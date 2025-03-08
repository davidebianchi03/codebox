import { Card, CardBody, CardHeader } from "reactstrap";
import { Workspace } from "../../types/workspace";
import { useCallback, useEffect, useState } from "react";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

interface Props {
  workspace: Workspace;
  fetchInterval: number;
}

export default function WorkspaceLogs({ workspace, fetchInterval }: Props) {
  const [logs, setLogs] = useState<string>("");

  const FetchLogs = useCallback(async () => {
    var [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${workspace.id}/logs`,
      "GET",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      setLogs(responseData.logs);
    }
  }, [workspace]);

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
        <CardHeader>
          <h3 className="mb-0">Logs</h3>
        </CardHeader>
        <CardBody>
          <div
            className="w-100 p-3"
            style={{
              backgroundColor: "var(--tblr-dark-bg-subtle)",
              borderRadius: 3,
              fontFamily: "Consolas",
            }}
          >
            {logs.split("\n").map((line, index) => (
              <p className="mb-0" key={index}>
                {line}
              </p>
            ))}
          </div>
        </CardBody>
      </Card>
    </>
  );
}
