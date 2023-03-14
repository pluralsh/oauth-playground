import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions =  {}
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Map: any;
  Time: any;
};

/** Input for adding a user to an organization as an administrator. */
export type Admin = {
  /** The ID of the user to add as an admin. */
  id: Scalars['ID'];
};

/** Representation a group of users. */
export type Group = {
  __typename?: 'Group';
  /** The users that are admins of the organization. */
  members?: Maybe<Array<User>>;
  /** The unique name of the group. */
  name: Scalars['String'];
  /** The organization that the group belongs to. */
  organization: Organization;
};

/** Representation of users and groups that are allowed to login with through OAuth2 Client. */
export type LoginBindings = {
  __typename?: 'LoginBindings';
  /** The groups that are allowed to login with this OAuth2 Client. */
  groups?: Maybe<Array<Group>>;
  /** The users that are allowed to login with this OAuth2 Client. */
  users?: Maybe<Array<User>>;
};

export type LoginBindingsInput = {
  /** The groups that are allowed to login with this OAuth2 Client. */
  groups?: Maybe<Array<Scalars['ID']>>;
  /** The users that are allowed to login with this OAuth2 Client. */
  users?: Maybe<Array<Scalars['ID']>>;
};

export type Mutation = {
  __typename?: 'Mutation';
  /** Create a new OAuth2 Client. */
  createOAuth2Client: OAuth2Client;
  /** Create a new user. */
  createUser: User;
  /** Delete a group. */
  deleteGroup: Group;
  /** Delete an OAuth2 Client. */
  deleteOAuth2Client: OAuth2Client;
  /** Delete an observability tenant. */
  deleteObservabilityTenant: ObservabilityTenant;
  /** Delete a user. */
  deleteUser: User;
  /** Create or update a group. */
  group: Group;
  /** Create or update an observability tenant. */
  observabilityTenant: ObservabilityTenant;
  /** Create a new organization. */
  organization: Organization;
  /** Update an OAuth 2 Client. */
  updateOAuth2Client: OAuth2Client;
};


export type MutationCreateOAuth2ClientArgs = {
  ClientSecretExpiresAt?: Maybe<Scalars['Int']>;
  allowedCorsOrigins?: Maybe<Array<Scalars['String']>>;
  audience?: Maybe<Array<Scalars['String']>>;
  authorizationCodeGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  authorizationCodeGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  authorizationCodeGrantRefreshTokenLifespan?: Maybe<Scalars['String']>;
  backChannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  backChannelLogoutUri?: Maybe<Scalars['String']>;
  clientCredentialsGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  clientName?: Maybe<Scalars['String']>;
  clientSecret?: Maybe<Scalars['String']>;
  clientUri?: Maybe<Scalars['String']>;
  contacts?: Maybe<Array<Scalars['String']>>;
  frontchannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  frontchannelLogoutUri?: Maybe<Scalars['String']>;
  grantTypes?: Maybe<Array<Scalars['String']>>;
  implicitGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  implicitGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  jwks?: Maybe<Scalars['Map']>;
  jwksUri?: Maybe<Scalars['String']>;
  jwtBearerGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  loginBindings?: Maybe<LoginBindingsInput>;
  logoUri?: Maybe<Scalars['String']>;
  metadata?: Maybe<Scalars['Map']>;
  policyUri?: Maybe<Scalars['String']>;
  postLogoutRedirectUris?: Maybe<Array<Scalars['String']>>;
  redirectUris?: Maybe<Array<Scalars['String']>>;
  responseTypes?: Maybe<Array<Scalars['String']>>;
  scope?: Maybe<Scalars['String']>;
  sectorIdentifierUri?: Maybe<Scalars['String']>;
  subjectType?: Maybe<Scalars['String']>;
  tokenEndpointAuthMethod?: Maybe<Scalars['String']>;
  tokenEndpointAuthSigningAlgorithm?: Maybe<Scalars['String']>;
  tosUri?: Maybe<Scalars['String']>;
  userinfoSignedResponseAlgorithm?: Maybe<Scalars['String']>;
};


export type MutationCreateUserArgs = {
  email: Scalars['String'];
  name?: Maybe<NameInput>;
};


export type MutationDeleteGroupArgs = {
  name: Scalars['String'];
};


export type MutationDeleteOAuth2ClientArgs = {
  clientId: Scalars['String'];
};


export type MutationDeleteObservabilityTenantArgs = {
  name: Scalars['String'];
};


