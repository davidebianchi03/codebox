import React, { useCallback, useState } from "react";
import { Badge, Button, Input, Modal, ModalBody, ModalFooter, ModalHeader } from "reactstrap";

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
                isOpen={isOpen}
                toggle={() => {
                    HandleCloseModal(false);
                }}
                centered
                size="lg"
                modalClassName="modal-blur"
                fade
            >
                <ModalHeader
                    toggle={() => {
                        HandleCloseModal(false);
                    }}
                >
                    Delete Runner
                </ModalHeader>
                <ModalBody>
                    <p>
                        Are you sure that you want to delete this runner?
                        Workspaces using this runner will be stopped,
                        and users will have to select another runner to start them.
                        This could cause data loss.
                    </p>
                    <p>
                        Enter <Badge color="accent">confirm</Badge> to confirm that you want to delete the runner
                    </p>
                    <Input
                        placeholder="confirm"
                        onChange={(e) => setConfirmText(e.target.value)}
                        value={confirmText}
                    />
                </ModalBody>
                <ModalFooter>
                    <Button
                        color="accent"
                        onClick={(e) => {
                            e.preventDefault();
                            HandleCloseModal(false);
                        }}
                    >
                        Cancel
                    </Button>
                    <Button
                        color="danger"
                        disabled={confirmText !== "confirm"}
                        onClick={(e) => {
                            e.preventDefault();
                            HandleCloseModal(true);
                        }}
                    >
                        Delete
                    </Button>
                </ModalFooter>
            </Modal>
        </React.Fragment>
    )
}