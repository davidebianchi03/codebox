import React, { useCallback, useEffect, useState } from "react";
import { Workspace } from "../../types/workspace";
import { APIDeleteWorkspace, APIRetrieveWorkspaceById, APIStartWorkspace, APIStopWorkspace } from "../../api/workspace";
import { toast } from "react-toastify";
import Swal from "sweetalert2";
import { GetBeautyNameForStatus, GetWorkspaceStatusColor } from "../../common/workspace";
import { WorkspaceSelectRunnerModal } from "./WorkspaceSelectRunnerModal";
import { Button, Dropdown } from "react-bootstrap";

interface WorkspaceStatusDropdownProps {
    workspace: Workspace;
    onStatusChange: () => void;
}

export function WorkspaceStatusDropdown({
    workspace,
    onStatusChange
}: WorkspaceStatusDropdownProps) {
    const [statusName, setStatusName] = useState<string>(GetBeautyNameForStatus(workspace.status));
    const [showSelectRunnerModal, setShowSelectRunnerModal] = useState<boolean>(false);

    const HandleStartWorkspace = useCallback(async () => {
        // check if workspace has a runner assigned
        const w = await APIRetrieveWorkspaceById(workspace.id);
        if (w) {
            if (w.runner == null) {
                setShowSelectRunnerModal(true);
            } else {
                if (await APIStartWorkspace(workspace.id)) {
                    onStatusChange();
                } else {
                    toast.error(`Failed to start workspace, try again later`);
                }
            }
        } else {
            toast.error(
                `Failed to fetch workspace details, try again later`
            );
        }
    }, [onStatusChange, workspace]);

    const HandleStopWorkspace = useCallback(async () => {
        if (await APIStopWorkspace(workspace.id)) {
            onStatusChange();
        } else {
            toast.error(`Failed to stop workspace, try again later`);
        }
    }, [onStatusChange, workspace]);

    const HandleDeleteWorkspace = useCallback(async (force: boolean) => {
        if (
            (
                await Swal.fire({
                    title: "Delete workspace",
                    icon: "warning",
                    text: `
                  Are you sure that you want to delete this workspace?
                  ${force && (`
                    Force-deleting a workspace may result in orphaned containers if runner errors, 
                    including connection issues or container removal failures, are encountered
                  `)}
                `,
                    showCancelButton: true,
                    reverseButtons: true,
                    confirmButtonText: "Delete",
                    customClass: {
                        popup: "bg-dark text-light",
                        cancelButton: "btn btn-accent",
                        confirmButton: "btn btn-warning",
                    },
                })
            ).isConfirmed
        ) {
            if (await APIDeleteWorkspace(workspace.id, force)) {
                onStatusChange();
            } else {
                toast.error(
                    `Failed to delete workspace, try again later`
                );
            }
        }
    }, [onStatusChange, workspace]);

    useEffect(() => {
        setStatusName(GetBeautyNameForStatus(workspace.status));
    }, [workspace.status]);

    return (
        <React.Fragment>
            {
                workspace.status != "starting" &&
                    workspace.status != "stopping" &&
                    workspace.status != "deleting" ? (
                    <Dropdown>
                        <Dropdown.Toggle variant={GetWorkspaceStatusColor(workspace.status)} id="dropdown-basic">
                            {statusName}
                        </Dropdown.Toggle>
                        <Dropdown.Menu>
                            <Dropdown.Item onClick={() => {
                                if (
                                    workspace.status === "running" ||
                                    workspace.status === "error"
                                ) {
                                    HandleStopWorkspace();
                                } else {
                                    HandleStartWorkspace();
                                }
                            }}>
                                {workspace?.status === "running" ||
                                    workspace?.status === "error"
                                    ? "Stop workspace"
                                    : "Start workspace"}
                            </Dropdown.Item>
                            <Dropdown.Item onClick={() => {
                                HandleDeleteWorkspace(false);
                            }}>
                                Delete workspace
                            </Dropdown.Item>
                            {workspace.status === "error" && (
                                <Dropdown.Item onClick={() => {
                                    HandleDeleteWorkspace(true);
                                }}>
                                    Force delete workspace
                                </Dropdown.Item>
                            )}
                        </Dropdown.Menu>
                    </Dropdown>
                ) : (
                    <Button variant={GetWorkspaceStatusColor(workspace.status)}>
                        {statusName}
                    </Button>
                )}
            <WorkspaceSelectRunnerModal
                isOpen={showSelectRunnerModal}
                onClose={async (updated) => {
                    setShowSelectRunnerModal(false);
                    if (updated) {
                        HandleStartWorkspace();
                    }
                }}
                workspace={workspace}
            />
        </React.Fragment>
    )
}