export type MutationDeleteUserArgs = {
  id: Scalars['ID'];
};


export type MutationGroupArgs = {
  members?: Maybe<Array<Scalars['String']>>;
  name: Scalars['String'];
};


export type MutationObservabilityTenantArgs = {
  editors?: Maybe<ObservabilityTenantEditorsInput>;
  name: Scalars['String'];
  viewers?: Maybe<ObservabilityTenantViewersInput>;
};


export type MutationOrganizationArgs = {
  admins: Array<Scalars['String']>;
  name: Scalars['String'];
};


export type MutationUpdateOAuth2ClientArgs = {
  ClientSecretExpiresAt?: Maybe<Scalars['Int']>;
  allowedCorsOrigins?: Maybe<Array<Scalars['String']>>;
  audience?: Maybe<Array<Scalars['String']>>;
  authorizationCodeGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  authorizationCodeGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  authorizationCodeGrantRefreshTokenLifespan?: Maybe<Scalars['String']>;
  backChannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  backChannelLogoutUri?: Maybe<Scalars['String']>;
  clientCredentialsGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  clientId: Scalars['String'];
  clientName?: Maybe<Scalars['String']>;
  clientSecret?: Maybe<Scalars['String']>;
  clientUri?: Maybe<Scalars['String']>;
  contacts?: Maybe<Array<Scalars['String']>>;
  frontchannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  frontchannelLogoutUri?: Maybe<Scalars['String']>;
  grantTypes?: Maybe<Array<Scalars['String']>>;
  implicitGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  implicitGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  jwks?: Maybe<Scalars['Map']>;
  jwksUri?: Maybe<Scalars['String']>;
  jwtBearerGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  loginBindings?: Maybe<LoginBindingsInput>;
  logoUri?: Maybe<Scalars['String']>;
  metadata?: Maybe<Scalars['Map']>;
  policyUri?: Maybe<Scalars['String']>;
  postLogoutRedirectUris?: Maybe<Array<Scalars['String']>>;
  redirectUris?: Maybe<Array<Scalars['String']>>;
  responseTypes?: Maybe<Array<Scalars['String']>>;
  scope?: Maybe<Scalars['String']>;
  sectorIdentifierUri?: Maybe<Scalars['String']>;
  subjectType?: Maybe<Scalars['String']>;
  tokenEndpointAuthMethod?: Maybe<Scalars['String']>;
  tokenEndpointAuthSigningAlgorithm?: Maybe<Scalars['String']>;
  tosUri?: Maybe<Scalars['String']>;
  userinfoSignedResponseAlgorithm?: Maybe<Scalars['String']>;
};

/** The first and last name of a user. */
export type Name = {
  __typename?: 'Name';
  /** The user's first name. */
  first?: Maybe<Scalars['String']>;
  /** The user's last name. */
  last?: Maybe<Scalars['String']>;
};

export type NameInput = {
  /** The user's first name. */
  first?: Maybe<Scalars['String']>;
  /** The user's last name. */
  last?: Maybe<Scalars['String']>;
};

