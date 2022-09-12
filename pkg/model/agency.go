package model

import "encoding/json"

type Agency struct {
	Address   string
	Port      string
	Anonymous string
	Type      string
	Location  string
}

func (t Agency) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t Agency) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}
