import React, { useCallback, useEffect, useState } from "react";
import { PageWithCardLayout } from "../../layouts/PageWithCardLayout";
import { Button, Form } from "react-bootstrap";
import { APICanResetPassword, APIInitialUserExists, RetrieveCurrentUserDetails } from "../../api/common";
import { useNavigate } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { APIRequestPasswordReset } from "../../api/auth";
import { NonFieldError } from "../../components/NonFieldError";

export function PasswordResetPage() {
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

    const checkCanResetPassword = useCallback(async () => {
        const canReset = await APICanResetPassword();
        if (!canReset) {
            navigate("/");
            return;
        }
    }, [navigate]);

    useEffect(() => {
        checkUserState();
        checkCanResetPassword();
    }, [checkUserState, checkCanResetPassword]);

    const validation = useFormik({
        initialValues: {
            email: "",
        },
        validateOnBlur: false,
        validateOnChange: false,
        validationSchema: Yup.object({
            email: Yup.string()
                .required("A valid email address is required")
                .email("A valid email address is required"),
        }),
        onSubmit: async (values) => {
            if(await APIRequestPasswordReset(values.email)) {
                navigate("/password-reset/sent");
                setNonFieldError("");
            } else {
                setNonFieldError("An error occurred while sending the reset instructions. Please try again later.");
            }
        }
    })

    return (
        <React.Fragment>
            <PageWithCardLayout
                title="Reset your password"
            >
                <form onSubmit={validation.handleSubmit}>
                    <Form.Group className="mb-3" controlId="formBasicEmail">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            name="email"
                            value={validation.values.email}
                            onChange={validation.handleChange}
                            onBlur={validation.handleBlur}
                            isInvalid={!!validation.errors.email}
                        />
                        <Form.Control.Feedback type="invalid">
                            {validation.errors.email}
                        </Form.Control.Feedback>
                        <Form.Text className="text-muted">
                            We'll send you an email with instructions to reset your password.
                        </Form.Text>
                    </Form.Group>
                    {nonFieldError && <NonFieldError error={nonFieldError} />}
                    <Button variant="light" type="submit" className="w-100">
                        Send reset instructions
                    </Button>
                </form>
            </PageWithCardLayout>
        </React.Fragment>
    );
}
