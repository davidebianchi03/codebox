import React, { useCallback, useEffect, useState } from "react";
import { ContainerPort, Workspace, WorkspaceContainer } from "../../types/workspace";
import { APIListWorkspaceContainerPorts } from "../../api/workspace";
import { ExposedPortsDropdown } from "./ExposedPortsDropdown";
import { WorkspaceContainerService } from "./WorkspaceContainerService";
import VsCodeIcon from "../../assets/images/vscode.png";
import TerminalIcon from "../../assets/images/terminal.png";
import PublicPortIcon from "../../assets/images/earth.png";
import PrivatePortIcon from "../../assets/images/padlock.png";
import { Col, Row } from "react-bootstrap";

interface SelectedContainerDetailsProps {
    workspace: Workspace;
    container: WorkspaceContainer;
}

export function SelectedContainerDetails({
    workspace,
    container
}: SelectedContainerDetailsProps) {

    const [containerPorts, setContainerPorts] = useState<ContainerPort[]>([]);

    const FetchSelectedContainerPorts = useCallback(async () => {
        const ports = await APIListWorkspaceContainerPorts(workspace.id, container.container_name);
        if (ports) {
            setContainerPorts(ports);
        }
    }, [container, workspace]);

    useEffect(() => {
        FetchSelectedContainerPorts();
    }, [FetchSelectedContainerPorts])

    return (
        <React.Fragment>
            <h4 className="d-flex justify-content-end mt-1">
                <ExposedPortsDropdown
                    onChange={() => {
                        FetchSelectedContainerPorts();
                    }}
                    workspace={workspace}
                    container={container}
                />
            </h4>
            {container && (
                <React.Fragment>
                    <WorkspaceContainerService
                        icon={VsCodeIcon}
                        title="Visual Studio Code"
                        description="Open container in visual studio code"
                        url={
                            `vscode://davidebianchi.codebox-remote/open?workspace_id=${workspace.id}` +
                            `&container_name=${container.container_name}` +
                            `&server_hostname=${import.meta.env.VITE_SERVER_URL === "" ?
                                window.location.host : new URL(import.meta.env.VITE_SERVER_URL).hostname}`
                        }
                    />
                    <WorkspaceContainerService
                        icon={TerminalIcon}
                        title="Terminal"
                        description="Open terminal"
                        url={`${import.meta.env.VITE_SERVER_URL}/views/workspace/${workspace.id}/container/${container.container_name}/terminal`}
                    />
                </React.Fragment>
            )}
            <div>
                {containerPorts.length > 0 && (
                    <Row>
                        {containerPorts.map((port, index) => (
                            <Col md={12} className="my-1" key={index} >
                                <WorkspaceContainerService
                                    icon={
                                        port.public
                                            ? PublicPortIcon
                                            : PrivatePortIcon
                                    }
                                    title={port.service_name}
                                    description={`Port: ${port.port_number}`}
                                    url={port.port_url}
                                />
                            </Col>
                        ))}
                    </Row>
                )}
            </div>
        </React.Fragment >
    )
}
