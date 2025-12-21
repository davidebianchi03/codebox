import { useCallback, useEffect, useState } from "react";
import { WorkspaceType } from "../../types/workspace";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { Runner } from "../../types/runner";
import { ToastContainer, toast } from "react-toastify";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import { ListRunners } from "../../api/runner";
import { APICreateWorkspace, APIListWorkspacesTypes } from "../../api/workspace";
import { APIListTemplates, APIRetrieveTemplateLatestVersion } from "../../api/templates";
import { EnvEditor } from "../../components/EnvEditor";
import { Card, Col, Container, FormGroup, Row, Form, Button } from "react-bootstrap";

export default function CreateWorkspace() {
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const [runners, setRunners] = useState<Runner[]>([]);
  const [templates, setTemplates] = useState<WorkspaceTemplate[]>([]);
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const validation = useFormik({
    initialValues: {
      workspaceName: searchParams.get("name") || "",
      workspaceType: "",
      configSource: "",
      runner: -1,
      gitRepositoryURL: "",
      gitRefName: "",
      configFilesPath: "",
      environment: "",
      template: "",
    },
    validationSchema: Yup.object({
      workspaceName: Yup.string()
        .required("Workspace name is required")
        .matches(
          /^\w+$/,
          "Workspace name can only contain letters, numbers and underscores"
        ),
      workspaceType: Yup.string().required("Workspace type is required"),
      configSource: Yup.string().required("Configuration source is required"),
      runner: Yup.number().min(0, "Runner is required"),
      gitRepositoryURL: Yup.string().when("configSource", {
        is: "git",
        then: (schema) => schema.required("Repository URL is required"),
      }),
      configFilesPath: Yup.string().when("configSource", {
        is: "git",
        then: (schema) => schema.required("Config file path is required"),
      }),
      template: Yup.number().when("configSource", {
        is: "template",
        then: (schema) => schema.required("Template is required").min(0, "Template is required"),
      }),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      var template = templates.find((t) => t.id === parseInt(values.template));
      var templateVersion: WorkspaceTemplateVersion | null = null;

      if (values.configSource === "template") {
        if (!template) {
          toast.error("Template not found");
          return;
        }

        templateVersion = await APIRetrieveTemplateLatestVersion(template.id) || null;
        if (!templateVersion) {
          toast.error("There are no versions available for this template");
          return;
        }
      }

      const w = await APICreateWorkspace(
        values.workspaceName,
        values.workspaceType,
        parseInt(values.runner.toString()),
        values.configSource,
        values.gitRepositoryURL,
        values.gitRefName,
        values.configFilesPath,
        values.environment.split("\n").filter((line) => line.trim() !== ""),
        templateVersion ? templateVersion.id : 0,
      );

      if (w) {
        navigate(`/workspaces/${w.id}`);
      } else {
        toast.error(
          `Failed to create workspace, try again later`
        );
      }
    },
  });

  const FetchWorkspaceTypes = useCallback(async () => {
    const wt = await APIListWorkspacesTypes();
    if (wt) {
      setWorkspaceTypes(wt);
    }
  }, []);

  const FetchRunners = useCallback(async () => {
    const r = await ListRunners();
    if (r) {
      setRunners(r);
    }
  }, []);

  const FetchTemplates = useCallback(async () => {
    const t = await APIListTemplates();
    if (t) {
      setTemplates(t);
    }
  }, []);

  useEffect(() => {
    FetchWorkspaceTypes();
    FetchRunners();
    FetchTemplates();
  }, [FetchWorkspaceTypes, FetchRunners, FetchTemplates]);

  return (
    <>
      <Container className="pb-4 mt-4">
        <div className="row g-2 align-items-center">
          <div className="col">
            <div className="page-pretitle">Workspaces</div>
            <h2 className="page-title">Create new workspace</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <div className="btn-list"></div>
          </div>
        </div>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            validation.handleSubmit();
            return false;
          }}
        >
          <Row className="mt-3">
            <Col md={6}>
              <Card>
                <Card.Body>
                  <FormGroup>
                    <Form.Label>Workspace name</Form.Label>
                    <Form.Control
                      type="text"
                      placeholder="my awesome workspace"
                      name="workspaceName"
                      onChange={validation.handleChange}
                      value={validation.values.workspaceName}
                      isInvalid={validation.errors.workspaceName !== undefined}
                    />
                    <Form.Control.Feedback>{validation.errors.workspaceName}</Form.Control.Feedback>
                  </FormGroup>
                  <FormGroup>
                    <Form.Label>Workspace type</Form.Label>
                    <select
                      className={`form-control ${validation.errors.workspaceType ? "is-invalid" : ""}`}
                      name="workspaceType"
                      onChange={(e) => {
                        var workspaceType = workspaceTypes.find(
                          (t) => t.id === e.target.value
                        );
                        if (workspaceType) {
                          validation.setFieldValue(
                            "configFilesPath",
                            workspaceType.config_files_default_path
                          );
                        }

                        if (e.target.value === "") {
                          validation.setFieldValue("configSource", "");
                          validation.setFieldValue("runner", "");
                          validation.setFieldValue("gitRepositoryURL", "");
                          validation.setFieldValue("gitRefName", "");
                          validation.setFieldValue("configFilesPath", "");
                        } else {
                          validation.setFieldValue("configSource", workspaceType?.supported_config_sources[0] || "");
                        }

                        validation.handleChange(e);
                      }}
                      value={validation.values.workspaceType}
                    >
                      <option value="">Select workspace type</option>
                      {workspaceTypes.map((wt) => (
                        <option value={wt.id} key={wt.id}>
                          {wt.name}
                        </option>
                      ))}
                    </select>
                    <Form.Control.Feedback>{validation.errors.workspaceType}</Form.Control.Feedback>
                  </FormGroup>
                  <FormGroup>
                    <Form.Label>Runner</Form.Label>
                    <select
                      className={`form-control ${validation.errors.runner ? "is-invalid" : ""
                        }`}
                      name="runner"
                      disabled={validation.values.workspaceType === ""}
                      onChange={validation.handleChange}
                      value={validation.values.runner}
                    >
                      {validation.values.workspaceType === "" && (
                        <option value="-1">Select a workspace type before</option>
                      )}
                      {validation.values.workspaceType !== "" &&
                        runners.length > 0 && (
                          <>
                            <option value="-1">Select a runner</option>
                            {runners.map((r) => (
                              <option value={r.id} key={r.id}>
                                {r.name}
                              </option>
                            ))}
                          </>
                        )}
                      {validation.values.workspaceType !== "" &&
                        runners.length === 0 && (
                          <option value="-1">No runner available</option>
                        )}
                    </select>
                    <span className="text-warning mt-1">
                      {(() => {
                        var selectedRunner = runners.find(
                          (r) =>
                            r.id === parseInt(validation.values.runner.toString())
                        );

                        if (!selectedRunner) {
                          return "";
                        }

                        if (
                          new Date().getTime() -
                          new Date(selectedRunner.last_contact).getTime() >
                          5 * 60 * 1000
                        ) {
                          return "Warning: last contact with this runner was more than 5 minutes ago.";
                        }

                        return "";
                      })()}
                    </span>
                    <Form.Control.Feedback>{validation.errors.runner}</Form.Control.Feedback>
                  </FormGroup>
                </Card.Body>
              </Card>
            </Col>
            <Col md={6}>
              <Card>
                <Card.Body>
                  <FormGroup>
                    <Form.Label>Config source</Form.Label>
                    <select
                      className={`form-control`}
                      name="configSource"
                      disabled={validation.values.workspaceType === ""}
                      value={validation.values.configSource}
                      onChange={(e) => {
                        if (e.target.value !== "git") {
                          validation.setFieldValue("gitRepositoryURL", "");
                          validation.setFieldValue("gitRefName", "");
                          validation.setFieldValue("configFilesPath", "");
                        }
                        validation.handleChange(e);
                      }}
                      defaultValue={""}
                    >
                      {validation.values.workspaceType === "" ? (
                        <option value={""}>Select a workspace type before</option>
                      ) : (
                        <option value={""}>Select config source</option>
                      )}
                      {workspaceTypes.find(
                        (t) => t.id === validation.values.workspaceType
                      ) !== undefined ? (
                        <>
                          {validation.values.workspaceType !== "" &&
                            workspaceTypes
                              .find((t) => t.id === validation.values.workspaceType)
                              ?.supported_config_sources.map((s, index) => (
                                <option value={s} key={index}>
                                  {s[0].toUpperCase() + s.substring(1)}
                                </option>
                              ))}
                        </>
                      ) : null}
                    </select>
                  </FormGroup>
                  {validation.values.configSource !== "" &&
                    (validation.values.configSource === "git" ? (
                      <>
                        <FormGroup>
                          <Form.Label>Repository URL</Form.Label>
                          <Form.Control
                            name="gitRepositoryURL"
                            placeholder="git@example.com/my-awesome-project"
                            value={validation.values.gitRepositoryURL}
                            onChange={validation.handleChange}
                            isInvalid={
                              validation.errors.gitRepositoryURL !== undefined
                            }
                          />
                          <Form.Control.Feedback>
                            {validation.errors.gitRepositoryURL}
                          </Form.Control.Feedback>
                        </FormGroup>
                        <FormGroup>
                          <Form.Label>Ref Name</Form.Label>
                          <Form.Control
                            name="gitRefName"
                            placeholder="refs/heads/main"
                            value={validation.values.gitRefName}
                            onChange={validation.handleChange}
                            isInvalid={validation.errors.gitRefName !== undefined}
                          />
                          <Form.Control.Feedback>
                            {validation.errors.gitRefName}
                          </Form.Control.Feedback>
                        </FormGroup>
                        <FormGroup>
                          <Form.Label>Config files path</Form.Label>
                          <Form.Control
                            name="configFilesPath"
                            placeholder={
                              workspaceTypes.find(
                                (t) => t.id === validation.values.workspaceType
                              )?.config_files_default_path
                            }
                            value={validation.values.configFilesPath}
                            onChange={validation.handleChange}
                            isInvalid={
                              validation.errors.configFilesPath !== undefined
                            }
                          />
                          <Form.Control.Feedback>
                            {validation.errors.configFilesPath}
                          </Form.Control.Feedback>
                        </FormGroup>
                      </>
                    ) : (
                      <>
                        <FormGroup>
                          <Form.Label>Template</Form.Label>
                          <Form.Control
                            name="template"
                            type="select"
                            value={validation.values.template}
                            onChange={validation.handleChange}
                            isInvalid={
                              validation.errors.template !== undefined
                            }
                          >
                            <option value={""}>Select a template</option>
                            {
                              templates.filter(
                                (template) => template.type === validation.values.workspaceType
                              ).map((template, index) => (
                                <option key={index} value={template.id}>{template.name}</option>
                              ))
                            }
                          </Form.Control>
                          <Form.Control.Feedback>
                            {validation.errors.template}
                          </Form.Control.Feedback>
                        </FormGroup>
                      </>
                    ))}
                </Card.Body>
              </Card>
            </Col>
          </Row>
          <Row className="mt-3">
            <Col md={12}>
              <Card>
                <Card.Body>
                  <FormGroup>
                    <p>
                      <Form.Label className="mb-0">Environment</Form.Label>
                      <small className="text-muted">
                        Define environment variables
                      </small>
                    </p>
                    <EnvEditor
                      value={validation.values.environment}
                      onChange={(value) => {
                        validation.setFieldValue("environment", value)
                      }}
                    />
                    <Form.Control.Feedback>{validation.errors.environment}</Form.Control.Feedback>
                  </FormGroup>
                </Card.Body>
              </Card>
            </Col>
          </Row>
          <div className="d-flex justify-content-end mt-4">
            <Link to={"/"}
              className="btn btn-accent me-2">
              Cancel
            </Link>
            <Button variant="light" type="submit">
              Create workspace
            </Button>
          </div>
        </form>
        <ToastContainer
          toastClassName={"bg-dark"}
        />
      </Container>
    </>
  );
}
