package clients

import (
	"context"
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

type ObservabilityTenantPermission string

const (
	ObservabilityTenantPermissionView ObservabilityTenantPermission = "viewers"
	ObservabilityTenantPermissionEdit ObservabilityTenantPermission = "editors"
)

func (c *ClientWrapper) MutateObservabilityTenant(ctx context.Context, name string, viewers *model.ObservabilityTenantViewersInput, editors *model.ObservabilityTenantEditorsInput) (*model.ObservabilityTenant, error) {

	// TODO: figure out which members to add or remove
	log := c.Log.WithName("ObservabilityTenant").WithValues("Name", name)

	// TODO: figure out how to distinguish between creating or updating a group
	// updating a group would require that we first check if it exists and if a user is allowed to update it
	// creating a group would require that we first check if it exists and if a user is allowed to create it

	tenantpExists, err := c.ObservabilityTenantExistsInKeto(ctx, name)
	if err != nil {
		log.Error(err, "Failed to check if observability tenant already exists in keto")
		return nil, err
	}

	if !tenantpExists {
		err := c.CreateObservabilityTenantInKeto(ctx, name)
		if err != nil {
			log.Error(err, "Failed to create observability tenant in keto")
			return nil, err
		}
	}

	viewUsersToAdd, viewUsersToRemove, viewGroupsToAdd, viewGroupsToRemove, viewClientsToAdd, viewClientsToRemove, err := c.OsTenantChangeset(ctx, name, viewers, nil, ObservabilityTenantPermissionView)
	if err != nil {
		log.Error(err, "Failed to get observability tenant changeset")
		return nil, err
	}

	if err := c.AddUsersToTenantInKeto(ctx, name, viewUsersToAdd, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to add users as viewers to observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveUsersFromTenantInKeto(ctx, name, viewUsersToRemove, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to remove users as viewers from observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.AddGroupsToTenantInKeto(ctx, name, viewGroupsToAdd, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to add groups as viewers to observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveGroupsFromTenantInKeto(ctx, name, viewGroupsToRemove, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to remove groups as viewers from observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.AddOauthClientsToTenantInKeto(ctx, name, viewClientsToAdd, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to add clients as viewers to observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveOauthClientsFromTenantInKeto(ctx, name, viewClientsToRemove, ObservabilityTenantPermissionView); err != nil {
		log.Error(err, "Failed to remove clients as viewers from observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	editUsersToAdd, editUsersToRemove, editGroupsToAdd, editGroupsToRemove, _, _, err := c.OsTenantChangeset(ctx, name, nil, editors, ObservabilityTenantPermissionEdit)
	if err != nil {
		log.Error(err, "Failed to get observability tenant changeset")
		return nil, err
	}

	if err := c.AddUsersToTenantInKeto(ctx, name, editUsersToAdd, ObservabilityTenantPermissionEdit); err != nil {
		log.Error(err, "Failed to add users as editors to observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveUsersFromTenantInKeto(ctx, name, editUsersToRemove, ObservabilityTenantPermissionEdit); err != nil {
		log.Error(err, "Failed to remove users as editors from observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.AddGroupsToTenantInKeto(ctx, name, editGroupsToAdd, ObservabilityTenantPermissionEdit); err != nil {
		log.Error(err, "Failed to add groups as editors to observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	if err := c.RemoveGroupsFromTenantInKeto(ctx, name, editGroupsToRemove, ObservabilityTenantPermissionEdit); err != nil {
		log.Error(err, "Failed to remove groups as editors from observability tenant in keto")
		// return nil, err // TODO: add some way to wrap errors
	}

	return &model.ObservabilityTenant{
		Name: name,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that checks if an observability tenant exists in keto
func (c *ClientWrapper) ObservabilityTenantExistsInKeto(ctx context.Context, name string) (bool, error) {
	log := c.Log.WithName("ObservabilityTenantExistsInKeto").WithValues("Name", name)

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    px.Ptr(name),
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

// function that creates an observability tenant in keto
func (c *ClientWrapper) CreateObservabilityTenantInKeto(ctx context.Context, name string) error {
	log := c.Log.WithName("CreateObservabilityTenantInKeto").WithValues("Name", name)

	tenantTuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, tenantTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating group in keto")
	return nil
}

// function that determines which users or groups to add or remove from the observability tenant of an oauth2 client
func (c *ClientWrapper) OsTenantChangeset(ctx context.Context, name string, viewers *model.ObservabilityTenantViewersInput, editors *model.ObservabilityTenantEditorsInput, permission ObservabilityTenantPermission) (usersToAdd []string, usersToRemove []string, groupsToAdd []string, groupsToRemove []string, clientsToAdd []string, clientsToRemove []string, err error) {
	var currentUsers []string
	var currentGroups []string
	var currentClients []string

	if permission == ObservabilityTenantPermissionView {
		currentUsers, currentGroups, currentClients, err = c.GetTenantViewersInKeto(ctx, name)
	} else if permission == ObservabilityTenantPermissionEdit {
		currentUsers, currentGroups, err = c.GetTenantEditorsInKeto(ctx, name)
	}

	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("failed to get current viewers: %w", err)
	}

	if permission == ObservabilityTenantPermissionView && viewers != nil {
		for _, user := range viewers.Users {
			if !stringContains(currentUsers, user) {
				usersToAdd = append(usersToAdd, user)
			}
		}

		for _, user := range currentUsers {
			if !stringContains(viewers.Users, user) {
				usersToRemove = append(usersToRemove, user)
			}
		}

		for _, group := range viewers.Groups {
			if !stringContains(currentGroups, group) {
				groupsToAdd = append(groupsToAdd, group)
			}
		}

		for _, group := range currentGroups {
			if !stringContains(viewers.Groups, group) {
				groupsToRemove = append(groupsToRemove, group)
			}
		}

		for _, client := range viewers.Oauth2Clients {
			if !stringContains(currentClients, client) {
				clientsToAdd = append(clientsToAdd, client)
			}
		}

		for _, client := range currentClients {
			if !stringContains(viewers.Oauth2Clients, client) {
				clientsToRemove = append(clientsToRemove, client)
			}
		}
	} else if permission == ObservabilityTenantPermissionEdit && editors != nil {
		for _, user := range editors.Users {
			if !stringContains(currentUsers, user) {
				usersToAdd = append(usersToAdd, user)
			}
		}

		for _, user := range currentUsers {
			if !stringContains(editors.Users, user) {
				usersToRemove = append(usersToRemove, user)
			}
		}

		for _, group := range editors.Groups {
			if !stringContains(currentGroups, group) {
				groupsToAdd = append(groupsToAdd, group)
			}
		}

		for _, group := range currentGroups {
			if !stringContains(editors.Groups, group) {
				groupsToRemove = append(groupsToRemove, group)
			}
		}
	}

	return usersToAdd, usersToRemove, groupsToAdd, groupsToRemove, clientsToAdd, clientsToRemove, nil
}

// function that gets the current viewers of an observability tenant from keto
func (c *ClientWrapper) GetTenantViewersInKeto(ctx context.Context, name string) (users []string, groups []string, clients []string, err error) {
	log := c.Log.WithName("GetTenantViewersInKeto").WithValues("Name", name)

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    px.Ptr(name),
		Relation:  px.Ptr(string(ObservabilityTenantPermissionView)),
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		log.Error(err, "Failed to query tuples")
		return nil, nil, nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			users = append(users, subjectSet.Object)
		} else if subjectSet.Namespace == "Group" && subjectSet.Object != "" {
			groups = append(groups, subjectSet.Object)
		} else if subjectSet.Namespace == "OAuth2Client" && subjectSet.Object != "" {
			clients = append(clients, subjectSet.Object)
		} else {
			continue
		}
	}

	return users, groups, clients, nil
}

// function that gets the current editors of an observability tenant from keto
func (c *ClientWrapper) GetTenantEditorsInKeto(ctx context.Context, name string) (users []string, groups []string, err error) {
	log := c.Log.WithName("GetTenantEditorsInKeto").WithValues("Name", name)

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    px.Ptr(name),
		Relation:  px.Ptr(string(ObservabilityTenantPermissionEdit)),
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		log.Error(err, "Failed to query tuples")
		return nil, nil, fmt.Errorf("failed to query tuples: %w", err)
	}

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

	return users, groups, nil
}

// function that adds users to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddUsersToTenantInKeto(ctx context.Context, name string, users []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddUsersToTenantInKeto").WithValues("Name", name)

	for _, user := range users {
		err := c.AddUserToTenantInKeto(ctx, name, user, permission)
		if err != nil {
			log.Error(err, "Failed to add user to observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success adding users to observability tenant")
	return nil
}

// function that adds a user to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddUserToTenantInKeto(ctx context.Context, name string, user string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddUserToTenantInKeto").WithValues("Name", name, "User", user)

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"User",
			user,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating tuple in keto")
	return nil
}

// function that removes users from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveUsersFromTenantInKeto(ctx context.Context, name string, users []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveUsersFromTenantInKeto").WithValues("Name", name)

	for _, user := range users {
		err := c.RemoveUserFromTenantInKeto(ctx, name, user, permission)
		if err != nil {
			log.Error(err, "Failed to remove user from observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success removing users from observability tenant")
	return nil
}

// function that removes a user from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveUserFromTenantInKeto(ctx context.Context, name string, user string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveUserFromTenantInKeto").WithValues("Name", name, "User", user)

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"User",
			user,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success deleting tuple in keto")
	return nil
}

// function that adds groups to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddGroupsToTenantInKeto(ctx context.Context, name string, groups []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddGroupsToTenantInKeto").WithValues("Name", name)

	for _, group := range groups {
		err := c.AddGroupToTenantInKeto(ctx, name, group, permission)
		if err != nil {
			log.Error(err, "Failed to add group to observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success adding groups to observability tenant")
	return nil
}

// function that adds a group to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddGroupToTenantInKeto(ctx context.Context, name string, group string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddGroupToTenantInKeto").WithValues("Name", name, "Group", group)

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"Group",
			group,
			"members",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating tuple in keto")
	return nil
}

// function that removes groups from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveGroupsFromTenantInKeto(ctx context.Context, name string, groups []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveGroupsFromTenantInKeto").WithValues("Name", name)

	for _, group := range groups {
		err := c.RemoveGroupFromTenantInKeto(ctx, name, group, permission)
		if err != nil {
			log.Error(err, "Failed to remove group from observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success removing groups from observability tenant")
	return nil
}

// function that removes a group from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveGroupFromTenantInKeto(ctx context.Context, name string, group string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveGroupFromTenantInKeto").WithValues("Name", name, "Group", group)

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"Group",
			group,
			"members",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success deleting tuple in keto")
	return nil
}

// function that adds an oauth client to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddOauthClientsToTenantInKeto(ctx context.Context, name string, oauthClients []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddOauthClientsToTenantInKeto").WithValues("Name", name)

	if permission == ObservabilityTenantPermissionEdit {
		err := fmt.Errorf("oauth clients do not support edit permission")
		log.Error(err, "Failed to add oauth clients from observability tenant")
		return err
	}

	for _, oauthClient := range oauthClients {
		err := c.AddOauthClientToTenantInKeto(ctx, name, oauthClient, permission)
		if err != nil {
			log.Error(err, "Failed to add oauth client to observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success adding oauth clients to observability tenant")
	return nil
}

// function that adds an oauth client to the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) AddOauthClientToTenantInKeto(ctx context.Context, name string, oauthClient string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("AddOauthClientToTenantInKeto").WithValues("Name", name, "OAuth2Client", oauthClient)

	if permission == ObservabilityTenantPermissionEdit {
		err := fmt.Errorf("oauth clients do not support edit permission")
		log.Error(err, "Failed to add oauth clients from observability tenant")
		return err
	}

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"OAuth2Client",
			oauthClient,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating tuple in keto")
	return nil
}

// function that removes oauth clients from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveOauthClientsFromTenantInKeto(ctx context.Context, name string, oauthClients []string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveOauthClientsFromTenantInKeto").WithValues("Name", name)

	if permission == ObservabilityTenantPermissionEdit {
		err := fmt.Errorf("oauth clients do not support edit permission")
		log.Error(err, "Failed to remove oauth clients from observability tenant")
		return err
	}

	for _, oauthClient := range oauthClients {
		err := c.RemoveOauthClientFromTenantInKeto(ctx, name, oauthClient, permission)
		if err != nil {
			log.Error(err, "Failed to remove oauth client from observability tenant")
			// return err // TODO: add some way to wrap errors
			continue
		}
	}

	log.Info("Success removing oauth clients from observability tenant")
	return nil
}

// function that removes an oauth client from the viewers or editors of an observability tenant in keto
func (c *ClientWrapper) RemoveOauthClientFromTenantInKeto(ctx context.Context, name string, oauthClient string, permission ObservabilityTenantPermission) error {
	log := c.Log.WithName("RemoveOauthClientFromTenantInKeto").WithValues("Name", name, "OAuth2Client", oauthClient)

	if permission == ObservabilityTenantPermissionEdit {
		err := fmt.Errorf("oauth clients do not support edit permission")
		log.Error(err, "Failed to remove oauth clients from observability tenant")
		return err
	}

	tuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  string(permission),
		Subject: rts.NewSubjectSet(
			"OAuth2Client",
			oauthClient,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, tuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success deleting tuple in keto")
	return nil
}

// function that lists all observability tenants in keto
func (c *ClientWrapper) ListTenantsInKeto(ctx context.Context) ([]*model.ObservabilityTenant, error) {
	log := c.Log.WithName("ListTenantsInKeto")

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    nil,
		Relation:  px.Ptr("organizations"),
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputTenants []*model.ObservabilityTenant

	for _, tuple := range respTuples {
		if tuple.Object != "" {
			tenant, err := c.GetTenantFromKeto(ctx, tuple.Object)
			if err != nil {
				continue
			}
			outputTenants = append(outputTenants, tenant)
		} else {
			continue
		}
	}

	log.Info("Success listing observability tenants in keto")
	return outputTenants, nil
}

// function that gets an observability tenant from keto
func (c *ClientWrapper) GetTenantFromKeto(ctx context.Context, name string) (*model.ObservabilityTenant, error) {
	log := c.Log.WithName("GetTenantFromKeto").WithValues("Name", name)

	if name == "" {
		return nil, fmt.Errorf("observability tenant name cannot be empty")
	}

	// check if group exists in keto
	exists, err := c.ObservabilityTenantExistsInKeto(ctx, name)
	if err != nil {
		log.Error(err, "Failed to check if observability tenant exists in keto")
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("observability tenant does not exist in keto")
	}

	return &model.ObservabilityTenant{
		Name: name,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that gets the viewers of an observability tenant from keto
func (c *ClientWrapper) GetViewersOfTenantFromKeto(ctx context.Context, name string) (*model.ObservabilityTenantViewers, error) {
	log := c.Log.WithName("GetViewersOfTenantFromKeto").WithValues("Name", name)
	// TODO: dedupe with GetTenantViewersInKeto since they are almost identical

	if name == "" {
		return nil, fmt.Errorf("observability tenant name cannot be empty")
	}

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    px.Ptr(name),
		Relation:  px.Ptr(string(ObservabilityTenantPermissionView)),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputViewers *model.ObservabilityTenantViewers
	var outputUserIds []*model.User
	var outputGroupNames []*model.Group
	var outputOauthClientIds []*model.OAuth2Client

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			outputUserIds = append(outputUserIds, &model.User{ID: subjectSet.Object})
		} else if subjectSet.Namespace == "Group" && subjectSet.Object != "" {
			outputGroupNames = append(outputGroupNames, &model.Group{Name: subjectSet.Object})
		} else if subjectSet.Namespace == "OAuth2Client" && subjectSet.Object != "" {
			outputOauthClientIds = append(outputOauthClientIds, &model.OAuth2Client{ClientID: &subjectSet.Object})
		} else {
			continue
		}
	}

	if len(outputUserIds) > 0 || len(outputGroupNames) > 0 || len(outputOauthClientIds) > 0 {
		outputViewers = &model.ObservabilityTenantViewers{
			Users:         outputUserIds,
			Groups:        outputGroupNames,
			Oauth2Clients: outputOauthClientIds,
		}
	}

	log.Info("Success getting viewers of observability tenant from keto")
	return outputViewers, nil
}

// function that gets the editors of an observability tenant from keto
func (c *ClientWrapper) GetEditorsOfTenantFromKeto(ctx context.Context, name string) (*model.ObservabilityTenantEditors, error) {
	log := c.Log.WithName("GetEditorsOfTenantFromKeto").WithValues("Name", name)
	// TODO: dedupe with GetTenantEditorsInKeto since they are almost identical

	if name == "" {
		return nil, fmt.Errorf("observability tenant name cannot be empty")
	}

	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Object:    px.Ptr(name),
		Relation:  px.Ptr(string(ObservabilityTenantPermissionEdit)),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputEditors *model.ObservabilityTenantEditors
	var outputUserIds []*model.User
	var outputGroupNames []*model.Group

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			outputUserIds = append(outputUserIds, &model.User{ID: subjectSet.Object})
		} else if subjectSet.Namespace == "Group" && subjectSet.Object != "" {
			outputGroupNames = append(outputGroupNames, &model.Group{Name: subjectSet.Object})
		} else {
			continue
		}
	}

	if len(outputUserIds) > 0 || len(outputGroupNames) > 0 {
		outputEditors = &model.ObservabilityTenantEditors{
			Users:  outputUserIds,
			Groups: outputGroupNames,
		}
	}

	log.Info("Success getting editors of observability tenant from keto")
	return outputEditors, nil
}

// function that gets user objects from a list of user ids
func (c *ClientWrapper) GetObservabilityTenantUsers(ctx context.Context, users []*model.User) ([]*model.User, error) {
	log := c.Log.WithName("GetObservabilityTenantUsers")
	//TODO: dedupe with GetOAuth2ClientUserLoginBindings

	var output []*model.User

	for _, inUser := range users {
		user, err := c.GetUserFromId(ctx, inUser.ID)
		if err != nil {
			log.Error(err, "failed to get user", "ID", inUser.ID)
			continue
		}
		output = append(output, user)
	}
	return output, nil
}

// function that gets group objects from a list of group names
func (c *ClientWrapper) GetObservabilityTenantGroups(ctx context.Context, groups []*model.Group) ([]*model.Group, error) {
	log := c.Log.WithName("GetObservabilityTenantGroups")
	// TODO: dedupe with GetOAuth2ClientGroupLoginBindings

	var output []*model.Group

	for _, inGroup := range groups {
		group, err := c.GetGroupFromName(ctx, inGroup.Name)
		if err != nil {
			log.Error(err, "failed to get group", "Name", inGroup.Name)
			continue
		}
		output = append(output, group)
	}
	return output, nil
}

// function that gets oauth2 client objects from a list of oauth2 client ids
func (c *ClientWrapper) GetObservabilityTenantOauth2Clients(ctx context.Context, clients []*model.OAuth2Client) ([]*model.OAuth2Client, error) {
	log := c.Log.WithName("GetObservabilityTenantOauth2Clients")
	// TODO: turn this into a more generic function that can be used for all the Get*From* functions

	var output []*model.OAuth2Client

	for _, inClient := range clients {
		client, err := c.GetOAuth2Client(ctx, *inClient.ClientID)
		if err != nil {
			log.Error(err, "failed to get oauth2 client", "ClientID", inClient.ClientID)
			continue
		}
		output = append(output, client)
	}
	return output, nil
}

// function that deletes an observability tenant from keto
func (c *ClientWrapper) DeleteObservabilityTenantInKeto(ctx context.Context, name string) (*model.ObservabilityTenant, error) {
	log := c.Log.WithName("DeleteObservabilityTenantInKeto").WithValues("Name", name)

	if name == "" {
		return nil, fmt.Errorf("observability tenant name cannot be empty")
	}

	// delete the relation tuple for the tenant
	tenantTuple := &rts.RelationTuple{
		Namespace: "ObservabilityTenant",
		Object:    name,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}
	err := c.KetoClient.DeleteTuple(ctx, tenantTuple)
	if err != nil {
		return nil, fmt.Errorf("failed to delete relation tuple: %w", err)
	}

	log.Info("Success deleting observability tenant in keto")
	return &model.ObservabilityTenant{
		Name: name,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}
