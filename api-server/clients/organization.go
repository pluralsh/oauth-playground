package clients

import (
	"context"

	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

func (c *ClientWrapper) UpdateOrganization(ctx context.Context, name string) (*model.Organization, error) {

	// TODO: figure out which admins to add or remove
	// TODO: create separate functions for adding and removing admins from an organization
	// log := c.Log.WithName("Organization").WithValues("Name", name)

	// var adminTuples []*rts.RelationTuple

	// for _, admin := range admins {
	// 	adminTuples = append(adminTuples, &rts.RelationTuple{
	// 		Namespace: "Organization",
	// 		Object:    org.ID.String(),
	// 		Relation:  "admins",
	// 		Subject: rts.NewSubjectSet(
	// 			"User",
	// 			admin.ID,
	// 			"",
	// 		),
	// 	})
	// }

	// err = r.C.KetoClient.CreateTuples(ctx, adminTuples)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create tuple: %w", err)
	// }

	// query := rts.RelationQuery{
	// 	Namespace: px.Ptr("Organization"),
	// 	Object:    px.Ptr(org.ID.String()),
	// 	Relation:  px.Ptr("admins"),
	// 	Subject:   nil,
	// }

	// respTuples, err := r.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query tuples: %w", err)
	// }

	// var outputAdmins []*model.User

	// for _, tuple := range respTuples {
	// 	subjectSet := tuple.Subject.GetSet()
	// 	if subjectSet.Namespace == "User" {
	// 		user, err := r.C.GetUserFromId(ctx, subjectSet.Object)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		outputAdmins = append(outputAdmins, user)
	// 	} else {
	// 		continue
	// 	}

	// }

	// return &model.Organization{
	// 	// ID:     org.ID.String(),
	// 	Name:   &org.Name,
	// 	Admins: outputAdmins,
	// }, nil
	return nil, nil
}