/** Representation of the information about an OAuth2 Client sourced from Hydra. */
export type OAuth2Client = {
  __typename?: 'OAuth2Client';
  /** OAuth 2.0 Client Secret Expires At. The field is currently not supported and its value is always 0. */
  ClientSecretExpiresAt?: Maybe<Scalars['Int']>;
  /** OAuth 2.0 Client Allowed CORS Origins. AllowedCORSOrigins is an array of allowed CORS origins. If the array is empty, the value of the first element is considered valid. */
  allowedCorsOrigins?: Maybe<Array<Scalars['String']>>;
  /** OAuth 2.0 Client Audience. Audience is an array of URLs that the OAuth 2.0 Client is allowed to request tokens for. */
  audience?: Maybe<Array<Scalars['String']>>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  authorizationCodeGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  authorizationCodeGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  authorizationCodeGrantRefreshTokenLifespan?: Maybe<Scalars['String']>;
  /** OpenID Connect Back-Channel Logout Session Required  Boolean value specifying whether the RP requires that a sid (session ID) Claim be included in the Logout Token to identify the RP session with the OP when the backchannel_logout_uri is used. If omitted, the default value is false. */
  backChannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  /** OpenID Connect Back-Channel Logout URI. RP URL that will cause the RP to log itself out when sent a Logout Token by the OP. */
  backChannelLogoutUri?: Maybe<Scalars['String']>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  clientCredentialsGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client ID. The ID is autogenerated and immutable. */
  clientId?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Name. The human-readable name of the client to be presented to the end-user during authorization. */
  clientName?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Secret. The secret will be included in the create request as cleartext, and then never again. The secret is kept in hashed format and is not recoverable once lost. */
  clientSecret?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client URI. ClientURI is a URL string of a web page providing information about the client. If present, the server SHOULD display this URL to the end-user in a clickable fashion. */
  clientUri?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Contacts. Contacts is an array of strings representing ways to contact people responsible for this client, typically email addresses. */
  contacts?: Maybe<Array<Scalars['String']>>;
  /** OAuth 2.0 Client Creation Date. CreatedAt returns the timestamp of the client's creation. */
  createdAt?: Maybe<Scalars['Time']>;
  /** OpenID Connect Front-Channel Logout Session Required. Boolean value specifying whether the RP requires that iss (issuer) and sid (session ID) query parameters be included to identify the RP session with the OP when the frontchannel_logout_uri is used. If omitted, the default value is false. */
  frontchannelLogoutSessionRequired?: Maybe<Scalars['Boolean']>;
  /** OpenID Connect Front-Channel Logout URI. RP URL that will cause the RP to log itself out when rendered in an iframe by the OP. */
  frontchannelLogoutUri?: Maybe<Scalars['String']>;
  grantTypes?: Maybe<Array<Scalars['String']>>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  implicitGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  implicitGrantIdTokenLifespan?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client JSON Web Key Set. Client's JSON Web Key Set [JWK] document, passed by value. The semantics of the jwks parameter are the same as the jwks_uri parameter, other than that the JWK Set is passed by value, rather than by reference. This parameter is intended only to be used by Clients that, for some reason, are unable to use the jwks_uri parameter, for instance, by native applications that might not have a location to host the contents of the JWK Set. If a Client can use jwks_uri, it MUST NOT use jwks. One significant downside of jwks is that it does not enable key rotation (which jwks_uri does, as described in Section 10 of OpenID Connect Core 1.0 [OpenID.Core]). The jwks_uri and jwks parameters MUST NOT be used together. */
  jwks?: Maybe<Scalars['Map']>;
  /** OAuth 2.0 Client JSON Web Key Set URI. Client's JSON Web Key Set [JWK] document URI, passed by reference. The semantics of the jwks_uri parameter are the same as the jwks parameter, other than that the JWK Set is passed by reference, rather than by value. The jwks_uri and jwks parameters MUST NOT be used together. */
  jwksUri?: Maybe<Scalars['String']>;
  /** Specify a time duration in milliseconds, seconds, minutes, hours. For example, 1h, 1m, 1s, 1ms. */
  jwtBearerGrantAccessTokenLifespan?: Maybe<Scalars['String']>;
  /** The users and groups that are allowed to login with this OAuth2 Client. */
  loginBindings?: Maybe<LoginBindings>;
  /** OAuth 2.0 Client Logo URI. A URL string referencing the client's logo. */
  logoUri?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Metadata. Metadata is a map of key-value pairs that contain additional information about the client. */
  metadata?: Maybe<Scalars['Map']>;
  /** The organization that owns this OAuth2 Client. */
  organization: Organization;
  /** OAuth 2.0 Client Owner. Owner is a string identifying the owner of the OAuth 2.0 Client. */
  owner?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Policy URI. PolicyURI is a URL string that points to a human-readable privacy policy document that describes how the deployment organization collects, uses, retains, and discloses personal data. */
  policyUri?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Post Logout Redirect URIs. PostLogoutRedirectUris is an array of allowed URLs to which the RP is allowed to redirect the End-User's User Agent after a logout has been performed. */
  postLogoutRedirectUris?: Maybe<Array<Scalars['String']>>;
  /** OAuth 2.0 Client Redirect URIs. RedirectUris is an array of allowed redirect URLs for the OAuth 2.0 Client. */
  redirectUris?: Maybe<Array<Scalars['String']>>;
  /** OAuth 2.0 Client Response Types. ResponseTypes is an array of the OAuth 2.0 response type strings that the client can use at the Authorization Endpoint. */
  responseTypes?: Maybe<Array<Scalars['String']>>;
  /** OAuth 2.0 Client Scope. Scope is a string containing a space-separated list of scope values (as described in Section 3.3 of OAuth 2.0 [RFC6749]) that the client can use when requesting access tokens. */
  scope?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Sector Identifier URI. SectorIdentifierURI is a URL string using the https scheme referencing a file with a single JSON array of redirect_uri values. */
  sectorIdentifierUri?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Subject Type. SubjectType requested for responses to this Client. The subject_types_supported Discovery parameter contains a list of the supported subject_type values for this server. Valid types include pairwise and public. */
  subjectType?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Token Endpoint Auth Method. TokenEndpointAuthMethod is the requested Client Authentication method for the Token Endpoint. The token_endpoint_auth_methods_supported Discovery parameter contains a list of the authentication methods supported by this server. Valid types include client_secret_post, client_secret_basic, private_key_jwt, and none. */
  tokenEndpointAuthMethod?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Token Endpoint Auth Signing Algorithm. TokenEndpointAuthSigningAlgorithm is the requested Client Authentication signing algorithm for the Token Endpoint. The token_endpoint_auth_signing_alg_values_supported Discovery parameter contains a list of the supported signing algorithms for the token endpoint. */
  tokenEndpointAuthSigningAlgorithm?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Terms of Service URI. A URL string pointing to a human-readable terms of service document for the client that describes a contractual relationship between the end-user and the client that the end-user accepts when authorizing the client. */
  tosUri?: Maybe<Scalars['String']>;
  /** OAuth 2.0 Client Updated Date. UpdatedAt returns the timestamp of the client's last update. */
  updatedAt?: Maybe<Scalars['Time']>;
  /** OpenID Connect Userinfo Signed Response Algorithm. UserInfoSignedResponseAlg is a string containing the JWS signing algorithm (alg) parameter required for signing UserInfo Responses. The value none MAY be used, which indicates that the UserInfo Response will not be signed. The alg value RS256 MUST be used unless support for RS256 has been explicitly disabled. If support for RS256 has been disabled, the value none MUST be used. */
  userinfoSignedResponseAlgorithm?: Maybe<Scalars['String']>;
};

