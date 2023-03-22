import { LoginFlow, UpdateLoginFlowBody, OAuth2Client, OAuth2ConsentRequest } from "@ory/client"
import { UserAuthCard, UserConsentCard } from "@ory/elements"
import { useCallback, useEffect, useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { CircularProgress } from '@mui/material'
import { sdk, sdkError } from "../apis/ory"
import { useOAuth2ConsentRequestQuery, useAcceptOAuth2ConsentRequestMutation } from "../generated/graphql"

export const Consent = (): JSX.Element => {
  // const [flow, setFlow] = useState<LoginFlow | null>(null)
  const [searchParams, setSearchParams] = useSearchParams()
  const win: Window = window;

  const navigate = useNavigate()

  const challenge = searchParams.get("consent_challenge")

  if (!challenge) {
    return <div>There is no consent challenge</div>
  }

  const { data } = useOAuth2ConsentRequestQuery({
    variables: {
      challenge: challenge
    }
  })

  const [mutation, { loading, error }] = useAcceptOAuth2ConsentRequestMutation({
    variables: {
      challenge,
      grantScope: data?.oauth2ConsentRequest?.requestedScope || ['profile', 'openid'],
      remember: data?.oauth2ConsentRequest?.skip,
      // rememberFor: 3600,
      // session: // TODO: need to parse using the subject and scopes. See https://github.com/ory/kratos-selfservice-ui-node/pull/248/files#diff-f55c47595a4b4dc1dc448defc15f0157e124c1f8241c25474835948ca51be903R24
    },
    onCompleted: ({ acceptOAuth2ConsentRequest: { redirectTo } }) => {
      win.location = redirectTo
    },
  })

  if (data?.oauth2ConsentRequest?.skip) {
    mutation(
      {
        variables: {
          challenge,
          grantScope: data?.oauth2ConsentRequest?.requestedScope,
          remember: data?.oauth2ConsentRequest?.skip,

        },
      },
    )
  }


  // we check if the flow is set, if not we show a loading indicator
  return data?.oauth2ConsentRequest ? (
    <UserConsentCard
      csrfToken="csrfToken"
      consent={data.oauth2ConsentRequest as OAuth2ConsentRequest}
      cardImage={data?.oauth2ConsentRequest?.client?.logoUri || "/logo192.png"}
      client_name="Ory Kratos"
      requested_scope={data?.oauth2ConsentRequest?.requestedScope || []}
      client={data?.oauth2ConsentRequest?.client as OAuth2Client}
      action={(process.env.BASE_URL || "") + "/consent"}
      />
  ) : (
    <CircularProgress />
  )
}
