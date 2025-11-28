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
import { faCopy, faXmark } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";
import { Link } from "react-router-dom";
import { AdminListRunners } from "../../api/runner";
import DataTable from "../../components/DataTable";
import React from "react";
import { useSelector } from "react-redux";
import { RootState } from "../../redux/store";

export function AdminRunners() {
  const [runners, setRunners] = useState<RunnerAdmin[]>([]);
  const [showCreateRunnerModal, setCreateRunnerModal] =
    useState<boolean>(false);
  const [runnerToken, setRunnerToken] = useState<string>("");
  const [runnerId, setRunnerID] = useState<string>("");
  const settings = useSelector((state: RootState) => state.settings);

  const FetchRunners = useCallback(async () => {
    const r = await AdminListRunners();
    if (r) {
      setRunners(r);
    }
  }, []);
  useEffect(() => {
    FetchRunners();
  }, [FetchRunners]);

  return (
    <Container>
      <div className="row g-2 align-items-center mb-4">
        <div className="col">
          <h2 className="mb-0 mt-2">Runners</h2>
          <p className="text-muted">List of all runners</p>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <Button color="primary" onClick={() => setCreateRunnerModal(true)}>
            Add new runner
          </Button>
        </div>
      </div>
      {runnerToken.length > 0 && runnerId.length > 0 && (
        <Row>
          <Col md={12}>
            <Card style={{ borderRadius: 5 }}>
              <CardBody className="bg-primary" style={{ borderRadius: 5 }}>
                <div className="d-flex justify-content-between">
                  <h3>Runner has been created</h3>
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
                    style={{ background: "var(--tblr-primary-darken)" }}
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
                    style={{ background: "var(--tblr-primary-darken)" }}
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
                      <b>{runner.name}</b>
                      {runner.deletion_in_progress && (
                        <React.Fragment>
                          <Badge color="orange" className="text-white">
                            Deletion in progress
                            <Spinner size="sm" />
                          </Badge>
                        </React.Fragment>
                      )}
                      {settings.recommended_runner_version !== runner.version && (
                        <React.Fragment>
                          <Badge color="warning" className="text-white">
                            Version mismatch
                          </Badge>
                        </React.Fragment>
                      )}
                    </Link>
                  ),
                },
                {
                  label: "Type",
                  key: "type",
                },
                {
                  label: "Last contact",
                  key: "last_contact",
                  render: (_, runner: RunnerAdmin) => (
                    runner.last_contact ? new Date(runner.last_contact).toLocaleString() : "Never"
                  ),
                },
                {
                  label: "Status",
                  key: "_",
                  render: (_, runner: RunnerAdmin) => (
                    <React.Fragment>
                      {new Date(runner.last_contact) > new Date(Date.now() - 5 * 60 * 1000)
                        ? (
                          <React.Fragment>
                            <span className="text-success pe-1">●</span>
                            Online
                          </React.Fragment>
                        ) : (
                          <React.Fragment>
                            <span className="text-danger pe-1">●</span>
                            Offline
                          </React.Fragment>
                        )}
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
