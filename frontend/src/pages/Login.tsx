import { useEffect, useState } from 'react';
import {
  SelfServiceLoginFlow,
  SubmitSelfServiceLoginFlowBody
} from '@ory/client';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import {
  Avatar,
  Button,
  CssBaseline,
  TextField,
  Grid,
  Box,
  Typography,
  Container
} from '@mui/material';
import { LockOutlined } from '@mui/icons-material';
import { handleFlowError } from '../pkg/errors';
import ory from '../apis/ory';
import { Flow } from '../pkg/ui';
import { AxiosError } from 'axios';

function Login() {
  const [flow, setFlow] = useState<SelfServiceLoginFlow>();
  const [searchParams, setSearchParams] = useSearchParams();

  const flowId = searchParams.get('flow');
  const returnTo = searchParams.get('return_to');
  const refresh = searchParams.get('refresh');
  const aal = searchParams.get('aal');

  const navigate = useNavigate();

  useEffect(() => {
    if (flow) {
      return;
    }

    if (flowId) {
      ory
        .getSelfServiceLoginFlow(String(flowId))
        .then(({ data }: any) => {
          setFlow(data);
        })
        .catch(handleFlowError(navigate, 'login', setFlow));

      return;
    }

    ory
      .initializeSelfServiceLoginFlowForBrowsers(
        Boolean(refresh),
        aal ? String(aal) : undefined,
        returnTo ? String(returnTo) : undefined
      )
      .then(({ data }: any) => {
        setFlow(data);
      })
      .catch(handleFlowError(navigate, 'login', setFlow));
  }, [flowId, aal, refresh, returnTo, flow]);

  const onSubmit = (values: SubmitSelfServiceLoginFlowBody) => {
    navigate(`/login?flow=${flow?.id}`);

    ory
      .submitSelfServiceLoginFlow(String(flow?.id), values, undefined)
      .then(res => {
        if (flow?.return_to) {
          window.location.href = flow?.return_to;
          return;
        }

        window.location.replace('/');
      })
      .catch(handleFlowError(navigate, 'login', setFlow))
      .catch((err: AxiosError<SelfServiceLoginFlow>) => {
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
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Box
        sx={{
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center'
        }}
      >
        <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
          <LockOutlined />
        </Avatar>
        <Typography component="h1" variant="h5" mb={1}>
          Sign in
        </Typography>
        <Flow onSubmit={onSubmit} flow={flow} />
        <Grid container>
          <Grid item xs>
            <Link to="/recovery">
              <Typography variant="body2">Forgot password?</Typography>
            </Link>
          </Grid>
          <Grid item>
            <Link to="/registration">
              <Typography variant="body2">
                Don't have an account? Sign Up
              </Typography>
            </Link>
          </Grid>
        </Grid>
      </Box>
    </Container>
  );
}

export default Login;
