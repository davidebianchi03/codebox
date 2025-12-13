import React, { useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { Button, Card, Col, Container, FormFeedback, FormGroup, Input, Label, Row } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { useFormik } from "formik";
import * as Yup from "yup";
import { APIVerifyEmailAddress, APIVerifyEmailCode } from "../../api/auth";

export function VerifyEmailPage() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    const validation = useFormik({
        initialValues: {
            code: ""
        },
        validationSchema: Yup.object({
            code: Yup.string().required("Verification code is required")
        }),
        validateOnBlur: false,
        validateOnChange: false,
        onSubmit: async (values) => {
            const result = await APIVerifyEmailAddress(values.code);
            if (result === APIVerifyEmailCode.SUCCESS) {
                navigate("/login");
                return;
            } else if (result === APIVerifyEmailCode.EMAIL_ALREADY_VERIFIED) {
                navigate("/login");
                return;
            } else if (result === APIVerifyEmailCode.INVALID_CODE) {
                validation.setFieldError("code", "Invalid verification code");
                return;
            } else {
                validation.setFieldError("code", "An unknown error occurred. Please try again.");
                return;
            }
        }
    })

    useEffect(() => {
        const code = searchParams.get("code");
        if (code) {
            validation.setFieldValue("code", code);
        }
    }, []);


    return (
        <React.Fragment>
            <div className="page page-center">
                <Container className="w-100 d-flex justify-content-center">
                    <div className="d-flex flex-column align-items-center">
                        <div className="d-flex align-items-center justify-content-center">
                            <img src={CodeboxLogo} alt="logo" width={185} />
                        </div>
                        <Row className="d-flex flex-column align-items-center mt-5">
                            <Col md={8} style={{ minWidth: 350 }}>
                                <Card body className="text-center">
                                    <h2>Verify your email address</h2>
                                    <p>
                                        Enter the verification code below
                                    </p>
                                    <form onSubmit={validation.handleSubmit}>
                                        <FormGroup>
                                            <Label>Verification code</Label>
                                            <Input
                                                type="text"
                                                name="code"
                                                value={validation.values.code}
                                                onChange={validation.handleChange}
                                                invalid={!!validation.errors.code}
                                            />
                                            <FormFeedback>
                                                {validation.errors.code}
                                            </FormFeedback>
                                        </FormGroup>
                                        <Button type="submit" color="primary" className="w-100">
                                            Submit
                                        </Button>
                                    </form>
                                </Card>
                            </Col>
                        </Row>
                    </div>
                </Container>
            </div>
        </React.Fragment>
    );
}
