import React from "react";
import { Container } from "react-bootstrap";

export default function LoadingFallback() {
    return (
        <React.Fragment>
            <Container
                style={{ minHeight: "500px", paddingTop: "2rem", paddingBottom: "2rem" }}
            >
            </Container>
        </React.Fragment>
    );
};
