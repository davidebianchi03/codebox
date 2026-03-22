import React from "react";
import { Container } from "reactstrap";
import { Navbar } from "./Navbar";
import { Footer } from "./Footer";
import { AdminAnalyticsModal } from "../components/AdminAnalyticsModal";

interface NavbarLayoutProps {
    children: React.ReactNode;
}

export function NavbarLayout({ children }: NavbarLayoutProps) {
    return (
        <React.Fragment>
            <Navbar />
            <Container style={{ minHeight: "calc(100vh - 170px)" }}>
                {children}
            </Container>
            <Footer />
            <AdminAnalyticsModal />
        </React.Fragment>
    );
}
