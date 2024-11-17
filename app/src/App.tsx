import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import LoginPage from './pages/Login';
import HomePage from "./pages/Home";
import WorkspaceDetails from "./pages/WorkspaceDetails/WorkspaceDetails";
import PageNotFound from "./pages/PageNotFound";
import CreateWorkspace from "./pages/CreateWorkspace";
import Profile from "./pages/Profile";

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/login/' element={<LoginPage />} />
        <Route path='/' element={<HomePage />} />
        <Route path='/workspaces/:workspaceId' element={<WorkspaceDetails />}  />
        <Route path='/create-workspace' element={<CreateWorkspace />}  />
        <Route path='/profile' element={<Profile />}  />
        <Route path='*' element={<PageNotFound />}  />
      </Routes>
    </Router>
  );
}

export default App;
