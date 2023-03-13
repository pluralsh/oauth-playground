package clients

import (
	"context"
	"encoding/json"
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	kratos "github.com/ory/kratos-client-go"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

// function that gets a user from the Kratos API
func (c *ClientWrapper) GetUserFromId(ctx context.Context, id string) (*model.User, error) {
	log := c.Log.WithName("User").WithValues("ID", id)

	user, resp, err := c.KratosClient.IdentityApi.GetIdentity(ctx, id).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to get user")
		return nil, err
	}

	outUser, err := c.UnmarshalUserTraits(user)
	if err != nil {
		log.Error(err, "Error when unmarshalling user")
		return nil, err
	}

	log.Info("Success getting User")

	return outUser, nil
}

// function that will list all users using the kratos api
func (c *ClientWrapper) ListUsers(ctx context.Context) ([]*model.User, error) {
	log := c.Log.WithName("ListUsers")
	users, resp, err := c.KratosClient.IdentityApi.ListIdentities(ctx).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var output []*model.User

	for _, user := range users {

		user, err := c.UnmarshalUserTraits(&user)
		if err != nil {
			log.Error(err, "Error when unmarshalling user", "ID", user.ID)
			continue
		}

		output = append(output, user)

	}
	return output, nil
}

// UnmarshalUserTraits unmarshals the user traits from the Kratos Identity.
// It expects that the user traits are in the format of the model.User struct.

func (c *ClientWrapper) UnmarshalUserTraits(user *kratos.Identity) (*model.User, error) {
	log := c.Log.WithName("UnmarshalUserTraits")

	outUser := &model.User{}

	byteData, err := json.Marshal(user.Traits)
	if err != nil {
		log.Error(err, "Error when marshalling user traits")
		return nil, err
	}

	err = json.Unmarshal(byteData, outUser)
	if err != nil {
		log.Error(err, "Error when unmarshalling user traits")
		return nil, err
	}

	// the unmarshal function does not set the ID
	outUser.ID = user.Id

	return outUser, nil
}

func (c *ClientWrapper) CreateUser(ctx context.Context, email string, name *model.NameInput) (*model.User, error) {
	log := c.Log.WithName("CreateUser").WithValues("Email", email)

	// test := "password12345"

	traits := map[string]interface{}{
		"email": email,
	}

	if name != nil {
		traits["name"] = make(map[string]interface{})
		if name.First != nil {
			traits["name"].(map[string]interface{})["first"] = *name.First
		}
		if name.Last != nil {
			traits["name"].(map[string]interface{})["last"] = *name.Last
		}
	}

	kratosUser, resp, err := c.KratosClient.IdentityApi.CreateIdentity(ctx).CreateIdentityBody(
		kratos.CreateIdentityBody{
			SchemaId: "person",
			Traits:   traits,
			// VerifiableAddresses: []kratos.VerifiableIdentityAddress{
			// 	{
			// 		Value: email,
			// 		Via:   "email",
			// 	},
			// },
			// Credentials: &kratos.IdentityWithCredentials{
			// 	Password: &kratos.IdentityWithCredentialsPassword{
			// 		Config: &kratos.IdentityWithCredentialsPasswordConfig{
			// 			Password: &test,
			// 		},
			// 	},
			// },
		},
	).Execute()

	if err != nil || resp.StatusCode != 201 {
		log.Error(err, "failed to create user")
		return nil, err
	}

	// TODO: this is giving a 404 not found error
	// recoveryLink, err := c.CreateRecoveryLinkForIdentity(ctx, kratosUser.Id)
	// if err != nil {
	// 	log.Error(err, "failed to create recovery link for user")
	// }

	exits, err := c.UserExistsInKeto(ctx, kratosUser.Id)
	if err != nil {
		log.Error(err, "failed to check if user exists in keto")
		// return nil, err
	}

	if !exits {
		err = c.CreateUserInKeto(ctx, kratosUser.Id)
		if err != nil {
			log.Error(err, "failed to create user in keto")
			return nil, err
		}
	}

	outUser, err := c.UnmarshalUserTraits(kratosUser)
	if err != nil {
		log.Error(err, "failed to unmarshal user traits")
		return nil, err
	}

	// TODO: reenable once it is working
	// outUser.RecoveryLink = recoveryLink

	log.Info("Success creating User")
	return outUser, nil
}

