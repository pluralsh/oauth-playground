package clients

import (
	"context"
	"fmt"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

func (c *ClientWrapper) UpdateOrganization(ctx context.Context, name string, admins []string) (*model.Organization, error) {

	// TODO: figure out which admins to add or remove
	// TODO: create separate functions for adding and removing admins from an organization
	log := c.Log.WithName("Organization").WithValues("Name", name)

	_, err := c.OrganizationExistsInKeto(ctx, name)
	if err != nil {
		log.Error(err, "Failed to check if organization already exists in keto")
		return nil, err
	}

	//TODO: this doesn't seem to work and blocks updating with the first admin
	// if !exists {
	// 	log.Error(nil, "Organization does not exist in keto. Having multiple organizations is not yet supported.")
	// 	return nil, fmt.Errorf("Organization does not exist in keto. Having multiple organizations is not yet supported.")
	// }

	toAdd, toRemove, err := c.OrgAdminChangeset(ctx, name, admins)
	if err != nil {
		log.Error(err, "Failed to get organization admin changeset")
		return nil, err
	}

	for _, admin := range toAdd {
		err := c.AddAdminToOrganization(ctx, name, admin)
		if err != nil {
			log.Error(err, "Failed to add user to organization admins in keto", "User", admin)
			// TODO: add some way to wrap errors
			continue
		}
	}

	for _, admin := range toRemove {
		err := c.RemoveAdminFromOrganization(ctx, name, admin)
		if err != nil {
			log.Error(err, "Failed to remove user from organization admins in keto", "User", admin)
			// TODO: add some way to wrap errors
			continue
		}
	}

	return &model.Organization{
		Name: name,
	}, nil
}

// function that determines which admins to add or remove from an organization
func (c *ClientWrapper) OrgAdminChangeset(ctx context.Context, orgName string, admins []string) (toAdd []string, toRemove []string, err error) {
	currentAdmins, err := c.GetOrganizationAdmins(ctx, orgName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current admins: %w", err)
	}

	for _, admin := range admins {
		if !userIdInListOfUsers(currentAdmins, admin) {
			toAdd = append(toAdd, admin)
		}
	}

	for _, admin := range currentAdmins {
		if !stringContains(admins, admin.ID) {
			toRemove = append(toRemove, admin.ID)
		}
	}

	return toAdd, toRemove, nil
}

// function that adds an admin to an organization in keto
func (c *ClientWrapper) AddAdminToOrganization(ctx context.Context, orgName string, adminId string) error {
	adminTuple := &rts.RelationTuple{
		Namespace: "Organization",
		Object:    orgName,
		Relation:  "admins",
		Subject: rts.NewSubjectSet(
			"User",
			adminId,
			"",
		),
	}

	err := c.KetoClient.CreateTuple(ctx, adminTuple)
	if err != nil {
		return fmt.Errorf("failed to create tuple: %w", err)
	}

	return nil
}

// function that removes an admin from an organization in keto
func (c *ClientWrapper) RemoveAdminFromOrganization(ctx context.Context, orgName string, adminId string) error {
	adminTuple := &rts.RelationTuple{
		Namespace: "Organization",
		Object:    orgName,
		Relation:  "admins",
		Subject: rts.NewSubjectSet(
			"User",
			adminId,
			"",
		),
	}

	err := c.KetoClient.DeleteTuple(ctx, adminTuple)
	if err != nil {
		return fmt.Errorf("failed to delete tuple: %w", err)
	}

	return nil
}

// function that returns all admins for an organization
func (c *ClientWrapper) GetOrganizationAdmins(ctx context.Context, orgName string) ([]*model.User, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Object:    px.Ptr(orgName),
		Relation:  px.Ptr("admins"),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputAdmins []*model.User

	for _, tuple := range respTuples {
		subjectSet := tuple.Subject.GetSet()
		if subjectSet.Namespace == "User" {
			user, err := c.GetUserFromId(ctx, subjectSet.Object)
			if err != nil {
				continue
			}
			outputAdmins = append(outputAdmins, user)
		} else {
			continue
		}

	}

	return outputAdmins, nil
}

// function that checks if an organization exists in keto
func (c *ClientWrapper) OrganizationExistsInKeto(ctx context.Context, orgName string) (bool, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Object:    px.Ptr(orgName),
		Relation:  nil,
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return false, fmt.Errorf("failed to query tuples: %w", err)
	}

	if len(respTuples) == 0 {
		return false, nil
	}

	return true, nil
}

// function that lists all organizations in keto
func (c *ClientWrapper) ListOrganizations(ctx context.Context) ([]*model.Organization, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Object:    nil,
		Relation:  px.Ptr("admins"),
		Subject:   nil,
	}

	respTuples, err := c.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to query tuples: %w", err)
	}

	var outputOrgs []*model.Organization

	for _, tuple := range respTuples {
		outputOrgs = append(outputOrgs, &model.Organization{
			Name: tuple.Object,
		})
	}

	return outputOrgs, nil
}

// function that lists all organizations in keto
func (c *ClientWrapper) GetOrganization(ctx context.Context, orgName string) (*model.Organization, error) {
	log := c.Log.WithName("Organization").WithValues("Name", orgName)
	exists, err := c.OrganizationExistsInKeto(ctx, orgName)
	if err != nil {
		log.Error(err, "Failed to check if organization exists in keto")
		return nil, err
	}

	if !exists {
		log.Error(nil, "Organization does not exist in keto. Having multiple organizations is not yet supported.")
		return nil, fmt.Errorf("Organization does not exist in keto. Having multiple organizations is not yet supported.")
	}

	currentAdmins, err := c.GetOrganizationAdmins(ctx, orgName)
	if err != nil {
		return nil, fmt.Errorf("failed to get current admins: %w", err)
	}

	return &model.Organization{
		Name:   orgName,
		Admins: currentAdmins,
	}, nil
}
