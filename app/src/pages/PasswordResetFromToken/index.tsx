import React, { useCallback, useEffect, useState } from "react";
import { PageWithCardLayout } from "../../layouts/PageWithCardLayout";
import { useFormik } from "formik";
import { Button, Form } from "react-bootstrap";
import { NonFieldError } from "../../components/NonFieldError";
import { APICanResetPasswordCode, APIResetPasswordFromToken } from "../../api/auth";
import { useNavigate } from "react-router-dom";
import * as Yup from "yup";
import { APIInitialUserExists, RetrieveCurrentUserDetails } from "../../api/common";
import { toast } from "react-toastify";

export function PasswordResetFromTokenPage() {
    const navigate = useNavigate();
    const [nonFieldError, setNonFieldError] = useState<string>("");

    const checkUserState = useCallback(async () => {
        const user = await RetrieveCurrentUserDetails();
        if (user) {
            navigate("/");
            return;
        }

        const exists = await APIInitialUserExists();
        if (!exists) navigate("/signup");
    }, [navigate]);

    useEffect(() => {
        checkUserState();
    }, [checkUserState]);

    useEffect(() => {
        // check that token is present in query params
        const token = new URLSearchParams(window.location.search).get("token");
        if (!token) {
            navigate("/login");
        }
    });

    const validation = useFormik({
        initialValues: {
            password: "",
            confirmPassword: "",
        },
        validateOnChange: false,
        validateOnBlur: false,
        validationSchema: Yup.object({
            password: Yup.string()
                .required("A password is required")
                .test({
                    name: "password",
                    exclusive: false,
                    params: {},
                    message: "The password must be at least 10 characters long and include at least one uppercase letter and one special symbol (!_-,.?!).",
                    test: (value, context) => {
                        if (value.length < 10) {
                            return false;
                        }
                        const hasUppercase = /[A-Z]/.test(value);
                        const hasSpecialSymbol = /[!_\-,.?]/.test(value);
                        return hasUppercase && hasSpecialSymbol;
                    },
                }),
            confirmPassword: Yup.string()
                .required("Confirm the password")
                .test({
                    name: "confirmPassword",
                    exclusive: false,
                    params: {},
                    message: "Passwords do not match",
                    test: (value, context) => value === context.parent.password,
                }),
        }),
        onSubmit: async (values) => {
            const code = await APIResetPasswordFromToken(
                new URLSearchParams(window.location.search).get("token") || "",
                values.password,
            );

            if (code === APICanResetPasswordCode.SUCCESS) {
                toast.success("Password reset successfully. Please log in with your new password.");
                setNonFieldError("");
                navigate("/login");
            } else if (code === APICanResetPasswordCode.INVALID_TOKEN) {
                setNonFieldError("The password reset link is invalid or has expired.");
            } else if (code === APICanResetPasswordCode.PASSWORD_RESET_NOT_AVAILABLE) {
                setNonFieldError("Password reset is not available. Please contact support.");
            } else {
                setNonFieldError("An unknown error occurred. Please try again later.");
            }
        }
    });

    return (
        <React.Fragment>
            <PageWithCardLayout
                title="Password Reset"
            >
                <form onSubmit={validation.handleSubmit}>
                    <Form.Group>
                        <Form.Label>New Password</Form.Label>
                        <Form.Control
                            type="password"
                            name="password"
                            value={validation.values.password}
                            onChange={validation.handleChange}
                            isInvalid={!!validation.errors.password}
                            placeholder="Enter the new password"
                        />
                        <Form.Control.Feedback type="invalid">
                            {validation.errors.password}
                        </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group className="mt-3">
                        <Form.Label>Confirm New Password</Form.Label>
                        <Form.Control
                            type="password"
                            name="confirmPassword"
                            value={validation.values.confirmPassword}
                            onChange={validation.handleChange}
                            isInvalid={!!validation.errors.confirmPassword}
                            placeholder="Confirm the new password"
                        />
                        <Form.Control.Feedback type="invalid">
                            {validation.errors.confirmPassword}
                        </Form.Control.Feedback>
                    </Form.Group>
                    {nonFieldError && <NonFieldError error={nonFieldError} />}
                    <Button variant="light" className="w-100 mt-4" type="submit">
                        Reset Password
                    </Button>
                </form>
            </PageWithCardLayout>
        </React.Fragment>
    );
}
