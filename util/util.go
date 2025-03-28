package util

import (
	"encoding/json"
	"strconv"
)

// money represented in cents
const dollar_to_cents uint8 = 100

type StringUInt32 uint32

func (s *StringUInt32) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return err
	}
	*s = StringUInt32(uint32(value * float64(dollar_to_cents)))
	return nil
}

type StringUInt64 uint64

func (s *StringUInt64) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*s = StringUInt64(uint64(value))
	return nil
}

type StringFloat32 float32

func (s *StringFloat32) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	*s = StringFloat32(value)
	return nil
}
