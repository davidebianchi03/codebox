import { Card, CardBody, CardHeader, Col, Row } from "reactstrap";
import { Workspace } from "../../types/workspace";
import React, { useCallback, useEffect, useRef, useState } from "react";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

interface Props {
  workspace: Workspace;
  fetchInterval: number;
}

export default function WorkspaceLogs({ workspace, fetchInterval }: Props) {
  const [logs, setLogs] = useState<string>("");
  const logsContainerRef = useRef<any>(null);

  const FetchLogs = useCallback(async () => {
    var [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${workspace.id}/logs`,
      "GET",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      var scrollToBottom = (responseData.logs as string).length !== logs.length;
      setLogs(responseData.logs);

      // if length of logs has changed scroll to the last row of logs
      if (scrollToBottom) {
        var logsContainer = logsContainerRef.current;
        if (logsContainer) {
          logsContainer.scrollTop = logsContainer.scrollHeight;
        }
      }
    }
  }, [workspace, logs.length]);

  useEffect(() => {
    FetchLogs();
    const interval = setInterval(FetchLogs, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchLogs, fetchInterval]);

  return (
    <>
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
    </>
  );
}
