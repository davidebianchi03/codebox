import React, { Component, ReactNode, useEffect, useState } from "react";
import { Navigate, Params, RouteProps, useLocation, useNavigate, useParams } from "react-router-dom"
import BasePage from "./base/Base";
import { Http } from "../api/http";
import { RequestStatus } from "../api/types";
import Card from "../theme/components/card/Card";
import Button from "../theme/components/button/Button";
import { RetrieveBeautyNameForStatus, RetrieveColorForWorkspaceStatus } from "../utils/workspaceStatus";

interface WorkspaceDetailsProps {

}

interface WorkspaceDetails {
    id?: number
    name?: string
    status?: string
    type?: string
}

export default function WorkspaceDetails(props: WorkspaceDetailsProps) {

    const navigate = useNavigate();
    const [workspaceDetails, setWorkspaceDetails] = useState<WorkspaceDetails>({});
    const [workspaceLogs, setWorkspaceLogs] = useState<string>("");

    let params = useParams();
    let workspaceId = params.workspaceId;
    if (workspaceId === undefined) {
        navigate("/");
    }

    // check that users are authenticated
    useEffect(() => {
        (async () => {
            let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/whoami`, "GET", null);
            if (status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
                navigate("/login")
            }
        })();
    }, []);

    // retrieve workspace details
    const UpdateWorkspaceDetails = async () => {
        let [status, statusCode, data, errorDescription] = await Http.Request(`${Http.GetServerURL()}/api/v1/workspace/${workspaceId}`, "GET", null);
        if (status === RequestStatus.OK) {
            setWorkspaceDetails(data);
        }
    };

    // retrieve workspace logs
    const retrieveWorkspaceLogs = async () => {
        let [status, statusCode, data, errorDescription] = await Http.Request(`${Http.GetServerURL()}/api/v1/workspace/${workspaceId}/logs`, "GET", null);
        if (status === RequestStatus.OK) {
            setWorkspaceLogs(data.logs.replaceAll("\r", "\n"));
        }
    };

    // use effect to check if is ne
    let updateWorkspaceLogsInterval: NodeJS.Timer | null = null;
    useEffect(() => {
        if (workspaceDetails.status === "creating" || workspaceDetails.status === "starting" || workspaceDetails.status === "stopping") {
            retrieveWorkspaceLogs();
            // updateWorkspaceLogsInterval = setInterval(retrieveWorkspaceLogs, 800);
        } else {
            // if (updateWorkspaceLogsInterval !== null) {
            //     clearInterval(updateWorkspaceLogsInterval);
            // }
        }
    }, [workspaceDetails]);


    useEffect(() => {
        UpdateWorkspaceDetails();
        retrieveWorkspaceLogs();
        setInterval((UpdateWorkspaceDetails), 5000);
    }, []);

    let borderColorCssVar = RetrieveColorForWorkspaceStatus(workspaceDetails.status)
    return (
        <BasePage>
            <Card style={{
                width: "90%",
                minWidth: "450px",
                margin: "auto",
                marginTop: "40pt",
                marginBottom: "30pt",
                paddingTop: "10pt",
                border: `solid var(${borderColorCssVar}) 1px`,
                boxShadow: `-1px 0px 15px -5px var(${borderColorCssVar})`,
            }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div style={{ display: "flex", flexDirection: "column", alignItems: "start", justifyContent: "center" }}>
                        <h3 style={{ marginBottom: 0, marginTop: 0 }}>{workspaceDetails.name}</h3>
                        <small style={{ color: "var(--grey-300)" }}>{workspaceDetails.type}</small>
                    </div>
                    <div style={{
                        background: `var(${borderColorCssVar})`,
                        fontSize: "11pt",
                        padding: "5px 10px",
                        borderRadius: "15px",
                        minWidth: "50px",
                        textAlign: "center",
                    }}>
                        {RetrieveBeautyNameForStatus(workspaceDetails.status)}
                    </div>
                </div>
                <div style={{
                    border: "solid var(--background-divider) 1px",
                    marginTop: "30px",
                    // padding: "10pt",
                    borderRadius: "5pt",
                }}>
                    <div style={{ padding: "10pt" }}>
                        <b>Logs</b>
                    </div>
                    <div style={{
                        padding: "10pt",
                        paddingBottom: "10pt",
                        background: "var(--background-color-paper)",
                        maxHeight: "200px",
                        overflowY: "auto",
                        fontFamily: "Consolas, monaco, monospace",
                        fontSize: "14px",
                        whiteSpace: "pre-wrap"
                    }}>
                        {workspaceLogs}
                    </div>
                </div>
            </Card>
        </BasePage>
    );
}