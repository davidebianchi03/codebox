import React, { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { toast, ToastContainer } from "react-toastify";
import { Workspace } from "../../types/workspace";
import { Button, Col, Container, Row } from "reactstrap";
import WorkspaceLogs from "./WorkspaceLogs";
import WorkspaceContainers from "./WorkspaceContainers";
import Swal from "sweetalert2";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCloudArrowUp, faGear } from "@fortawesome/free-solid-svg-icons";
import { WorkspaceSettingsModal } from "./WorkspaceSettingsModal";
import { APIDeleteWorkspace, APIRetrieveWorkspaceById, APIStartWorkspace, APIStopWorkspace, APIUpdateWorkspaceConfig } from "../../api/workspace";
import { APIRetrieveTemplateById, APIRetrieveTemplateLatestVersion } from "../../api/templates";
import { WorkspaceTemplate } from "../../types/templates";
import { WorkspaceSelectRunnerModal } from "./WorkspaceSelectRunnerModal";
import { WorkspaceStatusDropdown } from "./WorkspaceStatusDropdown";

export default function WorkspaceDetails() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [workspace, setWorkspace] = useState<Workspace>();
  const [workspaceTemplate, setWorkspaceTemplate] = useState<WorkspaceTemplate | null>(null);
  const [fetchInterval, setFetchInterval] = useState(10000);
  const [showSettingsModal, setShowSettingsModal] = useState<boolean>(false);
  const [canUpdateConfigFiles, setCanUpdateConfigFiles] = useState<boolean>(false);
  const [showSelectRunnerModal, setShowSelectRunnerModal] = useState<boolean>(false);

  const FetchWorkspace = useCallback(async () => {
    if (id) {
      const w = await APIRetrieveWorkspaceById(parseInt(id));
      if (w) {
        setWorkspace(w);
      } else if (w === null) {
        navigate("/");
      } else {
        toast.error(
          `Failed to fetch workspace details, try again later`
        );
      }
    }
  }, [id, navigate]);

  const HandleUpdateConfigFiles = useCallback(async (interactive: boolean = false) => {
    if (id) {
      const workspace = await APIRetrieveWorkspaceById(parseInt(id));
      if (workspace) {
        if (workspace.runner == null) {
          setShowSelectRunnerModal(true);
        } else {
          let actionConfirmed = true;
          if (interactive) {
            actionConfirmed = (
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
                  cancelButton: "btn btn-accent",
                  confirmButton: "btn btn-primary",
                },
              })
            ).isConfirmed
          }

          if (actionConfirmed) {
            if (await APIUpdateWorkspaceConfig(parseInt(id))) {
              FetchWorkspace();
            } else {
              toast.error(
                `Failed to update workspace configuration, try again later`
              );
            }
          }
        }
      } else {
        toast.error(
          `Failed to fetch workspace details, try again later`
        );
      }
    }
  }, [FetchWorkspace, id]);

  const CheckNewTemplateVersionAvailable = useCallback(async () => {
    if (workspace) {
      if (workspace.template_version) {
        const latestTemplateVersion = await APIRetrieveTemplateLatestVersion(
          workspace?.template_version.template_id
        );
        if (latestTemplateVersion) {
          setCanUpdateConfigFiles(latestTemplateVersion.id !== workspace?.template_version.id);
        } else {
          setCanUpdateConfigFiles(false);
        }
      }
    }
  }, [workspace]);

  const FetchWorkspaceTemplate = useCallback(async () => {
    if (workspace?.template_version) {
      if (workspace?.template_version.template_id) {
        const template = await APIRetrieveTemplateById(workspace.template_version.template_id);
        if (template) {
          setWorkspaceTemplate(template);
        } else {
          setWorkspaceTemplate(null);
        }
      }
    }
  }, [workspace]);

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

    if (workspace?.config_source === "git") {
      setCanUpdateConfigFiles(true);
    } else {
      CheckNewTemplateVersionAvailable();
    }
  }, [CheckNewTemplateVersionAvailable, workspace]);

  useEffect(() => {
    FetchWorkspace();
    const interval = setInterval(FetchWorkspace, fetchInterval);
    return () => {
      clearInterval(interval);
    };
  }, [FetchWorkspace, fetchInterval]);

  useEffect(() => {
    FetchWorkspaceTemplate();
  }, [FetchWorkspaceTemplate]);

  return (
    <Container className="mt-4 mb-4 pb-4">
      <div className="row g-2 align-items-center mb-4">
        <div className="col">
          <div className="page-pretitle">Workspace</div>
          <h2 className="page-title mb-2">{workspace?.name}</h2>
          <p className="text-muted fs-6 fw-bolder mb-0">
            {workspace?.type.replaceAll("_", " ").toUpperCase()}
          </p>
          <p className="text-muted fs-6 fw-bolder mb-0 d-flex align-items-center">
            Config loaded from {workspace?.config_source === "git" ? "git repository" : "template"}&nbsp;
            {workspace?.config_source === "git" ? (
              <small className="text-muted fs-6 fw-bolder badge bg-dark">
                {workspace?.git_source?.repository_url}
              </small>
            ) : (
              workspaceTemplate && (
                <Link to={`/templates/${workspaceTemplate.id}`} className="text-decoration-none">
                  <small className="text-muted fs-6 fw-bolder badge bg-dark">
                    {workspaceTemplate?.name}
                  </small>
                </Link>
              )
            )}
          </p>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <div className="dropdown">
            {workspace?.status === "stopped" && (
              <React.Fragment>
                {canUpdateConfigFiles && (
                  <Button
                    color="accent"
                    className="me-1"
                    onClick={() => HandleUpdateConfigFiles(true)}
                  >
                    <FontAwesomeIcon icon={faCloudArrowUp} />
                    <span className="ms-2">
                      {workspace.config_source === "template" ? "Update template version" : "Update config files"}
                    </span>
                  </Button>
                )}
                <Button
                  color="accent"
                  className="me-3"
                  onClick={() => setShowSettingsModal(true)}
                >
                  <FontAwesomeIcon icon={faGear} />
                  <span className="ms-2">Settings</span>
                </Button>
              </React.Fragment>
            )}
            {workspace && (
              <WorkspaceStatusDropdown
                workspace={workspace}
                onStatusChange={FetchWorkspace}
              />
            )}
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
          <WorkspaceSelectRunnerModal
            isOpen={showSelectRunnerModal}
            onClose={async (updated) => {
              setShowSelectRunnerModal(false);
              if (updated) {
                await FetchWorkspace();
                HandleUpdateConfigFiles(false);
              }
            }}
            workspace={workspace}
          />
        </>
      )}
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </Container>
  );
}
