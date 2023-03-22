import {
  BrowserRouter,
  Routes,
  Route,
  Navigate,
} from 'react-router-dom';
import Notes from './pages/Notes';
import { Login } from "./pages/Login"
import { Settings } from './pages/Settings';
// import './App.css';
import Layout from './components/Layout';
import Users from './pages/Users';
import Groups from './pages/Groups';
import { Verification } from "./pages/Verification"
import { Error } from "./pages/Error"
import { Registration } from './pages/Registration';
import { Recovery } from './pages/Recovery';
import { Consent } from './pages/Consent';

function App() {
    return (
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<Notes />} />
            <Route path="users" element={<Users />} />
            <Route path="groups" element={<Groups />} />
            <Route path="settings" element={<Settings />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
          <Route path="/registration" element={<Registration />} />
          <Route path="/login" element={<Login />} />
          <Route path="/verification" element={<Verification />} />
          <Route path="/error" element={<Error />} />
          <Route path="/recovery" element={<Recovery />} />
          <Route path="/consent" element={<Consent />} />
        </Routes>
      </BrowserRouter>
    )
}

export default App;
