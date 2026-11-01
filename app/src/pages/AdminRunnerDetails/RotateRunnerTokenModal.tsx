import { useEffect, useState } from "react";
import { faCopy } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { Button, Form, InputGroup, Modal, Spinner } from "react-bootstrap";
import { toast } from "react-toastify";
import { APIAdminRotateRunnerToken } from "../../api/admin";

interface RotateRunnerTokenModalProps {
  isOpen: boolean;
  runnerId: number;
  onClose: () => void;
}

export function RotateRunnerTokenModal({
  isOpen,
  runnerId,
  onClose,
}: RotateRunnerTokenModalProps) {
  const [isRotating, setIsRotating] = useState(false);
  const [newToken, setNewToken] = useState<string>();

  useEffect(() => {
    if (isOpen) {
      setNewToken(undefined);
    }
  }, [isOpen]);

  const rotateToken = async () => {
    setIsRotating(true);
    try {
      const result = await APIAdminRotateRunnerToken(runnerId);
      if (!result) {
        toast.error("Failed to rotate the runner token");
        return;
      }

      setNewToken(result.token);
    } finally {
      setIsRotating(false);
    }
  };

  const copyToken = async () => {
    if (!newToken) {
      return;
    }

    try {
      await navigator.clipboard.writeText(newToken);
      toast.success("Token copied to clipboard");
    } catch {
      toast.error("Failed to copy the token to clipboard");
    }
  };

  return (
    <Modal
      show={isOpen}
      onHide={() => {
        if (!isRotating) {
          onClose();
        }
      }}
      centered
      className="modal-blur"
      backdrop={isRotating ? "static" : true}
      keyboard={!isRotating}
    >
      <Modal.Header closeButton={!isRotating}>
        <Modal.Title>
          {newToken ? "Runner token rotated" : "Rotate runner token"}
        </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {newToken ? (
          <>
            <p>
              The previous token is no longer valid. Copy and save this new
              token now; it will not be shown again.
            </p>
            <InputGroup>
              <Form.Control
                readOnly
                value={newToken}
                aria-label="New runner token"
              />
              <Button
                variant="accent"
                onClick={copyToken}
                aria-label="Copy token"
                title="Copy token"
              >
                <FontAwesomeIcon icon={faCopy} />
              </Button>
            </InputGroup>
          </>
        ) : (
          <p>
            Rotating the token will immediately invalidate the current token.
            The connection between the runner and the server will be lost until the new token is used.
            <br/>
            <br/>
            Are you sure you want to continue?
          </p>
        )}
      </Modal.Body>
      <Modal.Footer>
        {newToken ? (
          <>
            <Button variant="accent" onClick={onClose}>
              Close
            </Button>
          </>
        ) : (
          <>
            <Button
              variant="accent"
              onClick={onClose}
              disabled={isRotating}
            >
              Cancel
            </Button>
            <Button
              variant="light"
              onClick={rotateToken}
              disabled={isRotating}
            >
              {isRotating && <Spinner size="sm" className="me-2" />}
              Rotate token
            </Button>
          </>
        )}
      </Modal.Footer>
    </Modal>
  );
}
