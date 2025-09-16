import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import { Link, useNavigate } from "react-router-dom";
import { InstanceSettings } from "../types/settings";
import { useCallback, useEffect, useState } from "react";
import { User } from "../types/user";
import { Logout, RetrieveInstanceSettings } from "../api/common";
import { UserDropdown } from "./UserDropdown";

interface Props {
  showLogo?: boolean;
}

export function Navbar({ showLogo = true }: Props) {
  return (
    <header className="navbar navbar-expand-md d-print-none">
      <div className="container-xl">
        {showLogo ? (
          <Link
            className="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3"
            to="/"
          >
            <img src={CodeboxLogo} alt="logo" width={120} />
          </Link>
        ) : (
          <div />
        )}
        <UserDropdown />
      </div>
    </header>
  );
}
