import { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  CardBody,
  Col,
  Container,
  Input,
  Label,
  Row,
  Spinner,
} from "reactstrap";
import { RunnerAdmin } from "../../types/runner";
import { CreateRunnerModal } from "./CreateRunnerModal";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy, faPencilAlt, faPlus, faXmark } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";
import { Link } from "react-router-dom";
import { AdminListRunners, AdminRetrieveRecommendedRunnerVersion } from "../../api/runner";
import DataTable from "../../components/DataTable";
import React from "react";
import { AdminRunnerStatusBadge } from "../../components/AdminRunnerStatusBadge";
import { TimeSince } from "../../common/time";
import { AdminPageHeader } from "../../components/AdminPageHeader";

export default function AdminRunners() {
  const [runners, setRunners] = useState<RunnerAdmin[]>([]);
  const [showCreateRunnerModal, setCreateRunnerModal] =
    useState<boolean>(false);
  const [runnerToken, setRunnerToken] = useState<string>("");
  const [runnerId, setRunnerID] = useState<string>("");
  const [recommendedRunnerVersion, setRecommendedRunnerVersion] = useState<string>("");

  const FetchRunners = useCallback(async () => {
    const r = await AdminListRunners();
    if (r) {
      setRunners(r);
    }
  }, []);

  const FetchRecommendedRunnerVersion = useCallback(async () => {
    const r = await AdminRetrieveRecommendedRunnerVersion();
    if (r) {
      setRecommendedRunnerVersion(r);
    }
  }, []);

  useEffect(() => {
    FetchRunners();
    FetchRecommendedRunnerVersion();
  }, [FetchRunners, FetchRecommendedRunnerVersion]);

  return (
    <Container>
      <AdminPageHeader
        title="Runners"
        description="Manage your runner instances"
        actions={
          <Button color="light" onClick={() => setCreateRunnerModal(true)}>
            <FontAwesomeIcon icon={faPlus} className="me-2" />
            Add runner
          </Button>
        }
      />
      {runnerToken.length > 0 && runnerId.length > 0 && (
        <Row>
          <Col md={12}>
            <Card style={{ borderRadius: 5 }}>
              <CardBody style={{ borderRadius: 5, background: "rgba(var(--tblr-primary-rgb), 0.5)" }}>
                <div className="d-flex justify-content-between">
                  <h3>New runner has been created</h3>
                  <Button
                    className="bg-transparent p-0 m-0"
                    style={{ height: 25 }}
                    onClick={() => {
                      setRunnerToken("");
                      setRunnerID("");
                    }}
                  >
                    <FontAwesomeIcon icon={faXmark} />
                  </Button>
                </div>
                <p>
                  Use the following ID and token to register the runner, the token will not
                  longer be visible
                </p>
                <div className="d-flex align-items-center">
                  <Label className="mt-2 me-6">
                    ID
                  </Label>
                  <Input
                    value={runnerId}
                    style={{ background: "transparent" }}
                    className="text-white"
                    disabled
                  />
                  <Button
                    className="bg-transparent"
                    onClick={() => {
                      navigator.clipboard.writeText(runnerId);
                      toast.info("Copied to clipboard");
                    }}
                  >
                    <FontAwesomeIcon icon={faCopy} />
                  </Button>
                </div>
                <div className="d-flex align-items-center mt-2">
                  <Label className="mt-2 me-3">
                    Token
                  </Label>
                  <Input
                    value={runnerToken}
                    style={{ background: "transparent" }}
                    className="text-white"
                    disabled
                  />
                  <Button
                    className="bg-transparent"
                    onClick={() => {
                      navigator.clipboard.writeText(runnerToken);
                      toast.info("Copied to clipboard");
                    }}
                  >
                    <FontAwesomeIcon icon={faCopy} />
                  </Button>
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>
      )}
      <Row className="mt-4">
        <Col md={12}>
          <Card body>
            <DataTable
              columns={[
                {
                  label: "Name",
                  key: "name",
                  render: (_, runner: RunnerAdmin) => (
                    <Link to={`/admin/runners/${runner.id}`} className="d-flex gap-2 align-items-center">
                      <div className="d-flex align-items-center gap-3">
                        <div
                          className="d-flex align-items-center justify-content-center text-white text-primary fw-bold"
                          style={{ width: 40, height: 40, borderRadius: 5, background: "rgba(var(--tblr-primary-rgb), 0.2)" }}
                        >
                          #{runner.id}
                        </div>
                        <div className="d-flex flex-column text-white">
                          <h4 className="mb-1">{runner.name}</h4>
                          {runner.deletion_in_progress ? (
                            <React.Fragment>
                              <Badge
                                className="text-warning bg-warning-lt"
                              >
                                Deletion in progress
                                <Spinner size="sm" className="ms-2" />
                              </Badge>
                            </React.Fragment>
                          ) : (
                            <small className="text-muted">
                              {runner.description || "No description provided"}
                            </small>
                          )}
                        </div>
                      </div>
                    </Link>
                  ),
                  sortable: true,
                },
                {
                  label: "Type",
                  key: "type",
                  sortable: true,
                  render: (_, runner: RunnerAdmin) => (
                    <React.Fragment>
                      <div className="d-flex align-items-center">
                        <Badge
                          className="text-primary bg-azure-lt"
                        >
                          {runner.type.at(0)?.toUpperCase() + runner.type.slice(1)}
                        </Badge>
                      </div>
                    </React.Fragment>
                  ),
                },
                {
                  label: "Status",
                  key: "status",
                  sortable: true,
                  render: (_, runner: RunnerAdmin) => (
                    <React.Fragment>
                      <div className="d-flex align-items-center">
                        <AdminRunnerStatusBadge runner={runner} />
                      </div>
                    </React.Fragment>
                  ),
                },
                {
                  label: "Version",
                  key: "version",
                  sortable: true,
                  render: (_, runner: RunnerAdmin) => (
                    <React.Fragment>
                      <div className="d-flex align-items-center gap-2">
                        {runner.version ? (
                          <React.Fragment>
                            <Badge
                              className="text-white"
                              style={{ background: "rgba(var(--tblr-white-rgb), 0.1) !important" }}
                            >
                              {runner.version}
                            </Badge>
                            {recommendedRunnerVersion !== runner.version && runner.version && (
                              <React.Fragment>
                                <Badge
                                  className="text-warning bg-warning-lt"
                                >
                                  !
                                </Badge>
                              </React.Fragment>
                            )}
                          </React.Fragment>) : (
                          <React.Fragment>
                            <Badge
                              className="text-white"
                              style={{ background: "rgba(var(--tblr-white-rgb), 0.1) !important" }}
                            >
                              -
                            </Badge>
                          </React.Fragment>
                        )}
                      </div>
                    </React.Fragment>
                  ),
                },
                {
                  label: "Last contact",
                  key: "last_contact",
                  sortable: true,
                  render: (_, runner: RunnerAdmin) => (

                    <React.Fragment>
                      <div className="d-flex align-items-center">
                        {runner.last_contact ? (
                          TimeSince(new Date(runner.last_contact))
                        ) : "Never"}
                      </div>
                    </React.Fragment>
                  ),
                },
                {
                  label: "Actions",
                  key: "actions",
                  sortable: false,
                  render: (_, runner: RunnerAdmin) => (
                    <React.Fragment>
                      <div className="d-flex align-items-center justify-content-end  gap-2">
                        <Link to={`/admin/runners/${runner.id}`} className="btn btn-accent px-2 text-muted">
                          <FontAwesomeIcon icon={faPencilAlt} />
                        </Link>
                      </div>
                    </React.Fragment>
                  ),
                },
              ]}
              data={runners}
            />
          </Card>
        </Col>
      </Row>
      <CreateRunnerModal
        isOpen={showCreateRunnerModal}
        onClose={(id, token) => {
          setCreateRunnerModal(false);
          FetchRunners();
          if (token && id) {
            setRunnerToken(token);
            setRunnerID(id);
          } else {
            setRunnerToken("");
            setRunnerID("");
          }
        }}
      />
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </Container>
  );
}
