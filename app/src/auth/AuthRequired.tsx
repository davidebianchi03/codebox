import { withRouter } from "../common/router";
import React, { useCallback, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { RetrieveCurrentUserDetails } from "../api/common";
import { useDispatch, useSelector } from 'react-redux'
import { setUser } from "../redux/slices/user";
import { RootState } from "../redux/store";

export type AuthRequiredProps = {
  children: React.ReactNode;
};

function AuthRequired({ children }: AuthRequiredProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const dispatch = useDispatch();  
  const currentUser = useSelector((state:RootState) => state.user);

  const WhoAmI = useCallback(async () => {
    const u = await RetrieveCurrentUserDetails();
    if (u) {
      dispatch(setUser(u));
    } else {
      navigate(`/login?next=${encodeURIComponent(location.pathname)}`);
    }
  }, [dispatch, navigate, location.pathname]);

  useEffect(() => {
    WhoAmI();
  }, [WhoAmI]);

  if (currentUser) {
    return (
      <React.Fragment>
        {children}
      </React.Fragment>
    );
  }

  return (
    <React.Fragment />
  );
}

export default withRouter(AuthRequired);
