package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/ory/oathkeeper/pipeline/authn"

	rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"
	px "github.com/ory/x/pointerx"
	ketocl "github.com/pluralsh/oauth-playground/keto/client"
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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/", LogAuthenticationSession)
	http.ListenAndServe(":8080", r)
}

func LogAuthenticationSession(w http.ResponseWriter, r *http.Request) {

	conndetails := ketocl.NewConnectionDetailsFromEnv()
	kcl, err := ketocl.NewGrpcClient(context.Background(), conndetails)
	if err != nil {
		panic("Encountered error: " + err.Error())
	}

	data := &AuthenticationSessionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// article := data.Subject
	// dbNewArticle(article)

	json, _ := json.Marshal(data.AuthenticationSession)
	log.Printf("Post body: %v", string(json))

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewAuthenticationSessionResponse(kcl, data))
}

func NewAuthenticationSessionResponse(kcl *ketocl.GrpcClient, session *AuthenticationSessionRequest) *AuthenticationSessionRequest {
	tenants, err := getUserTenants(kcl, session.Subject)
	if err != nil {
		log.Printf("Error getting tenants: %v", err)
	}
	log.Printf("Tenants: %v", tenants)

	session.Header = map[string][]string{
		"X-Scope-OrgID": {strings.Join(tenants, "|")},
	}
	return session
}

// Get all the ObservabilityTenants has permissions for based on direct bindings and group memberships
func getUserTenants(kcl *ketocl.GrpcClient, subject string) ([]string, error) {
	// get all the groups a user is a member of
	groups, err := getUserGroups(kcl, subject)
	if err != nil {
		return []string{}, err
	}

	// get all the tenants a user has permissions for
	tenants, err := getUserDirectTenants(kcl, subject)
	if err != nil {
		return []string{}, err
	}

	// get all the tenants a user has permissions for via group membership
	for _, group := range groups {
		groupTenants, err := getGroupTenants(kcl, group)
		if err != nil {
			return []string{}, err
		}
		tenants = append(tenants, groupTenants...)
	}

	// get all the tenants a user has permissions for via organization admin permissions
	orgs, err := isOrgAdmin(kcl, subject)
	if err != nil {
		return []string{}, err
	}
	for _, org := range orgs {
		orgTenants, err := getOrgTenants(kcl, org)
		if err != nil {
			return []string{}, err
		}
		tenants = append(tenants, orgTenants...)
	}

	return tenants, nil
}

// Check in which organizations a user is an admin
func isOrgAdmin(kcl *ketocl.GrpcClient, subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Organization"),
		Relation:  px.Ptr("admins"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
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
func getOrgTenants(kcl *ketocl.GrpcClient, org string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("Organization", org, ""),
	}
	respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
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
func getUserGroups(kcl *ketocl.GrpcClient, subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("Group"),
		Relation:  px.Ptr("members"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
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
func getUserDirectTenants(kcl *ketocl.GrpcClient, subject string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("User", subject, ""),
	}
	respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
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
func getGroupTenants(kcl *ketocl.GrpcClient, group string) ([]string, error) {
	query := rts.RelationQuery{
		Namespace: px.Ptr("ObservabilityTenant"),
		Subject:   rts.NewSubjectSet("Group", group, "members"),
	}
	respTuples, err := kcl.QueryAllTuples(context.Background(), &query, 100)
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
