package format

import "encoding/json"

// function that converts an interface{} to a map[string]interface{}
func ToMapStringInterface(i interface{}) (map[string]interface{}, error) {
	var output map[string]interface{}
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
