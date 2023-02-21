package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"

	kratosClient "github.com/ory/kratos-client-go"
)

const (
	ProjectKey         = "projects.platform.kubricks.io/name"
	ProjectNSKey       = "projects.platform.kubricks.io/namespace"
	ParentProjectKey   = "projects.platform.kubricks.io/parent-project"
	ParentProjectNSKey = "projects.platform.kubricks.io/parent-project-namespace"
	PersonalKey        = "environment.platform.kubricks.io/personal-namespace"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

type AccountType string

const (
	TypeServiceAccount AccountType = "ServiceAccount"
	TypeUser           AccountType = "User"
)

// A stand-in for our backed user object
type User struct {
	// Type          AccountType
	// Username      string
	Id    string
	Name  string
	Email string
	// Groups        []string
	// ProfileImage  *string
	// JWT           *string
	KratosSession *kratosClient.Session
	IsAdmin       bool
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(kratosClient *kratosClient.APIClient, log logr.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("ory_kratos_session")
			// Allow unauthenticated users in
			if err != nil || cookie == nil {
				next.ServeHTTP(w, r)
				return
			}

			// log.Info(fmt.Sprintf("Cookie: %s", cookie.String()))

			resp, req, err := kratosClient.FrontendApi.ToSession(context.Background()).Cookie(cookie.String()).Execute()
			if err != nil {
				// TODO: should we return here?
				log.Error(err, fmt.Sprintf("Error when calling `V0alpha2Api.ToSession``: %v\n", err))
				log.Error(err, fmt.Sprintf("Full HTTP response: %v\n", req))
			}
			// response from `ToSession`: Session
			log.Info(fmt.Sprintf("%s", resp.Identity.Traits))

			var email string
			var name string

			if val, ok := resp.Identity.Traits.(map[string]interface{})["email"]; ok {
				if foundEmail, ok := val.(string); ok {
					email = foundEmail
				} else {
					log.Error(err, "Error when parsing email")
				}
			}

			if val, ok := resp.Identity.Traits.(map[string]interface{})["name"]; ok {
				first, firstFound := val.(map[string]interface{})["first"]
				last, lastFound := val.(map[string]interface{})["last"]

				if firstName, ok := first.(string); ok {
					if lastName, ok := last.(string); ok {
						if firstFound && lastFound {
							name = firstName + " " + lastName
						}
					}
				}
			}

			user := &User{
				KratosSession: resp,
				Email:         email,
				Name:          name,
				Id:            resp.Identity.Id,
			}

			// user.Type = TypeUser
			// user.Username = strings.Replace(strings.Split(user.Email, "@")[0], ".", "-", -1)

			user.IsAdmin = false
			// for _, adminEmail := range kubricksConfig.Spec.Admins {
			// 	if adminEmail == email {
			// 		user.IsAdmin = true
			// 	}
			// }

			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			log.Info("Success auth", "user", user.Email)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}
