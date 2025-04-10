import {
  Card,
  CardBody,
  CardHeader,
  Col,
  Row,
  Table,
  Tooltip,
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
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

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
  const [warningTooltipIsOpen, setWarningTooltipIsOpen] = useState(false);

  const FetchSelectedContainer = useCallback(
    async (containerName: string) => {
      // fetch container details
      var [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${
          workspace.id
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
        `${Http.GetServerURL()}/api/v1/workspace/${
          workspace.id
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
      }
    } else {
      setContainers([]);
    }
    setLoading(false);
  }, [
    workspace,
    selectedContainer,
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
          <CardHeader>
            <h3 className="mb-0">Containers</h3>
          </CardHeader>
          <CardBody>
            {containers.length > 0 ? (
              <>
                <Row>
                  <Col sm={4}>
                    <Table>
                      <tbody>
                        {containers.map((c) => (
                          <tr key={c.container_id}>
                            <th
                              style={{
                                cursor: "pointer",
                                background: `${
                                  selectedContainer?.container_name ===
                                  c.container_name
                                    ? "rgba(var(--tblr-primary-rgb), 0.2)"
                                    : "transparent"
                                }`,
                                borderRadius: "4px",
                              }}
                              onClick={() => {
                                FetchSelectedContainer(c.container_name);
                                FetchSelectedContainerPorts(c.container_name);
                              }}
                            >
                              <span>{c.container_name}</span>
                            </th>
                          </tr>
                        ))}
                      </tbody>
                    </Table>
                  </Col>
                  <Col sm={8}>
                    {selectedContainer && (
                      <>
                        <h3>
                          {selectedContainer.container_name}
                          {new Date().getTime() -
                            new Date(
                              selectedContainer.agent_last_contact
                            ).getTime() >
                          5 * 60 * 1000 ? (
                            <span className="ms-2 text-warning">
                              <Tooltip
                                isOpen={warningTooltipIsOpen}
                                target="warningIcon"
                                toggle={() =>
                                  setWarningTooltipIsOpen(!warningTooltipIsOpen)
                                }
                              >
                                Last contact with the agent running in this
                                container was more than 5 minutes ago.
                              </Tooltip>
                              <FontAwesomeIcon
                                id="warningIcon"
                                icon={faTriangleExclamation}
                              />
                            </span>
                          ) : null}
                          <br />
                          <small className="text-muted">
                            {selectedContainer.container_image}
                          </small>
                        </h3>
                        <div>
                          <Row
                            className="border rounded align-items-center"
                            style={{ width: 300, cursor: "pointer" }}
                            onClick={() => {
                              window.location.href = `vscode://davidebianchi.codebox-remote/open?workspace_id=${workspace.id}&container_name=${selectedContainer.container_name}&server_hostname=${settings?.server_hostname}`;
                            }}
                          >
                            <Col sm={2}>
                              <img src={VsCodeIcon} alt="" width={25} />
                            </Col>
                            <Col sm={10}>
                              <h4 className="mb-0">Visual studio code</h4>
                              <small className="text-muted">
                                Open workspace in visual studio code
                              </small>
                            </Col>
                          </Row>
                        </div>
                      </>
                    )}
                    <h4 className="mt-3">Exposed ports</h4>
                    <div>
                      {selectedContainerExposedPorts.length > 0 ? (
                        <>
                          {selectedContainerExposedPorts.map((port) => (
                            <Row
                              key={port.port_number}
                              className="border rounded align-items-center"
                              style={{ width: 300, cursor: "pointer" }}
                              onClick={() => {
                                var portUrl = `${window.location.protocol}//${settings?.server_hostname}/api/v1/workspace/${workspace.id}/container/${selectedContainer?.container_name}/forward-http/${port.port_number}?path=%2F`;
                                if (settings?.use_subdomains) {
                                  portUrl = `${window.location.protocol}//codebox--${workspace.id}--${selectedContainer?.container_name}--${port.port_number}.${settings?.server_hostname}`;
                                }
                                window.open(portUrl, "_blank")?.focus();
                              }}
                            >
                              <Col sm={2}>
                                <img
                                  src={
                                    port.public
                                      ? PublicPortIcon
                                      : PrivatePortIcon
                                  }
                                  alt=""
                                  width={25}
                                />
                              </Col>
                              <Col sm={10}>
                                <h4 className="mb-0">{port.service_name}</h4>
                                <small className="text-muted">
                                  {port.port_number}
                                </small>
                              </Col>
                            </Row>
                          ))}
                        </>
                      ) : (
                        <>
                          <h5>This container has no exposed ports</h5>
                        </>
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
        </Card>
      )}
    </>
  );
}
