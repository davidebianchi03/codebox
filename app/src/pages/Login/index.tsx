import React, { useEffect, useState } from "react";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { APIInitialUserExists, APISignUpOpen, RetrieveCurrentUserDetails } from "../../api/common";
import { APILogin, APILoginCode } from "../../api/auth";
import { NonFieldError } from "../../components/NonFieldError";
import { FieldError } from "../../components/FieldError";
import { Button, Card, Container, Form } from "react-bootstrap";

export default function LoginPage() {
  const [form, setForm] = useState({
    email: "",
    password: "",
    rememberMe: false,
  });
  const [errors, setErrors] = useState({
    email: "",
    password: "",
    nonFieldError: "",
  });
  const [signupOpen, setSignupOpen] = useState<boolean>(false);

  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const updateField = (key: string, value: any) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  useEffect(() => {
    const checkUserState = async () => {
      const user = await RetrieveCurrentUserDetails();
      if (user) {
        navigate("/");
        return;
      }

      const exists = await APIInitialUserExists();
      if (!exists) navigate("/signup");
    };

    const checkSignupOpen = async () => {
      setSignupOpen(await APISignUpOpen());
    }

    checkUserState();
    checkSignupOpen();
  }, [navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const { email, password, rememberMe } = form;

    if (!email || !password) {
      const errs = {
        email: "",
        password: "",
        nonFieldError: "",
      };

      if (!email) {
        errs.email = "Email is required";
      }
      if (!password) {
        errs.password = "Password is required";
      }

      setErrors(errs);
      return;
    }

    const { code, token } = await APILogin(email, password, rememberMe);
    if (code === APILoginCode.SUCCESS) {
      if (token) {
        setErrors({
          email: "",
          password: "",
          nonFieldError: "",
        });
        navigate(searchParams.get("next") || "/");
      } else {

        // clear errors
        setErrors({
          email: "",
          password: "",
          nonFieldError: "Wrong email or password",
        });
      }
    } else if (code === APILoginCode.EMAIL_NOT_VERIFIED) {
      navigate("/email-not-verified");
    } else if (code === APILoginCode.RATELIMIT) {
      setErrors({
        email: "",
        password: "",
        nonFieldError: "Too many requests, try again later",
      });
    } else {
      setErrors({
        email: "",
        password: "",
        nonFieldError: "Wrong email or password",
      });
    }
  };

  return (
    <div className="page page-center">
      <Container className="container-tight py-4">
        <div className="text-center mb-4">
          <div className="navbar-brand navbar-brand-autodark">
            <img src={CodeboxLogo} alt="logo" width={185} />
          </div>
        </div>

        <Card className="card-md">
          <Card.Body>
            <h2 className="h2 text-center mb-4">Login to your account</h2>

            <form onSubmit={handleSubmit}>
              <div className="mb-3">
                <label className="form-label">Email</label>
                <Form.Control
                  autoFocus
                  type="text"
                  placeholder="email@example.com"
                  value={form.email}
                  autoComplete="email"
                  onChange={(e) => updateField("email", e.target.value)}
                  isInvalid={!!errors.email}
                />
                <FieldError error={errors.email} />
              </div>

              <div className="mb-2">
                <label className="form-label">Password</label>
                <Form.Control
                  type="password"
                  placeholder="password"
                  value={form.password}
                  autoComplete="current-password"
                  onChange={(e) => updateField("password", e.target.value)}
                  isInvalid={!!errors.password}
                />
                <FieldError error={errors.password} />
              </div>

              {errors.nonFieldError && (
                <NonFieldError error={errors.nonFieldError} />
              )}

              <Form.Group className="d-flex align-items-center">
                <input
                  type="checkbox"
                  className="form-check-input form-check-input-light"
                  id="remember_me"
                  checked={form.rememberMe}
                  onChange={(e) => updateField("rememberMe", e.target.checked)}
                />
                <Form.Label for="remember_me" className="mt-2 ms-2" style={{ userSelect: "none" }}>
                  Remember me
                </Form.Label>
              </Form.Group>

              <div className="d-flex justify-content-between">
                <Button variant="light" className="w-75 mx-auto mt-5" type="submit">
                  Login
                </Button>
              </div>
            </form>
            {signupOpen && (
              <React.Fragment>
                <div className="hr-text">or</div>
                <div className="text-center fs-5">
                  Don't have an account? <Link to="/signup">Sign Up</Link>
                </div>
              </React.Fragment>
            )}
          </Card.Body>
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
  );
}
