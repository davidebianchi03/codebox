import { useFormik } from "formik";
import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, CardHeader, Col, Container, FormGroup, Input, Label, Row } from "reactstrap";
import { APIAdminRetrieveInstanceSettings, APIAdminUpdateInstanceSettings } from "../../api/common";
import { toast, ToastContainer } from "react-toastify";

export function AdminInstanceSettingsPage() {
    const [loading, setLoading] = useState<boolean>(true);

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
                fetchSettings();
            } else {
                toast.error("An error occured, please try again later");
            }
        }
    })

    const fetchSettings = useCallback(async () => {
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
        setLoading(false);
    }, []);

    useEffect(() => {
        fetchSettings();
    }, [fetchSettings]);

    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">Codebox Settings</h2>
                    <span className="text-muted">Configure codebox</span>
                </div>
                <form onSubmit={validation.handleSubmit}>
                    <Row className="mt-4">
                        <Col>
                            <Card>
                                <CardHeader>
                                    <h2 className="mb-0">Sign In/Sign Up</h2>
                                </CardHeader>
                                <CardBody>
                                    <FormGroup>
                                        <div className="d-flex">
                                            <Input
                                                type="checkbox"
                                                id="signUpOpen"
                                                name="signUpOpen"
                                                checked={validation.values.signUpOpen}
                                                onChange={validation.handleChange}
                                            />
                                            <Label for="signUpOpen" className="ms-3">
                                                Sign Up Open
                                            </Label>
                                        </div>
                                        <p className="mb-0">
                                            <small className="text-muted">
                                                Allow users to register. If this feature is enabled, anyone can create an account.
                                                You can restrict sign-ups to a specific group of users by using the whitelist below.
                                                You have to configure email service to enable this setting.
                                            </small>
                                        </p>
                                    </FormGroup>
                                    <FormGroup>
                                        <div className="d-flex">
                                            <Input
                                                type="checkbox"
                                                id="signUpRestricted"
                                                name="signUpRestricted"
                                                checked={validation.values.signUpRestricted}
                                                onChange={validation.handleChange}
                                                disabled={!validation.values.signUpOpen}
                                            />
                                            <Label for="signUpRestricted" className="ms-3">
                                                Sign Up Restricted
                                            </Label>
                                        </div>
                                        <p className="mb-0">
                                            <small className="text-muted">
                                                To restrict sign-ups, only users whose email addresses match the regular expression
                                                (regex) below will be allowed to create an account.
                                            </small>
                                        </p>
                                    </FormGroup>
                                    <FormGroup>
                                        <Label for="allowedEmailRegex">
                                            Allowed Email Addresses Regex
                                        </Label>
                                        <Input
                                            type="textarea"
                                            id="allowedEmailRegex"
                                            name="allowedEmailRegex"
                                            value={validation.values.allowedEmailRegex}
                                            onChange={validation.handleChange}
                                            disabled={validation.values.signUpRestricted || !validation.values.signUpOpen}
                                            placeholder="e.g. ^.*@example\.com$"
                                        />
                                        <p className="mb-0">
                                            <small className="text-muted">
                                                To ensure security, only users with verified email addresses that
                                                successfully match the regular expression (regex) defined below are
                                                allowed to sign-up and access Codebox; verification is required before signing in.
                                                Enter one regex per line.
                                            </small>
                                        </p>
                                    </FormGroup>
                                    <FormGroup>
                                        <Label for="blacklistedEmailRegex">
                                            Blacklisted Email Addresses Regex
                                        </Label>
                                        <Input
                                            type="textarea"
                                            id="blacklistedEmailRegex"
                                            name="blacklistedEmailRegex"
                                            value={validation.values.blacklistedEmailRegex}
                                            onChange={validation.handleChange}
                                            placeholder="e.g. ^.*@example\.com$"
                                        />
                                        <p className="mb-0">
                                            <small className="text-muted">
                                                Users whose emails match this regex cannot sign up. Enter one regex per line.
                                            </small>
                                        </p>
                                    </FormGroup>
                                </CardBody>
                            </Card>
                            <Card className="mt-3">
                                <CardBody className="d-flex justify-content-end">
                                    <Button color="accent" className="me-2" onClick={(e) => {
                                        e.preventDefault();
                                        validation.resetForm();
                                        fetchSettings();
                                    }}>
                                        Discard
                                    </Button>
                                    <Button color="primary" type="submit">
                                        Save
                                    </Button>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </form>
                <ToastContainer
                    toastClassName={"bg-dark"}
                />
            </Container>
        </React.Fragment >
    )
}