/** Representation a tenant in the Grafana observability stack where metrics, logs and traces can be sent to or retrieved from. */
export type ObservabilityTenant = {
  __typename?: 'ObservabilityTenant';
  /** The users and groups that can edit a tenant to add users, groups or oauth2 clients to it. */
  editors?: Maybe<ObservabilityTenantEditors>;
  /** The unique name of the tenant. */
  name: Scalars['String'];
  /** The organization that the tenant belongs to. */
  organization: Organization;
  /** The users that are admins of the organization. */
  viewers?: Maybe<ObservabilityTenantViewers>;
};

/** Representation of the users and groups that can edit a tenant. */
export type ObservabilityTenantEditors = {
  __typename?: 'ObservabilityTenantEditors';
  /** The groups that can edit a tenant. */
  groups?: Maybe<Array<Group>>;
  /** The users that can edit a tenant. */
  users?: Maybe<Array<User>>;
};

export type ObservabilityTenantEditorsInput = {
  /** The names of groups that can edit a tenant. */
  groups?: Maybe<Array<Scalars['String']>>;
  /** The IDs of users that can edit a tenant. */
  users?: Maybe<Array<Scalars['String']>>;
};

/** Representation of the users, groups and oauth2 clients that can view or send data a tenant. */
export type ObservabilityTenantViewers = {
  __typename?: 'ObservabilityTenantViewers';
  /** The groups that can view a tenant. */
  groups?: Maybe<Array<Group>>;
  /** The oauth2 clients that can send data a tenant. */
  oauth2Clients?: Maybe<Array<OAuth2Client>>;
  /** The users that can view a tenant. */
  users?: Maybe<Array<User>>;
};

export type ObservabilityTenantViewersInput = {
  /** The names of groups that can view a tenant. */
  groups?: Maybe<Array<Scalars['String']>>;
  /** The clientIDs oauth2 clients that can send data a tenant. */
  oauth2Clients?: Maybe<Array<Scalars['String']>>;
  /** The IDs of users that can view a tenant. */
  users?: Maybe<Array<Scalars['String']>>;
};

/** Representation an Organization in the auth stack. */
export type Organization = {
  __typename?: 'Organization';
  /** The users that are admins of the organization. */
  admins?: Maybe<Array<User>>;
  /** The unique name of the organization. */
  name: Scalars['String'];
};

