import "./WorkspaceDetails.css"
import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom"
import BasePage from "../base/Base";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import Card from "../../theme/components/card/Card";
import { RetrieveBeautyNameForStatus, RetrieveColorForWorkspaceStatus } from "../../utils/workspaceStatus";
import EarthIcon from "../../assets/images/earth.png";
import LockIcon from "../../assets/images/padlock.png";

interface WorkspaceDetailsProps {

}

interface ForwardedPortsDetails {
    port_number: number
    active: boolean
    connection_type: string
    public: boolean
    url?: string
}

interface ContainerDetails {
    id?: number
    type?: string
    name?: string
    container_user?: string
    container_status?: string
    agent_status?: string
    can_connect_remote_developing?: boolean
    workspace_path_in_container?: string
    forwarded_ports?: Array<ForwardedPortsDetails> | null
}

interface WorkspaceDetails {
    id?: number
    name?: string
    status?: string
    type?: string
    containers?: Array<ContainerDetails> | null
}

export default function WorkspaceDetails(props: WorkspaceDetailsProps) {

    const navigate = useNavigate();
    const [workspaceDetails, setWorkspaceDetails] = useState<WorkspaceDetails>({});
    const [workspaceLogs, setWorkspaceLogs] = useState<string>("");
    const [selectedContainerIndex, setSelectedContainerIndex] = useState<number>(0);
    const [selectedContainer, setSelectedContainer] = useState<ContainerDetails>({});
    const [instanceSettings, setInstanceSettings] = useState<any>({});

    var logsContainerRef = useRef<null | HTMLDivElement>(null);


    let params = useParams();
    let workspaceId = params.workspaceId;
    if (workspaceId === undefined) {
        navigate("/");
    }

    // retrieve workspace details
    useEffect(() => {
        if (workspaceDetails.containers) {
            if (selectedContainerIndex < 0 || selectedContainerIndex > workspaceDetails.containers?.length) {
                setSelectedContainerIndex(0);
                return;
            }

            var selectedContainer = workspaceDetails.containers[selectedContainerIndex];
            if (instanceSettings && selectedContainer.forwarded_ports) {
                for (let i = 0; i < selectedContainer.forwarded_ports.length; i++) {
                    if (selectedContainer.forwarded_ports[i].connection_type === "http" && instanceSettings.server_hostname) {
                        if(instanceSettings.use_subdomains) {
                            selectedContainer.forwarded_ports[i].url = `http://codebox--w${workspaceDetails.id}--c${selectedContainer.name}--p${selectedContainer.forwarded_ports[i].port_number}.${instanceSettings.server_hostname}`;
                        } else {
                            selectedContainer.forwarded_ports[i].url = `http://${instanceSettings.server_hostname}/api/v1/workspace/${workspaceDetails.id}/container/${selectedContainer.id}/forward/${selectedContainer.forwarded_ports[i].port_number}`;
                        }
                    }
                }
            }

            setSelectedContainer(selectedContainer);
        }
    }, [workspaceDetails, selectedContainerIndex, instanceSettings]);

    const UpdateWorkspaceDetails = async () => {
        let [, , is] = await Http.Request(`${Http.GetServerURL()}/api/v1/instance-settings`, "GET", null);
        setInstanceSettings(is);
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
            logsContainerRef.current?.scroll({ top: logsContainerRef.current?.scrollHeight });
        }
    };

    // use effect to check if is ne
    let updateWorkspaceLogsInterval: NodeJS.Timer | null = null;
    useEffect(() => {
        if (workspaceDetails.status === "creating" || workspaceDetails.status === "starting" || workspaceDetails.status === "stopping") {
            retrieveWorkspaceLogs();
            updateWorkspaceLogsInterval = setInterval(retrieveWorkspaceLogs, 800);
        } else {
            if (updateWorkspaceLogsInterval !== null) {
                clearInterval(updateWorkspaceLogsInterval);
            }
        }

        return () => {
            if (updateWorkspaceLogsInterval !== null) {
                clearInterval(updateWorkspaceLogsInterval);
            }
        }
    }, [workspaceDetails]);


    var updatedWorkspaceDetailsInterval: NodeJS.Timer | null = null;
    useEffect(() => {
        UpdateWorkspaceDetails();
        retrieveWorkspaceLogs();
        if (updatedWorkspaceDetailsInterval !== null) {
            clearInterval(updatedWorkspaceDetailsInterval);
        }
        updatedWorkspaceDetailsInterval = setInterval((UpdateWorkspaceDetails), 5000);

        return () => {
            if (updatedWorkspaceDetailsInterval !== null) {
                clearInterval(updatedWorkspaceDetailsInterval);
            }
        }
    }, []);

    const StartWorkspace = async () => {
        await Http.Request(`${Http.GetServerURL()}/api/v1/workspace/${workspaceId}/start`, "POST", null);
        UpdateWorkspaceDetails();
    }

    const StopWorkspace = async () => {
        await Http.Request(`${Http.GetServerURL()}/api/v1/workspace/${workspaceId}/stop`, "POST", null);
        UpdateWorkspaceDetails();
    }

    const OpenRemoteDeveloping = () => {
        if (
            workspaceDetails.id !== undefined &&
            selectedContainer.id !== undefined &&
            selectedContainer.container_user !== undefined &&
            selectedContainer.workspace_path_in_container !== undefined
        ) {
            var urlQueryParams = new URLSearchParams();
            urlQueryParams.set("workspace_id", workspaceDetails.id?.toString());
            urlQueryParams.set("container_id", selectedContainer.id?.toString());
            urlQueryParams.set("container_port", "2222");
            urlQueryParams.set("container_user", selectedContainer.container_user);
            urlQueryParams.set("workspace_path", encodeURI(selectedContainer.workspace_path_in_container));
            document.location.href = `vscode://davidebianchi.codebox-remote/open?${urlQueryParams.toString()}`;
        }
    }

    let borderColorCssVar = RetrieveColorForWorkspaceStatus(workspaceDetails.status)
    return (
        <BasePage authRequired={true}>
            <Card style={{
                width: "90%",
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
                    <div style={{ display: "flex", alignItems: "center" }}>
                        {
                            workspaceDetails.status !== "creating" && workspaceDetails.status !== "starting" && workspaceDetails.status !== "stopping" ?
                                <div
                                    style={{
                                        marginRight: "10pt",
                                        padding: "5pt 10pt",
                                        border: "solid var(--background-divider) 1.5px",
                                        borderRadius: "5px",
                                        cursor: "pointer"
                                    }}
                                    onClick={() => {
                                        if (workspaceDetails.status === "running" || workspaceDetails.status === "error") {
                                            StopWorkspace();
                                        } else {
                                            StartWorkspace();
                                        }
                                    }}
                                >
                                    {
                                        workspaceDetails.status === "running" || workspaceDetails.status === "error" ?
                                            "Stop workspace" :
                                            "Start workspace"
                                    }
                                </div> : null
                        }
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
                </div>
                <div style={{
                    border: "solid var(--background-divider) 1px",
                    marginTop: "30px",
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
                    }}
                        className="workspace-logs-container"
                        ref={logsContainerRef}>
                        {workspaceLogs}
                    </div>
                </div>
            </Card>
            <Card
                style={{
                    width: "90%",
                    margin: "auto",
                    marginTop: "40pt",
                    marginBottom: "30pt",
                    paddingTop: "10pt",
                }}
            >
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div style={{ display: "flex", flexDirection: "column", alignItems: "start", justifyContent: "center" }}>
                        <h3 style={{ marginBottom: 0, marginTop: 0 }}>Containers</h3>
                    </div>
                </div>
                {workspaceDetails.containers !== null ?
                    <div style={{ display: "flex" }} className="workspace-containers">
                        <ul>
                            {
                                workspaceDetails.containers?.map((container, index) => (
                                    <li style={index === 0 ? { borderTop: "none" } : {}} onClick={() => { setSelectedContainerIndex(index); UpdateWorkspaceDetails(); }} key={container.id}>
                                        {container.name}
                                    </li>
                                ))
                            }
                        </ul>
                        <div style={{ marginLeft: "10pt" }}>
                            <div>
                                <h4 style={{ marginBottom: 0 }}>{selectedContainer.name}</h4>
                                <small style={{ color: "var(--grey-300)" }}>{selectedContainer.type === "docker_container" ? "Docker container" : "N/A"}</small>
                            </div>
                            <div>
                                <div>
                                    <h5 style={{ marginBottom: "8pt" }}>Forwarded ports</h5>
                                </div>
                                <div style={{
                                    display: "flex",
                                    flexWrap: "wrap"
                                }}>
                                    {
                                        selectedContainer.forwarded_ports?.map((port) => (
                                            <a style={{
                                                display: "flex",
                                                flexDirection: "row",
                                                alignItems: "center",
                                                border: "solid var(--background-divider) 1px",
                                                minWidth: "150px",
                                                padding: "4pt 7pt",
                                                borderRadius: "4pt",
                                                margin: "2pt"
                                            }}
                                                key={port.port_number}
                                            >
                                                {
                                                    port.public ?
                                                        <img alt="Public access" src={EarthIcon} width={"20px"} height={"20px"} />
                                                        :
                                                        <img alt="Authentication required" src={LockIcon} width={"20px"} height={"20px"} />
                                                }

                                                {port.port_number === 2222 ?
                                                    <span
                                                        style={{
                                                            display: "flex",
                                                            flexDirection: "column",
                                                            flexWrap: "wrap",
                                                            marginLeft: "5pt",
                                                            cursor: "pointer",
                                                        }}
                                                        onClick={() => OpenRemoteDeveloping()}
                                                    >
                                                        <span>
                                                            {port.port_number}
                                                            <small style={{ fontSize: "9pt", marginLeft: "4pt" }}>(remote developing)</small>
                                                        </span>
                                                        <small style={{ fontSize: "8pt", color: "var(--grey-300)", }}>
                                                            {port.connection_type === "ws" ? "TCP over WS" : "HTTP"}
                                                        </small>
                                                    </span>
                                                    :
                                                    <span
                                                        style={{
                                                            display: "flex",
                                                            flexDirection: "column",
                                                            flexWrap: "wrap",
                                                            marginLeft: "5pt",
                                                            cursor: port.url && port.url !== "" ? "pointer": "default",
                                                        }}
                                                        onClick={() => {
                                                            if (port.url) {
                                                                if (port.url !== "") {
                                                                    window.open(port.url, "_blank")?.focus();
                                                                }
                                                            }
                                                        }}
                                                    >
                                                        <span>
                                                            {port.port_number}
                                                        </span>
                                                        <small style={{ fontSize: "8pt", color: "var(--grey-300)", }}>
                                                            {port.connection_type === "ws" ? "TCP over WS" : "HTTP"}
                                                        </small>
                                                    </span>
                                                }
                                            </a>
                                        ))
                                    }
                                </div>
                            </div>
                        </div>
                    </div>
                    :
                    <div style={{ textAlign: "center" }}>
                        <h5 style={{ margin: "10pt" }}>No running containers
                            {
                                workspaceDetails.status === "stopped" ?
                                    <span>
                                        ,&nbsp;
                                        <a style={{
                                            textDecoration: "underline",
                                            cursor: "pointer"
                                        }}
                                            onClick={() => StartWorkspace()}
                                        >
                                            start workspace
                                        </a>
                                    </span>
                                    : null
                            }
                        </h5>
                    </div>
                }
            </Card>
        </BasePage>
    );
}