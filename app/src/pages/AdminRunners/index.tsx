import { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Col,
  Container,
  Row,
  Spinner,
} from "reactstrap";
import { RunnerAdmin } from "../../types/runner";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPencilAlt, faPlus } from "@fortawesome/free-solid-svg-icons";
import { ToastContainer } from "react-toastify";
import { Link, useNavigate } from "react-router-dom";
import { AdminListRunners, AdminRetrieveRecommendedRunnerVersion } from "../../api/runner";
import DataTable from "../../components/DataTable";
import React from "react";
import { AdminRunnerStatusBadge } from "../../components/AdminRunnerStatusBadge";
import { TimeSince } from "../../common/time";
import { AdminPageHeader } from "../../components/AdminPageHeader";

export default function AdminRunners() {
  const [runners, setRunners] = useState<RunnerAdmin[]>([]);
  const [recommendedRunnerVersion, setRecommendedRunnerVersion] = useState<string>("");
  const navigate = useNavigate();

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
          <Button color="light" onClick={() => navigate("/admin/runners/create")}>
            <FontAwesomeIcon icon={faPlus} className="me-2" />
            Add runner
          </Button>
        }
      />
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
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </Container>
  );
}
