import { AdminRunners } from "../pages/AdminRunners";
import CliLogin from "../pages/CliLogin";
import CreateWorkspace from "../pages/CreateWorkspace";
import HomePage from "../pages/Home";
import LoginPage from "../pages/Login";
import Profile from "../pages/Profile";
import SignUpPage from "../pages/SignUp";
import WorkspaceDetails from "../pages/WorkspaceDetails";
import { Route } from "./types";

export const PublicRoutes: Route[] = [
  {
    path: "login",
    element: <LoginPage />,
  },
  {
    path: "signup",
    element: <SignUpPage />,
  },
];

export const AuthProtectedRoutes: Route[] = [
  {
    path: "",
    element: <HomePage />,
  },
  {
    path: "/create-workspace",
    element: <CreateWorkspace />,
  },
  {
    path: "/workspaces/:id",
    element: <WorkspaceDetails />,
  },
  {
    path: "/profile",
    element: <Profile />,
  },
  {
    path: "/cli-login",
    element: <CliLogin />,
    showNavbar: false,
  },
];

export const SuperUserRoutes: Route[] = [
  {
    path: "/admin/runners",
    element: <AdminRunners />,
  },
];
