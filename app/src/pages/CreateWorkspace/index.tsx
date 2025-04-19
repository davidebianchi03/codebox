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
import { Workspace, WorkspaceType } from "../../types/workspace";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Link, useNavigate } from "react-router-dom";
import { Runner } from "../../types/runner";
import { ToastContainer, toast } from "react-toastify";

export default function CreateWorkspace() {
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const [runners, setRunners] = useState<Runner[]>([]);
  const navigate = useNavigate();

  const validation = useFormik({
    initialValues: {
      workspaceName: "",
      workspaceType: "",
      configSource: "",
      runner: -1,
      gitRepositoryURL: "",
      gitRefName: "",
      configFilesPath: "",
      environment: "",
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
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      var data = {
        name: values.workspaceName,
        type: values.workspaceType,
        runner_id: parseInt(values.runner.toString()),
        config_source: values.configSource,
        git_repo_url: values.gitRepositoryURL,
        git_ref_name: values.gitRefName,
        config_source_path: values.configFilesPath,
        environment_variables: values.environment.split("\n"),
      };

      var [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace`,
        "POST",
        JSON.stringify(data),
        "application/json"
      );
      if (status === RequestStatus.OK && statusCode === 201) {
        var workspace = responseData as Workspace;
        navigate(`/workspaces/${workspace.id}`);
      } else {
        toast.error(
          `Failed to create workspace, received status ${statusCode}`
        );
      }
    },
  });

  const FetchWorkspaceTypes = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace-types`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setWorkspaceTypes(responseData as WorkspaceType[]);
    }
  }, []);

  const FetchRunners = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/runners`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setRunners(responseData as Runner[]);
    }
  }, []);

  useEffect(() => {
    FetchWorkspaceTypes();
    FetchRunners();
  }, [FetchWorkspaceTypes, FetchRunners]);

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
                      <div className="mb-3">TODO</div>
                    ))}
                </CardBody>
              </Card>
            </Col>
          </Row>
          <Row className="mt-3">
            <Card>
              <CardBody>
                <FormGroup>
                  <p>
                    <Label className="mb-0">Environment</Label>
                    <small className="text-muted">
                      Define environment variables, one per line, using the format
                      'KEY=VALUE'
                    </small>
                  </p>
                  <textarea
                    className={`form-control ${!!validation.errors.environment ? "is-invalid" : ""
                      }`}
                    rows={10}
                    placeholder="VAR1=VALUE1"
                    name="environment"
                    onChange={validation.handleChange}
                  ></textarea>
                  <FormFeedback>{validation.errors.environment}</FormFeedback>
                </FormGroup>
              </CardBody>
            </Card>
          </Row>
          <div className="d-flex justify-content-end mt-4">
            <Link to={"/"} 
              className="btn btn-outline-muted text-white me-1">
              Cancel
            </Link>
            <Button color="primary" type="submit">
              Create workspace
            </Button>
          </div>
        </form>
        <ToastContainer />
      </Container>
    </>
  );
}
