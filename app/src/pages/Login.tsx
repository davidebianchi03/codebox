import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, Container, Input } from "reactstrap";
import LogoSquare from "../assets/images/logo-square.png";
import { useNavigate, useSearchParams } from "react-router-dom";
import { LoginStatus, RequestStatus } from "../api/types";
import { Http } from "../api/http";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [searchParams] = useSearchParams();

  const navigate = useNavigate();

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
  }, [IsAuthenticated]);

  const SubmitLoginForm = async (event: any) => {
    event.preventDefault();

    // validate fields
    if (email === "" || password === "") {
      setError("Missing email or password");
      return;
    }

    // process login
    let [status, jwtToken, expirationDate] = await Http.Login(email, password);
    if (status === LoginStatus.OK) {
      setError("");
      document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()};domain=${
        window.location.hostname
      }`;
      document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()};domain=.${
        window.location.hostname
      }`;
      if (process.env.NODE_ENV === "development") {
        document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()};domain=${
          new URL(Http.GetServerURL()).hostname
        }`;
      }

      navigate(searchParams.get("next") || "/");
    } else {
      document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=${window.location.hostname}`;
      document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=.${window.location.hostname}`;
      if (process.env.NODE_ENV === "development") {
        document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;domain=${
          new URL(Http.GetServerURL()).hostname
        }`;
      }
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
              <img src={LogoSquare} alt="logo" width={50} />
              <h2 style={{ fontSize: "20pt", marginTop: "10px" }}>Codebox</h2>
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
                <p className="text-danger text-center">{error}</p>
                <div className="d-flex justify-content-between">
                  <Button color="primary w-75 mx-auto" type="submit">
                    Login
                  </Button>
                </div>
              </form>
            </CardBody>
          </Card>
          <div className="text-center text-secondary mt-3">
            Don't have account yet?{" "}
            <a href="./sign-up.html" tabIndex={-1}>
              Sign up
            </a>
          </div>
          <div className="d-flex flex-column justify-content-between mt-2">
            <p className="w-100 text-center mb-0">
              <small className="text-muted">
                &copy; codebox {new Date().getFullYear()}
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
