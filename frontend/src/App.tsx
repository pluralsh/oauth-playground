import { useEffect, useState } from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate
} from 'react-router-dom';
import Notes from './pages/Notes';
import Login from './pages/Login';
import Settings from './pages/Settings';
import './App.css';
import Layout from './components/Layout';
import Users from './pages/Users';
import Groups from './pages/Groups';
import Registration from './pages/Registration';
import Recovery from './pages/Recovery';
import { Session } from '@ory/client';
import ory from './apis/ory';
import { CircularProgress } from '@mui/material';

function App() {
  const [status, setStatus] = useState('idle');
  const [session, setSession] = useState<Session | undefined>();

  const isLoading = status === 'loading';
  const isSuccess = status === 'success';
  const isError = status === 'error';

  useEffect(() => {
    setStatus('loading');
    ory
      .toSession()
      .then(({ data }) => {
        setSession(data);
        setStatus('success');
      })
      .catch(err => {
        setStatus('error');
        console.log(err.message);
      });
  }, []);

  if (isLoading) {
    return <CircularProgress />;
  }

  if (isSuccess || isError) {
    return session ? (
      <Router>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<Notes />} />
            {/* <Route path="storage" element={<Storage />} /> */}
            <Route path="settings" element={<Settings />} />
            {/* <Route path="workspaces" element={<Workspaces />}>
              <Route path=":name" element={<WorkspaceDetails />} />
              <Route path="create" element={<CreateWorkspaceForm />} />
            </Route> */}
            <Route path="users" element={<Users />} />
            <Route path="groups" element={<Groups />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </Router>
    ) : (
      <Router>
        <Routes>
          <Route path="/registration" element={<Registration />} />
          <Route path="/login" element={<Login />} />
          <Route path="/recovery" element={<Recovery />} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </Router>
    );
  }

  return null;
}

export default App;
