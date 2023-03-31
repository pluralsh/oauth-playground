package clients

import (
	"context"
	"fmt"
	"net/http"

	hydra "github.com/ory/hydra-client-go/v2"
	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/api-server/format"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

type HydraOperation string

const (
	HydraOperationCreate HydraOperation = "create"
	HydraOperationUpdate HydraOperation = "update"
)

// ListOAuth2Clients is the resolver for the listOAuth2Clients field.
func (c *ClientWrapper) ListOAuth2Clients(ctx context.Context) ([]*model.OAuth2Client, error) {
	log := c.Log.WithName("ListOAuth2Clients")
	clients, resp, err := c.HydraClient.OAuth2Api.ListOAuth2Clients(ctx).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to list oauth2 clients")
		return nil, fmt.Errorf("failed to list oauth2 clients: %w", err)
	}
	var output []*model.OAuth2Client
	for _, client := range clients {
		output = append(output, format.HydraOAuth2ClientToGraphQL(client))
	}
	return output, nil
}

// GetOAuth2Client is the resolver for the listOAuth2Clients field.
func (c *ClientWrapper) GetOAuth2Client(ctx context.Context, id string) (*model.OAuth2Client, error) {
	log := c.Log.WithName("ListOAuth2Clients")
	client, resp, err := c.HydraClient.OAuth2Api.GetOAuth2Client(ctx, id).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to list oauth2 clients")
		return nil, fmt.Errorf("failed to list oauth2 clients: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("client not found")
	}
	return format.HydraOAuth2ClientToGraphQL(*client), nil
}

// function that gets all users that can login to the oauth2 client
func (c *ClientWrapper) GetOAuth2ClientUserLoginBindings(ctx context.Context, obj *model.LoginBindings) ([]*model.User, error) {
	log := c.Log.WithName("GetOAuth2ClientUserLoginBindings")

	var output []*model.User

	for _, inUser := range obj.Users {
		user, err := c.GetUserFromId(ctx, inUser.ID)
		if err != nil {
			log.Error(err, "failed to get user", "ID", inUser.ID)
			continue
		}
		output = append(output, user)
	}
	return output, nil
}

// function that gets all groups that can login to the oauth2 client
func (c *ClientWrapper) GetOAuth2ClientGroupLoginBindings(ctx context.Context, obj *model.LoginBindings) ([]*model.Group, error) {
	log := c.Log.WithName("GetOAuth2ClientUserLoginBindings")

	var output []*model.Group

	for _, inGroup := range obj.Groups {
		group, err := c.GetGroupFromName(ctx, inGroup.Name)
		if err != nil {
			log.Error(err, "failed to get group", "Name", inGroup.Name)
			continue
		}
		output = append(output, group)
	}
	return output, nil
}

// function that gets all login bindings for an oauth2 client from keto
func (c *ClientWrapper) GetOAuth2ClientLoginBindings(ctx context.Context, id string) (*model.LoginBindings, error) {
	log := c.Log.WithName("GetOAuth2ClientLoginBindings").WithValues("ID", id)

	query := rts.RelationQuery{
		Namespace: px.Ptr("OAuth2Client"),
		Object:    px.Ptr(id),
		Relation:  px.Ptr("login"),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputLoginBindings *model.LoginBindings
	var outputUserBindingIds []*model.User
	var outputGroupBindingNames []*model.Group

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			outputUserBindingIds = append(outputUserBindingIds, &model.User{ID: subjectSet.Object})
		} else if subjectSet.Namespace == "Group" && subjectSet.Object != "" {
			outputGroupBindingNames = append(outputGroupBindingNames, &model.Group{Name: subjectSet.Object})
		} else {
			continue
		}

	}

	if len(outputUserBindingIds) > 0 || len(outputGroupBindingNames) > 0 {
		outputLoginBindings = &model.LoginBindings{
			Users:  outputUserBindingIds,
			Groups: outputGroupBindingNames,
		}
	}

	log.Info("Success getting group members in keto")
	return outputLoginBindings, nil
}

