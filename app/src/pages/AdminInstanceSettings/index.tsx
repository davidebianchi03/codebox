import { useFormik } from "formik";
import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, Col, Container, Row } from "reactstrap";
import { APIAdminEmailServiceConfigured, APIAdminRetrieveInstanceSettings, APIAdminUpdateInstanceSettings } from "../../api/common";
import { toast, ToastContainer } from "react-toastify";
import { SignSignUpCard } from "./SignSignUpCard";
import { SettingsFormPlaceholder } from "./SettingsFormPlaceholder";

export function AdminInstanceSettingsPage() {
    const [loading, setLoading] = useState<boolean>(true);
    const [emailServiceConfigured, setEmailServiceConfigured] = useState<boolean>(true);

    const validation = useFormik({
        initialValues: {
            signUpOpen: false,
            signUpRestricted: false,
            allowedEmailRegex: "",
            blacklistedEmailRegex: "",
        },
        validateOnBlur: false,
        validateOnChange: false,
        onSubmit: async (values) => {
            const r = await APIAdminUpdateInstanceSettings(
                values.signUpOpen,
                values.signUpRestricted,
                values.allowedEmailRegex,
                values.blacklistedEmailRegex,
            )

            if (r) {
                toast.success("Settings have been updated");
                fetchConfig();
            } else {
                toast.error("An error occured, please try again later");
            }
        }
    })

    const fetchConfig = useCallback(async () => {
        setLoading(true);
        const s = await APIAdminRetrieveInstanceSettings();
        if (s) {
            validation.setValues({
                signUpOpen: s.is_signup_open,
                signUpRestricted: s.is_signup_restricted,
                allowedEmailRegex: s.allowed_emails_regex,
                blacklistedEmailRegex: s.blocked_emails_regex,
            })
        } else {
            toast.error("Failed to fetch settings, try again later");
        }
        
        const configured = await APIAdminEmailServiceConfigured();
        setEmailServiceConfigured(configured);
        
        setLoading(false);
    }, []);

    useEffect(() => {
        fetchConfig();
    }, [fetchConfig]);

    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">Codebox Settings</h2>
                    <span className="text-muted">Configure codebox</span>
                </div>
                {loading ? (
                    <React.Fragment>
                        <SettingsFormPlaceholder />
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <form onSubmit={validation.handleSubmit}>
                            <Row className="mt-4">
                                <Col>
                                    {emailServiceConfigured ? (
                                        <React.Fragment>
                                            <SignSignUpCard validation={validation} />
                                            <Card className="mt-3">
                                                <CardBody className="d-flex justify-content-end">
                                                    <Button color="accent" className="me-2" onClick={(e) => {
                                                        e.preventDefault();
                                                        validation.resetForm();
                                                        fetchConfig();
                                                    }}>
                                                        Discard
                                                    </Button>
                                                    <Button color="light" type="submit">
                                                        Save
                                                    </Button>
                                                </CardBody>
                                            </Card>
                                        </React.Fragment>
                                    ) : (
                                        <React.Fragment>
                                            <div className="alert alert-warning">
                                                Email server is not configured. Configure email server to enable instance settings.
                                            </div>
                                        </React.Fragment>
                                    )}
                                </Col>
                            </Row>
                        </form>
                    </React.Fragment>
                )}
                <ToastContainer
                    toastClassName={"bg-dark"}
                />
            </Container>
        </React.Fragment >
    )
}
