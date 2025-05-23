import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, Container, Input } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";
import { useNavigate, useSearchParams } from "react-router-dom";
import { LoginStatus, RequestStatus } from "../../api/types";
import { Http } from "../../api/http";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [searchParams] = useSearchParams();

  const navigate = useNavigate();

  const CheckIfInitialUserExists = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/initial-user-exists`,
      "GET",
      null,
      "application/json",
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      if (!(responseBody as any).exists) {
        navigate("/signup");
        return;
      }
    }
  }, [navigate]);

  const IsAuthenticated = useCallback(async () => {
    // redirect to home if user is already authenticated
    let [status, statusCode] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/user-details`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      navigate("/");
      return;
    }
  }, [navigate]);

  useEffect(() => {
    IsAuthenticated();
    CheckIfInitialUserExists();
  }, [IsAuthenticated, CheckIfInitialUserExists]);

  const SubmitLoginForm = async (event: any) => {
    event.preventDefault();

    // validate fields
    if (email === "" || password === "") {
      setError("Missing email or password");
      return;
    }

    // process login
    let [status] = await Http.Login(email, password);
    if (status === LoginStatus.OK) {
      setError("");
      navigate(searchParams.get("next") || "/");
    } else {
      if (status === LoginStatus.INVALID_CREDENTIALS) {
        setError("Invalid credentials");
      } else {
        setError("Unknown error, check that server is reachable");
      }
    }
  };

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
                <div className="mb-4">
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
                version: {process.env.REACT_APP_VERSION}
              </small>
            </p>
          </div>
        </Container>
      </div>
    </React.Fragment>
  );
}
