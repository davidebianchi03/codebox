import { withRouter } from "../common/router";
import { Http } from "../api/http";
import { useCallback, useEffect, useState } from "react";
import { RequestStatus } from "../api/types";
import { useNavigate } from "react-router-dom";
import { User } from "../types/user";
import { Navbar } from "./Navbar";

type Props = {
  children: string | JSX.Element | JSX.Element[] | (() => JSX.Element);
};

function AuthRequired({ children }: Props) {
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
      navigate("/login");
    }
  }, [navigate]);

  useEffect(() => {
    WhoAmI();
  }, [WhoAmI]);

  return (
    <>
      {user && <Navbar user={user} />}
      {children}
    </>
  );
}

export default withRouter(AuthRequired);
