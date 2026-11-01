import React, { useCallback, useState } from "react";
import { Badge, Button, Form, Modal } from "react-bootstrap";

interface ConfirmDeleteRunnerModalProps {
    isOpen: boolean;
    onClose: (deleteConfirmed: boolean) => void;
}

export function ConfirmDeleteRunnerModal({
    isOpen,
    onClose,
}: ConfirmDeleteRunnerModalProps) {
    const [confirmText, setConfirmText] = useState<string>("");

    const HandleCloseModal = useCallback((deleteConfirmed: boolean) => {
        onClose(deleteConfirmed);
    }, [onClose]);

    return (
        <React.Fragment>
            <Modal
                show={isOpen}
                onHide={() => HandleCloseModal(false)}
                centered
                size="lg"
                className="modal-blur"
            >
                <Modal.Header closeButton>
                    <Modal.Title>Delete Runner</Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <p>
                        Are you sure that you want to delete this runner?
                        Workspaces using this runner will be stopped,
                        and users will have to select another runner to start them.
                        This could cause data loss.
                    </p>
                    <p>
                        Enter <Badge bg="accent">confirm</Badge> to confirm that you want to delete the runner
                    </p>
                    <Form.Control
                        placeholder="confirm"
                        onChange={(e) => setConfirmText(e.target.value)}
                        value={confirmText}
                    />
                </Modal.Body>
                <Modal.Footer>
                    <Button
                        variant="accent"
                        onClick={() => HandleCloseModal(false)}
                    >
                        Cancel
                    </Button>
                    <Button
                        variant="danger"
                        disabled={confirmText !== "confirm"}
                        onClick={() => HandleCloseModal(true)}
                    >
                        Delete
                    </Button>
                </Modal.Footer>
            </Modal>
        </React.Fragment>
    )
}