import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import DefaultAvatar from "../assets/images/default-avatar.png";
import sha256 from "crypto-js/sha256";
import { Http } from "../api/http";
import { Link, useNavigate } from "react-router-dom";
import { InstanceSettings } from "../types/settings";
import { useCallback, useEffect, useState } from "react";
import { RequestStatus } from "../api/types";
import { User } from "../types/user";

interface Props {
  user: User;
}

export function Navbar({ user }: Props) {
  const navigate = useNavigate();
  const [settings, setSettings] = useState<InstanceSettings | null>(null);

  const HandleLogout = (e: any) => {
    e.preventDefault();
    Http.Request(`${Http.GetServerURL()}/api/v1/auth/logout`, "POST", null);
    document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=${window.location.hostname}`;
    document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=.${window.location.hostname}`;
    if (process.env.NODE_ENV === "development") {
      document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=${new URL(Http.GetServerURL()).hostname
        }`;
    }
    navigate("/login");
  };

  const FetchSettings = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/instance-settings`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setSettings(responseBody as InstanceSettings);
    }
  }, []);

  useEffect(() => {
    FetchSettings();
  }, [FetchSettings]);

  return (
    <header className="navbar navbar-expand-md d-print-none">
      <div className="container-xl">
        <Link
          className="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3"
          to="/"
        >
          <img src={CodeboxLogo} alt="logo" width={120} />
        </Link>
        <div className="nav-item dropdown">
          <span
            className="nav-link d-flex lh-1 p-0 px-2"
            data-bs-toggle="dropdown"
            aria-label="Open user menu"
          >
            <img
              className="avatar avatar-sm"
              src={
                settings?.use_gravatar && user
                  ? `https://www.gravatar.com/avatar/${sha256(user?.email)}`
                  : DefaultAvatar
              }
              alt="avatar"
            />
            <div className="d-none d-xl-block ps-2">
              <div>
                <b>
                  {user?.first_name} {user?.last_name}
                </b>
              </div>
              <div className="mt-1 small text-secondary">{user?.email}</div>
            </div>
          </span>
          <div className="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
            {/* <span className="dropdown-header" style={{ fontSize: "11px" }}>
              User
            </span> */}
            <Link to="/profile" className="dropdown-item">
              Profile
            </Link>
            {user?.is_superuser && (
              <Link to="/admin" className="dropdown-item">
                Admin Area
              </Link>
            )}
            {/* <div className="dropdown-divider"></div> */}
            <Link to="/" className="dropdown-item" onClick={HandleLogout}>
              Logout
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}
