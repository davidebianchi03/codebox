import React from "react";
import { Link } from "react-router-dom";
import { Container } from "reactstrap";
import CodeboxLogo from "../../assets/images/codebox-logo-white.png";

export default function NotFound() {
  return (
    <React.Fragment>
      <div className="page page-center">
        <Container className="w-100 d-flex justify-content-center">
          <div className="mb-5">
            <div className="section-header">
              <div className="d-flex align-items-center justify-content-center">
                <img src={CodeboxLogo} alt="logo" width={185} />
              </div>
              <h1 className="section-title section-title-lg mt-5">
                Oooops! Page Not Found
              </h1>
              <p className="section-description text-secondary text-center">
                This page doesn't exist or was removed!
                <br />
                We suggest you back home
              </p>
              <div className="mt-5 btn-list d-flex justify-content-center">
                <Link to="/" className="btn btn-accent">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="icon icon-tabler icon-tabler-chevron-left icon"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <polyline points="15 6 9 12 15 18"></polyline>
                  </svg>
                  <span className="btn-text">Back to home</span>
                </Link>
              </div>
            </div>
          </div>
        </Container>
      </div>
    </React.Fragment>
  );
}
