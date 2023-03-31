package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/ory/oathkeeper/pipeline/authn"
	"github.com/pluralsh/oauth-playground/api-server/utils"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
)

type AuthenticationSessionRequest struct {
	*authn.AuthenticationSession

	// User *UserPayload `json:"user,omitempty"`

	// ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *AuthenticationSessionRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Subject == "" {
		return errors.New("missing required Article fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	// a.ProtectedID = ""                                 // unset the protected ID
	// a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

func (a *AuthenticationSessionRequest) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// a.Elapsed = 10
	return nil
}

func (h *Handler) HydrateObservabilityTenants(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithName("HydrateObservabilityTenants")
	data := &AuthenticationSessionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// article := data.Subject
	// dbNewArticle(article)

	json, _ := json.Marshal(data.AuthenticationSession)
	log.Info("Post", "body", string(json))

	render.Status(r, http.StatusOK) // TODO: we should probably do error checking above so we can return a 400 if something goes wrong
	render.Render(w, r, h.NewAuthenticationSessionResponse(data))
}

func (h *Handler) NewAuthenticationSessionResponse(session *AuthenticationSessionRequest) *AuthenticationSessionRequest {
	log := h.Log.WithName("NewAuthenticationSessionResponse")
	tenants, err := h.getUserTenants(session.Subject)
	if err != nil {
		log.Error(err, "Error getting tenants")
	}
	log.Info("Success getting tenants", "Tenants", tenants)

	session.Header = map[string][]string{
		"X-Scope-OrgID": {strings.Join(tenants, "|")},
	}
	return session
}

// Get all the ObservabilityTenants has permissions for based on direct bindings and group memberships
func (h *Handler) getUserTenants(subject string) ([]string, error) {
	// get all the groups a user is a member of
	groups, err := h.getUserGroups(subject)
	if err != nil {
		return []string{}, err
	}

	// get all the tenants a user has permissions for
	tenants, err := h.getUserDirectTenants(subject)
	if err != nil {
		return []string{}, err
	}

	// get all the tenants a user has permissions for via group membership
	for _, group := range groups {
		groupTenants, err := h.getGroupTenants(group)
		if err != nil {
			return []string{}, err
		}
		tenants = append(tenants, groupTenants...)
	}

	// get all the tenants a user has permissions for via organization admin permissions
	orgs, err := h.isOrgAdmin(subject)
	if err != nil {
		return []string{}, err
	}
	for _, org := range orgs {
		orgTenants, err := h.getOrgTenants(org)
		if err != nil {
			return []string{}, err
		}
		tenants = append(tenants, orgTenants...)
	}

	return utils.DedupeList(tenants), nil
}

// Check in which organizations a user is an admin
func (h *Handler) isOrgAdmin(subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Relation:  px.Ptr("admins"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := h.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return []string{}, err
	}

	output := []string{}
	for _, tuple := range respTuples {
		// likely unnecessary but just in case
		if tuple.Namespace == "Organization" && tuple.Relation == "admins" && tuple.Object != "" {
			output = append(output, tuple.Object)
		}
	}

	return output, nil
}

// Get all the ObservabilityTenants that belong to an organization
func (h *Handler) getOrgTenants(org string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("Organization", org, ""),
	}
	respTuples, err := h.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return []string{}, err
	}

	output := []string{}
	for _, tuple := range respTuples {
		// likely unnecessary but just in case
		if tuple.Namespace == "ObservabilityTenant" && tuple.Relation == "organizations" && tuple.Object != "" {
			output = append(output, tuple.Object)
		}
	}

	return output, nil
}

// Get the groups a user is a member of
func (h *Handler) getUserGroups(subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Group"),
		Relation:  px.Ptr("members"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := h.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return []string{}, err
	}

	output := []string{}
	for _, tuple := range respTuples {
		// likely unnecessary but just in case
		if tuple.Namespace == "Group" && tuple.Relation == "members" && tuple.Object != "" {
			output = append(output, tuple.Object)
		}
	}

	return output, nil
}

// Get the ObservabilityTenants a user has permissions for
func (h *Handler) getUserDirectTenants(subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := h.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return []string{}, err
	}

	output := []string{}
	for _, tuple := range respTuples {
		// likely unnecessary but just in case
		if tuple.Namespace == "ObservabilityTenant" && tuple.Object != "" {
			output = append(output, tuple.Object)
		}
	}

	return output, nil
}

// Get the ObservabilityTenants a group has permissions for
func (h *Handler) getGroupTenants(group string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("Group", group, "members"),
	}
	respTuples, err := h.C.KetoClient.QueryAllTuples(context.Background(), &query, 100)
	if err != nil {
		return []string{}, err
	}

	output := []string{}
	for _, tuple := range respTuples {
		// likely unnecessary but just in case
		if tuple.Namespace == "ObservabilityTenant" && tuple.Object != "" {
			output = append(output, tuple.Object)
		}
	}

	return output, nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
