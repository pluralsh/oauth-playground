package format

import (
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/pluralsh/oauth-playground/api-server/graph/model"
)

// function to convert a hydra OAuth2Client to a GraphQL OAuth2Client
func HydraOAuth2ClientToGraphQL(oauth2Client hydra.OAuth2Client) *model.OAuth2Client {

	jwks, _ := ToMapStringInterface(oauth2Client.Jwks)

	metadata, _ := ToMapStringInterface(oauth2Client.Metadata)

	return &model.OAuth2Client{
		AllowedCorsOrigins: oauth2Client.AllowedCorsOrigins,
		Audience:           oauth2Client.Audience,
		AuthorizationCodeGrantAccessTokenLifespan:  oauth2Client.AuthorizationCodeGrantAccessTokenLifespan,
		AuthorizationCodeGrantIDTokenLifespan:      oauth2Client.AuthorizationCodeGrantIdTokenLifespan,
		AuthorizationCodeGrantRefreshTokenLifespan: oauth2Client.AuthorizationCodeGrantRefreshTokenLifespan,
		BackChannelLogoutSessionRequired:           oauth2Client.BackchannelLogoutSessionRequired,
		BackChannelLogoutURI:                       oauth2Client.BackchannelLogoutUri,
		ClientCredentialsGrantAccessTokenLifespan:  oauth2Client.ClientCredentialsGrantAccessTokenLifespan,
		ClientID:                          oauth2Client.ClientId,
		ClientName:                        oauth2Client.ClientName,
		ClientSecret:                      oauth2Client.ClientSecret,
		ClientSecretExpiresAt:             oauth2Client.ClientSecretExpiresAt,
		ClientURI:                         oauth2Client.ClientUri,
		Contacts:                          oauth2Client.Contacts,
		CreatedAt:                         oauth2Client.CreatedAt,
		FrontchannelLogoutSessionRequired: oauth2Client.FrontchannelLogoutSessionRequired,
		FrontchannelLogoutURI:             oauth2Client.FrontchannelLogoutUri,
		GrantTypes:                        oauth2Client.GrantTypes,
		ImplicitGrantAccessTokenLifespan:  oauth2Client.ImplicitGrantAccessTokenLifespan,
		ImplicitGrantIDTokenLifespan:      oauth2Client.ImplicitGrantIdTokenLifespan,
		Jwks:                              jwks,
		JwksURI:                           oauth2Client.JwksUri,
		JwtBearerGrantAccessTokenLifespan: oauth2Client.JwtBearerGrantAccessTokenLifespan,
		LogoURI:                           oauth2Client.LogoUri,
		Metadata:                          metadata,
		Owner:                             oauth2Client.Owner,
		PolicyURI:                         oauth2Client.PolicyUri,
		PostLogoutRedirectUris:            oauth2Client.PostLogoutRedirectUris,
		RedirectUris:                      oauth2Client.RedirectUris,
		ResponseTypes:                     oauth2Client.ResponseTypes,
		Scope:                             oauth2Client.Scope,
		SectorIdentifierURI:               oauth2Client.SectorIdentifierUri,
		SubjectType:                       oauth2Client.SubjectType,
		TokenEndpointAuthMethod:           oauth2Client.TokenEndpointAuthMethod,
		TokenEndpointAuthSigningAlgorithm: oauth2Client.TokenEndpointAuthSigningAlg,
		TosURI:                            oauth2Client.TosUri,
		UpdatedAt:                         oauth2Client.UpdatedAt,
		UserinfoSignedResponseAlgorithm:   oauth2Client.UserinfoSignedResponseAlg,
	}
}

func GraphQLNewOAuth2ClientToHydra(allowedCorsOrigins []string, audience []string, authorizationCodeGrantAccessTokenLifespan *string, authorizationCodeGrantIDTokenLifespan *string, authorizationCodeGrantRefreshTokenLifespan *string, backChannelLogoutSessionRequired *bool, backChannelLogoutURI *string, clientCredentialsGrantAccessTokenLifespan *string, clientName string, clientSecret string, clientSecretExpiresAt *int64, clientURI *string, contacts []string, frontchannelLogoutSessionRequired *bool, frontchannelLogoutURI *string, grantTypes []string, implicitGrantAccessTokenLifespan *string, implicitGrantIDTokenLifespan *string, jwks map[string]interface{}, jwksURI *string, jwtBearerGrantAccessTokenLifespan *string, logoURI *string, metadata map[string]interface{}, policyURI *string, postLogoutRedirectUris []string, redirectUris []string, responseTypes []string, scope *string, sectorIdentifierURI *string, subjectType *string, tokenEndpointAuthMethod *string, tokenEndpointAuthSigningAlgorithm *string, tosURI *string, userinfoSignedResponseAlgorithm *string, loginBindings *model.LoginBindingsInput) hydra.OAuth2Client {
	return hydra.OAuth2Client{
		AllowedCorsOrigins: allowedCorsOrigins,
		Audience:           audience,
		AuthorizationCodeGrantAccessTokenLifespan:  authorizationCodeGrantAccessTokenLifespan,
		AuthorizationCodeGrantIdTokenLifespan:      authorizationCodeGrantIDTokenLifespan,
		AuthorizationCodeGrantRefreshTokenLifespan: authorizationCodeGrantRefreshTokenLifespan,
		BackchannelLogoutSessionRequired:           backChannelLogoutSessionRequired,
		BackchannelLogoutUri:                       backChannelLogoutURI,
		ClientCredentialsGrantAccessTokenLifespan:  clientCredentialsGrantAccessTokenLifespan,
		ClientName:                        &clientName,
		ClientSecret:                      &clientSecret,
		ClientSecretExpiresAt:             clientSecretExpiresAt,
		ClientUri:                         clientURI,
		Contacts:                          contacts,
		FrontchannelLogoutSessionRequired: frontchannelLogoutSessionRequired,
		FrontchannelLogoutUri:             frontchannelLogoutURI,
		GrantTypes:                        grantTypes,
		ImplicitGrantAccessTokenLifespan:  implicitGrantAccessTokenLifespan,
		ImplicitGrantIdTokenLifespan:      implicitGrantIDTokenLifespan,
		Jwks:                              jwks,
		JwksUri:                           jwksURI,
		JwtBearerGrantAccessTokenLifespan: jwtBearerGrantAccessTokenLifespan,
		LogoUri:                           logoURI,
		Metadata:                          metadata,
		PolicyUri:                         policyURI,
		PostLogoutRedirectUris:            postLogoutRedirectUris,
		RedirectUris:                      redirectUris,
		ResponseTypes:                     responseTypes,
		Scope:                             scope,
		SectorIdentifierUri:               sectorIdentifierURI,
		SubjectType:                       subjectType,
		TokenEndpointAuthMethod:           tokenEndpointAuthMethod,
		TokenEndpointAuthSigningAlg:       tokenEndpointAuthSigningAlgorithm,
		TosUri:                            tosURI,
		UserinfoSignedResponseAlg:         userinfoSignedResponseAlgorithm,
	}
}
