import React, { useCallback, useEffect, useState } from "react";
import { Badge, Button, Card, Col, Container, Row, Spinner } from "react-bootstrap";
import { APIAdminEmailServiceConfigured } from "../../api/common";
import { ToastContainer } from "react-toastify";
import { AdminSendTestEmailResponse, APIAdminSendTestEmail } from "../../api/admin";

export function AdminEmailSenderPage() {

    const [loading, setLoading] = useState<boolean>(true);
    const [sending, setSending] = useState<boolean>(false);
    const [status, setStatus] = useState<"error" | "success">();
    const [description, setDescription] = useState<string>("");
    const [emailServiceConfigured, setEmailServiceConfigured] = useState<boolean>(true);

    const fetchConfig = useCallback(async () => {
        setLoading(true);
        const configured = await APIAdminEmailServiceConfigured();
        setEmailServiceConfigured(configured);
        setLoading(false);
    }, []);

    const handleSendTestEmail = useCallback(async () => {
        setSending(true);
        const r = await APIAdminSendTestEmail();
        setDescription(r.description);
        if (r.response === AdminSendTestEmailResponse.SUCCESS) {
            setStatus("success");
        } else {
            setStatus("error");
        }

        setSending(false);
    }, []);

    useEffect(() => {
        fetchConfig();
    }, [fetchConfig]);

    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">Email Sender</h2>
                    <p className="text-muted">
                        Email Sender is the SMTP service used to send emails to users.
                        You may choose not to configure it, but it is recommended.
                        Some features (such as user sign-up) are disabled if the
                        email service is not configured. You must configure
                        it using stack environment variables. Refer to the &nbsp;
                        <a href="https://codebox4073715.gitlab.io/codebox/" target="_blank">documentation</a>
                        &nbsp; for configuration details.
                    </p>
                </div>

                {loading ? (
                    <React.Fragment>

                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <Row className="mt-4">
                            <Col>
                                {emailServiceConfigured ? (
                                    <React.Fragment>
                                        <Card>
                                            <Card.Header>
                                                <div className="d-flex flex-column">
                                                    <h3 className="mb-0">Test Email</h3>
                                                    <small className="mb-0 text-muted">
                                                        Send a test email to ensure the email service is properly
                                                        configured and functioning.
                                                        The test email will be sent to your email address.
                                                    </small>
                                                </div>
                                            </Card.Header>
                                            <Card.Body>
                                                {status === "success" && (
                                                    <React.Fragment>
                                                        <div className="mb-2">
                                                            <Badge
                                                                bg="success"
                                                                className="text-white"
                                                            >
                                                                Success
                                                            </Badge>
                                                        </div>
                                                    </React.Fragment>
                                                )}
                                                {status === "error" && (
                                                    <React.Fragment>
                                                        <div className="mb-2">
                                                            <Badge
                                                                bg="danger"
                                                                className="text-white"
                                                            >
                                                                Error
                                                            </Badge>
                                                        </div>
                                                    </React.Fragment>
                                                )}
                                                {description && (
                                                    <React.Fragment>
                                                        <div className="bg-dark px-3 py-3 my-3 rounded" style={{ fontFamily: "Consolas" }}>
                                                            {description}
                                                        </div>
                                                    </React.Fragment>
                                                )}
                                                <Button
                                                    variant="light"
                                                    onClick={handleSendTestEmail}
                                                    disabled={sending}
                                                >
                                                    Send Test Email
                                                    {sending && (
                                                        <React.Fragment>
                                                            <Spinner
                                                                size="sm"
                                                                className="ms-2"
                                                            />
                                                        </React.Fragment>
                                                    )}
                                                </Button>
                                                {sending && (
                                                    <React.Fragment>
                                                        <p className="mt-2 text-yellow">
                                                            Sending email, please do not refresh the page...
                                                        </p>
                                                    </React.Fragment>
                                                )}
                                            </Card.Body>
                                        </Card>
                                    </React.Fragment>
                                ) : (
                                    <React.Fragment>
                                        <div className="alert alert-warning">
                                            Email server is not configured. Configure email server to show this section.
                                        </div>
                                    </React.Fragment>
                                )}
                            </Col>
                        </Row>
                    </React.Fragment>
                )}
                <ToastContainer
                    toastClassName={"bg-dark"}
                />
            </Container>
        </React.Fragment>
    )
}
