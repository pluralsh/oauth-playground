import { AxiosError } from 'axios';
import { useState, useEffect, DependencyList } from 'react';
import ory from '../apis/ory';

// Returns a function which will log the user out
export function useLogoutHandler(deps?: DependencyList) {
  const [logoutToken, setLogoutToken] = useState<string>('');

  useEffect(() => {
    ory
      .createSelfServiceLogoutFlowUrlForBrowsers()
      .then(({ data }) => {
        setLogoutToken(data.logout_token);
      })
      .catch((err: AxiosError) => {
        switch (err.response?.status) {
          case 401:
            // do nothing, the user is not logged in
            return;
        }

        // Something else happened!
        return Promise.reject(err);
      });
  }, deps);

  return () => {
    if (logoutToken) {
      ory
        .submitSelfServiceLogoutFlow(logoutToken)
        .then(() => window.location.replace('/login'));
    }
  };
}