func (c *ClientWrapper) DeleteUser(ctx context.Context, id string) (*model.User, error) {
	log := c.Log.WithName("DeleteUser").WithValues("ID", id)

	resp, err := c.KratosClient.IdentityApi.DeleteIdentity(ctx, id).Execute()

	if err != nil || resp.StatusCode != 204 {
		log.Error(err, "failed to delete user")
		return nil, err
	}

	exits, err := c.UserExistsInKeto(ctx, id)
	if err != nil {
		log.Error(err, "failed to check if user exists in keto")
		// return nil, err
	}

	if exits {
		err = c.DeleteUserInKeto(ctx, id)
		if err != nil {
			log.Error(err, "failed to delete user in keto")
			return nil, err
		}
	}

	log.Info("Success deleting User")
	return &model.User{
		ID:    id,
		Email: "deleted",
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that checks if a user exists in keto
func (c *ClientWrapper) UserExistsInKeto(ctx context.Context, id string) (bool, error) {
	log := c.Log.WithName("UserExistsInKeto").WithValues("ID", id)

	query := rts.RelationQuery{
		Namespace: px.Ptr("User"),
		Object:    px.Ptr(id),
		Relation:  px.Ptr("organizations"),
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		log.Error(err, "Failed to query tuples")
		return false, fmt.Errorf("failed to query tuples: %w", err)
	}

	if len(respTuples) == 0 {
		return false, nil
	}

	return true, nil
}

// function that create a recovery link for a user
func (c *ClientWrapper) CreateRecoveryLinkForIdentity(ctx context.Context, id string) (*string, error) {
	log := c.Log.WithName("CreateRecoveryLinkForIdentity").WithValues("ID", id)

	link, resp, err := c.KratosClient.IdentityApi.CreateRecoveryLinkForIdentity(ctx).CreateRecoveryLinkForIdentityBody(
		kratos.CreateRecoveryLinkForIdentityBody{
			IdentityId: id,
		},
	).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to create recovery link for identity")
		return nil, err
	}

	log.Info("Success creating recovery link for identity")
	return &link.RecoveryLink, nil
}

// function that creates a user in keto
func (c *ClientWrapper) CreateUserInKeto(ctx context.Context, id string) error {
	log := c.Log.WithName("CreateUserInKeto").WithValues("ID", id)

	userTuple := &rts.RelationTuple{
		Namespace: "User",
		Object:    id,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, userTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating user in keto")
	return nil
}

// function that deletes a user from keto
func (c *ClientWrapper) DeleteUserInKeto(ctx context.Context, id string) error {
	log := c.Log.WithName("DeleteUserInKeto").WithValues("ID", id)

	userTuple := &rts.RelationTuple{
		Namespace: "User",
		Object:    id,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, userTuple)
	if err != nil {
		log.Error(err, "failed to delete tuple")
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success deleting user in keto")
	return nil
}

// function that will get all the groups a user is in
func (c *ClientWrapper) GetUserGroups(ctx context.Context, id string) ([]*model.Group, error) {
	log := c.Log.WithName("GetUserGroups").WithValues("ID", id)

	respTuples, err := c.KetoClient.QueryAllTuples(ctx, &rts.RelationQuery{
		Namespace: px.Ptr("Group"),
		Relation:  px.Ptr("members"),
		Subject:   rts.NewSubjectSet("User", id, ""),
	}, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get tuples: %w", err)
	}

	var groups []*model.Group

	for _, tuple := range respTuples {
		group, err := c.GetGroupFromName(ctx, tuple.Object)
		if err != nil {
			log.Error(err, "failed to get group", "Name", tuple.Object)
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}
