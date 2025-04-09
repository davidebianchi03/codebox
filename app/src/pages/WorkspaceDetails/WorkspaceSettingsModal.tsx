import React from "react";
import {
  Button,
  FormFeedback,
  Input,
  Label,
  Modal,
  ModalBody,
  ModalHeader,
} from "reactstrap";
import { Workspace } from "../../types/workspace";
import { useFormik } from "formik";
import * as Yup from "yup";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  workspace: Workspace;
}

export function WorkspaceSettingsModal({ isOpen, onClose, workspace }: Props) {
  const handleCloseModal = () => {
    onClose();
  };

  const validation = useFormik({
    initialValues: {
      gitRepoUrl: workspace.git_source?.repository_url,
      gitRefName: workspace.git_source?.ref_name,
      configSourcePath: workspace.git_source?.config_file_relative_path,
      environemntVariables: workspace.environment_variables,
    },
    validationSchema: Yup.object({
      gitRepoUrl: Yup.string().required("This field is required"),
      configSourcePath: Yup.string().required("This field is required"),
    }),
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      console.log(values);
    },
  });

  /* GitRepoUrl           *string `json:"git_repo_url"`
      GitRefName           *string `json:"git_ref_name"`
      ConfigSourcePath     *string `json:"config_source_path"`
      EnvironmentVariables *string `json:"environment_variables"`
      UpdateConfig   */

  return (
    <React.Fragment>
      <Modal
        centered
        isOpen={isOpen}
        toggle={handleCloseModal}
        modalClassName="modal-blur"
      >
        <ModalHeader toggle={handleCloseModal}>Settings</ModalHeader>
        <ModalBody>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              validation.handleSubmit();
              return false;
            }}
          >
            {workspace.config_source === "git" && (
              <React.Fragment>
                <div className="mb-3">
                  <Label>Git repository url</Label>
                  <Input
                    name="gitRepoUrl"
                    placeholder="git@example.com"
                    onChange={validation.handleChange}
                    value={validation.values.gitRepoUrl}
                    invalid={!!validation.errors.gitRepoUrl}
                  />
                  <FormFeedback>{validation.errors.gitRepoUrl}</FormFeedback>
                </div>
                <div className="mb-3">
                  <Label>Git ref name</Label>
                  <Input
                    name="gitRefName"
                    placeholder="refs/heads/main"
                    onChange={validation.handleChange}
                    value={validation.values.gitRefName}
                    invalid={!!validation.errors.gitRefName}
                  />
                  <FormFeedback>{validation.errors.gitRefName}</FormFeedback>
                </div>
                <div className="mb-3">
                  <Label>Config source path</Label>
                  <Input
                    name="configSourcePath"
                    placeholder="docker-compose.yml"
                    onChange={validation.handleChange}
                    value={validation.values.configSourcePath}
                    invalid={!!validation.errors.configSourcePath}
                  />
                  <FormFeedback>
                    {validation.errors.configSourcePath}
                  </FormFeedback>
                </div>
              </React.Fragment>
            )}
            <div className="mb-3">
              <Label>Environment variables</Label>
              <Input
                name="environemntVariables"
                type="textarea"
                placeholder="Define environment variables, one per line, using the format 'KEY=VALUE'"
                onChange={validation.handleChange}
                value={validation.values.environemntVariables}
                invalid={!!validation.errors.environemntVariables}
              />
              <FormFeedback>
                {validation.errors.environemntVariables}
              </FormFeedback>
            </div>
            <small className="text-muted">
              Workspace configuration won't be automatically updated, click on
              'Update config files to update it'
            </small>
            <div className="d-flex align-items-center justify-content-end">
              <Button
                color="outline-secondary"
                onClick={(e) => {
                  e.preventDefault();
                  handleCloseModal();
                }}
              >
                Close
              </Button>
              <Button color="primary" className="ms-2" type="submit">
                Save
              </Button>
            </div>
          </form>
        </ModalBody>
      </Modal>
    </React.Fragment>
  );
}
