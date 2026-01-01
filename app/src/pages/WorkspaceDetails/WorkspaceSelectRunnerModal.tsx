import React, { useCallback, useEffect, useState } from "react";
import {
  Button,
  FormFeedback,
  FormGroup,
  Input,
  Label,
  Modal,
  ModalBody,
  ModalHeader,
  Spinner,
} from "reactstrap";
import { Workspace } from "../../types/workspace";
import { toast, ToastContainer } from "react-toastify";
import { Runner } from "../../types/runner";
import { ListRunners } from "../../api/runner";
import { APISetWorkspaceRunner } from "../../api/workspace";

interface Props {
  isOpen: boolean;
  onClose: (updated: boolean) => void;
  workspace: Workspace;
}

export function WorkspaceSelectRunnerModal({ isOpen, onClose, workspace }: Props) {
  const [runners, setRunners] = useState<Runner[]>([]);
  const [selectedRunner, setSelectedRunner] = useState<Runner>();
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  const FetchRunners = useCallback(async () => {
    const r = await ListRunners();
    if (r) {
      setRunners(r);
    }
  }, []);

  const handleCloseModal = useCallback((updated: boolean) => {
    onClose(updated);
    setSelectedRunner(undefined);
    setError("");
  }, [onClose]);

  const handleSubmit = useCallback(async () => {
    if (selectedRunner !== undefined) {
      setError("");
      setSaving(true);
      // TODO: check that the runner is compatible
      if (!(await APISetWorkspaceRunner(
        workspace.id,
        selectedRunner.id,
      ))) {
        toast.error("Failed to update the workspace");
      } else {
        handleCloseModal(true);
      }
      setSaving(false);
    } else {
      setError("Please select a runner to continue");
    }
  }, [handleCloseModal, selectedRunner, workspace]);

  useEffect(() => {
    FetchRunners();
  }, [FetchRunners]);

  return (
    <React.Fragment>
      <Modal
        centered
        isOpen={isOpen}
        toggle={() => handleCloseModal(false)}
        modalClassName="modal-blur"
      >
        <ModalHeader toggle={() => handleCloseModal(false)}>Select a runner</ModalHeader>
        <ModalBody>
          <p>No runner is selected. Select a runner to start the workspace.</p>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleSubmit();
              return false;
            }}
          >
            <FormGroup>
              <Label>Runner</Label>
              <Input
                type="select"
                onChange={
                  (e) => {
                    if (parseInt(e.target.value) > -1) {
                      setSelectedRunner(runners.find(r => r.id === parseInt(e.target.value)));
                      setError("");
                    } else {
                      setSelectedRunner(undefined)
                    }
                  }
                }
                value={selectedRunner ? selectedRunner.id : "-1"}
                invalid={error !== ""}
              >
                <option value={-1}>---------</option>
                {runners.map(r => (
                  <option value={r.id} key={r.id}>
                    {r.name}
                  </option>
                ))}
              </Input>
              <FormFeedback>{error}</FormFeedback>
              {selectedRunner && (
                new Date().getTime() - new Date(selectedRunner.last_contact).getTime() > 5 * 1000 * 60 && (
                  <span className="text-warning">
                    Warning: last contact with this runner was more than 5 minutes ago.
                  </span>
                )
              )}
            </FormGroup>
            <div className="d-flex align-items-center justify-content-end">
              <Button
                color="accent"
                onClick={(e) => {
                  e.preventDefault();
                  handleCloseModal(false);
                }}
              >
                Close
              </Button>
              <Button
                color="light"
                className="ms-2"
                type="submit"
                disabled={saving}
              >
                {saving && <Spinner size="sm" className="me-2" />}
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
