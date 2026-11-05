import React, { useCallback, useEffect, useState } from "react";
import { RunnerAdmin, RunnerType } from "../../types/runner";
import { toast, ToastContainer } from "react-toastify";
import { useNavigate, useParams } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong, faFileLines } from "@fortawesome/free-solid-svg-icons";
import {
  AdminDeleteRunner,
  AdminListRunners,
  AdminRetrieveRecommendedRunnerVersion,
  AdminRetrieveRunnerById,
  AdminUpdateRunner,
  ListRunnerTypes
} from "../../api/runner";
import { ConfirmDeleteRunnerModal } from "./ConfirmDeleteRunnerModal";
import { RotateRunnerTokenModal } from "./RotateRunnerTokenModal";
import { RunnerLogsModal } from "./RunnerLogsModal";
import { Button, Card, Col, Container, Form, Row, Spinner } from "react-bootstrap";
import { AdminRunnerStatusBadge, IsRunnerOnline } from "../../components/AdminRunnerStatusBadge";
import { AdminPageHeader } from "../../components/AdminPageHeader";

export default function AdminRunnerDetails() {
  const [runner, setRunner] = useState<RunnerAdmin>();
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const [showConfirmDeleteModal, setShowConfirmDeleteModal] = useState<boolean>(false);
  const [showRotateTokenModal, setShowRotateTokenModal] = useState<boolean>(false);
  const [showRunnerLogsModal, setShowRunnerLogsModal] = useState<boolean>(false);
  const [recommendedRunnerVersion, setRecommendedRunnerVersion] = useState<string>("");

  const { id } = useParams();
  const navigate = useNavigate();

  const FetchRunnerTypes = useCallback(async () => {
    const rt = await ListRunnerTypes();
    if (rt) {
      setRunnerTypes(rt);
    }
  }, []);

  const FetchRecommendedRunnerVersion = useCallback(async () => {
    const r = await AdminRetrieveRecommendedRunnerVersion();
    if (r) {
      setRecommendedRunnerVersion(r);
    }
  }, []);

  const validation = useFormik({
    initialValues: {
      runnerName: "",
      description: "",
      runnerType: "",
      usePublicUrl: false,
      publicUrl: "",
    },
    validationSchema: Yup.object({
      runnerName: Yup.string()
        .required("Runner name is required")
        .test(
          "Another runner with the same name already exists",
          async (value) => {
            // TODO: runner by id
            const runners = await AdminListRunners();
            if (runners) {
              return (
                runners.find(
                  (r) => r.name === value && r.id !== runner?.id
                ) === undefined
              );
            }

            return false;
          }
        ),
      runnerType: Yup.string().required("Runner type is required"),
      publicUrl: Yup.string().when("usePublicUrl", {
        is: true,
        then: (schema) =>
          schema
            .required("Public url is required")
            .test(
              "Another runner with the same public url already exists",
              async (value) => {
                // TODO: runner by name
                const runners = await AdminListRunners();
                if (!runners) {
                  toast.error(`Failed to update runner, try again later`);
                  return false;
                }
                return (
                  runners.find(
                    (r) =>
                      r.public_url === value &&
                      r.use_public_url &&
                      r.id !== runner?.id
                  ) === undefined
                );
              }
            ),
      }),
    }),
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      if (id) {
        const r = await AdminUpdateRunner(
          parseInt(id),
          values.runnerName,
          values.runnerType,
          values.description,
          values.usePublicUrl,
          values.publicUrl
        );
        if (r) {
          setRunner(r);
          toast.success("Runner has been updated");
        } else {
          validation.resetForm();
          FetchRunner();
          toast.error("An error occured while updating runner details");
        }
      }
    },
  });

  const FetchRunner = useCallback(async () => {
    if (id) {
      const r = await AdminRetrieveRunnerById(parseInt(id));
      if (r) {
        setRunner(r);
      } else {
        navigate("/admin/runners");
      }
    }
  }, [id, navigate]);

  const DeleteRunner = useCallback(async () => {
    if (id) {
      if (await AdminDeleteRunner(parseInt(id))) {
        navigate("/admin/runners");
      } else {
        toast.error("Failed to delete the runner")
      }
    }
  }, [id, navigate]);

  useEffect(() => {
    FetchRunner();
    FetchRunnerTypes();
    FetchRecommendedRunnerVersion();
  }, [FetchRunner, FetchRunnerTypes, FetchRecommendedRunnerVersion]);

  useEffect(() => {
    validation.setValues({
      runnerName: runner?.name || "",
      description: runner?.description || "",
      runnerType: runner?.type || "",
      usePublicUrl: runner?.use_public_url || false,
      publicUrl: runner?.public_url || "",
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [runner]);

  return (
    <React.Fragment>
      <Container>
        <AdminPageHeader
          title="Edit runner"
          description="Manage runner configuration, connectivity and health."
          actions={
            <Button
              variant="accent"
              onClick={(e) => {
                e.preventDefault();
                navigate("/admin/runners");
              }}
            >
              <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
              Back to runners
            </Button>
          }
        />
        {runner && (
          <React.Fragment>
            {runner.version && recommendedRunnerVersion !== runner.version && (
              <React.Fragment>
                <Row className="mb-4">
                  <Col>
                    <Card
                      body
                      className="bg-warning-lt border-warning-lt"
                      style={{ borderLeft: "5px solid var(--tblr-warning)" }}
                    >
                      <div className="d-flex justify-content-center align-items-start gap-3">
                        <div>
                        </div>
                        <div>
                          <h3 className="mb-1">Runner version mismatch</h3>
                          <p className="text-white mb-1">
                            This runner is running version <b>{runner.version}</b>,
                            but the recommended version is <b>{recommendedRunnerVersion}</b>.
                            Please consider updating the runner to ensure compatibility and access to the latest features and a complete compatibility.
                          </p>
                        </div>
                      </div>
                    </Card>
                  </Col>
                </Row>
              </React.Fragment>
            )}
          </React.Fragment>
        )}
        <form
          onSubmit={(e) => {
            e.preventDefault();
            validation.handleSubmit();
            return false;
          }}
        >
          <Row>
            <Col md={12}>
              <Card>
                <Card.Header className="d-flex flex-column align-items-start">
                  <h3 className="mb-1">Runner Configuration</h3>
                  <small className="text-muted">Basic details and how Codebox connects to this runner.</small>
                </Card.Header>
                <Card.Body>
                  <h4 className="text-uppercase text-muted">
                    <b>General</b>
                  </h4>
                  <div className="mb-3">
                    <Form.Label>Runner name</Form.Label>
                    <Form.Control
                      name="runnerName"
                      value={validation.values.runnerName}
                      onChange={validation.handleChange}
                      isInvalid={!!validation.errors.runnerName}
                      disabled={runner?.deletion_in_progress}
                    />
                    <small className="text-muted">
                      A unique name to identify this runner in your infrastructure.
                    </small>
                    <Form.Control.Feedback>
                      {validation.errors.runnerName}
                    </Form.Control.Feedback>
                  </div>
                  <div className="mb-3">
                    <Form.Label>Runner type</Form.Label>
                    <Form.Select
                      name="runnerType"
                      value={validation.values.runnerType}
                      onChange={validation.handleChange}
                      isInvalid={!!validation.errors.runnerType}
                      disabled={runner?.deletion_in_progress}
                    >
                      <option value={""}>Select runner type</option>
                      {runnerTypes.map((t) => (
                        <option value={t.id} key={t.id}>
                          {t.name}
                        </option>
                      ))}
                    </Form.Select>
                    <small className="text-muted">
                      Supports {" "}
                      {runnerTypes
                        .find((r) => r.id === validation.values.runnerType)
                        ?.supported_types.map((st) => st.name)
                        .join(", ")}.
                    </small>
                    <Form.Control.Feedback type="invalid">
                      {validation.errors.runnerType}
                    </Form.Control.Feedback>
                  </div>
                  <div className="mb-3">
                    <Form.Label>Description</Form.Label>
                    <Form.Control
                      as="textarea"
                      rows={3}
                      placeholder="Describe this runner"
                      name="description"
                      value={validation.values.description}
                      onChange={validation.handleChange}
                      disabled={runner?.deletion_in_progress}
                    />
                  </div>
                  <hr className="my-4" />
                  <h4 className="text-uppercase text-muted">
                    <b>Connectivity</b>
                  </h4>
                  <div
                    className="d-flex justify-content-between align-items-center mb-3 p-3 rounded-2 bg-azure-lt"
                  >
                    <div>
                      <h3 className="mb-1">Direct Connection</h3>
                      <small className="text-muted">
                        Connect directly to this runner without using tunneling.
                      </small>
                    </div>
                    <Form.Switch
                      name="usePublicUrl"
                      onChange={validation.handleChange}
                      checked={validation.values.usePublicUrl}
                      disabled={runner?.deletion_in_progress}
                    />
                  </div>
                  {validation.values.usePublicUrl && (

                    <div className="mb-3">
                      <Form.Label>Runner URL</Form.Label>
                      <Form.Control
                        placeholder="http://my-host.example.com:12345"
                        name="publicUrl"
                        onChange={validation.handleChange}
                        value={validation.values.publicUrl}
                        isInvalid={!!validation.errors.publicUrl}
                        disabled={!validation.values.usePublicUrl || runner?.deletion_in_progress}
                      />
                      <Form.Control.Feedback>
                        {validation.errors.publicUrl}
                      </Form.Control.Feedback>
                    </div>
                  )}
                </Card.Body>
              </Card>
            </Col>
          </Row>
          <Row className="mt-4">
            <Col md={6}>
              <Card className="h-100">
                <Card.Header className="d-flex flex-column align-items-start">
                  <h3 className="mb-1">Runner health</h3>
                  <small className="text-muted">Live connection information</small>
                </Card.Header>
                <Card.Body>
                  <div>
                    <div className="d-flex justify-content-between align-items-center">
                      <p className="mb-0 fw-bold text-muted">Connection</p>
                      <div>
                        {runner && <AdminRunnerStatusBadge runner={runner} />}
                      </div>
                    </div>
                  </div>
                  <div className="mt-3">
                    <div className="d-flex justify-content-between align-items-center">
                      <p className="mb-0 fw-bold text-muted">Last contact</p>
                      <div>
                        {runner?.last_contact ? new Date(runner.last_contact).toLocaleString() : "Never"}
                      </div>
                    </div>
                  </div>
                  <div className="mt-3">
                    <div className="d-flex justify-content-between align-items-center">
                      <p className="mb-0 fw-bold text-muted">Runner ID</p>
                      <div>
                        {runner?.id || "Unknown"}
                      </div>
                    </div>
                  </div>
                  <div className="mt-3">
                    <div className="d-flex justify-content-between align-items-center">
                      <p className="mb-0 fw-bold text-muted">Version</p>
                      <div>
                        {runner?.version || "Unknown"}
                      </div>
                    </div>
                  </div>
                  <div
                    className="d-flex justify-content-between p-3 rounded-2 mb-3 align-items-center gap-3 bg-azure-lt mt-3"
                  >
                    <div>
                      <h4 className="mb-1">Runner logs</h4>
                      <small className="text-muted">
                        Inspect runner logs
                      </small>
                    </div>
                    {runner && (
                      <React.Fragment>
                        {IsRunnerOnline(runner) ? (
                          <Button
                            variant="accent"
                            onClick={() => setShowRunnerLogsModal(true)}
                          >
                            <FontAwesomeIcon icon={faFileLines} className="me-2" />
                            View logs
                          </Button>
                        ) : (
                          <React.Fragment>
                            <Button
                              variant="outline-light"
                              disabled
                            >
                              Runner is offline
                            </Button>
                          </React.Fragment>
                        )}
                      </React.Fragment>
                    )}
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col md={6}>
              <Card className="h-100">
                <Card.Header className="d-flex flex-column align-items-start">
                  <h3 className="mb-1">Security & maintenance</h3>
                  <small className="text-muted">Manage runner credentials and lifecycle.</small>
                </Card.Header>
                <Card.Body>
                  {runner?.deletion_in_progress ? (
                    <React.Fragment>
                      <div className="d-flex justify-content-end mt-4">
                        <div className="btn btn-orange text-white">
                          Deletion in progress
                          <Spinner size="sm" className="ms-2" />
                        </div>
                      </div>
                    </React.Fragment>
                  ) : (
                    <React.Fragment>
                      <div
                        className="d-flex justify-content-between p-3 rounded-2 mb-3 align-items-center gap-3 bg-azure-lt"
                      >
                        <div>
                          <h4 className="mb-1">Rotate runner token</h4>
                          <small className="text-muted">
                            Rotate the runner token to invalidate the current token and generate a new one.
                          </small>
                        </div>
                        <div>
                          <Button
                            variant="accent"
                            className="me-2"
                            onClick={() => setShowRotateTokenModal(true)}
                          >
                            Rotate
                          </Button>
                        </div>
                      </div>
                      <div
                        className="d-flex justify-content-between p-3 rounded-2 mb-3 align-items-center gap-3 bg-danger-lt"
                      >
                        <div>
                          <h4 className="mb-1 ">Remove runner</h4>
                          <small className="text-muted">
                            Disconnect and delete this runner.
                          </small>
                        </div>
                        <div>
                          <Button
                            variant="outline-danger"
                            className="me-2"
                            onClick={(e) => {
                              e.preventDefault();
                              setShowConfirmDeleteModal(true);
                            }}
                          >
                            Delete
                          </Button>
                        </div>
                      </div>
                    </React.Fragment>
                  )}
                </Card.Body>
              </Card>
            </Col>
          </Row>
          <Row className="mt-4" style={{ position: "sticky", bottom: 0, zIndex: 100 }}>
            <Col>
              <Card body className="d-flex justify-content-end align-items-end gap-2">
                <Button
                  variant="accent"
                  className="me-2"
                  onClick={(e) => {
                    e.preventDefault();
                    FetchRunner();
                  }}
                >
                  Cancel
                </Button>
                <Button type="submit" variant="light">
                  Save
                </Button>
              </Card>
            </Col>
          </Row>
        </form>
        <ConfirmDeleteRunnerModal
          isOpen={showConfirmDeleteModal}
          onClose={(deleteConfirmed) => {
            setShowConfirmDeleteModal(false);
            if (deleteConfirmed) {
              DeleteRunner();
            }
          }}
        />
        {runner && (
          <RotateRunnerTokenModal
            isOpen={showRotateTokenModal}
            runnerId={runner.id}
            onClose={() => setShowRotateTokenModal(false)}
          />
        )}
        {runner && (
          <RunnerLogsModal
            isOpen={showRunnerLogsModal}
            runnerId={runner.id}
            runnerName={runner.name}
            onClose={() => setShowRunnerLogsModal(false)}
          />
        )}
        <ToastContainer
          toastClassName={"bg-dark"}
        />
      </Container>
    </React.Fragment>
  );
}
