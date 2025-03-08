import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import LoginPage from './pages/Login';
import HomePage from "./pages/Home";
import "@tabler/core/dist/css/tabler.min.css";
import 'bootstrap/dist/js/bootstrap.js';
import AuthRequired from "./pages/AuthRequired";
import CreateWorkspace from "./pages/CreateWorkspace";
import WorkspaceDetails from "./pages/WorkspaceDetails";


export default function App() {
  return (
    <Router>
      <Routes>
        <Route path='' element={<AuthRequired><HomePage /></AuthRequired>} />
        <Route path='/create-workspace' element={<AuthRequired><CreateWorkspace /></AuthRequired>} />
        <Route path='/workspaces/:id' element={<AuthRequired><WorkspaceDetails /></AuthRequired>} />
      </Routes>
      <Routes>
        <Route path='login' element={<LoginPage />} />
      </Routes>
    </Router>
  );
}
