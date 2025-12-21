import React from "react";
import { Link } from "react-router-dom";
import { Card, Col, Container, Row } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";

export function EmailNotVerifiedPage() {
    return (
        <React.Fragment>
            <div className="page page-center">
                <Container className="w-100 d-flex justify-content-center">
                    <div className="d-flex flex-column align-items-center">
                        <div className="d-flex align-items-center justify-content-center">
                            <img src={CodeboxLogo} alt="logo" width={185} />
                        </div>
                        <Row className="d-flex flex-column align-items-center mt-5">
                            <Col md={8}>
                                <Card body className="text-center">
                                    <h2>Email address not verified</h2>
                                    <p>
                                        Your email address has not yet been verified.
                                        Check your email for the verification link.
                                    </p>
                                    <Link to="/login" className="btn btn-light">
                                        Back to login
                                    </Link>
                                </Card>
                            </Col>
                        </Row>
                    </div>
                </Container>
            </div>
        </React.Fragment>
    );
}
