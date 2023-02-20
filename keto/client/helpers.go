//source & inspiration: https://github.com/ory/keto/blob/6c0e1ba87f4d3a355cebd0ea77f28319be2dd606/cmd/relationtuple/output.go

package client

import (
	"encoding/json"
	"sort"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
)

type (
	Collection struct {
		relations []*rts.RelationTuple
	}
	OutputTuple struct {
		*rts.RelationTuple
	}
)

func NewCollection(rels []*rts.RelationTuple) (*Collection, error) {
	r := &Collection{relations: rels}
	//for i, rel := range rels {
	//	var err error
	//	r.relations[i], err = &rts.RelationTuple{}).FromDataProvider(rel)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	return r, nil
}

func NewAPICollection(rels []*rts.RelationTuple) *Collection {
	return &Collection{relations: rels}
}

func (r *Collection) Header() []string {
	return []string{
		"NAMESPACE",
		"OBJECT",
		"RELATION NAME",
		"SUBJECT",
	}
}

func (r *Collection) Table() [][]string {
	ir := r.relations

	data := make([][]string, len(ir))
	for i, rel := range ir {
		var sub string
		if rel.Subject != nil {
			sub = rel.Subject.String()
		} else {
			sub = ""
		}

		data[i] = []string{rel.Namespace, rel.Object, rel.Relation, sub}
	}

	return data
}

func (r *Collection) Interface() interface{} {
	return r.relations
}

func (r *Collection) MarshalJSON() ([]byte, error) {
	ir := r.relations
	return json.Marshal(ir)
}

func (r *Collection) UnmarshalJSON(raw []byte) error {
	return json.Unmarshal(raw, &r.relations)
}

func (r *Collection) Len() int {
	return len(r.relations)
}

func (r *Collection) IDs() []string {
	ts := r.relations
	ids := make([]string, len(ts))
	for i, rt := range ts {
		ids[i] = rt.String()
	}
	return ids
}

func (r *OutputTuple) Header() []string {
	return []string{
		"NAMESPACE",
		"OBJECT ID",
		"RELATION NAME",
		"SUBJECT",
	}
}

func (r *OutputTuple) Columns() []string {
	return []string{
		r.Namespace,
		r.Object,
		r.Relation,
		r.Subject.String(),
	}
}

// function to sort the relation tuples by namespace, object, relation, subject
func SortRelationTuples(relations []*rts.RelationTuple) {
	sort.Slice(relations, func(i, j int) bool {
		if relations[i].Namespace < relations[j].Namespace {
			return true
		}
		if relations[i].Namespace > relations[j].Namespace {
			return false
		}
		if relations[i].Object < relations[j].Object {
			return true
		}
		if relations[i].Object > relations[j].Object {
			return false
		}
		if relations[i].Relation < relations[j].Relation {
			return true
		}
		if relations[i].Relation > relations[j].Relation {
			return false
		}
		if relations[i].Subject == nil {
			return true
		}
		if relations[j].Subject == nil {
			return false
		}
		return relations[i].Subject.String() < relations[j].Subject.String()
	})
}

// function to sort a list of SubjectTrees by nodetype, subject, children and tuple
func SortSubjectTrees(subjectTrees []*rts.SubjectTree) {
	sort.Slice(subjectTrees, func(i, j int) bool {
		if subjectTrees[i].NodeType < subjectTrees[j].NodeType {
			return true
		}
		if subjectTrees[i].NodeType > subjectTrees[j].NodeType {
			return false
		}
		if subjectTrees[i].Subject == nil {
			return true
		}
		if subjectTrees[j].Subject == nil {
			return false
		}
		if subjectTrees[i].Children == nil {
			return true
		}
		if subjectTrees[j].Children == nil {
			return false
		}
		if len(subjectTrees[i].Children) < len(subjectTrees[j].Children) {
			return true
		}
		if len(subjectTrees[i].Children) > len(subjectTrees[j].Children) {
			return false
		}
		if subjectTrees[i].Tuple == nil {
			return true
		}
		if subjectTrees[j].Tuple == nil {
			return false
		}
		return subjectTrees[i].Tuple.String() < subjectTrees[j].Tuple.String()
	})
}
