package meta

import "encoding/json"

type Meta struct {
	Total   int
	Removed int
	Limit   int
	Offset  int
}

func (m Meta) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
