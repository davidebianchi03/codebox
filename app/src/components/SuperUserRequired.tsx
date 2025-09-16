import { withRouter } from "../common/router";
import { Navbar } from "./Navbar";
import { Container } from "reactstrap";
import React from "react";
import { SuperUserSidebar } from "./SuperUserSidebar";

type Props = {
  children: any;
  showNavbar?: boolean;
};

function AuthRequired({ children, showNavbar = true }: Props) {
  return (
    <React.Fragment>
      <SuperUserSidebar />
      <div className="page-wrapper">
        <div className="superuser-navbar">
          <Navbar showLogo={false} />
        </div>
        <Container className="mt-4 mb-4">
          {children}
        </Container>
      </div>
    </React.Fragment>
  );
}

export default withRouter(AuthRequired);
