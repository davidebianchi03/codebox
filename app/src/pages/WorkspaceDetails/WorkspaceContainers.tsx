import {
  Badge,
  Card,
  CardBody,
  CardHeader,
  Col,
  Row,
} from "reactstrap";
import {
  ContainerPort,
  Workspace,
  WorkspaceContainer,
} from "../../types/workspace";
import React, { useCallback, useEffect, useState } from "react";
import VsCodeIcon from "../../assets/images/vscode.png";
import TerminalIcon from "../../assets/images/terminal.png";
import PublicPortIcon from "../../assets/images/earth.png";
import PrivatePortIcon from "../../assets/images/padlock.png";
import { APIListWorkspaceContainers, APIListWorkspaceContainerPorts, APIRetrieveWorkspaceContainer } from "../../api/workspace";
import { ExposedPortsDropdown } from "./ExposedPortsDropdown";
import { WorkspaceContainerService } from "./WorkspaceContainerService";

interface Props {
  workspace: Workspace;
  fetchInterval: number;
}

export default function WorkspaceContainers({
  workspace,
  fetchInterval,
}: Props) {
  const [containers, setContainers] = useState<WorkspaceContainer[]>([]);
  const [selectedContainer, setSelectedContainer] =
    useState<WorkspaceContainer | null>(null);
  const [selectedContainerExposedPorts, setSelectedContainerExposedPorts] =
    useState<ContainerPort[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  const FetchSelectedContainer = useCallback(
    async (containerName: string) => {
      const c = await APIRetrieveWorkspaceContainer(workspace.id, containerName);
      if (c) {
        setSelectedContainer(c);
      }
    },
    [workspace]
  );

  const FetchSelectedContainerPorts = useCallback(
    async (containerName: string) => {
      const ports = await APIListWorkspaceContainerPorts(workspace.id, containerName);
      if (ports) {
        setSelectedContainerExposedPorts(ports);
      }
    },
    [workspace]
  );

  const FetchContainers = useCallback(async () => {
    const c = await APIListWorkspaceContainers(workspace.id);
    if (c) {
      setContainers(c);
      if (selectedContainer === null && c.length > 0) {
        FetchSelectedContainer(c[0].container_name);
        FetchSelectedContainerPorts(c[0].container_name);
      }

      if (selectedContainer !== null && c.length === 0) {
        setSelectedContainer(null);
      } else {
        var sc = c.find(
          (container) =>
            container.container_id === selectedContainer?.container_id
        );
        if (sc) {
          setSelectedContainer(sc);
        }
      }
    } else {
      setContainers([]);
    }
    setLoading(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    workspace,
    FetchSelectedContainer,
    FetchSelectedContainerPorts,
  ]);


  useEffect(() => {
    FetchContainers();
    const interval = setInterval(FetchContainers, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchContainers, fetchInterval]);

  return (
    <>
      {!loading && (
        <Card>
          <CardHeader className="border-0">
            <h3 className="mb-0">Containers</h3>
          </CardHeader>
          <CardBody className="pt-0">
            {containers.length > 0 ? (
              <>
                <Row className="justify-content-between">
                  <Col sm={4}>
                    <div className="d-flex flex-column">
                      {containers.map((c) => (
                        <div
                          key={c.container_id}
                          className={`my-1 py-2 px-2 w-100 ${c.container_id === selectedContainer?.container_id ? "border rounded" : ""}`}
                          style={{
                            cursor: "pointer",
                            borderRadius: "7px",
                          }}
                          onClick={() => {
                            FetchSelectedContainer(c.container_name);
                            FetchSelectedContainerPorts(c.container_name);
                          }}
                        >
                          <div className="d-flex justify-content-between">
                            <h4 className="mb-0">{c.container_name}</h4>
                            {new Date().getTime() -
                              new Date(
                                selectedContainer?.agent_last_contact || ""
                              ).getTime() >
                              5 * 60 * 1000 ? (
                              <Badge
                                color="warning"
                                className="text-white"
                                style={{ opacity: c.container_id === selectedContainer?.container_id ? "1" : "0.7" }}
                              >
                                Not connected
                              </Badge>
                            ) : (
                              <Badge
                                color="success"
                                className="text-white"
                                style={{ opacity: c.container_id === selectedContainer?.container_id ? "1" : "0.7" }}
                              >
                                Connected
                              </Badge>
                            )}
                          </div>
                          <small className="text-muted">{c.container_image}</small>
                        </div>
                      ))}
                    </div>
                  </Col>
                  <Col sm={7} className="ms-3">
                    <h4 className="d-flex justify-content-end mt-1">
                      {selectedContainer && (
                        <ExposedPortsDropdown
                          onChange={() => {
                            FetchSelectedContainerPorts(selectedContainer.container_name);
                          }}
                          workspace={workspace}
                          container={selectedContainer}
                        />
                      )}
                    </h4>
                    {selectedContainer && (
                      <React.Fragment>
                        <WorkspaceContainerService
                          icon={VsCodeIcon}
                          title="Visual Studio Code"
                          description="Open container in visual studio code"
                          url={
                            `vscode://davidebianchi.codebox-remote/open?workspace_id=${workspace.id}` +
                            `&container_name=${selectedContainer.container_name}` +
                            `&server_hostname=${import.meta.env.VITE_SERVER_URL === "" ?
                              window.location.host : new URL(import.meta.env.VITE_SERVER_URL).hostname}`
                          }
                        />
                        <WorkspaceContainerService
                          icon={TerminalIcon}
                          title="Terminal"
                          description="Open terminal"
                          url={`${import.meta.env.VITE_SERVER_URL}/views/workspace/${workspace.id}/container/${selectedContainer.container_name}/terminal`}
                        />
                      </React.Fragment>
                    )}
                    <div>
                      {selectedContainerExposedPorts.length > 0 && (
                        <Row>
                          {selectedContainerExposedPorts.map((port) => (
                            <Col md={12} className="my-1">
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
                  </Col>
                </Row>
              </>
            ) : (
              <>
                {workspace.status === "running" ? (
                  <h4 className="text-center mb-3">No containers found</h4>
                ) : (
                  <h4 className="text-center mb-3">Workspace is not running</h4>
                )}
              </>
            )}
          </CardBody>
        </Card >
      )
      }
    </>
  );
}
