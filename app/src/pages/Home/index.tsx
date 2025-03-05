import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, CardBody, Container, Input, Label } from "reactstrap";
import { Workspace } from "../../types/workspace";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

export default function HomePage() {
  const [searchText, setSearchText] = useState<string>("");
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);

  const FetchWorkspaces = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace`,
      "GET",
      null
    );
    if (status === RequestStatus.OK) {
      setWorkspaces(responseData as Workspace[]);
    } else {
      console.log(`Error: received ${statusCode} from server`);
    }
  }, []);

  useEffect(() => {
    FetchWorkspaces();
  }, [FetchWorkspaces]);

  return (
    <Container className="pb-4">
      <div className="mt-5 w-100 d-flex justify-content-end">
        <Button color="primary">Create workspace</Button>
      </div>
      <Card className="my-5">
        <CardBody>
          <Label>Filter workspaces:</Label>
          <Input
            placeholder="my awesome workspace"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
        </CardBody>
      </Card>
      <Card className="my-1">
        <CardBody>
          {workspaces.length &&
          workspaces.filter((w) => w.name.indexOf(searchText) >= 0).length > 0 ? (
            <>
              {workspaces.map((workspace: Workspace) => (
                <>
                  {workspace.name.indexOf(searchText) >= 0 && (
                    <>
                      <div className="d-flex align-items-center">
                        <div
                          style={{
                            width: 50,
                            height: 50,
                            fontSize: 20,
                            opacity: 0.5,
                            borderRadius: 4,
                          }}
                          className="bg-primary text-white d-flex align-items-center justify-content-center"
                        >
                          {workspace.name[0].toUpperCase()}
                        </div>
                        <div className="ms-4">
                          <h3 className="mb-0">{workspace.name}</h3>
                          <small className="text-muted">{workspace.type}</small>
                        </div>
                      </div>
                      <hr className="my-3" />
                    </>
                  )}
                </>
              ))}
            </>
          ) : (
            <p className="text-center">
              {workspaces.length > 0 ? (
                <>
                  No workspaces found matching '{searchText}'' {" "}
                  <span>
                    <u>create it</u>
                  </span>
                </>
              ) : (
                <>
                  No workspaces found,{" "}
                  <span>
                    <u>create your first workspace</u>
                  </span>
                </>
              )}
            </p>
          )}
        </CardBody>
      </Card>
    </Container>
  );
}
