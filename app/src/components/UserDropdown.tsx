import { faBorderTopLeft, faGears, faRightFromBracket, faUser, faUserSecret } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button, Col, Row } from "reactstrap";
import DefaultAvatar from "../assets/images/default-avatar.png";
import sha256 from "crypto-js/sha256";
import { Logout, RetrieveInstanceSettings } from "../api/common";
import { InstanceSettings } from "../types/settings";
import { useSelector } from "react-redux";
import { RootState } from "../redux/store";
import { StopImpersonation } from "../api/users";
import { toast } from "react-toastify";

export function UserDropdown() {
    const navigate = useNavigate();
    const [settings, setSettings] = useState<InstanceSettings | null>(null);
    const user = useSelector((state: RootState) => state.user);

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

    const HandleStopImpersonation = useCallback(async() => {
        if(user.impersonated) {
            if(await StopImpersonation()) {
                // trigger a complete reload of the page
                window.location.href = `/admin/users/${user.email}`
            } else {
                toast.error(`Failed to stop to impersonate ${user.email}`);
            }
        } 
    }, [user.email, user.impersonated]);

    useEffect(() => {
        FetchSettings();
    }, [FetchSettings]);

    return (
        <React.Fragment>
            <div className="d-flex">
                {user.impersonated && (
                    <Button 
                        className="mx-1 px-2 text-warning btn btn-outline-warning"
                        onClick={HandleStopImpersonation}
                        title="Stop impersonating"
                    >
                        <FontAwesomeIcon icon={faUserSecret} />
                    </Button>
                )}
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
        </React.Fragment>
    )
}