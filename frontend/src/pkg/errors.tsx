import { AxiosError } from 'axios';
import { Dispatch, SetStateAction } from 'react';
import { NavigateFunction } from 'react-router-dom';

// A small function to help us deal with errors coming from fetching a flow.
function handleGetFlowError<S>(
  navigate: NavigateFunction,
  flowType: 'login' | 'registration' | 'settings' | 'recovery' | 'verification',
  resetFlow: Dispatch<SetStateAction<S | undefined>>
) {
  return async (err: any) => {
    switch (err.response?.data.error?.id) {
      case 'session_aal2_required':
        // 2FA is enabled and enforced, but user did not perform 2fa yet!
        console.log('session_aal2_required');
        window.location.href = err.response?.data.redirect_browser_to;
        return;
      case 'session_already_available':
        // User is already signed in, let's redirect them home!
        console.log('session_already_available');
        navigate('/');
        return;
      case 'session_refresh_required':
        // We need to re-authenticate to perform this action
        console.log('session_refresh_required');
        window.location.href = err.response?.data.redirect_browser_to;
        return;
      case 'self_service_flow_return_to_forbidden':
        // The flow expired, let's request a new one.
        console.log('self_service_flow_return_to_forbidden');
        resetFlow(undefined);
        navigate('/' + flowType);
        return;
      case 'self_service_flow_expired':
        // The flow expired, let's request a new one.
        console.log('self_service_flow_expired');
        resetFlow(undefined);
        navigate('/' + flowType);
        return;
      case 'security_csrf_violation':
        // A CSRF violation occurred. Best to just refresh the flow!
        console.log('security_csrf_violation');
        resetFlow(undefined);
        navigate('/' + flowType);
        return;
      case 'security_identity_mismatch':
        // The requested item was intended for someone else. Let's request a new flow...
        console.log('security_identity_mismatch');
        resetFlow(undefined);
        navigate('/' + flowType);
        return;
      case 'browser_location_change_required':
        // Ory Kratos asked us to point the user to this URL.
        console.log('browser_location_change_required');
        window.location.href = err.response.data.redirect_browser_to;
        return;
    }

    switch (err.response?.status) {
      case 410:
        // The flow expired, let's request a new one.
        console.log('expired_flow');
        resetFlow(undefined);
        navigate('/' + flowType);
        return;
    }

    // We are not able to handle the error? Return it.
    return Promise.reject(err);
  };
}

export const handleFlowError = handleGetFlowError;
