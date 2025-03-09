import { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { Workspace } from "../../types/workspace";
import { Col, Container, Row } from "reactstrap";
import WorkspaceLogs from "./WorkspaceLogs";
import WorkspaceContainers from "./WorkspaceContainers";

export default function WorkspaceDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [workspace, setWorkspace] = useState<Workspace>();

  const FetchWorkspace = useCallback(async () => {
    var [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${id}`,
      "GET",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      setWorkspace(responseData as Workspace);
    } else if (statusCode === 404) {
      navigate("/");
    } else {
      toast.error(
        `Failed to fetch workspace details, received status ${statusCode}`
      );
    }
  }, [id, navigate]);

  useEffect(() => {
    FetchWorkspace();
  }, [FetchWorkspace]);

  return (
    <Container className="mt-4 mb-4">
      <div className="row g-2 align-items-center mb-4">
        <div className="col">
          <div className="page-pretitle">Workspaces</div>
          <h2 className="page-title">{workspace?.name}</h2>
        </div>
        <div className="col-auto ms-auto d-print-none"></div>
      </div>
      {workspace && (
        <>
          <Row>
            <Col md={12}>
              <WorkspaceLogs workspace={workspace} fetchInterval={10000}/>
            </Col>
          </Row>
          <Row className="mt-4">
            <Col md={12}>
              <WorkspaceContainers workspace={workspace} fetchInterval={10000}/>
            </Col>
          </Row>
        </>
      )}
    </Container>
  );
}
