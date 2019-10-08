package parameters

import (
	"encoding/json"
	"sync/atomic"
)

type Int8Type struct {
	value atomic.Value
}

type Int8JSON struct {
	Value uint8 `json:"value"`
}

func NewInt8() *Int8Type {
	self := Int8Type{}
	self.value.Store(uint8(0))
	return &self
}

func (i *Int8Type) Type() string {
	return UInt8
}

func (i *Int8Type) Get() interface{} {
	return i.value.Load().(uint8)
}

func (i *Int8Type) Set(value interface{}) {
	i.value.Store(value.(*uint8))
}

func (i *Int8Type) MarshalJSON() ([]byte, error) {
	data := Int8JSON{
		Value: i.value.Load().(uint8),
	}

	return json.Marshal(data)
}

func (i *Int8Type) UnmarshalJSON(data []byte) error {
	input := Int8JSON{}

	err := json.Unmarshal(data, &input)
	if err != nil {
		return err
	}

	i.Set(input.Value)

	return nil
}