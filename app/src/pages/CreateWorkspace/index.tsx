import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Col,
  Container,
  FormFeedback,
  FormGroup,
  Input,
  Label,
  Row,
} from "reactstrap";
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
                <CardBody>
                  <FormGroup>
                    <Label>Workspace name</Label>
                    <Input
                      type="text"
                      placeholder="my awesome workspace"
                      name="workspaceName"
                      onChange={validation.handleChange}
                      value={validation.values.workspaceName}
                      invalid={validation.errors.workspaceName !== undefined}
                    />
                    <FormFeedback>{validation.errors.workspaceName}</FormFeedback>
                  </FormGroup>
                  <FormGroup>
                    <Label>Workspace type</Label>
                    <select
                      className={`form-control ${validation.errors.workspaceType ? "is-invalid" : ""
                        }`}
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
                    <FormFeedback>{validation.errors.workspaceType}</FormFeedback>
                  </FormGroup>
                  <FormGroup>
                    <Label>Runner</Label>
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
                    <FormFeedback>{validation.errors.runner}</FormFeedback>
                  </FormGroup>
                </CardBody>
              </Card>
            </Col>
            <Col md={6}>
              <Card>
                <CardBody>
                  <FormGroup>
                    <Label>Config source</Label>
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
                          <Label>Repository URL</Label>
                          <Input
                            name="gitRepositoryURL"
                            placeholder="git@example.com/my-awesome-project"
                            value={validation.values.gitRepositoryURL}
                            onChange={validation.handleChange}
                            invalid={
                              validation.errors.gitRepositoryURL !== undefined
                            }
                          />
                          <FormFeedback>
                            {validation.errors.gitRepositoryURL}
                          </FormFeedback>
                        </FormGroup>
                        <FormGroup>
                          <Label>Ref Name</Label>
                          <Input
                            name="gitRefName"
                            placeholder="refs/heads/main"
                            value={validation.values.gitRefName}
                            onChange={validation.handleChange}
                            invalid={validation.errors.gitRefName !== undefined}
                          />
                          <FormFeedback>
                            {validation.errors.gitRefName}
                          </FormFeedback>
                        </FormGroup>
                        <FormGroup>
                          <Label>Config files path</Label>
                          <Input
                            name="configFilesPath"
                            placeholder={
                              workspaceTypes.find(
                                (t) => t.id === validation.values.workspaceType
                              )?.config_files_default_path
                            }
                            value={validation.values.configFilesPath}
                            onChange={validation.handleChange}
                            invalid={
                              validation.errors.configFilesPath !== undefined
                            }
                          />
                          <FormFeedback>
                            {validation.errors.configFilesPath}
                          </FormFeedback>
                        </FormGroup>
                      </>
                    ) : (
                      <>
                        <FormGroup>
                          <Label>Template</Label>
                          <Input
                            name="template"
                            type="select"
                            value={validation.values.template}
                            onChange={validation.handleChange}
                            invalid={
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
                          </Input>
                          <FormFeedback>
                            {validation.errors.template}
                          </FormFeedback>
                        </FormGroup>
                      </>
                    ))}
                </CardBody>
              </Card>
            </Col>
          </Row>
          <Row className="mt-3">
            <Col md={12}>
              <Card>
                <CardBody>
                  <FormGroup>
                    <p>
                      <Label className="mb-0">Environment</Label>
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
                    <FormFeedback>{validation.errors.environment}</FormFeedback>
                  </FormGroup>
                </CardBody>
              </Card>
            </Col>
          </Row>
          <div className="d-flex justify-content-end mt-4">
            <Link to={"/"}
              className="btn btn-accent me-2">
              Cancel
            </Link>
            <Button color="primary" type="submit">
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
