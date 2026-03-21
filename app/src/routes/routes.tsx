import { AdminAnalyticsPage } from "../pages/AdminAnalytics";
import { Route } from "./types";
import { lazy } from "react";

const LoginPage = lazy(() => import("../pages/Login"));
const SignUpPage = lazy(() => import("../pages/SignUp"));
const HomePage = lazy(() => import("../pages/Home"));
const CreateWorkspace = lazy(() => import("../pages/CreateWorkspace"));
const WorkspaceDetails = lazy(() => import("../pages/WorkspaceDetails"));
const Profile = lazy(() => import("../pages/Profile"));
const CliLogin = lazy(() => import("../pages/CliLogin"));
const TemplatesList = lazy(() => import("../pages/TemplatesList"));
const TemplateDetailsPage = lazy(() => import("../pages/TemplateDetails"));
const AdminDashboard = lazy(() => import("../pages/AdminDashboard"));
const AdminUsersList = lazy(() => import("../pages/AdminUsersList"));
const AdminUserDetails = lazy(() => import("../pages/AdminUserDetails"));
const AdminRunners = lazy(() => import("../pages/AdminRunners"));
const AdminRunnerDetails = lazy(() => import("../pages/AdminRunnerDetails"));
const AdminAuthenticationSettingsPage = lazy(() => import("../pages/AdminAuthenticationSettings"));
const AdminEmailSenderPage = lazy(() => import("../pages/AdminEmailSender"));
const CLIDownloadPage = lazy(() => import("../pages/CLIDownload"));
const CreditsPage = lazy(() => import("../pages/Credits"));
const VerifyEmailPage = lazy(() => import("../pages/VerifyEmail"));
const EmailNotVerifiedPage = lazy(() => import("../pages/EmailNotVerified"));
const PasswordResetPage = lazy(() => import("../pages/PasswordReset"));
const PasswordResetSentPage = lazy(() => import("../pages/PasswordResetSent"));
const PasswordResetFromTokenPage = lazy(() => import("../pages/PasswordResetFromToken"));
const TemplateVersionEditor = lazy(() => import("../pages/TemplateVersionEditor"));


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
  {
    path: "password-reset",
    element: <PasswordResetPage />,
  },
  {
    path: "password-reset/sent",
    element: <PasswordResetSentPage />,
  },
  {
    path: "password-reset/from-token",
    element: <PasswordResetFromTokenPage />,
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
    path: "/admin/authentication-settings",
    element: <AdminAuthenticationSettingsPage />,
  },
  {
    path: "/admin/email-sender",
    element: <AdminEmailSenderPage />,
  },
  {
    path: "/admin/analytics",
    element: <AdminAnalyticsPage />,
  },
];
