package clients

import (
	"context"

	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

func (c *ClientWrapper) GetUserFromId(ctx context.Context, id string) (*model.User, error) {
	log := c.Log.WithName("User").WithValues("ID", id)

	user, resp, err := c.KratosClient.IdentityApi.GetIdentity(ctx, id).Execute()
	if err != nil || resp.StatusCode != 200 {
		log.Error(err, "failed to get user")
		return nil, err
	}

	var email string
	var name string

	if val, ok := user.Traits.(map[string]interface{})["email"]; ok {
		if foundEmail, ok := val.(string); ok {
			email = foundEmail
		} else {
			log.Error(err, "Error when parsing email")
		}
	}

	if val, ok := user.Traits.(map[string]interface{})["name"]; ok {
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

	log.Info("Success getting User")
	return &model.User{
		ID:    user.Id,
		Email: email,
		Name:  &name,
	}, nil
}
