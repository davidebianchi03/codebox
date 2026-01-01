import React, { useEffect } from "react";
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
import { toast, ToastContainer } from "react-toastify";
import { APIUpdateWorkspace } from "../../api/workspace";
import { EnvEditor } from "../../components/EnvEditor";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  workspace: Workspace;
}

export function WorkspaceSettingsModal({ isOpen, onClose, workspace }: Props) {
  const handleCloseModal = () => {
    validation.resetForm();
    onClose();
  };

  const validation = useFormik({
    initialValues: {
      gitRepoUrl: "",
      gitRefName: "",
      configSourcePath: "",
      environmentVariables: "",
    },
    validationSchema: Yup.object({
      gitRepoUrl: Yup.string().required("This field is required"),
      configSourcePath: Yup.string().required("This field is required"),
    }),
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      const r = await APIUpdateWorkspace(
        workspace.id,
        values.gitRepoUrl,
        values.gitRefName,
        values.configSourcePath,
        values.environmentVariables !== "" ? values.environmentVariables.split("\n") : [],
      );

      if (r) {
        toast.success("Workspace has been updated");
        handleCloseModal();
      } else {
        toast.error("Failed to update workspace");
      }
    },
  });

  useEffect(() => {
    validation.setValues({
      gitRepoUrl: workspace.git_source?.repository_url || "",
      gitRefName: workspace.git_source?.ref_name || "",
      configSourcePath: workspace.git_source?.config_file_relative_path || "",
      environmentVariables: workspace.environment_variables.join("\n") || "",
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen]);

  return (
    <React.Fragment>
      <Modal
        centered
        isOpen={isOpen}
        toggle={handleCloseModal}
        modalClassName="modal-blur"
        size="xl"
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
              <EnvEditor
                value={validation.values.environmentVariables}
                onChange={(value) => {
                  validation.setFieldValue("environmentVariables", value);
                }}
                invalid={!!validation.errors.environmentVariables}
              />
              <FormFeedback>
                {validation.errors.environmentVariables}
              </FormFeedback>
            </div>
            <small className="text-muted">
              Workspace configuration won't be automatically updated, click on
              'Update config files to update it'
            </small>
            <div className="d-flex align-items-center justify-content-end">
              <Button
                color="accent"
                onClick={(e) => {
                  e.preventDefault();
                  handleCloseModal();
                }}
              >
                Close
              </Button>
              <Button color="light" className="ms-2" type="submit">
                Save
              </Button>
            </div>
          </form>
        </ModalBody>
        <ToastContainer
          toastClassName={"bg-dark"}
        />
      </Modal>
    </React.Fragment>
  );
}
