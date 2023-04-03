package custom

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/prometheus/common/model"
)

func MarshalDuration(d model.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, d.String())
	})
}

func UnmarshalDuration(v interface{}) (model.Duration, error) {
	switch v := v.(type) {
	case string:
		return model.ParseDuration(v)
	default:
		return 0, fmt.Errorf("%T is not a string", v)
	}
}
