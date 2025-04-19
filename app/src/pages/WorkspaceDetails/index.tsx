import React, { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { Workspace } from "../../types/workspace";
import { Button, Col, Container, Row } from "reactstrap";
import WorkspaceLogs from "./WorkspaceLogs";
import WorkspaceContainers from "./WorkspaceContainers";
import Swal from "sweetalert2";
import {
  GetBeautyNameForStatus,
  GetWorkspaceStatusColor,
} from "../../common/workspace";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp, faGear } from "@fortawesome/free-solid-svg-icons";
import { WorkspaceSettingsModal } from "./WorkspaceSettingsModal";

export default function WorkspaceDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [workspace, setWorkspace] = useState<Workspace>();
  const [fetchInterval, setFetchInterval] = useState(10000);
  const [showSettingsModal, setShowSettingsModal] = useState<boolean>(false);

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

  const HandleStartWorkspace = async () => {
    var [status, statusCode] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${id}/start`,
      "POST",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      FetchWorkspace();
    } else {
      toast.error(`Failed to start workspace, received status ${statusCode}`);
    }
  };

  const HandleStopWorkspace = async () => {
    var [status, statusCode] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace/${id}/stop`,
      "POST",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      FetchWorkspace();
    } else {
      toast.error(`Failed to stop workspace, received status ${statusCode}`);
    }
  };

  const HandleDeleteWorkspace = async (force: boolean) => {
    if (
      (
        await Swal.fire({
          title: "Delete workspace",
          icon: "warning",
          text: `
            Are you sure that you want to delete this workspace?
            ${force && (`
              Force-deleting a workspace may result in orphaned containers if runner errors, 
              including connection issues or container removal failures, are encountered
            `)}
          `,
          showCancelButton: true,
          reverseButtons: true,
          confirmButtonText: "Delete",
          customClass: {
            popup: "bg-dark text-light",
            cancelButton: "btn btn-light",
            confirmButton: "btn btn-warning",
          },
        })
      ).isConfirmed
    ) {
      var [status, statusCode] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${id}`,
        "DELETE",
        JSON.stringify({
          skip_errors: force,
        })
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        FetchWorkspace();
      } else {
        toast.error(
          `Failed to delete workspace, received status ${statusCode}`
        );
      }
    }
  };

  const HandleUpdateConfigFiles = useCallback(async () => {
    if (
      (
        await Swal.fire({
          title: "Update configuration files",
          text: `
            Updating configuration files to the latest version may cause data loss. 
            Are you sure you want to proceed?
          `,
          icon: "warning",
          showCancelButton: true,
          reverseButtons: true,
          cancelButtonText: "Cancel",
          confirmButtonText: "Update",
          customClass: {
            popup: "bg-dark text-light",
            cancelButton: "btn btn-light",
            confirmButton: "btn btn-primary",
          },
        })
      ).isConfirmed
    ) {
      var [status, statusCode] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${id}/update-config`,
        "POST",
        null
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        FetchWorkspace();
      } else {
        toast.error(
          `Failed to update workspace configuration, received status ${statusCode}`
        );
      }
    }
  }, [FetchWorkspace, id]);

  useEffect(() => {
    if (
      workspace?.status === "creating" ||
      workspace?.status === "starting" ||
      workspace?.status === "stopping" ||
      workspace?.status === "deleting"
    ) {
      setFetchInterval(800);
    } else {
      setFetchInterval(10000);
    }
  }, [workspace]);

  useEffect(() => {
    FetchWorkspace();
    const interval = setInterval(FetchWorkspace, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchWorkspace, fetchInterval]);

  return (
    <Container className="mt-4 mb-4">
      <div className="row g-2 align-items-center mb-4">
        <div className="col">
          <div className="page-pretitle">Workspace</div>
          <h2 className="page-title">{workspace?.name}</h2>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <div className="dropdown">
            {workspace?.status === "stopped" && (
              <React.Fragment>
                <Button
                  color="outline-white"
                  className="me-1"
                  onClick={HandleUpdateConfigFiles}
                >
                  <FontAwesomeIcon icon={faCloudArrowUp} />
                  <span className="ms-2">Update config files</span>
                </Button>
                <Button
                  color="outline-white"
                  className="me-3"
                  onClick={() => setShowSettingsModal(true)}
                >
                  <FontAwesomeIcon icon={faGear} />
                  <span className="ms-2">Settings</span>
                </Button>
              </React.Fragment>
            )}
            <button
              className={`btn btn-${GetWorkspaceStatusColor(
                workspace?.status
              )} dropdown-toggle`}
              type="button"
              data-bs-toggle="dropdown"
              aria-haspopup="true"
              aria-expanded="false"
            >
              {GetBeautyNameForStatus(workspace?.status)}
            </button>
            <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
              {/* <span className="dropdown-item">TODO: Update workspace</span> */}
              <span
                className="dropdown-item"
                onClick={() => {
                  if (
                    workspace?.status === "running" ||
                    workspace?.status === "error"
                  ) {
                    HandleStopWorkspace();
                  } else {
                    HandleStartWorkspace();
                  }
                }}
              >
                {workspace?.status === "running" ||
                workspace?.status === "error"
                  ? "Stop workspace"
                  : "Start workspace"}
              </span>
              <span
                className="dropdown-item"
                onClick={() => {
                  HandleDeleteWorkspace(false);
                }}
              >
                Delete workspace
              </span>
              {workspace?.status === "error" && (
                <span
                  className="dropdown-item"
                  onClick={() => {
                    HandleDeleteWorkspace(true);
                  }}
                >
                  Force delete workspace
                </span>
              )}
            </div>
          </div>
        </div>
      </div>
      {workspace && (
        <>
          <Row>
            <Col md={12}>
              <WorkspaceLogs
                workspace={workspace}
                fetchInterval={fetchInterval}
              />
            </Col>
          </Row>
          <Row className="mt-4">
            <Col md={12}>
              <WorkspaceContainers
                workspace={workspace}
                fetchInterval={fetchInterval}
              />
            </Col>
          </Row>
          <WorkspaceSettingsModal
            isOpen={showSettingsModal}
            onClose={() => {
              setShowSettingsModal(false);
              FetchWorkspace();
            }}
            workspace={workspace}
          />
        </>
      )}
    </Container>
  );
}
