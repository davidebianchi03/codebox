import React, { useCallback } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button, Col, Row } from "reactstrap";
import { Logout } from "../api/common";
import { useSelector } from "react-redux";
import { RootState } from "../redux/store";
import { StopImpersonation } from "../api/users";
import { toast } from "react-toastify";
import { Dropdown } from "react-bootstrap";
import { ShieldCogIcon, TemplateIcon, TerminalIcon, UserIcon, LogoutIcon, GhostOffIcon } from "../icons/Tabler";

export function UserDropdown() {
    const navigate = useNavigate();
    const user = useSelector((state: RootState) => state.user);

    const HandleLogout = (e: any) => {
        e.preventDefault();
        Logout();
        navigate("/login");
    };

    const HandleStopImpersonation = useCallback(async () => {
        if (user.impersonated) {
            if (await StopImpersonation()) {
                // trigger a complete reload of the page
                window.location.href = `/admin/users/${user.email}`
            } else {
                toast.error(`Failed to stop to impersonate ${user.email}`);
            }
        }
    }, [user.email, user.impersonated]);

    return (
        <React.Fragment>
            <div className="d-flex user-dropdown">
                {user.impersonated && (
                    <Button
                        className="mx-1 px-2 text-warning btn btn-outline-warning"
                        onClick={HandleStopImpersonation}
                        title="Stop impersonating"
                    >
                        <GhostOffIcon />
                    </Button>
                )}
                <Dropdown>
                    <Dropdown.Toggle
                        variant="link"
                        className="d-flex lh-1 p-0 px-2 text-white"
                    >
                        <div
                            className="avatar avatar-sm bg-azure-lt text-white d-flex align-items-center justify-content-center"
                            style={{ width: '35px', height: '35px' }}
                        >
                            {user?.first_name.length > 0 && user.first_name[0]}
                            &nbsp;
                            {user?.last_name.length > 0 && user.last_name[0]}
                        </div>
                        <div className="d-none d-xl-block ps-2 text-start">
                            <div>
                                <b>
                                    {user?.first_name} {user?.last_name}
                                </b>
                            </div>
                            <div className="mt-1 small text-secondary">{user?.email}</div>
                        </div>
                    </Dropdown.Toggle>
                    <Dropdown.Menu>
                        <Dropdown.Item as={Link} to="/profile">
                            <Row>
                                <Col md={4}>
                                    <UserIcon />
                                </Col>
                                <Col md={8}>
                                    Profile
                                </Col>
                            </Row>
                        </Dropdown.Item>
                        {user?.is_superuser && (
                            <Dropdown.Item as={Link} to="/admin">
                                <Row>
                                    <Col md={4} className="pe-0">
                                        <ShieldCogIcon />
                                    </Col>
                                    <Col md={8} className="ps-0">
                                        Admin Area
                                    </Col>
                                </Row>
                            </Dropdown.Item>
                        )}
                        <Dropdown.Item as={Link} to="/templates">
                            <Row>
                                <Col md={4} className="pe-0">
                                    <TemplateIcon />
                                </Col>
                                <Col md={8} className="ps-0">
                                    Templates
                                </Col>
                            </Row>
                        </Dropdown.Item>
                        <Dropdown.Item as={Link} to="/cli">
                            <Row>
                                <Col md={6}>
                                    <TerminalIcon />
                                </Col>
                                <Col md={6}>
                                    CLI
                                </Col>
                            </Row>
                        </Dropdown.Item>
                        <Dropdown.Item as={Link} to="/" onClick={HandleLogout}>
                            <Row>
                                <Col md={4}>
                                    <LogoutIcon />
                                </Col>
                                <Col md={8}>
                                    Logout
                                </Col>
                            </Row>
                        </Dropdown.Item>
                    </Dropdown.Menu>
                </Dropdown>
            </div>
        </React.Fragment>
    )
}