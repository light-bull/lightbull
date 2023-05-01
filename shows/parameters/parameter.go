package parameters

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

// Parameter is an effect parameter
type Parameter struct {
	// ID is the globally unique UUID for this parameter
	ID uuid.UUID

	// Key is the id that is unique for a single effect
	Key string

	// Name is the nice name for the UI
	Name string

	// current and default value of parameter, they have to be the same DataType
	cur DataType
	def DataType

	// linked parameters that will change the value together with this one
	// warning: can have loops, you have to check this when iterating the links
	linkedParameters                 []*Parameter
	linkedParametersDuringUnmarshall []uuid.UUID

	// FIXME: mux!
}

// NewParameter returns a new parameter of the specified data type (or nil)
func NewParameter(key string, datatype string, name string) *Parameter {
	parameter := Parameter{}

	parameter.ID = uuid.New() // TODO: make sure that unique
	parameter.Key = key
	parameter.Name = name

	if datatype == Color {
		parameter.cur = NewColor()
		parameter.def = NewColor()
	} else if datatype == Percent {
		parameter.cur = NewPercent()
		parameter.def = NewPercent()
	} else if datatype == IntegerGreaterOrEqualZero {
		parameter.cur = NewIntegerGreaterZero()
		parameter.def = NewIntegerGreaterZero()
	} else if datatype == Boolean {
		parameter.cur = NewBooleanType()
		parameter.def = NewBooleanType()
	} else {
		return nil
	}

	return &parameter
}

// MarshalJSON is there to implement the `json.Marshaller` interface.
func (parameter *Parameter) MarshalJSON() ([]byte, error) {
	type format struct {
		ID             uuid.UUID   `json:"id"`
		Key            string      `json:"key"`
		Name           string      `json:"name"` // will be ignored for deserialization
		Type           string      `json:"type"` // will be ignored for deserialization
		Default        DataType    `json:"default"`
		Current        DataType    `json:"current"`
		LinkParameters []uuid.UUID `json:"linkedParameters"`
	}

	data := format{
		ID:             parameter.ID,
		Key:            parameter.Key,
		Name:           parameter.Name,
		Type:           parameter.cur.Type(),
		Current:        parameter.cur,
		Default:        parameter.def,
		LinkParameters: make([]uuid.UUID, len(parameter.linkedParameters)),
	}

	for i, linkedParameter := range parameter.linkedParameters {
		data.LinkParameters[i] = linkedParameter.ID
	}

	return json.Marshal(data)
}

// UnmarshalJSON is there to implement the `json.Unmarshaller` interface.
func (parameter *Parameter) UnmarshalJSON(data []byte) error {
	type format struct {
		ID             uuid.UUID        `json:"id"`
		Key            string           `json:"key"`
		Current        *json.RawMessage `json:"current"`
		Default        *json.RawMessage `json:"default"`
		LinkParameters []uuid.UUID      `json:"linkedParameters"`
	}

	dataMap := format{}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	parameter.ID = dataMap.ID
	parameter.Key = dataMap.Key
	parameter.linkedParametersDuringUnmarshall = dataMap.LinkParameters

	if dataMap.Current != nil {
		err = parameter.SetFromJSON(*dataMap.Current)
		if err != nil {
			return err
		}
	}

	if dataMap.Default != nil {
		err = parameter.SetDefaultFromJSON(*dataMap.Default)
		if err != nil {
			return err
		}
	}

	return nil
}

// FillLinkedParametersAfterUnmarshall resolves the UUIDs of linked parameters and puts the real parameters into the list
func (parameter *Parameter) FillLinkedParametersAfterUnmarshall(mapping map[uuid.UUID]*Parameter) error {
	parameter.linkedParameters = make([]*Parameter, len(parameter.linkedParametersDuringUnmarshall))

	for i, parameterId := range parameter.linkedParametersDuringUnmarshall {
		linkedParameter, ok := mapping[parameterId]
		if !ok {
			return errors.New("unknown UUID in linkedParameters")
		}
		parameter.linkedParameters[i] = linkedParameter
	}

	return nil
}

// Get returns the currently set value
func (parameter *Parameter) Get() interface{} {
	return parameter.cur.Get()
}

// SetFromJSON sets a new value from the JSON data
func (parameter *Parameter) SetFromJSON(data []byte) error {
	err := parameter.cur.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	parameter.updateLinkedParameters()

	return nil
}

// SetDefaultFromJSON sets a new default value from the JSON data
func (parameter *Parameter) SetDefaultFromJSON(data []byte) error {
	err := parameter.def.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	parameter.updateLinkedParameters()

	return nil
}

// SetDefault sets the current value as default
// TODO: remove?
func (parameter *Parameter) SetDefault() {
	parameter.def.Set(parameter.cur.Get())

	parameter.updateLinkedParameters()
}

// RestoreDefault sets the current value back to the default value
// TODO: remove?
func (parameter *Parameter) RestoreDefault() {
	parameter.cur.Set(parameter.def.Get())

	parameter.updateLinkedParameters()
}

// AddLink adds a new link
// warning: no locking on parameter level - you have to take care on your own!
func (parameter *Parameter) AddLink(otherParameter *Parameter) {
	// check if parameter is already in list
	for _, link := range parameter.linkedParameters {
		if link.ID == otherParameter.ID {
			return
		}
	}

	// otherwise add it
	parameter.linkedParameters = append(parameter.linkedParameters, otherParameter)

	// make sure that all parameters have the same value
	parameter.updateLinkedParameters()
}

// DeleteLink removed a link between parameters
// warning: no locking on parameter level - you have to take care on your own!
func (parameter *Parameter) DeleteLink(otherParameter *Parameter) {
	for pos, cur := range parameter.linkedParameters {
		if otherParameter.ID == cur.ID {
			parameter.linkedParameters = append(parameter.linkedParameters[:pos], parameter.linkedParameters[pos+1:]...)
			break

		}
	}
}

// updateLinkedParameters updates all linked parameters (current + default value)
func (parameter *Parameter) updateLinkedParameters() {
	allLinks := parameter.getAllLinkedParameters()
	for _, link := range allLinks {
		link.cur.Set(parameter.cur.Get())
		link.def.Set(parameter.def.Get())
	}
}

// getAllLinkedParameters returns an unique list of all linked parameters
func (parameter *Parameter) getAllLinkedParameters() []*Parameter {
	// we use a map to temporary store all linked parameters as we can easily add elements and check if they are included
	queuedLinks := make(map[*Parameter]bool)
	allLinks := make(map[*Parameter]bool)

	// fill lists with the links that we have in this parameter
	for _, link := range parameter.linkedParameters {
		queuedLinks[link] = true
		allLinks[link] = true
	}

	// iterate queue until we do not find new parameters
	for len(queuedLinks) > 0 {
		for currentLink, _ := range queuedLinks {
			// handle the linked parameters of the current parameter
			for _, link := range currentLink.linkedParameters {
				// we already have the parameter in our list -> ignore it
				if _, found := allLinks[link]; found {
					continue
				}

				// new parameter found -> add it to queue and list of all parameters
				queuedLinks[link] = true
				allLinks[link] = true
			}

			// remove ffrom queue
			delete(queuedLinks, currentLink)
		}
	}

	// convert to output format
	result := make([]*Parameter, len(allLinks))
	i := 0
	for link, _ := range allLinks {
		result[i] = link
		i++
	}
	return result
}
