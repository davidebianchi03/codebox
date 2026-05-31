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
import React, { Suspense } from "react";
import { EmptyLayout } from "./layouts/EmptyLayout";
import LoadingFallback from "./components/LoadingFallback";

axios.defaults.withCredentials = true;
axios.defaults.baseURL = import.meta.env.VITE_SERVER_URL;

export default function App() {
  console.log(
    `
   ██████╗ ██████╗ ██████╗ ███████╗██████╗  ██████╗ ██╗  ██╗
  ██╔════╝██╔═══██╗██╔══██╗██╔════╝██╔══██╗██╔═══██╗╚██╗██╔╝
  ██║     ██║   ██║██║  ██║█████╗  ██████╔╝██║   ██║ ╚███╔╝
  ██║     ██║   ██║██║  ██║██╔══╝  ██╔══██╗██║   ██║ ██╔██╗
  ╚██████╗╚██████╔╝██████╔╝███████╗██████╔╝╚██████╔╝██╔╝ ██╗
   ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═════╝  ╚═════╝ ╚═╝  ╚═╝

   Copyright (c) ${new Date().getFullYear()} Davide Bianchi
   Codebox is licensed under MIT license
   Version: ${import.meta.env.VITE_APP_VERSION}
    `
  );

  return (
    <Router>
      <Routes>
        {PublicRoutes.map((r, i) => (
          <Route
            path={r.path}
            element={
              <Suspense fallback={<LoadingFallback />}>
                {r.element}
              </Suspense>
            }
            key={i}
          />
        ))}
        {AuthProtectedRoutes.map((r, i) => (
          <Route
            key={i}
            path={r.path}
            element={
              <Suspense fallback={<LoadingFallback />}>
                <AuthRequired showNavbar={r.showNavbar}>
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
                </AuthRequired>
              </Suspense>
            }
          />
        ))}
        {SuperUserRoutes.map((r, i) => {
          return (
            <Route
              key={i}
              path={r.path}
              element={
                <Suspense fallback={<LoadingFallback />}>
                  <SuperUserRequired showNavbar={r.showNavbar}>
                    <SidebarLayout sidebarItems={SuperUserSidebarItems}>
                      {r.element}
                    </SidebarLayout>
                  </SuperUserRequired>
                </Suspense>
              }
            />
          );
        })}
        <Route element={<NotFound />} path="*" />
      </Routes>
    </Router>
  );
}
