import { AdminDashboard } from "../pages/AdminDashboard";
import { AdminInstanceSettingsPage } from "../pages/AdminInstanceSettings";
import { AdminRunnerDetails } from "../pages/AdminRunnerDetails";
import { AdminRunners } from "../pages/AdminRunners";
import { AdminUserDetails } from "../pages/AdminUserDetails";
import { AdminUsersList } from "../pages/AdminUsersList";
import { CLIDownloadPage } from "../pages/CLIDownload";
import CliLogin from "../pages/CliLogin";
import CreateWorkspace from "../pages/CreateWorkspace";
import { CreditsPage } from "../pages/Credits";
import { EmailNotVerifiedPage } from "../pages/EmailNotVerified";
import HomePage from "../pages/Home";
import LoginPage from "../pages/Login";
import Profile from "../pages/Profile";
import SignUpPage from "../pages/SignUp";
import { TemplateDetailsPage } from "../pages/TemplateDetails";
import TemplatesList from "../pages/TemplatesList";
import { TemplateVersionEditor } from "../pages/TemplateVersionEditor";
import { VerifyEmailPage } from "../pages/VerifyEmail";
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
  {
    path: "email-not-verified",
    element: <EmailNotVerifiedPage />,
  },
  {
    path: "verify-email",
    element: <VerifyEmailPage />,
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
  {
    path: "/templates",
    element: <TemplatesList />,
  },
  {
    path: "/templates/:id",
    element: <TemplateDetailsPage />,
  },
  {
    path: "/templates/:templateId/versions/:versionId/editor",
    element: <TemplateVersionEditor />, // TODO: protect view, pnly template managers and admin can view this page
    showNavbar: false,
  },
  {
    path: "/credits",
    element: <CreditsPage />,
  },
  {
    path: "/cli",
    element: <CLIDownloadPage />,
  },
];

export const SuperUserRoutes: Route[] = [
  {
    path: "/admin",
    element: <AdminDashboard />,
  },
  {
    path: "/admin/users",
    element: <AdminUsersList />,
  },
  {
    path: "/admin/users/:email",
    element: <AdminUserDetails />,
  },
  {
    path: "/admin/runners",
    element: <AdminRunners />,
  },
  {
    path: "/admin/runners/:id",
    element: <AdminRunnerDetails />,
  },
  {
    path: "/admin/settings",
    element: <AdminInstanceSettingsPage />,
  },
];
