import React, { useCallback, useEffect } from "react";
import {
  Button,
  Card,
  CardBody,
  Container,
  FormFeedback,
  Input,
} from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { useNavigate } from "react-router-dom";
import { LoginStatus, RequestStatus } from "../../api/types";
import { Http } from "../../api/http";
import { useFormik } from "formik";
import * as Yup from "yup";
import { toast, ToastContainer } from "react-toastify";
import { RetrieveCurrentUserDetails } from "../../api/common";

export default function SignUpPage() {

  const navigate = useNavigate();

  const CheckIfInitialUserExists = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/initial-user-exists`,
      "GET",
      null,
      "application/json",
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      if ((responseBody as any).exists) {
        navigate("/");
        return;
      }
    }
  }, [navigate]);

  const IsAuthenticated = useCallback(async () => {
    // redirect to home if user is already authenticated
    const user = await RetrieveCurrentUserDetails();
    if (user !== undefined) {
      navigate("/");
    }
  }, [navigate]);

  useEffect(() => {
    IsAuthenticated();
    CheckIfInitialUserExists();
  }, [IsAuthenticated, CheckIfInitialUserExists]);

  var validation = useFormik({
    initialValues: {
      email: "",
      firstName: "",
      lastName: "",
      password: "",
      confirmPassword: "",
    },
    validationSchema: Yup.object({
      email: Yup.string()
        .required("A valid email address is required")
        .email("A valid email address is required"),
      firstName: Yup.string().required("First name is required"),
      lastName: Yup.string().required("First name is required"),
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
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      let [status, code] = await Http.SignUp(values.email, values.password, values.firstName, values.lastName)
      if (status !== LoginStatus.OK || code < 200 || code > 299) {
        toast.error(`Error, recived status ${code}`);
      } else {
        navigate("/login")
      }
    },
  });

  return (
    <React.Fragment>
      <div className="page page-center">
        <Container className="container-tight py-4">
          <div className="text-center mb-4">
            <div className="navbar-brand navbar-brand-autodark">
              <img src={CodeboxLogo} alt="logo" width={185} />
            </div>
          </div>
          <Card className="card-md">
            <CardBody>
              <h2 className="h2 text-center mb-4">Sign-up</h2>
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  validation.handleSubmit();
                  return false;
                }}
              >
                <div className="mb-3">
                  <label className="form-label">First name</label>
                  <Input
                    autoFocus
                    type="text"
                    placeholder="John"
                    name="firstName"
                    value={validation.values.firstName}
                    onChange={validation.handleChange}
                    invalid={validation.errors.firstName !== undefined}
                  />
                  <FormFeedback>{validation.errors.firstName}</FormFeedback>
                </div>
                <div className="mb-3">
                  <label className="form-label">Last name</label>
                  <Input
                    autoFocus
                    type="text"
                    placeholder="Doe"
                    name="lastName"
                    value={validation.values.lastName}
                    onChange={validation.handleChange}
                    invalid={validation.errors.lastName !== undefined}
                  />
                  <FormFeedback>{validation.errors.lastName}</FormFeedback>
                </div>
                <div className="mb-3">
                  <label className="form-label">Email</label>
                  <Input
                    autoFocus
                    type="text"
                    placeholder="email@example.com"
                    name="email"
                    value={validation.values.email}
                    onChange={validation.handleChange}
                    invalid={validation.errors.email !== undefined}
                  />
                  <FormFeedback>{validation.errors.email}</FormFeedback>
                </div>
                <div className="mb-4">
                  <label className="form-label">Password</label>
                  <Input
                    type="password"
                    placeholder="password"
                    name="password"
                    value={validation.values.password}
                    onChange={validation.handleChange}
                    invalid={validation.errors.password !== undefined}
                  />
                  <FormFeedback>{validation.errors.password}</FormFeedback>
                </div>
                <div className="mb-4">
                  <label className="form-label">Confirm password</label>
                  <Input
                    type="password"
                    placeholder="confirm password"
                    name="confirmPassword"
                    value={validation.values.confirmPassword}
                    onChange={validation.handleChange}
                    invalid={validation.errors.confirmPassword !== undefined}
                  />
                  <FormFeedback>
                    {validation.errors.confirmPassword}
                  </FormFeedback>
                </div>
                <div className="d-flex justify-content-between">
                  <Button color="primary w-75 mx-auto" type="submit">
                    Sign up
                  </Button>
                </div>
              </form>
            </CardBody>
          </Card>
          <div className="d-flex flex-column justify-content-between mt-2">
            <p className="w-100 text-center mb-0">
              <small className="text-muted">
                &copy;codebox {new Date().getFullYear()}
              </small>
            </p>
            <p className="w-100 text-center">
              <small className="text-muted">
                version: {import.meta.env.VITE_APP_VERSION}
              </small>
            </p>
          </div>
        </Container>
      </div>
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </React.Fragment>
  );
}
