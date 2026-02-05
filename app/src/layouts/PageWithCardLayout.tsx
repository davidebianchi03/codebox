import React from "react";
import { Card, Container } from "react-bootstrap";
import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import { ToastContainer } from "react-toastify";

export interface PageWithCardLayoutProps {
    title: string;
    children?: React.ReactNode;
}

export function PageWithCardLayout({ title, children }: PageWithCardLayoutProps) {
    return (
        <React.Fragment>
            <div className="page page-center">
                <Container className="container-tight py-4">
                    <div className="text-center mb-4">
                        <div className="navbar-brand navbar-brand-autodark">
                            <img src={CodeboxLogo} alt="logo" width={185} />
                        </div>
                    </div>

                    <Card className="card-md">
                        <Card.Body>
                            <h2 className="h2 text-center mb-4">{title}</h2>
                            {children}
                        </Card.Body>
                    </Card>

                    <div className="d-flex flex-column justify-content-between mt-2">
                        <p className="w-100 text-center mb-0">
                            <small className="text-muted">
                                &copy;&nbsp;
                                <a href="https://github.com/davidebianchi03/codebox" target="_blank">Codebox</a>
                                &nbsp;{new Date().getFullYear()}
                            </small>
                        </p>
                        <p className="w-100 text-center">
                            <small className="text-muted">
                                version: {import.meta.env.VITE_APP_VERSION}
                            </small>
                        </p>
                    </div>
                </Container>
            </div>

            <ToastContainer
                toastClassName={"bg-dark"}
            />
        </React.Fragment>
    )
}
