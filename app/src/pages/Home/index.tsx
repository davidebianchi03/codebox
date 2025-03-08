import React, { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  CardBody,
  Container,
  Input,
  Label,
} from "reactstrap";
import { Workspace, WorkspaceType } from "../../types/workspace";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import {
  GetBeautyNameForStatus,
  GetWorkspaceStatusColor,
} from "../../common/workspace";
import { Link, useNavigate } from "react-router-dom";

export default function HomePage() {
  const [searchText, setSearchText] = useState<string>("");
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const navigate = useNavigate();

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

  const FetchWorkspaceTypes = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace-types`,
      "GET",
      null
    );
    if (status === RequestStatus.OK) {
      setWorkspaceTypes(responseData as WorkspaceType[]);
    } else {
      console.log(`Error: received ${statusCode} from server`);
    }
  }, []);

  useEffect(() => {
    FetchWorkspaces();
    FetchWorkspaceTypes();
  }, [FetchWorkspaces, FetchWorkspaceTypes]);

  return (
    <Container className="pb-4 mt-4">
      <div className="row g-2 align-items-center">
        <div className="col">
          <div className="page-pretitle">Overview</div>
          <h2 className="page-title">Workspaces</h2>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <div className="btn-list">
            <Button
              color="primary"
              onClick={() => navigate("/create-workspace")}
            >
              Create workspace
            </Button>
          </div>
        </div>
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
          workspaces.filter((w) => w.name.indexOf(searchText) >= 0).length >
            0 ? (
            <>
              {workspaces
                .sort((w1: any, w2: any) => {
                  var d1 = new Date(w1.updated_at);
                  var d2 = new Date(w2.updated_at);
                  if (d1 < d2) return 1;
                  else if (d1 > d2) return -1;
                  else return 0;
                })
                .map((workspace: Workspace) => (
                  <div key={workspace.id}>
                    {workspace.name.indexOf(searchText) >= 0 && (
                      <>
                        <div className="d-flex align-items-center justify-content-between">
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
                              <h3 className="mb-0">
                                <Link to={`/workspaces/${workspace.id}`}>{workspace.name}</Link>
                              </h3>
                              <small className="text-muted">
                                {(() => {
                                  var prettyType = "Unknown type";
                                  var workspaceType = workspaceTypes.find(
                                    (wt: WorkspaceType) =>
                                      wt.id === workspace.type
                                  );
                                  if (workspaceType) {
                                    prettyType = workspaceType.name;
                                  }
                                  return prettyType;
                                })()}
                              </small>
                            </div>
                          </div>
                          <div className="d-flex flex-column align-items-end">
                            <Badge
                              color={GetWorkspaceStatusColor(workspace.status)}
                              className="text-white mb-2"
                              style={{ fontSize: 11 }}
                            >
                              {GetBeautyNameForStatus(workspace.status)}
                            </Badge>
                            <p
                              className="mb-0 text-muted"
                              style={{ fontSize: 12 }}
                            >
                              <small>
                                Last activity{" "}
                                {new Date(
                                  workspace.updated_at
                                ).toLocaleString()}
                              </small>
                            </p>
                          </div>
                        </div>
                        <hr className="my-3" />
                      </>
                    )}
                  </div>
                ))}
            </>
          ) : (
            <p className="text-center">
              {workspaces.length > 0 ? (
                <>
                  No workspaces found matching '{searchText}''{" "}
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
