import React, { useCallback, useEffect, useState } from "react";
import { Workspace, WorkspaceType } from "../../types/workspace";
import { Link, useNavigate } from "react-router-dom";
import { APIListWorkspaces, APIListWorkspacesTypes } from "../../api/workspace";
import { WorkspaceItem } from "./WorkspaceItem";
import { Button, Card, Container, Form } from "react-bootstrap";

export default function HomePage() {
  const navigate = useNavigate();

  const [searchText, setSearchText] = useState<string>("");
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);

  const FetchWorkspaces = useCallback(async () => {
    const w = await APIListWorkspaces();
    if (w) {
      setWorkspaces(w);
    }
  }, []);

  const FetchWorkspaceTypes = useCallback(async () => {
    const wt = await APIListWorkspacesTypes();
    if (wt) {
      setWorkspaceTypes(wt);
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
              variant="light"
              onClick={() => navigate("/create-workspace")}
            >
              Create workspace
            </Button>
          </div>
        </div>
      </div>
      <Card className="my-5">
        <Card.Body>
          <Form.Label>Filter workspaces:</Form.Label>
          <Form.Control
            placeholder="my awesome workspace"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
        </Card.Body>
      </Card>
      <Card className="my-1">
        <Card.Body>
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
                      <React.Fragment>
                        <WorkspaceItem workspace={workspace} workspaceTypes={workspaceTypes} />
                      </React.Fragment>
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
                    <Link to={`/create-workspace?name=${searchText}`}>create it</Link>
                  </span>
                </>
              ) : (
                <>
                  No workspaces found,{" "}
                  <span>
                    <u
                      onClick={() => navigate("/create-workspace")}
                      style={{ cursor: "pointer" }}
                    >
                      create your first workspace
                    </u>
                  </span>
                </>
              )}
            </p>
          )}
        </Card.Body>
      </Card>
    </Container>
  );
}
