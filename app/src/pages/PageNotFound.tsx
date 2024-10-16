import React, { Component, ReactNode, useEffect, useState } from "react";
import { Navigate, Params, RouteProps, useLocation, useNavigate, useParams } from "react-router-dom"
import BasePage from "./base/Base";
import { Http } from "../api/http";
import { RequestStatus } from "../api/types";
import Card from "../theme/components/card/Card";
import Button from "../theme/components/button/Button";
import { RetrieveBeautyNameForStatus, RetrieveColorForWorkspaceStatus } from "../utils/workspaceStatus";
import { Link } from "react-router-dom";

import CodeboxLogoWhite from "../assets/images/logo-white.png";


export default function PageNotFound(props: any) {
    return (
        <div style={{
            width: "100%",
            height: "100%",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
        }}>
            <img src={CodeboxLogoWhite} alt="Codebox log" width={"200px"} />
            <h3 style={{ marginBottom: 0 }}>Oops....</h3>
            <h4 style={{ color: "var(--grey-500)" }}>Page not found</h4>
            <Link to={"/"} style={{ color: "var(--blue)" }}>Go back to home</Link>
        </div>
    );
}