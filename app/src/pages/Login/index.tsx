import React, { useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Container,
  FormGroup,
  Input,
  Label,
} from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { APIInitialUserExists, APISignUpOpen, RetrieveCurrentUserDetails } from "../../api/common";
import { APILogin, APILoginCode } from "../../api/auth";

export default function LoginPage() {
  const [form, setForm] = useState({
    email: "",
    password: "",
    rememberMe: false,
  });
  const [error, setError] = useState("");
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
      setError("Missing email or password");
      return;
    }

    const { code, token } = await APILogin(email, password, rememberMe);
    if (code === APILoginCode.SUCCESS) {
      if (token) {
        setError("");
        navigate(searchParams.get("next") || "/");
      } else {
        setError("Invalid credentials");
      }
    } else if (code === APILoginCode.EMAIL_NOT_VERIFIED) {
      navigate("/email-not-verified");
    } else {
      setError("Invalid credentials");
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
          <CardBody>
            <h2 className="h2 text-center mb-4">Login to your account</h2>

            <form onSubmit={handleSubmit}>
              <div className="mb-3">
                <label className="form-label">Email</label>
                <Input
                  autoFocus
                  type="text"
                  placeholder="email@example.com"
                  value={form.email}
                  autoComplete="email"
                  onChange={(e) => updateField("email", e.target.value)}
                />
              </div>

              <div className="mb-2">
                <label className="form-label">Password</label>
                <Input
                  type="password"
                  placeholder="password"
                  value={form.password}
                  autoComplete="current-password"
                  onChange={(e) => updateField("password", e.target.value)}
                />
              </div>

              {error && (
                <p
                  className="alert border-0 d-flex justify-content-center"
                  style={{ background: "rgba(var(--tblr-danger-rgb), 0.8)" }}
                >
                  {error}
                </p>
              )}

              <FormGroup className="d-flex align-items-center">
                <Input
                  type="checkbox"
                  id="remember_me"
                  checked={form.rememberMe}
                  onChange={(e) => updateField("rememberMe", e.target.checked)}
                />
                <Label for="remember_me" className="mt-2 ms-2">
                  Remember me
                </Label>
              </FormGroup>

              <div className="d-flex justify-content-between">
                <Button color="primary w-75 mx-auto" type="submit">
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
  );
}
