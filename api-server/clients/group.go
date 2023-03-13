package clients

import (
	"context"
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

func (c *ClientWrapper) MutateGroup(ctx context.Context, name string, members []string) (*model.Group, error) {

	// TODO: figure out which members to add or remove
	log := c.Log.WithName("Group").WithValues("Name", name)

	// TODO: figure out how to distinguish between creating or updating a group
	// updating a group would require that we first check if it exists and if a user is allowed to update it
	// creating a group would require that we first check if it exists and if a user is allowed to create it

	groupExists, err := c.GroupExistsInKeto(ctx, name)
	if err != nil {
		log.Error(err, "Failed to check if group already exists in keto")
		return nil, err
	}

	if !groupExists {
		err := c.CreateGroupInKeto(ctx, name)
		if err != nil {
			log.Error(err, "Failed to create group in keto")
			return nil, err
		}
	}

	toAdd, toRemove, err := c.GroupChangeset(ctx, name, members)
	if err != nil {
		log.Error(err, "Failed to get group changeset")
		return nil, err
	}

	for _, member := range toAdd {
		err := c.AddUserToGroupInKeto(ctx, name, member)
		if err != nil {
			log.Error(err, "Failed to add user to group in keto", "User", member)
			// TODO: add some way to wrap errors
			continue
		}
	}

	for _, member := range toRemove {
		err := c.RemoveUserFromGroupInKeto(ctx, name, member)
		if err != nil {
			log.Error(err, "Failed to remove user from group in keto", "User", member)
			// TODO: add some way to wrap errors
			continue
		}
	}

	return &model.Group{
		Name: name,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that determines which users to add or remove from a group
func (c *ClientWrapper) GroupChangeset(ctx context.Context, groupName string, members []string) (toAdd []string, toRemove []string, err error) {
	currentMembers, err := c.GetGroupMembersInKeto(ctx, groupName)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current members: %w", err)
	}

	for _, member := range members {
		if !userIdInListOfUsers(currentMembers, member) {
			toAdd = append(toAdd, member)
		}
	}

	for _, member := range currentMembers {
		if !stringContains(members, member.ID) {
			toRemove = append(toRemove, member.ID)
		}
	}

	return toAdd, toRemove, nil
}

// function that checks if a user ID is in a []*model.User
func userIdInListOfUsers(users []*model.User, userId string) bool {
	for _, user := range users {
		if user.ID == userId {
			return true
		}
	}
	return false
}

// function that checks if a string is in a []string
func stringContains(list []string, s string) bool {
	for _, u := range list {
		if u == s {
			return true
		}
	}
	return false
}

// function that checks if a user is part of a group
func (c *ClientWrapper) IsUserInGroup(ctx context.Context, groupName string, userId string) (bool, error) {
	log := c.Log.WithName("IsUserInGroup").WithValues("Name", groupName)

	userTuple := &rts.RelationTuple{
		Namespace: "Group",
		Object:    groupName,
		Relation:  "members",
		Subject: rts.NewSubjectSet(
			"User",
			userId,
			"",
		),
	}

	_, err := c.KetoClient.Check(ctx, userTuple)
	if err != nil {
		return false, fmt.Errorf("failed to check tuple: %w", err)
	}

	log.Info("Success checking if user is in group")
	return true, nil
}

// function that creates a group in keto
func (c *ClientWrapper) CreateGroupInKeto(ctx context.Context, name string) error {
	log := c.Log.WithName("CreateGroupInKeto").WithValues("Name", name)

	groupTuple := &rts.RelationTuple{
		Namespace: "Group",
		Object:    name,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, groupTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success creating group in keto")
	return nil
}

// func that adds a user to a group in keto
func (c *ClientWrapper) AddUserToGroupInKeto(ctx context.Context, groupName string, userId string) error {
	log := c.Log.WithName("AddUserToGroupInKeto").WithValues("Name", groupName)

	userTuple := &rts.RelationTuple{
		Namespace: "Group",
		Object:    groupName,
		Relation:  "members",
		Subject: rts.NewSubjectSet(
			"User",
			userId,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, userTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	log.Info("Success adding user to group in keto")
	return nil
}

// function that removes a user from a group in keto
func (c *ClientWrapper) RemoveUserFromGroupInKeto(ctx context.Context, groupName string, userId string) error {
	log := c.Log.WithName("RemoveUserFromGroupInKeto").WithValues("Name", groupName)

	userTuple := &rts.RelationTuple{
		Namespace: "Group",
		Object:    groupName,
		Relation:  "members",
		Subject: rts.NewSubjectSet(
			"User",
			userId,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, userTuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	log.Info("Success removing user from group in keto")
	return nil
}

// function that gets all members of a group in keto
func (c *ClientWrapper) GetGroupMembersInKeto(ctx context.Context, groupName string) ([]*model.User, error) {
	log := c.Log.WithName("GetGroupMembersInKeto").WithValues("Name", groupName)

	query := rts.RelationQuery{
		Namespace: px.Ptr("Group"),
		Object:    px.Ptr(groupName),
		Relation:  px.Ptr("members"),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputMembers []*model.User

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" && subjectSet.Object != "" {
			user, err := c.GetUserFromId(ctx, subjectSet.Object) // TODO: it might be better to split this off into a separate function
			if err != nil {
				continue
			}
			outputMembers = append(outputMembers, user)
		} else {
			continue
		}

	}

	log.Info("Success getting group members in keto")
	return outputMembers, nil
}

// function that gets a group from keto
func (c *ClientWrapper) GetGroupFromName(ctx context.Context, groupName string) (*model.Group, error) {
	log := c.Log.WithName("GetGroupFromName").WithValues("Name", groupName)

	if groupName == "" {
		return nil, fmt.Errorf("group name cannot be empty")
	}

	// check if group exists in keto
	exists, err := c.GroupExistsInKeto(ctx, groupName)
	if err != nil {
		log.Error(err, "Failed to check if group exists in keto")
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("group does not exist in keto")
	}

	return &model.Group{
		Name: groupName,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}

// function that checks if a group exists in keto
func (c *ClientWrapper) GroupExistsInKeto(ctx context.Context, groupName string) (bool, error) {
	log := c.Log.WithName("GroupExistsInKeto").WithValues("Name", groupName)

	query := rts.RelationQuery{
		Namespace: px.Ptr("Group"),
		Object:    px.Ptr(groupName),
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

// function that lists all groups in keto
func (c *ClientWrapper) ListGroupsInKeto(ctx context.Context) ([]*model.Group, error) {
	log := c.Log.WithName("ListGroupsInKeto")

	query := rts.RelationQuery{
		Namespace: px.Ptr("Group"),
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

	var outputGroups []*model.Group

	for _, tuple := range respTuples {
		if tuple.Object != "" {
			group, err := c.GetGroupFromName(ctx, tuple.Object)
			if err != nil {
				continue
			}
			outputGroups = append(outputGroups, group)
		} else {
			continue
		}
	}

	log.Info("Success listing groups in keto")
	return outputGroups, nil
}

// funtion that deletes a group in keto
func (c *ClientWrapper) DeleteGroup(ctx context.Context, groupName string) (*model.Group, error) {
	log := c.Log.WithName("DeleteGroup").WithValues("Name", groupName)

	tuple := &rts.RelationTuple{
		Namespace: "Group",
		Object:    groupName,
		Relation:  "organizations",
		Subject: rts.NewSubjectSet(
			"Organization",
			"main", //TODO: decide whether to hardcode this or not
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, tuple)
	if err != nil {
		log.Error(err, "Failed to delete tuple")
		return nil, err
	}

	log.Info("Success deleting group in keto")
	return &model.Group{
		Name: groupName,
		Organization: &model.Organization{
			Name: "main", //TODO: decide whether to hardcode this or not
		},
	}, nil
}
