import React from "react";
import { Container } from "reactstrap";
import { Navbar } from "./Navbar";

interface NavbarLayoutProps {
    children: React.ReactNode;
}

export function NavbarLayout({ children }: NavbarLayoutProps) {
    return (
        <React.Fragment>
            <Navbar />
            <Container>
                {children}
            </Container>
        </React.Fragment>
    );
}
