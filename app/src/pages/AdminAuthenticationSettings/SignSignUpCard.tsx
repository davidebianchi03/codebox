import React from "react";
import { Card, CardBody, CardHeader, FormGroup, Input, Label } from "reactstrap";

interface SignSignUpCardProps {
    validation: any
}

export function SignSignUpCard({ validation }: SignSignUpCardProps) {
    return (
        <React.Fragment>
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
                                className="form-check-input-light"
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
                                id="usersMustBeApproved"
                                name="usersMustBeApproved"
                                className="form-check-input-light"
                                checked={validation.values.usersMustBeApproved}
                                onChange={validation.handleChange}
                                disabled={!validation.values.signUpOpen}
                            />
                            <Label for="usersMustBeApproved" className="ms-3">
                                Users must be approved
                            </Label>
                        </div>
                        <p className="mb-0">
                            <small className="text-muted">
                                If this settings is enabled new users cannot login until they are
                                approved by an administrator.
                            </small>
                        </p>
                    </FormGroup>
                    <FormGroup>
                        <Label for="autoApprovedUsersRegex">
                            Auto Approved Email Regex
                        </Label>
                        <Input
                            type="textarea"
                            id="autoApprovedUsersRegex"
                            name="autoApprovedUsersRegex"
                            value={validation.values.autoApprovedUsersRegex}
                            onChange={validation.handleChange}
                            disabled={!validation.values.signUpOpen}
                            placeholder="e.g. ^.*@example\.com$"
                        />
                        <p className="mb-0">
                            <small className="text-muted">
                                If new users require administrator approval,
                                you can define a list of regular expressions 
                                (regex) to automatically approve specific 
                                email domains or addresses. Enter one regex per line.
                            </small>
                        </p>
                    </FormGroup>
                    <FormGroup>
                        <div className="d-flex">
                            <Input
                                type="checkbox"
                                id="signUpRestricted"
                                name="signUpRestricted"
                                className="form-check-input-light"
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
                            disabled={!validation.values.signUpRestricted || !validation.values.signUpOpen}
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
        </React.Fragment>
    )
}
