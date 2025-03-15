import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import "@tabler/core/dist/css/tabler.min.css";
import "bootstrap/dist/js/bootstrap.js";
import AuthRequired from "./components/AuthRequired";
import NotFound from "./pages/NotFound";
import SuperUserRequired from "./components/SuperUserRequired";
import {
  AuthProtectedRoutes,
  PublicRoutes,
  SuperUserRoutes,
} from "./routes/routes";

export default function App() {
  return (
    <Router>
      <Routes>
        {PublicRoutes.map((r, i) => (
          <Route path={r.path} element={r.element} key={`pu-${i}`} />
        ))}
        {AuthProtectedRoutes.map((r, i) => (
          <Route
            path={r.path}
            element={
              <AuthRequired showNavbar={r.showNavbar} key={`ar-${i}`}>
                {r.element}
              </AuthRequired>
            }
          />
        ))}
        {SuperUserRoutes.map((r, i) => {
          return (
            <Route
              path={r.path}
              element={
                <SuperUserRequired showNavbar={r.showNavbar} key={`su-${i}`}>
                  {r.element}
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
