import { withRouter } from "../common/router";
import React, { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { User } from "../types/user";
import { Navbar } from "./Navbar";
import { RetrieveCurrentUserDetails } from "../api/common";

type Props = {
  children: string | JSX.Element | JSX.Element[] | (() => JSX.Element);
  showNavbar?: boolean;
};

function TemplateManagerRequired({ children, showNavbar = true }: Props) {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  const WhoAmI = useCallback(async () => {
    const user = await RetrieveCurrentUserDetails();
    if (user) {
      setUser(user);
    } else {
      navigate(`/login`);
    }
  }, [navigate]);

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

export default withRouter(TemplateManagerRequired);