export type Query = {
  __typename?: 'Query';
  getObservabilityTenant: ObservabilityTenant;
  /** Get a user by ID. */
  getUser: User;
  /** Get a list of all users. */
  listGroups?: Maybe<Array<Group>>;
  /** Get a list of all OAuth2 Clients. */
  listOAuth2Clients: Array<OAuth2Client>;
  /** Get a list of all users. */
  listObservabilityTenants: Array<ObservabilityTenant>;
  /** Get a list of all users. */
  listOrganizations: Array<Organization>;
  /** Get a list of all users. */
  listUsers: Array<User>;
  /** Get a single OAuth2 Client by ID. */
  oAuth2Client?: Maybe<OAuth2Client>;
};


export type QueryGetObservabilityTenantArgs = {
  name: Scalars['String'];
};


export type QueryGetUserArgs = {
  id: Scalars['ID'];
};


export type QueryOAuth2ClientArgs = {
  clientId: Scalars['ID'];
};

/** Representation of the information about a user sourced from Kratos. */
export type User = {
  __typename?: 'User';
  /** The user's email address. */
  email: Scalars['String'];
  /** The groups the user belongs to. */
  groups?: Maybe<Array<Group>>;
  /** The unique ID of the user. */
  id: Scalars['ID'];
  /** The user's full name. */
  name?: Maybe<Name>;
  /** The organization the user belongs to. */
  organization: Organization;
  /** The link a user can use to recover their account. */
  recoveryLink?: Maybe<Scalars['String']>;
};

export type ListGroupsQueryVariables = Exact<{ [key: string]: never; }>;


export type ListGroupsQuery = { __typename?: 'Query', listGroups?: Array<{ __typename?: 'Group', name: string, members?: Array<{ __typename?: 'User', id: string, email: string, name?: { __typename?: 'Name', first?: string | null | undefined, last?: string | null | undefined } | null | undefined }> | null | undefined }> | null | undefined };

export type DeleteGroupMutationVariables = Exact<{
  name: Scalars['String'];
}>;


export type DeleteGroupMutation = { __typename?: 'Mutation', deleteGroup: { __typename?: 'Group', name: string } };


export const ListGroupsDocument = gql`
    query ListGroups {
  listGroups {
    name
    members {
      id
      email
      name {
        first
        last
      }
    }
  }
}
    `;

/**
 * __useListGroupsQuery__
 *
 * To run a query within a React component, call `useListGroupsQuery` and pass it any options that fit your needs.
 * When your component renders, `useListGroupsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useListGroupsQuery({
 *   variables: {
 *   },
 * });
 */
export function useListGroupsQuery(baseOptions?: Apollo.QueryHookOptions<ListGroupsQuery, ListGroupsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<ListGroupsQuery, ListGroupsQueryVariables>(ListGroupsDocument, options);
      }
export function useListGroupsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<ListGroupsQuery, ListGroupsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<ListGroupsQuery, ListGroupsQueryVariables>(ListGroupsDocument, options);
        }
export type ListGroupsQueryHookResult = ReturnType<typeof useListGroupsQuery>;
export type ListGroupsLazyQueryHookResult = ReturnType<typeof useListGroupsLazyQuery>;
export type ListGroupsQueryResult = Apollo.QueryResult<ListGroupsQuery, ListGroupsQueryVariables>;
export const DeleteGroupDocument = gql`
    mutation DeleteGroup($name: String!) {
  deleteGroup(name: $name) {
    name
  }
}
    `;
export type DeleteGroupMutationFn = Apollo.MutationFunction<DeleteGroupMutation, DeleteGroupMutationVariables>;

/**
 * __useDeleteGroupMutation__
 *
 * To run a mutation, you first call `useDeleteGroupMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteGroupMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteGroupMutation, { data, loading, error }] = useDeleteGroupMutation({
 *   variables: {
 *      name: // value for 'name'
 *   },
 * });
 */
export function useDeleteGroupMutation(baseOptions?: Apollo.MutationHookOptions<DeleteGroupMutation, DeleteGroupMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteGroupMutation, DeleteGroupMutationVariables>(DeleteGroupDocument, options);
      }
export type DeleteGroupMutationHookResult = ReturnType<typeof useDeleteGroupMutation>;
export type DeleteGroupMutationResult = Apollo.MutationResult<DeleteGroupMutation>;
export type DeleteGroupMutationOptions = Apollo.BaseMutationOptions<DeleteGroupMutation, DeleteGroupMutationVariables>;