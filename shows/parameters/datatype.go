package parameters

// DataType is a data type used for a parameter
type DataType interface {
	// Type returns a identifier like "color" or "bool"
	Type() string

	// Get the data
	Get() interface{}

	// Set the data
	Set(interface{})

	// MarshalJSON returns the data serialized as JSON (it implements the `json.Marshaller` interface)
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON loads the data from the JSON string (it implements the `json.Unmarshaller` interface)
	UnmarshalJSON(data []byte) error
}
