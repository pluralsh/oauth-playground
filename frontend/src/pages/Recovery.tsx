import { useState, useEffect } from 'react';
import {
  SelfServiceRecoveryFlow,
  SubmitSelfServiceRecoveryFlowBody
} from '@ory/client';
import { Flow } from '../pkg/ui';
import { handleFlowError } from '../pkg/errors';
import ory from '../apis/ory';
import { useNavigate, useSearchParams, Link } from 'react-router-dom';
import { Container, CssBaseline, Box, Typography, Grid } from '@mui/material';
import { AxiosError } from 'axios';

function Recovery() {
  const [flow, setFlow] = useState<SelfServiceRecoveryFlow>();

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
        .getSelfServiceRecoveryFlow(String(flowId))
        .then(({ data }) => {
          setFlow(data);
        })
        .catch(handleFlowError(navigate, 'recovery', setFlow));
      return;
    }

    ory
      .initializeSelfServiceRecoveryFlowForBrowsers(
        returnTo ? String(returnTo) : undefined
      )
      .then(({ data }: any) => {
        setFlow(data);
      })
      .catch(handleFlowError(navigate, 'recovery', setFlow))
      .catch((err: AxiosError<SelfServiceRecoveryFlow>) => {
        // If the previous handler did not catch the error it's most likely a form validation error
        if (err.response?.status === 400) {
          // Yup, it is!
          setFlow(err.response?.data);
          return;
        }

        return Promise.reject(err);
      });
  }, [flowId, returnTo, flow]);

  const onSubmit = (values: SubmitSelfServiceRecoveryFlowBody) => {
    navigate(`/recovery?flow=${flow?.id}`);
    ory
      .submitSelfServiceRecoveryFlow(String(flow?.id), values)
      .then(({ data }) => {
        if (flow?.return_to) {
          window.location.href = flow?.return_to;
          return;
        }

        window.location.replace('/login');
      })
      .catch(handleFlowError(navigate, 'recovery', setFlow))
      .catch((err: AxiosError<SelfServiceRecoveryFlow>) => {
        // If the previous handler did not catch the error it's most likely a form validation error
        if (err.response?.status === 400) {
          // Yup, it is!
          setFlow(err.response?.data);
          return;
        }

        throw err;
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
        <Typography component="h1" variant="h5" mb={1}>
          Recover your account
        </Typography>
        <Flow onSubmit={onSubmit} flow={flow} />
        <Grid container justifyContent="flex-end">
          <Grid item>
            <Link to="/login">
              <Typography variant="body2">Go back</Typography>
            </Link>
          </Grid>
        </Grid>
      </Box>
    </Container>
  );
}

export default Recovery;
