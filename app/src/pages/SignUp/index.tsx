import React, { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Container,
  FormFeedback,
  Input,
} from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { Link, useNavigate } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { ToastContainer } from "react-toastify";
import { APIInitialUserExists, APISignUpOpen, RetrieveCurrentUserDetails } from "../../api/common";
import { APISignUp, APISignUpCode } from "../../api/auth";
import { NonFieldError } from "../../components/NonFieldError";

export default function SignUpPage() {

  const navigate = useNavigate();
  const [firstUserExists, setFirstUserExists] = useState<boolean>(false);
  const [nonFieldError, setNonFieldError] = useState<string>("");

  const CheckCanCreateNewUser = useCallback(async () => {
    const initialUserExists = await APIInitialUserExists();
    setFirstUserExists(initialUserExists);
    const isSignupOpen = await APISignUpOpen();
    if (initialUserExists && !isSignupOpen) {
      navigate("/");
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
    CheckCanCreateNewUser();
  }, [IsAuthenticated, CheckCanCreateNewUser]);

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
      const signUpResult = await APISignUp(values.email, values.password, values.firstName, values.lastName);
      if (signUpResult === APISignUpCode.SUCCESS) {
        setNonFieldError("");
        navigate("/login");
      } else if (signUpResult === APISignUpCode.CANNOT_SIGNUP) {
        setNonFieldError(`
          Sign-up is disabled or your account cannot be created at this time. 
          Please contact the administrator for assistance.
        `);
      } else {
        setNonFieldError(`An unexpected error occured, try again later.`);
      }
    }
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
              {nonFieldError && (
                <NonFieldError error={nonFieldError} />
              )}
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
                    autocomplete="given-name"
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
                    autocomplete="family-name"
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
                    autocomplete="email"
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
                    autoComplete="new-password"
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
                    autoComplete="new-password"
                    value={validation.values.confirmPassword}
                    onChange={validation.handleChange}
                    invalid={validation.errors.confirmPassword !== undefined}
                  />
                  <FormFeedback>
                    {validation.errors.confirmPassword}
                  </FormFeedback>
                </div>
                <div className="d-flex justify-content-between">
                  <Button color="light" className="w-75 mx-auto" type="submit">
                    Sign up
                  </Button>
                </div>
              </form>
              {firstUserExists && (<React.Fragment>
                <div className="hr-text">or</div>
                <div className="text-center fs-5">
                  Already have an account? <Link to="/login">Login</Link>
                </div>
              </React.Fragment>)}
            </CardBody>
          </Card>
          <div className="d-flex flex-column justify-content-between mt-2">
            <p className="w-100 text-center mb-0">
              <small className="text-muted">
                &copy;&nbsp;
                <a href="https://gitlab.com/codebox4073715/codebox" target="_blank">Codebox</a>
                &nbsp;{new Date().getFullYear()}
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
