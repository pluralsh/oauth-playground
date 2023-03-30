package directives

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pluralsh/oauth-playground/api-server/handlers"
)

func (d *Directive) IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	userCtx := handlers.ForContext(ctx)
	if userCtx != nil {

		if userCtx.KratosSession != nil {

			session := userCtx.KratosSession

			if !session.GetActive() || session.ExpiresAt.Before(time.Now()) {
				return nil, fmt.Errorf("Access denied: User session not active or has expired.")
			}
		} else {
			// TODO: remove debug log message
			d.C.Log.Info("auth directive failed")
			return nil, fmt.Errorf("Access denied: User must be logged in.")
		}
	} else {
		// TODO: remove debug log message
		d.C.Log.Info("auth directive failed")
		return nil, fmt.Errorf("Access denied: User must be logged in.")
	}

	// TODO: remove debug log message
	// setupLog.Info("auth directive successful")
	d.C.Log.Info("auth directive successful")

	// let it pass through if user is authenticated
	return next(ctx)
}
