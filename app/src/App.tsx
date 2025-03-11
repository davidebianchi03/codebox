import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import LoginPage from "./pages/Login";
import HomePage from "./pages/Home";
import "@tabler/core/dist/css/tabler.min.css";
import "bootstrap/dist/js/bootstrap.js";
import AuthRequired from "./components/AuthRequired";
import CreateWorkspace from "./pages/CreateWorkspace";
import WorkspaceDetails from "./pages/WorkspaceDetails";
import Profile from "./pages/Profile";
import NotFound from "./pages/NotFound";
import { AdminRunners } from "./pages/AdminRunners";
import SuperUserRequired from "./components/SuperUserRequired";

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="login" element={<LoginPage />} />
        <Route
          path=""
          element={
            <AuthRequired>
              <HomePage />
            </AuthRequired>
          }
        />
        <Route
          path="/create-workspace"
          element={
            <AuthRequired>
              <CreateWorkspace />
            </AuthRequired>
          }
        />
        <Route
          path="/workspaces/:id"
          element={
            <AuthRequired>
              <WorkspaceDetails />
            </AuthRequired>
          }
        />
        <Route
          path="/profile"
          element={
            <AuthRequired>
              <Profile />
            </AuthRequired>
          }
        />
        <Route
          path="/admin/runners"
          element={
            <SuperUserRequired>
              <AdminRunners />
            </SuperUserRequired>
          }
        />
        <Route element={<NotFound />} path="*" />
      </Routes>
    </Router>
  );
}
