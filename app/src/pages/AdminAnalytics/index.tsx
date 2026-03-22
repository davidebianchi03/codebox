import React from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import { AdminAnalyticsConfig } from "../../components/AdminAnalyticsConfig";
import { ToastContainer } from "react-toastify";
import { AdminAnalyticsContentPreview } from "../../components/AdminAnalyticsContentPreview";

export function AdminAnalyticsPage() {
    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">Analytics</h2>
                    <h4>Help us improve Codebox</h4>
                    <p className="text-muted mb-1">
                        Codebox can send technical usage statistics to help us understand how the product is used and how to improve it.
                    </p>
                    <p className="text-muted mb-1">
                        This includes aggregated information such as:
                    </p>
                    <p className="text-muted mb-1">
                        <ul>
                            <li>number of users</li>
                            <li>number of workspaces and runners</li>
                            <li>server version and license type</li>
                        </ul>
                    </p>
                    <p className="text-muted mb-1">
                        We do not collect personal data or user content.
                    </p>
                    <p className="mb-1">
                        Analytics data is processed using a third-party service (PostHog).
                    </p>
                    <p className="text-muted mb-1">
                        Each installation is identified using a pseudonymous identifier that does not directly identify you.
                    </p>
                    <p className="text-muted mb-1">
                        Analytics are disabled by default and only enabled with your consent. You can change this setting at any time.
                    </p>
                </div>
                <Row className="mt-4">
                    <Col>
                        <Card body>
                            <AdminAnalyticsConfig />
                        </Card>
                    </Col>
                </Row>
                <p className="text-muted my-4">
                    Below is a real example of the data sent to our analytics server:
                </p>
                <Row className="mt-4">
                    <Col>
                        <AdminAnalyticsContentPreview />
                    </Col>
                </Row>
            </Container>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    )
}
