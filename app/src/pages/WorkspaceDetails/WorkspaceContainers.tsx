import {
  Badge,
  Button,
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
import { useCallback, useEffect, useState } from "react";
import VsCodeIcon from "../../assets/images/vscode.png";
import PublicPortIcon from "../../assets/images/earth.png";
import PrivatePortIcon from "../../assets/images/padlock.png";
import { InstanceSettings } from "../../types/settings";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { EditExposedPortsModal } from "./EditExposedPortsModal";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronRight } from "@fortawesome/free-solid-svg-icons";

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
  const [settings, setSettings] = useState<InstanceSettings>();
  const [showEditExposedPortsModal, setShowEditExposedPortaModal] =
    useState(false);

  const FetchSelectedContainer = useCallback(
    async (containerName: string) => {
      // fetch container details
      var [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${workspace.id
        }/container/${containerName}`,
        "GET",
        null
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        setSelectedContainer(responseData);
      } else {
        return;
      }
    },
    [workspace]
  );

  const FetchSelectedContainerPorts = useCallback(
    async (containerName: string) => {
      // fetch exposed ports
      var [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${workspace.id
        }/container/${containerName}/port`,
        "GET",
        null
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        setSelectedContainerExposedPorts(responseData);
      }
    },
    [workspace]
  );

  const FetchContainers = useCallback(async () => {
    var [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${workspace.id}/container`,
      "GET",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      setContainers(responseData);
      if (selectedContainer === null && responseData.length > 0) {
        FetchSelectedContainer(responseData[0].container_name);
        FetchSelectedContainerPorts(responseData[0].container_name);
      }

      if (selectedContainer !== null && responseData.length === 0) {
        setSelectedContainer(null);
      } else {
        var sc = (responseData as WorkspaceContainer[]).find(
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

  const FetchSettings = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/instance-settings`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setSettings(responseBody as InstanceSettings);
    }
  }, []);

  useEffect(() => {
    FetchSettings();
    FetchContainers();
    const interval = setInterval(FetchContainers, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchContainers, FetchSettings, fetchInterval]);

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
                      <Button
                        color="outline-muted"
                        className="text-white"
                        onClick={() => setShowEditExposedPortaModal(true)}
                      >
                        <span className="me-2">Edit exposed ports</span>
                      </Button>
                    </h4>
                    {selectedContainer && (
                      <>
                        <div style={{ marginTop: 5 }} className="my-1">
                          <div
                            className="d-flex border rounded align-items-center px-2"
                            style={{ cursor: "pointer", height: 50 }}
                            onClick={() => {
                              window.location.href = `vscode://davidebianchi.codebox-remote/open?workspace_id=${workspace.id}&container_name=${selectedContainer.container_name}&server_hostname=${settings?.server_hostname}`;
                            }}
                          >
                            <img src={VsCodeIcon} alt="vscode" width={25} className="me-3" />
                            <div className="d-flex justify-content-between align-items-center w-100 me-2">
                              <div className="d-flex align-items-center">
                                <h4 className="mb-0">Visual Studio Code</h4>
                                <span className="text-muted ms-5">
                                  Open container in visual studio code
                                </span>
                              </div>
                              <span className="text-muted">
                                <FontAwesomeIcon icon={faChevronRight} />
                              </span>
                            </div>
                          </div>
                        </div>
                      </>
                    )}
                    <div>
                      {selectedContainerExposedPorts.length > 0 && (
                        <Row>
                          {selectedContainerExposedPorts.map((port) => (
                            <Col md={12} className="my-1">
                              <div
                                key={port.port_number}
                                className="d-flex border rounded align-items-center px-2"
                                style={{ cursor: "pointer", height: 50 }}
                                onClick={() => {
                                  var portUrl = `${window.location.protocol}//${settings?.server_hostname}/api/v1/workspace/${workspace.id}/container/${selectedContainer?.container_name}/forward-http/${port.port_number}?path=%2F`;
                                  if (settings?.use_subdomains) {
                                    portUrl = `${window.location.protocol}//codebox--${workspace.id}--${selectedContainer?.container_name}--${port.port_number}.${settings?.server_hostname}`;
                                  }
                                  window.open(portUrl, "_blank")?.focus();
                                }}
                              >
                                <img
                                  src={
                                    port.public
                                      ? PublicPortIcon
                                      : PrivatePortIcon
                                  }
                                  className="me-3"
                                  alt=""
                                  width={25}
                                />
                                <div className="d-flex justify-content-between align-items-center w-100 me-2">
                                  <div className="d-flex align-items-center">
                                    <h4 className="mb-0">{port.service_name}</h4>
                                    <span className="text-muted ms-5">
                                      Port: {port.port_number}
                                    </span>
                                  </div>
                                  <span className="text-muted">
                                    <FontAwesomeIcon icon={faChevronRight} />
                                  </span>
                                </div>
                              </div>
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
                <h4 className="text-center">No containers found</h4>
              </>
            )}
          </CardBody>
          {workspace && selectedContainer && (
            <EditExposedPortsModal
              isOpen={showEditExposedPortsModal}
              onClose={() => {
                setShowEditExposedPortaModal(false);
                FetchSelectedContainerPorts(selectedContainer.container_name);
              }}
              workspace={workspace}
              container={selectedContainer}
            />
          )}
        </Card >
      )
      }
    </>
  );
}
