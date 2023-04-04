package custom

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MarshalDuration(d metav1.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(d.Duration.String()))
	})
}

func UnmarshalDuration(v interface{}) (metav1.Duration, error) {
	switch v := v.(type) {
	case string:
		duration, err := model.ParseDuration(v)
		return metav1.Duration{Duration: time.Duration(duration)}, err
	default:
		return metav1.Duration{Duration: 0}, fmt.Errorf("%T is not a string", v)
	}
}
