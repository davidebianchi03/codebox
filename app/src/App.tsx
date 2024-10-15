import {
  BrowserRouter as Router,
  Routes,
  Route,
} from "react-router-dom";
import LoginPage from './pages/Login';
import HomePage from "./pages/Home";
import WorkspaceDetails from "./pages/WorkspaceDetails";

function App() {
  return (
    <Router>
      <Routes>
        <Route path='/login/' element={<LoginPage />} />
        <Route path='/' element={<HomePage />} />
        <Route path='/workspaces/:workspaceId' element={<WorkspaceDetails />}  />
      </Routes>
    </Router>
  );
}

export default App;
