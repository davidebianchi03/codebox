import React, { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  Col,
  Container,
  FormFeedback,
  Input,
  Label,
  Row,
  Spinner,
} from "reactstrap";
import { RunnerAdmin, RunnerType } from "../../types/runner";
import { toast, ToastContainer } from "react-toastify";
import { useNavigate, useParams } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong } from "@fortawesome/free-solid-svg-icons";
import {
  AdminDeleteRunner,
  AdminListRunners,
  AdminRetrieveRecommendedRunnerVersion,
  AdminRetrieveRunnerById,
  AdminUpdateRunner,
  ListRunnerTypes
} from "../../api/runner";
import { ConfirmDeleteRunnerModal } from "./ConfirmDeleteRunnerModal";

export function AdminRunnerDetails() {
  const [runner, setRunner] = useState<RunnerAdmin>();
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const [showConfirmDeleteModal, setShowConfirmDeleteModal] = useState<boolean>(false);
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
          values.usePublicUrl,
          values.publicUrl
        );
        if (r) {
          setRunner(r);
        } else {
          validation.resetForm();
          FetchRunner();
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
      runnerType: runner?.type || "",
      usePublicUrl: runner?.use_public_url || false,
      publicUrl: runner?.public_url || "",
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [runner]);

  return (
    <React.Fragment>
      <Container>
        <Button
          color="accent"
          className="me-2 mb-4"
          onClick={(e) => {
            e.preventDefault();
            navigate("/admin/runners");
          }}
        >
          <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
          Back
        </Button>

        {runner && (
          <React.Fragment>
            {recommendedRunnerVersion !== runner.version && (
              <React.Fragment>
                <Row className="mb-4">
                  <Col>
                    <Card body className="bg-warning">
                      <div className="text-white">
                        <b>Warning:</b> This runner is running version <b>{runner.version}</b>, 
                        but the recommended version is <b>{recommendedRunnerVersion}</b>. 
                        Please consider updating the runner to ensure compatibility and access to the latest features.
                      </div>
                    </Card>
                  </Col>
                </Row>
              </React.Fragment>
            )}
          </React.Fragment>
        )}
        <Row>
          <Col md={12}>
            <Card body>
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  validation.handleSubmit();
                  return false;
                }}
              >
                <div className="mb-3">
                  <Label>Runner name</Label>
                  <Input
                    name="runnerName"
                    value={validation.values.runnerName}
                    onChange={validation.handleChange}
                    invalid={!!validation.errors.runnerName}
                    disabled={runner?.deletion_in_progress}
                  />
                  <FormFeedback>{validation.errors.runnerName}</FormFeedback>
                </div>
                <div className="mb-3">
                  <Label>Runner type</Label>
                  <select
                    name="runnerType"
                    className={`form-control ${validation.errors.runnerType ? "is-invalid" : ""
                      }`}
                    onChange={validation.handleChange}
                    value={validation.values.runnerType}
                    disabled={runner?.deletion_in_progress}
                  >
                    <option value={""}>Select runner type</option>
                    {runnerTypes.map((t) => (
                      <option value={t.id} key={t.id}>
                        {t.name}
                      </option>
                    ))}
                  </select>
                  <small className="text-muted">
                    Supported workspaces:{" "}
                    {runnerTypes
                      .find((r) => r.id === validation.values.runnerType)
                      ?.supported_types.map((st) => st.name)
                      .join(", ")}
                  </small>
                  <FormFeedback>{validation.errors.runnerType}</FormFeedback>
                </div>
                <div className="mb-3">
                  <label className="form-check">
                    <input
                      className="form-check-input"
                      type="checkbox"
                      name="usePublicUrl"
                      onClick={(e) => {
                        validation.handleChange(e);
                      }}
                      checked={validation.values.usePublicUrl}
                      disabled={runner?.deletion_in_progress}
                    />
                    <span className="form-check-label">Use public url</span>
                  </label>
                </div>
                <div className="mb-3">
                  <Label>Public url</Label>
                  <Input
                    placeholder="http://my-host.example.com:12345"
                    name="publicUrl"
                    onChange={validation.handleChange}
                    value={validation.values.publicUrl}
                    invalid={validation.errors.publicUrl ? true : false}
                    disabled={!validation.values.usePublicUrl || runner?.deletion_in_progress}
                  />
                  <FormFeedback>{validation.errors.publicUrl}</FormFeedback>
                </div>
                <div className="mb-3">
                  <Label>Status</Label>
                  {new Date(runner?.last_contact || "") > new Date(Date.now() - 5 * 60 * 1000)
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
                </div>
                <div className="mb-3">
                  <Label>Last contact</Label>
                  {runner?.last_contact ? new Date(runner.last_contact).toLocaleString() : "Never"}
                </div>
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
                    <div className="d-flex justify-content-end mt-4">
                      <Button
                        color="danger"
                        className="me-2"
                        onClick={(e) => {
                          e.preventDefault();
                          setShowConfirmDeleteModal(true);
                        }}
                      >
                        Delete
                      </Button>
                      <Button
                        color="accent"
                        className="me-2"
                        onClick={(e) => {
                          e.preventDefault();
                          FetchRunner();
                        }}
                      >
                        Cancel
                      </Button>
                      <Button type="submit" color="light">
                        Save
                      </Button>
                    </div>
                  </React.Fragment>
                )}
              </form>
            </Card>
          </Col>
        </Row>
        <ConfirmDeleteRunnerModal
          isOpen={showConfirmDeleteModal}
          onClose={(deleteConfirmed) => {
            setShowConfirmDeleteModal(false);
            if (deleteConfirmed) {
              DeleteRunner();
            }
          }}
        />
        <ToastContainer
          toastClassName={"bg-dark"}
        />
      </Container>
    </React.Fragment>
  );
}
