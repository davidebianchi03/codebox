import { withRouter } from "../common/router";
import { Http } from "../api/http";
import React, { useCallback, useEffect, useState } from "react";
import { RequestStatus } from "../api/types";
import { useLocation, useNavigate } from "react-router-dom";
import { User } from "../types/user";
import { Navbar } from "./Navbar";

type Props = {
  children: string | JSX.Element | JSX.Element[] | (() => JSX.Element);
  showNavbar?: boolean;
};

function AuthRequired({ children, showNavbar = true }: Props) {
  const location = useLocation();
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  const WhoAmI = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/user-details`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setUser(responseBody as User);
    } else {
      navigate(`/login?next=${encodeURIComponent(location.pathname)}`);
    }
  }, [navigate, location]);

  useEffect(() => {
    WhoAmI();
  }, [WhoAmI]);
  
  return (
    <React.Fragment>
      {user && (
        <>
          {showNavbar && <Navbar user={user} />}
          {children}
        </>
      )}
    </React.Fragment>
  );
}

export default withRouter(AuthRequired);
