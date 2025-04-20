import React, { useCallback, useEffect, useState } from "react";
import {
  ContainerPort,
  Workspace,
  WorkspaceContainer,
} from "../../types/workspace";
import {
  Button,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  Table,
} from "reactstrap";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast, ToastContainer } from "react-toastify";
import { EditExposedPortsAddPortModal } from "./EditExposedPortsAddPortModal";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  onChange: () => void;
  workspace: Workspace;
  container: WorkspaceContainer;
}

export function EditExposedPortsModal({
  isOpen,
  onClose,
  onChange,
  workspace,
  container,
}: Props) {
  const handleClose = useCallback(() => {
    onClose();
  }, [onClose]);

  const [containerExposedPorts, setContainerExposedPorts] = useState<
    ContainerPort[]
  >([]);
  const [showAddPortModal, setShowAddPortModal] = useState(false);

  const FetchSelectedContainerPorts = useCallback(
    async (containerName: string) => {
      // fetch exposed ports
      var [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${workspace.id
        }/container/${containerName}/port`,
        "GET",
        null
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        setContainerExposedPorts(responseData);
      } else {
        toast.error("Failed to fetch workspace container ports");
        setContainerExposedPorts([]);
      }
    },
    [workspace]
  );

  const handleDeletePort = useCallback(
    async (port: ContainerPort) => {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      var [status, statusCode] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/workspace/${workspace.id}/container/${container.container_name
        }/port/${port.port_number}`,
        "DELETE",
        null
      );

      if (statusCode !== 204) {
        toast.error("Failed to remove port");
      } else {
        FetchSelectedContainerPorts(container.container_name);
        onChange();
      }
    },
    [workspace.id, container.container_name, FetchSelectedContainerPorts]
  );

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
        <ModalHeader toggle={handleClose} className="border-0">
          Edit exposed ports
        </ModalHeader>
        <ModalBody className="pt-0">
          <Table striped bordered className="mb-0">
            <thead>
              <tr>
                <th>Port number</th>
                <th>Service name</th>
                <th>Public</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {containerExposedPorts.length === 0 ? (
                <React.Fragment>
                  <tr>
                    <td className="py-2 text-center" colSpan={4}>
                      There are no port exposed
                    </td>
                  </tr>
                </React.Fragment>
              ) : (
                containerExposedPorts.map((port) => (
                  <React.Fragment>
                    <tr>
                      <td style={{ width: 125, paddingTop: 12 }}>{port.port_number}</td>
                      <td style={{ paddingTop: 12 }}>{port.service_name}</td>
                      <td style={{ width: 100, paddingTop: 12 }}>{port.public ? "Yes" : "No"}</td>
                      <td style={{ width: 75 }}>
                        <Button
                          color="outline-danger"
                          size="sm"
                          style={{ width: 75 }}
                          onClick={() => handleDeletePort(port)}
                        >
                          Remove
                        </Button>
                      </td>
                    </tr>
                  </React.Fragment>
                ))
              )}
            </tbody>
          </Table>
          <div className="d-flex align-items-end justify-content-end mt-4">
            <Button color="primary" onClick={() => setShowAddPortModal(true)}>
              Add port
            </Button>
            <Button color="accent" onClick={() => handleClose()} className="ms-2">
              Close
            </Button>
          </div>
        </ModalBody>
      </Modal>
      <EditExposedPortsAddPortModal
        isOpen={showAddPortModal}
        onClose={() => {
          onChange();
          setShowAddPortModal(false);
          FetchSelectedContainerPorts(container.container_name);
        }}
        container={container}
        workspace={workspace}
      />
      <ToastContainer />
    </React.Fragment>
  );
}
