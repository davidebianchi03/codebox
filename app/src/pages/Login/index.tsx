import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, Container, FormGroup, Input, Label } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { useNavigate, useSearchParams } from "react-router-dom";
import { APIInitialUserExists, RetrieveCurrentUserDetails } from "../../api/common";
import { APILogin } from "../../api/auth";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(false);
  const [error, setError] = useState("");
  const [searchParams] = useSearchParams();

  const navigate = useNavigate();

  const CheckIfInitialUserExists = useCallback(async () => {
    if (!(await APIInitialUserExists())) {
      navigate("/signup");
    }
  }, [navigate]);

  const IsAuthenticated = useCallback(async () => {
    // redirect to home if user is already authenticated
    const user = await RetrieveCurrentUserDetails();
    if (user !== undefined) {
      navigate("/");
      return;
    }
  }, [navigate]);

  useEffect(() => {
    IsAuthenticated();
    CheckIfInitialUserExists();
  }, [IsAuthenticated, CheckIfInitialUserExists]);

  const SubmitLoginForm = useCallback(async (event: any) => {
    event.preventDefault();

    // validate fields
    if (email === "" || password === "") {
      setError("Missing email or password");
      return;
    }

    // process login
    const token = await APILogin(email, password, rememberMe);
    if (token) {
      setError("");
      navigate(searchParams.get("next") || "/");
    } else {
      setError("Invalid credentials");
    }
  }, [email, navigate, password, rememberMe, searchParams]);

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
              <h2 className="h2 text-center mb-4">Login to your account</h2>
              <form onSubmit={SubmitLoginForm}>
                <div className="mb-3">
                  <label className="form-label">Email</label>
                  <Input
                    autoFocus
                    type="text"
                    placeholder="email@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div className="mb-2">
                  <label className="form-label">Password</label>
                  <Input
                    type="password"
                    placeholder="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
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
                    id="remeber_me"
                    checked={rememberMe}
                    onChange={(e) => setRememberMe(e.target.checked)}
                  />
                  <Label for="remeber_me" className="mt-2 ms-2">
                    Remember me
                  </Label>
                </FormGroup>
                <div className="d-flex justify-content-between">
                  <Button color="primary w-75 mx-auto" type="submit">
                    Login
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
    </React.Fragment>
  );
}
