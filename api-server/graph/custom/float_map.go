package custom

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalFloatMap(val map[string]*float64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			panic(err)
		}
	})
}

// TODO: this unmarshaler is not working and is needed for input types
func UnmarshalFloatMap(v interface{}) (map[string]*float64, error) {
	if m, ok := v.(map[string]*float64); ok {
		return m, nil
	}
	return nil, fmt.Errorf("%T is not a map", v)
}
