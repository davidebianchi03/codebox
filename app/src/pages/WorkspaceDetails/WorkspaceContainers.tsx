import {
  Badge,
  Card,
  CardBody,
  CardHeader,
  Col,
  Row,
} from "reactstrap";
import {
  Workspace,
  WorkspaceContainer,
} from "../../types/workspace";
import React, { useCallback, useEffect, useState } from "react";
import { APIListWorkspaceContainers } from "../../api/workspace";
import { SelectedContainerDetails } from "./SelectedContainerDetails";

interface WorkspaceContainersProps {
  workspace: Workspace;
  fetchInterval: number;
}

export default function WorkspaceContainers({
  workspace,
  fetchInterval,
}: WorkspaceContainersProps) {
  const [containers, setContainers] = useState<WorkspaceContainer[]>([]);
  const [selectedContainerIndex, setSelectedContainerIndex] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);

  const FetchContainers = useCallback(async () => {
    const c = await APIListWorkspaceContainers(workspace.id);
    if (c) {
      setContainers(c);
      if (selectedContainerIndex >= c.length) {
        setSelectedContainerIndex(0);
      }
    } else {
      setContainers([]);
    }
    setLoading(false);
  }, [
    workspace,
  ]);

  useEffect(() => {
    FetchContainers();
    const interval = setInterval(FetchContainers, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchContainers, fetchInterval]);

  return (
    <React.Fragment>
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
                      {containers.map((c, index) => (
                        <div
                          key={c.container_id}
                          className={`my-1 py-2 px-2 w-100 ${selectedContainerIndex === index ? "border rounded" : ""}`}
                          style={{
                            cursor: "pointer",
                            borderRadius: "7px",
                          }}
                          onClick={() => {
                            setSelectedContainerIndex(index);
                          }}
                        >
                          <div className="d-flex justify-content-between">
                            <h4 className="mb-0">{c.container_name}</h4>
                            {new Date().getTime() -
                              new Date(
                                c?.agent_last_contact || ""
                              ).getTime() >
                              5 * 60 * 1000 ? (
                              <Badge
                                color="warning"
                                className="text-white"
                                style={{ opacity: c.container_id === c?.container_id ? "1" : "0.7" }}
                              >
                                Not connected
                              </Badge>
                            ) : (
                              <Badge
                                color="success"
                                className="text-white"
                                style={{ opacity: c.container_id === c?.container_id ? "1" : "0.7" }}
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
                    {selectedContainerIndex < containers.length && (
                      <SelectedContainerDetails
                        workspace={workspace}
                        container={containers[selectedContainerIndex]}
                      />
                    )}
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
    </React.Fragment>
  );
}
