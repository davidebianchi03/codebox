import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "@tabler/core/dist/css/tabler.min.css";
import "./assets/scss/custom.scss";
import AuthRequired from "./auth/AuthRequired";
import NotFound from "./pages/NotFound";
import SuperUserRequired from "./auth/SuperUserRequired";
import {
  AuthProtectedRoutes,
  PublicRoutes,
  SuperUserRoutes,
} from "./routes/routes";
import axios from "axios";
import { NavbarLayout } from "./layouts/NavbarLayout";
import { SidebarLayout } from "./layouts/SidebarLayout";
import { SuperUserSidebarItems } from "./layouts/SidebarItems";
import React from "react";
import { EmptyLayout } from "./layouts/EmptyLayout";
import { SettingsProvider } from "./components/SettingsProvider";

axios.defaults.withCredentials = true;
axios.defaults.baseURL = import.meta.env.VITE_SERVER_URL;

export default function App() {
  return (
    <Router>
      <Routes>
        {PublicRoutes.map((r, i) => (
          <Route path={r.path} element={r.element} key={i} />
        ))}
        {AuthProtectedRoutes.map((r, i) => (
          <Route
            key={i}
            path={r.path}
            element={
              <AuthRequired showNavbar={r.showNavbar}>
                <SettingsProvider>
                  {r.showNavbar === true || r.showNavbar === undefined ? (
                    <React.Fragment>
                      <NavbarLayout>
                        {r.element}
                      </NavbarLayout>
                    </React.Fragment>
                  ) : (
                    <React.Fragment>
                      <EmptyLayout>
                        {r.element}
                      </EmptyLayout>
                    </React.Fragment>
                  )}
                </SettingsProvider>
              </AuthRequired>
            }
          />
        ))}
        {SuperUserRoutes.map((r, i) => {
          return (
            <Route
              key={i}
              path={r.path}
              element={
                <SuperUserRequired showNavbar={r.showNavbar}>
                  <SettingsProvider>
                    <SidebarLayout sidebarItems={SuperUserSidebarItems}>
                      {r.element}
                    </SidebarLayout>
                  </SettingsProvider>
                </SuperUserRequired>
              }
            />
          );
        })}
        <Route element={<NotFound />} path="*" />
      </Routes>
    </Router>
  );
}
