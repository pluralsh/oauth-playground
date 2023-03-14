import {
  SelfServiceSettingsFlow,
  SubmitSelfServiceSettingsFlowBody
} from '@ory/client';
import { AxiosError } from 'axios';
import { ReactNode, useEffect, useState } from 'react';
import { Typography, Box } from '@mui/material';

import { Flow, Methods, Messages } from '../pkg/ui';
import { handleFlowError } from '../pkg/errors';
import ory from '../apis/ory';
import { useNavigate, useSearchParams } from 'react-router-dom';

interface Props {
  flow?: SelfServiceSettingsFlow;
  only?: Methods;
}

function SettingsCard({
  flow,
  only,
  children
}: Props & { children: ReactNode }) {
  if (!flow) {
    return null;
  }

  const nodes = only
    ? flow.ui.nodes.filter(({ group }) => group === only)
    : flow.ui.nodes;

  if (nodes.length === 0) {
    return null;
  }

  return <Box mb={2}>{children}</Box>;
}

function Settings() {
  const [flow, setFlow] = useState<SelfServiceSettingsFlow>();

  // Get ?flow=... from the URL
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const flowId = searchParams.get('flow');
  const returnTo = searchParams.get('return_to');

  useEffect(() => {
    if (flow) {
      return;
    }

    // If ?flow=.. was in the URL, we fetch it
    if (flowId) {
      ory
        .getSelfServiceSettingsFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleFlowError(navigate, 'settings', setFlow));
      return;
    }

    // Otherwise we initialize it
    ory
      .initializeSelfServiceSettingsFlowForBrowsers(
        returnTo ? String(returnTo) : undefined
      )
      .then(({ data }) => {
        setFlow(data);
      })
      .catch(handleFlowError(navigate, 'settings', setFlow));
  }, [flowId, returnTo, flow]);

  const onSubmit = (values: SubmitSelfServiceSettingsFlowBody) => {
    navigate(`/settings?flow=${flow?.id}`, undefined);
    // On submission, add the flow ID to the URL but do not navigate. This prevents the user loosing
    // his data when she/he reloads the page.

    ory
      .submitSelfServiceSettingsFlow(String(flow?.id), values, undefined)
      .then(({ data }) => {
        // The settings have been saved and the flow was updated. Let's show it to the user!
        setFlow(data);
      })
      .catch(handleFlowError(navigate, 'settings', setFlow))
      .catch((err: AxiosError<SelfServiceSettingsFlow>) => {
        // If the previous handler did not catch the error it's most likely a form validation error
        if (err.response?.status === 400) {
          // Yup, it is!
          setFlow(err.response?.data);
          return;
        }

        return Promise.reject(err);
      });
  };

  return (
    <Box sx={{ margin: 4, marginBottom: 0 }}>
      <SettingsCard only="profile" flow={flow}>
        <Typography variant="h4">Profile Settings</Typography>
        <Messages messages={flow?.ui.messages} />
        <Flow
          hideGlobalMessages
          onSubmit={onSubmit}
          only="profile"
          flow={flow}
        />
      </SettingsCard>
      <SettingsCard only="password" flow={flow}>
        <Typography variant="h4">Change Password</Typography>
        <Messages messages={flow?.ui.messages} />
        <Flow
          hideGlobalMessages
          onSubmit={onSubmit}
          only="password"
          flow={flow}
        />
      </SettingsCard>
      <SettingsCard only="oidc" flow={flow}>
        <Typography variant="h4">Manage Social Sign In</Typography>
        <Messages messages={flow?.ui.messages} />
        <Flow hideGlobalMessages onSubmit={onSubmit} only="oidc" flow={flow} />
      </SettingsCard>
      <SettingsCard only="lookup_secret" flow={flow}>
        <Typography variant="h4">Manage 2FA Backup Recovery Codes</Typography>
        <Messages messages={flow?.ui.messages} />
        <Typography paragraph>
          Recovery codes can be used in panic situations where you have lost
          access to your 2FA device.
        </Typography>
        <Flow
          hideGlobalMessages
          onSubmit={onSubmit}
          only="lookup_secret"
          flow={flow}
        />
      </SettingsCard>
      <SettingsCard only="totp" flow={flow}>
        <Typography variant="h4">Manage 2FA TOTP Authenticator App</Typography>
        <Typography paragraph>
          Add a TOTP Authenticator App to your account to improve your account
          security. Popular Authenticator Apps are{' '}
          <a href="https://www.lastpass.com" rel="noreferrer" target="_blank">
            LastPass
          </a>{' '}
          and Google Authenticator (
          <a
            href="https://apps.apple.com/us/app/google-authenticator/id388497605"
            target="_blank"
            rel="noreferrer"
          >
            iOS
          </a>
          ,{' '}
          <a
            href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en&gl=US"
            target="_blank"
            rel="noreferrer"
          >
            Android
          </a>
          ).
        </Typography>
        <Messages messages={flow?.ui.messages} />
        <Flow hideGlobalMessages onSubmit={onSubmit} only="totp" flow={flow} />
      </SettingsCard>
      <SettingsCard only="webauthn" flow={flow}>
        <Typography variant="h4">
          Manage Hardware Tokens and Biometrics
        </Typography>
        <Messages messages={flow?.ui.messages} />
        <Typography paragraph>
          Use Hardware Tokens (e.g. YubiKey) or Biometrics (e.g. FaceID,
          TouchID) to enhance your account security.
        </Typography>
        <Flow
          hideGlobalMessages
          onSubmit={onSubmit}
          only="webauthn"
          flow={flow}
        />
      </SettingsCard>
    </Box>
  );
}

export default Settings;
