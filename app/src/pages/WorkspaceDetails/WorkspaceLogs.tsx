import { Card, CardBody, CardHeader, Col, Row } from "reactstrap";
import { Workspace } from "../../types/workspace";
import React, { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";
import { APIRetrieveWorkspaceLogs } from "../../api/workspace";

interface Props {
  workspace: Workspace;
  fetchInterval: number;
}

export default function WorkspaceLogs({ workspace, fetchInterval }: Props) {
  const [logs, setLogs] = useState<string>("");
  const logsContainerRef = useRef<HTMLDivElement>(null);
  const latestRequestRef = useRef(0);

  const FetchLogs = useCallback(async () => {
    const requestId = ++latestRequestRef.current;
    const l = await APIRetrieveWorkspaceLogs(workspace.id);

    if (requestId === latestRequestRef.current && l !== undefined) {
      setLogs(l);
    }
  }, [workspace.id]);

  useEffect(() => {
    FetchLogs();
  }, [FetchLogs]);

  useLayoutEffect(() => {
    const logsContainer = logsContainerRef.current;
    if (logsContainer) {
      logsContainer.scrollTop = logsContainer.scrollHeight;
    }
  }, [logs]);

  useEffect(() => {
    if (fetchInterval > 0) {
      const interval = setInterval(FetchLogs, fetchInterval);

      return () => {
        clearInterval(interval);
      };
    }
  }, [FetchLogs, fetchInterval]);

  return (
    <React.Fragment>
      <Card>
        <CardHeader className="border-0">
          <h3 className="mb-0">Logs</h3>
        </CardHeader>
        <CardBody className="pt-0">
          <Row>
            <Col md={2} className="mb-2">
              <small className="text-muted">
                Last activity: {workspace.updated_at ? new Date(workspace.updated_at).toLocaleString() : "N/A"}
              </small>
            </Col>
            <Col md={10}>
              <div
                className="w-100 p-3"
                style={{
                  backgroundColor: "#000710",
                  borderRadius: 3,
                  fontFamily: "Consolas",
                  maxHeight: 250,
                  overflowY: "scroll"
                }}
                ref={logsContainerRef}
              >
                {logs.split("\n").map((line, index) => (
                  <p className="mb-0" key={index}>
                    {line}
                  </p>
                ))}
              </div>
            </Col>
          </Row>
        </CardBody>
      </Card>
    </React.Fragment>
  );
}
