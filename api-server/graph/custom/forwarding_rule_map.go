package custom

import (
	"encoding/json"
	"fmt"
	"io"

	observabilityv1alpha1 "github.com/pluralsh/trace-shield-controller/api/observability/v1alpha1"

	"github.com/99designs/gqlgen/graphql"
	// "github.com/pluralsh/oauth-playground/api-server/graph/model"
)

func MarshalForwardingRuleMap(val map[string]observabilityv1alpha1.ForwardingRule) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalForwardingRuleMap(v interface{}) (map[string]observabilityv1alpha1.ForwardingRule, error) {
	if m, ok := v.(map[string]observabilityv1alpha1.ForwardingRule); ok {
		return m, nil
	}
	return nil, fmt.Errorf("%T is not a map", v)
}
