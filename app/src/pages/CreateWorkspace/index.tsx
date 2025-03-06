import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Container,
  FormFeedback,
  Input,
  Label,
} from "reactstrap";
import { WorkspaceType } from "../../types/workspace";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Link } from "react-router-dom";
import { Runner } from "../../types/runner";

export default function CreateWorkspace() {
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const [runners, setRunners] = useState<Runner[]>([]);

  const validation = useFormik({
    initialValues: {
      workspaceName: "",
      workspaceType: "",
      configSource: "",
      runner: -1,
      gitRepositoryURL: "",
      configFilesPath: "",
    },
    validationSchema: Yup.object({
      workspaceName: Yup.string().required("Workspace name is required"),
      workspaceType: Yup.string().required("Workspace type is required"),
      configSource: Yup.string().required("Configuration source is required"),
      runner: Yup.number().min(0, "Runner is required"),
      gitRepositoryURL: Yup.string().when('pin', {
        is: "git",
        then: Yup.string().required(),
        otherwise: Yup.string(),
      }),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: (values) => {
      alert(JSON.stringify(values, null, 2));
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
        <Card className="mt-4">
          <CardBody>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                validation.handleSubmit();
                return false;
              }}
            >
              <div className="mb-3">
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
              </div>
              <div className="mb-3">
                <Label>Workspace type</Label>
                <select
                  className={`form-control ${
                    validation.errors.workspaceType ? "is-invalid" : ""
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
              </div>
              <div className="mb-3">
                <Label>Config source</Label>
                <select
                  className={`form-control`}
                  name="configSource"
                  disabled={validation.values.workspaceType === ""}
                  onChange={validation.handleChange}
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
              </div>
              {validation.values.configSource !== "" &&
                (validation.values.configSource === "git" ? (
                  <>
                    <div className="mb-3">
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
                    </div>
                    <div className="mb-3">
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
                    </div>
                  </>
                ) : (
                  <div className="mb-3">TODO</div>
                ))}
              <div className="mb-3">
                <Label>Runner</Label>
                <select
                  className={`form-control ${
                    validation.errors.runner ? "is-invalid" : ""
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
                <FormFeedback>{validation.errors.runner}</FormFeedback>
              </div>
              <div className="d-flex justify-content-end mt-4">
                <Link to={"/"} className="btn btn-outline-light me-1">
                  Cancel
                </Link>
                <Button color="primary" type="submit">
                  Create workspace
                </Button>
              </div>
            </form>
          </CardBody>
        </Card>
      </Container>
    </>
  );
}
