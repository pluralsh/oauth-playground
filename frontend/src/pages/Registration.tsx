import { useState, useEffect, Dispatch, SetStateAction } from 'react';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import {
  SelfServiceRegistrationFlow,
  SubmitSelfServiceRegistrationFlowBody
} from '@ory/client';
import {
  Avatar,
  CssBaseline,
  Grid,
  Box,
  Typography,
  Container
} from '@mui/material';
import { LockOutlined } from '@mui/icons-material';
import ory from '../apis/ory';
import { Flow } from '../pkg/ui';
import { handleFlowError } from '../pkg/errors';
import { AxiosError } from 'axios';

function Registration() {
  const [flow, setFlow] = useState<SelfServiceRegistrationFlow>();

  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();

  const flowId = searchParams.get('flow');
  const returnTo = searchParams.get('return_to');

  useEffect(() => {
    if (flow) {
      return;
    }

    if (flowId) {
      ory
        .getSelfServiceRegistrationFlow(String(flowId))
        .then(({ data }: any) => {
          setFlow(data);
        })
        .catch(handleFlowError(navigate, 'registration', setFlow));

      return;
    }

    ory
      .initializeSelfServiceRegistrationFlowForBrowsers(
        returnTo ? String(returnTo) : undefined
      )
      .then(({ data }: any) => {
        setFlow(data);
      })
      .catch(handleFlowError(navigate, 'registration', setFlow));
  }, [flowId, returnTo, flow]);

  const onSubmit = (values: SubmitSelfServiceRegistrationFlowBody) => {
    navigate(`/registration?flow=${flow?.id}`);
    ory
      .submitSelfServiceRegistrationFlow(String(flow?.id), values)
      .then(({ data }) => {
        if (flow?.return_to) {
          window.location.href = flow?.return_to;
          return;
        }

        window.location.replace('/');
      })
      .catch(handleFlowError(navigate, 'registration', setFlow))
      .catch((err: AxiosError<SelfServiceRegistrationFlow>) => {
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
          Sign up
        </Typography>
        <Flow onSubmit={onSubmit} flow={flow} />
        <Grid container justifyContent="flex-end">
          <Grid item>
            <Link to="/login">
              <Typography variant="body2">
                Already have an account? Sign in
              </Typography>
            </Link>
          </Grid>
        </Grid>
      </Box>
    </Container>
  );
}

export default Registration;