func (c *ClientWrapper) CreateOAuth2Client(ctx context.Context, mode HydraOperation, allowedCorsOrigins []string, audience []string, authorizationCodeGrantAccessTokenLifespan *string, authorizationCodeGrantIDTokenLifespan *string, authorizationCodeGrantRefreshTokenLifespan *string, backChannelLogoutSessionRequired *bool, backChannelLogoutURI *string, clientCredentialsGrantAccessTokenLifespan *string, clientID *string, clientName *string, clientSecret *string, clientSecretExpiresAt *int64, clientURI *string, contacts []string, frontchannelLogoutSessionRequired *bool, frontchannelLogoutURI *string, grantTypes []string, implicitGrantAccessTokenLifespan *string, implicitGrantIDTokenLifespan *string, jwks map[string]interface{}, jwksURI *string, jwtBearerGrantAccessTokenLifespan *string, logoURI *string, metadata map[string]interface{}, policyURI *string, postLogoutRedirectUris []string, redirectUris []string, responseTypes []string, scope *string, sectorIdentifierURI *string, subjectType *string, tokenEndpointAuthMethod *string, tokenEndpointAuthSigningAlgorithm *string, tosURI *string, userinfoSignedResponseAlgorithm *string, loginBindings *model.LoginBindingsInput) (*model.OAuth2Client, error) {
	log := c.Log.WithName("CreateOAuth2Client").WithValues("Name", clientName, "ID", clientID, "Mode", mode)

	client := format.GraphQLNewOAuth2ClientToHydra(allowedCorsOrigins, audience, authorizationCodeGrantAccessTokenLifespan, authorizationCodeGrantIDTokenLifespan, authorizationCodeGrantRefreshTokenLifespan, backChannelLogoutSessionRequired, backChannelLogoutURI, clientCredentialsGrantAccessTokenLifespan, clientID, clientName, clientSecret, clientSecretExpiresAt, clientURI, contacts, frontchannelLogoutSessionRequired, frontchannelLogoutURI, grantTypes, implicitGrantAccessTokenLifespan, implicitGrantIDTokenLifespan, jwks, jwksURI, jwtBearerGrantAccessTokenLifespan, logoURI, metadata, policyURI, postLogoutRedirectUris, redirectUris, responseTypes, scope, sectorIdentifierURI, subjectType, tokenEndpointAuthMethod, tokenEndpointAuthSigningAlgorithm, tosURI, userinfoSignedResponseAlgorithm, loginBindings)

	// TODO: add bound observability tenant as owner of the client so we can ensure that a client is only bound to one tenant

	var createdClient *hydra.OAuth2Client
	var resp *http.Response
	var err error

	if mode == HydraOperationCreate {
		createdClient, resp, err = c.HydraClient.OAuth2Api.CreateOAuth2Client(ctx).OAuth2Client(client).Execute()
		if err != nil || resp.StatusCode != 201 {
			log.Error(err, "failed to create oauth2 client")
			return nil, fmt.Errorf("failed to create oauth2 client: %w", err)
		}
	} else if mode == HydraOperationUpdate {
		if clientID == nil || *clientID == "" {
			return nil, fmt.Errorf("clientID is required for update")
		}
		createdClient, resp, err = c.HydraClient.OAuth2Api.SetOAuth2Client(ctx, *clientID).OAuth2Client(client).Execute()
		if err != nil || resp.StatusCode != 200 {
			log.Error(err, "failed to update oauth2 client")
			return nil, fmt.Errorf("failed to update oauth2 client: %w", err)
		}
	}

	if createdClient == nil {
		return nil, fmt.Errorf("failed to create oauth2 client")
	}

	exist, err := c.OAuth2ClientExistsInKeto(ctx, *createdClient.ClientId)
	if err != nil {
		log.Error(err, "failed to check if oauth2 client exists in keto")
	}
	if !exist {
		err = c.CreateOAuth2ClientInKeto(ctx, *createdClient.ClientId)
		if err != nil {
			log.Error(err, "failed to create oauth2 client in keto")
		}
	}

	usersToAdd, usersToRemove, groupsToAdd, groupsToRemove, err := c.LoginBindingsChangeset(ctx, *createdClient.ClientId, loginBindings)
	if err != nil {
		log.Error(err, "Failed to get oauth2 client changeset")
		return nil, err
	}

	if err := c.AddUsersToLoginBindings(ctx, *createdClient.ClientId, usersToAdd); err != nil {
		log.Error(err, "Failed to add users to oauth2 client bindings in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveUsersFromLoginBindings(ctx, *createdClient.ClientId, usersToRemove); err != nil {
		log.Error(err, "Failed to remove users from oauth2 client bindings in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.AddGroupsToLoginBindings(ctx, *createdClient.ClientId, groupsToAdd); err != nil {
		log.Error(err, "Failed to add groups to oauth2 client bindings in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveGroupsFromLoginBindings(ctx, *createdClient.ClientId, groupsToRemove); err != nil {
		log.Error(err, "Failed to remove groups from oauth2 client bindings in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	log.Info("Success creating oauth2 client in hydra")

	output := format.HydraOAuth2ClientToGraphQL(*createdClient)

	output.Organization = &model.Organization{
		Name: "main", //TODO: decide whether to hardcode this or not
	}

	return output, nil
}

// function that deletes an oauth2 client
func (c *ClientWrapper) DeleteOAuth2Client(ctx context.Context, clientID string) (*model.OAuth2Client, error) {
	log := c.Log.WithName("DeleteOAuth2Client").WithValues("ID", clientID)

	resp, err := c.HydraClient.OAuth2Api.DeleteOAuth2Client(ctx, clientID).Execute()
	if err != nil || resp.StatusCode != 204 {
		log.Error(err, "failed to delete oauth2 client")
		return nil, fmt.Errorf("failed to delete oauth2 client: %w", err)
	}

	exist, err := c.OAuth2ClientExistsInKeto(ctx, clientID)
	if err != nil {
		log.Error(err, "failed to check if oauth2 client exists in keto")
	}
	if exist {
		err = c.DeleteOAuth2ClientInKeto(ctx, clientID)
		if err != nil {
			log.Error(err, "failed to delete oauth2 client in keto")
		}
	}

	log.Info("Success deleting oauth2 client")

	return &model.OAuth2Client{
		ClientID: px.Ptr(clientID),
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that checks if an oauth2 client exists in keto
func (c *ClientWrapper) OAuth2ClientExistsInKeto(ctx context.Context, id string) (bool, error) {
	log := c.Log.WithName("OAuth2ClientExistsInKeto").WithValues("ID", id)

	query := rts.RelationQuery{
		Namespace: px.Ptr("OAuth2Client"),
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

// function that creates an oauth2client in keto
func (c *ClientWrapper) CreateOAuth2ClientInKeto(ctx context.Context, id string) error {
	log := c.Log.WithName("CreateOAuth2ClientInKeto").WithValues("ID", id)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    id,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating oauth2client in keto")
	return nil
}

// function that deletes an oauth2client in keto
func (c *ClientWrapper) DeleteOAuth2ClientInKeto(ctx context.Context, id string) error {
	log := c.Log.WithName("DeleteOAuth2ClientInKeto").WithValues("ID", id)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    id,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success deleting oauth2client in keto")
	return nil
}

// function that determines which users or groups to add or remove from the login bindings of an oauth2 client
func (c *ClientWrapper) LoginBindingsChangeset(ctx context.Context, clientId string, bindings *model.LoginBindingsInput) (usersToAdd []string, usersToRemove []string, groupsToAdd []string, groupsToRemove []string, err error) {
	currentUsers, currentGroups, err := c.GetLoginBindingsInKeto(ctx, clientId) // TODO: replace with GetOAuth2ClientLoginBindings

	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to get current members: %w", err)
	}

	for _, user := range bindings.Users {
		if !stringContains(currentUsers, user) {
			usersToAdd = append(usersToAdd, user)
		}
	}

	for _, user := range currentUsers {
		if !stringContains(bindings.Users, user) {
			usersToRemove = append(usersToRemove, user)
		}
	}

	for _, group := range bindings.Groups {
		if !stringContains(currentGroups, group) {
			groupsToAdd = append(groupsToAdd, group)
		}
	}

	for _, group := range currentGroups {
		if !stringContains(bindings.Groups, group) {
			groupsToRemove = append(groupsToRemove, group)
		}
	}

	return usersToAdd, usersToRemove, groupsToAdd, groupsToRemove, nil
}

// function that adds users to the login bindings of an oauth2 client
func (c *ClientWrapper) AddUsersToLoginBindings(ctx context.Context, clientID string, users []string) error {
	log := c.Log.WithName("AddUsersToLoginBindings").WithValues("ClientID", clientID, "Users", users)

	for _, user := range users {
		err := c.AddUserToLoginBindings(ctx, clientID, user)
		if err != nil {
			log.Error(err, "Failed to add user to login bindings")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success adding users to login bindings")
	return nil
}

// function that adds a user to the login bindings of an oauth2 client
func (c *ClientWrapper) AddUserToLoginBindings(ctx context.Context, clientID string, user string) error {
	log := c.Log.WithName("AddUserToLoginBindings").WithValues("ClientID", clientID, "User", user)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    clientID,
		Relation:  "login",
		Subject: rts.NewSubjectSet(
			"User",
			user,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success adding user to login bindings")
	return nil
}

// function that removes users from the login bindings of an oauth2 client
func (c *ClientWrapper) RemoveUsersFromLoginBindings(ctx context.Context, clientID string, users []string) error {
	log := c.Log.WithName("RemoveUsersFromLoginBindings").WithValues("ClientID", clientID, "Users", users)

	for _, user := range users {
		err := c.RemoveUserFromLoginBindings(ctx, clientID, user)
		if err != nil {
			log.Error(err, "Failed to remove user from login bindings")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success removing users from login bindings")
	return nil
}

// function that removes a user from the login bindings of an oauth2 client
func (c *ClientWrapper) RemoveUserFromLoginBindings(ctx context.Context, clientID string, user string) error {
	log := c.Log.WithName("RemoveUserFromLoginBindings").WithValues("ClientID", clientID, "User", user)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    clientID,
		Relation:  "login",
		Subject: rts.NewSubjectSet(
			"User",
			user,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success removing user from login bindings")
	return nil
}

// function that adds groups to the login bindings of an oauth2 client
func (c *ClientWrapper) AddGroupsToLoginBindings(ctx context.Context, clientID string, groups []string) error {
	log := c.Log.WithName("AddGroupsToLoginBindings").WithValues("ClientID", clientID, "Groups", groups)

	for _, group := range groups {
		err := c.AddGroupToLoginBindings(ctx, clientID, group)
		if err != nil {
			log.Error(err, "Failed to add group to login bindings")
			return err
		}
	}

	log.Info("Success adding groups to login bindings")
	return nil
}

// function that adds a group to the login bindings of an oauth2 client
func (c *ClientWrapper) AddGroupToLoginBindings(ctx context.Context, clientID string, group string) error {
	log := c.Log.WithName("AddGroupToLoginBindings").WithValues("ClientID", clientID, "Group", group)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    clientID,
		Relation:  "login",
		Subject: rts.NewSubjectSet(
			"Group",
			group,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success adding group to login bindings")
	return nil
}

// function that removes groups from the login bindings of an oauth2 client
func (c *ClientWrapper) RemoveGroupsFromLoginBindings(ctx context.Context, clientID string, groups []string) error {
	log := c.Log.WithName("RemoveGroupsFromLoginBindings").WithValues("ClientID", clientID, "Groups", groups)

	for _, group := range groups {
		err := c.RemoveGroupFromLoginBindings(ctx, clientID, group)
		if err != nil {
			log.Error(err, "Failed to remove group from login bindings")
			return err
		}
	}

	log.Info("Success removing groups from login bindings")
	return nil
}

// function that removes a group from the login bindings of an oauth2 client
func (c *ClientWrapper) RemoveGroupFromLoginBindings(ctx context.Context, clientID string, group string) error {
	log := c.Log.WithName("RemoveGroupFromLoginBindings").WithValues("ClientID", clientID, "Group", group)

	clientTuple := &rts.RelationTuple{
		Namespace: "OAuth2Client",
		Object:    clientID,
		Relation:  "login",
		Subject: rts.NewSubjectSet(
			"Group",
			group,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, clientTuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success removing group from login bindings")
	return nil
}

// function that get the current users and groups in the login bindings of an oauth2 client
func (c *ClientWrapper) GetLoginBindingsInKeto(ctx context.Context, clientID string) ([]string, []string, error) {
	log := c.Log.WithName("GetLoginBindingsInKeto").WithValues("ClientID", clientID)

	// TODO: remove since it is a duplicate of GetOAuth2ClientLoginBindings

	query := rts.RelationQuery{
		Namespace: px.Ptr("OAuth2Client"),
		Object:    px.Ptr(clientID),
		Relation:  px.Ptr("login"),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var users []string
	var groups []string
	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			users = append(users, subjectSet.Object)
		} else if subjectSet.Namespace == "Group" && subjectSet.Object != "" {
			groups = append(groups, subjectSet.Object)
		} else {
			continue
		}
	}

	log.Info("Success getting users in login bindings")
	return users, groups, nil
}
