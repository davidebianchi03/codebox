import { withRouter } from "../common/router";
import React, { useCallback, useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useDispatch, useSelector } from "react-redux";
import { RetrieveCurrentUserDetails } from "../api/common";
import { setUser } from "../redux/slices/user";
import { RootState } from "../redux/store";

export type SuperUserRequiredProps = {
  children: any;
  showNavbar?: boolean;
};

function SuperUserRequired({ children }: SuperUserRequiredProps) {
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
    if (currentUser.is_superuser) {
      return (
        <React.Fragment>
          {children}
        </React.Fragment>
      );
    }
  }

  return (
    <React.Fragment />
  );
}

export default withRouter(SuperUserRequired);
