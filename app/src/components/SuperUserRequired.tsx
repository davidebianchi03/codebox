import { withRouter } from "../common/router";
import { useCallback, useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { User } from "../types/user";
import { Navbar } from "./Navbar";
import { Badge, Card, CardBody, Col, Container, Row } from "reactstrap";
import React from "react";
import { RetrieveCurrentUserDetails } from "../api/common";

type Props = {
  children: any;
  showNavbar?: boolean;
};

function AuthRequired({ children, showNavbar = true }: Props) {
  const navigate = useNavigate();
  const location = useLocation();
  const [user, setUser] = useState<User | null>(null);

  const WhoAmI = useCallback(async () => {
    const user = await RetrieveCurrentUserDetails();
    if (user) {
      setUser(user);
    } else {
      navigate(`/login`);
    }
  }, [navigate]);

  useEffect(() => {
    WhoAmI();
  }, [WhoAmI]);

  return (
    <React.Fragment>
      {user && showNavbar && (
        <>
          <Navbar user={user} />
          <Container className="mt-4 mb-4">
            <div className="row g-2 align-items-center mb-4">
              <div className="col">
                <div className="page-pretitle">Codebox</div>
                <h2 className="page-title">Admin</h2>
              </div>
            </div>
            <Row>
              <Col md={12}>
                <Card>
                  <Row>
                    <Col md={3}>
                      <CardBody>
                        <h4 className="subheader">Settings</h4>
                        <div className="list-group list-group-transparent">
                          <Link
                            to="/admin"
                            className={`list-group-item list-group-item-action d-flex align-items-center ${location.pathname === "/admin" ? "active" : ""}`}
                          >
                            Dashboard
                          </Link>
                          <Link
                            to="/admin/users"
                            className={`list-group-item list-group-item-action d-flex align-items-center ${location.pathname.startsWith("/admin/users") ? "active" : ""}`}
                          >
                            Users
                          </Link>
                          <span
                            className="list-group-item list-group-item-action d-flex align-items-center justify-content-between"
                          >
                            Groups
                            <Badge color="orange" className="text-white ">
                              Coming soon
                            </Badge>
                          </span>
                          <Link
                            to="/admin/runners"
                            className={`list-group-item list-group-item-action d-flex align-items-center ${location.pathname.startsWith("/admin/runners") ? "active" : ""}`}
                          >
                            Runners
                          </Link>
                        </div>
                      </CardBody>
                    </Col>
                    <Col md={9}>
                      <CardBody>
                        {children}
                      </CardBody>
                    </Col>
                  </Row>
                </Card>
              </Col>
            </Row>
          </Container>
        </>)}
    </React.Fragment>
  );
}

export default withRouter(AuthRequired);
