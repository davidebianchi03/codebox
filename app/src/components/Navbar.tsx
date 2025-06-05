import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import DefaultAvatar from "../assets/images/default-avatar.png";
import sha256 from "crypto-js/sha256";
import { Link, useNavigate } from "react-router-dom";
import { InstanceSettings } from "../types/settings";
import { useCallback, useEffect, useState } from "react";
import { User } from "../types/user";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBorderTopLeft, faGears, faRightFromBracket, faUser } from "@fortawesome/free-solid-svg-icons";
import { Col, Row } from "reactstrap";
import { Logout, RetrieveInstanceSettings } from "../api/common";

interface Props {
  user: User;
}

export function Navbar({ user }: Props) {
  const navigate = useNavigate();
  const [settings, setSettings] = useState<InstanceSettings | null>(null);

  const HandleLogout = (e: any) => {
    e.preventDefault();
    Logout();
    navigate("/login");
  };

  const FetchSettings = useCallback(async () => {
    const s = await RetrieveInstanceSettings();
    if (s) {
      setSettings(s);
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
            <Link to="/profile" className="dropdown-item">
              <Row>
                <Col md={4}>
                  <FontAwesomeIcon icon={faUser} />
                </Col>
                <Col md={8}>
                  Profile
                </Col>
              </Row>
            </Link>
            {user?.is_superuser && (
              <Link to="/admin" className="dropdown-item">
                <Row>
                  <Col md={4} className="pe-0">
                    <FontAwesomeIcon icon={faGears} />
                  </Col>
                  <Col md={8} className="ps-0">
                    Admin Area
                  </Col>
                </Row>
              </Link>
            )}
            <Link to="/templates" className="dropdown-item">
              <Row>
                <Col md={4} className="pe-0">
                  <FontAwesomeIcon icon={faBorderTopLeft} />
                </Col>
                <Col md={8} className="ps-0">
                  Templates
                </Col>
              </Row>
            </Link>
            <Link to="/" className="dropdown-item" onClick={HandleLogout}>
              <Row>
                <Col md={4}>
                  <FontAwesomeIcon icon={faRightFromBracket} />
                </Col>
                <Col md={8}>
                  Logout
                </Col>
              </Row>
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
}
