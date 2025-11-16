import React, { useCallback, useState } from "react";
import { Workspace } from "../../types/workspace";
import { APIDeleteWorkspace, APIRetrieveWorkspaceById, APIStartWorkspace, APIStopWorkspace } from "../../api/workspace";
import { toast } from "react-toastify";
import Swal from "sweetalert2";
import { GetBeautyNameForStatus, GetWorkspaceStatusColor } from "../../common/workspace";
import { WorkspaceSelectRunnerModal } from "./WorkspaceSelectRunnerModal";

interface WorkspaceStatusDropdownProps {
    workspace: Workspace;
    onStatusChange: () => void;
}


export function WorkspaceStatusDropdown({
    workspace,
    onStatusChange
}: WorkspaceStatusDropdownProps) {
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

    return (
        <React.Fragment>
            <button
                className={`
                    btn btn-${GetWorkspaceStatusColor(workspace.status)} 
                    ${(
                        workspace.status != "starting" &&
                        workspace.status != "stopping" &&
                        workspace.status != "deleting" &&
                        "dropdown-toggle"
                    )}`}
                type="button"
                data-bs-toggle="dropdown"
                aria-haspopup="true"
                aria-expanded="false"
            >
                {GetBeautyNameForStatus(workspace.status)}
            </button>
            {(
                workspace.status != "starting" &&
                workspace.status != "stopping" &&
                workspace.status != "deleting" && (
                    <div className="dropdown-menu" aria-labelledby="dropdownMenuButton">
                        <span
                            className="dropdown-item"
                            onClick={() => {
                                if (
                                    workspace.status === "running" ||
                                    workspace.status === "error"
                                ) {
                                    HandleStopWorkspace();
                                } else {
                                    HandleStartWorkspace();
                                }
                            }}
                        >
                            {workspace?.status === "running" ||
                                workspace?.status === "error"
                                ? "Stop workspace"
                                : "Start workspace"}
                        </span>
                        <span
                            className="dropdown-item"
                            onClick={() => {
                                HandleDeleteWorkspace(false);
                            }}
                        >
                            Delete workspace
                        </span>
                        {workspace.status === "error" && (
                            <span
                                className="dropdown-item"
                                onClick={() => {
                                    HandleDeleteWorkspace(true);
                                }}
                            >
                                Force delete workspace
                            </span>
                        )}
                    </div>
                )
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
