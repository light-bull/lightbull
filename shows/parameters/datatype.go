package parameters

// DataType is a data type used for a parameter
type DataType interface {
	// Type returns a identifier like "color" or "bool"
	Type() string

	// Get the data
	Get() interface{}

	// Set the data
	Set(interface{})

	// ToJSON returns the data serialized as JSON
	ToJSON() []byte

	// UpdateFromJSON loads the data from the JSON string
	UpdateFromJSON([]byte) error
}
