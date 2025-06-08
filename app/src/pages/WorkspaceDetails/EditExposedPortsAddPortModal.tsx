import React, { useCallback, useEffect, useState } from "react";
import {
  ContainerPort,
  Workspace,
  WorkspaceContainer,
} from "../../types/workspace";
import {
  Button,
  FormFeedback,
  FormGroup,
  Input,
  Label,
  Modal,
  ModalBody,
  ModalHeader,
} from "reactstrap";
import { toast, ToastContainer } from "react-toastify";
import { useFormik } from "formik";
import * as Yup from "yup";
import { APICreateWorkspaceContainerPort, APIListWorkspaceContainerPorts } from "../../api/workspace";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  workspace: Workspace;
  container: WorkspaceContainer;
}

export function EditExposedPortsAddPortModal({
  isOpen,
  onClose,
  workspace,
  container,
}: Props) {
  const [containerExposedPorts, setContainerExposedPorts] = useState<
    ContainerPort[]
  >([]);

  const FetchSelectedContainerPorts = useCallback(
    async (containerName: string) => {
      const ports = await APIListWorkspaceContainerPorts(workspace.id, containerName);
      if (ports) {
        setContainerExposedPorts(ports);
      } else {
        toast.error("Failed to fetch workspace container ports");
        setContainerExposedPorts([]);
      }
    },
    [workspace]
  );

  const validation = useFormik({
    initialValues: {
      portNumber: 0,
      serviceName: "",
      public: false,
    },
    validateOnChange: false,
    validateOnBlur: false,
    validationSchema: Yup.object({
      portNumber: Yup.number()
        .required("This field is required")
        .min(1, "Min port number is 1")
        .max(65535, "Max port number is 65535")
        .test(
          "port_number",
          "Port is already exposed",
          (value) =>
            containerExposedPorts.find((port) => port.port_number === value) ===
            undefined
        ),
      serviceName: Yup.string()
        .required("This field is required")
        .test(
          "name_unique",
          "Another port with this name already exists",
          (value) =>
            containerExposedPorts.find(
              (port) => port.service_name === value
            ) === undefined
        ),
    }),
    onSubmit: async (values) => {
      const p = await APICreateWorkspaceContainerPort(
        workspace.id,
        container.container_name,
        values.portNumber,
        values.serviceName,
        values.public
      )
      if (p) {
        handleClose();
      } else {
        toast.error("Failed to add port");
      }
    },
  });

  const handleClose = useCallback(() => {
    validation.resetForm();
    onClose();
  }, [onClose, validation]);

  useEffect(() => {
    if (isOpen) {
      FetchSelectedContainerPorts(container.container_name);
    }
  }, [isOpen, FetchSelectedContainerPorts, container]);

  return (
    <React.Fragment>
      <Modal
        isOpen={isOpen}
        toggle={handleClose}
        centered
        size="lg"
        modalClassName="modal-blur"
      >
        <ModalHeader toggle={handleClose}>Edit exposed ports</ModalHeader>
        <ModalBody>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              validation.handleSubmit();
            }}
          >
            <FormGroup>
              <Label>Port number</Label>
              <Input
                name="portNumber"
                type="number"
                min={1}
                max={65535}
                value={validation.values.portNumber > 0 ? validation.values.portNumber : ""}
                onChange={validation.handleChange}
                invalid={!!validation.errors.portNumber}
                placeholder="Port number"
              />
              <FormFeedback>{validation.errors.portNumber}</FormFeedback>
            </FormGroup>
            <FormGroup>
              <Label>Service name</Label>
              <Input
                name="serviceName"
                type="text"
                value={validation.values.serviceName}
                onChange={validation.handleChange}
                invalid={!!validation.errors.serviceName}
                placeholder="My awesome service"
              />
              <FormFeedback>{validation.errors.serviceName}</FormFeedback>
            </FormGroup>
            <FormGroup className="d-flex">
              <Input
                name="public"
                type="checkbox"
                id="public"
                checked={validation.values.public}
                onChange={validation.handleChange}
              />
              <Label className="ms-1" for="public">Is Public</Label>
            </FormGroup>
            <div className="d-flex align-items-center justify-content-end">
              <Button
                color="accent"
                onClick={(e) => {
                  e.preventDefault();
                  handleClose();
                }}
              >
                Close
              </Button>
              <Button color="primary" className="ms-1" type="submit">
                Add port
              </Button>
            </div>
          </form>
        </ModalBody>
      </Modal>
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </React.Fragment>
  );
}
