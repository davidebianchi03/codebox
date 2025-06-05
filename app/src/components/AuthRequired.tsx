import { withRouter } from "../common/router";
import React, { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { User } from "../types/user";
import { Navbar } from "./Navbar";
import { RetrieveCurrentUserDetails } from "../api/common";

export type AuthRequiredProps = {
  children: string | JSX.Element | JSX.Element[] | (() => JSX.Element);
  showNavbar?: boolean;
};

function AuthRequired({ children, showNavbar = true }: AuthRequiredProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  const WhoAmI = useCallback(async () => {
    const user = await RetrieveCurrentUserDetails();
    if (user) {
      setUser(user);
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
