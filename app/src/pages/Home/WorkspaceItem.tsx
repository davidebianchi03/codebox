import React from "react";
import { Workspace, WorkspaceType } from "../../types/workspace";
import { GetBeautyNameForStatus, GetWorkspaceStatusColor } from "../../common/workspace";
import { OverlayTrigger, Popover } from "react-bootstrap";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTriangleExclamation } from "@fortawesome/free-solid-svg-icons";

export interface WorkspaceItemProps {
    workspace: Workspace;
    workspaceTypes: WorkspaceType[];
}

export function WorkspaceItem({
    workspace, workspaceTypes
}: WorkspaceItemProps) {
    return (
        <React.Fragment>
            <div className="d-flex align-items-center justify-content-between">
                <div className="d-flex align-items-center">
                    <div
                        style={{
                            width: 50,
                            height: 50,
                            fontSize: 20,
                            opacity: 0.5,
                            borderRadius: 4,
                        }}
                        className="bg-primary text-white d-flex align-items-center justify-content-center"
                    >
                        {workspace.name[0].toUpperCase()}
                    </div>
                    <div className="ms-4">
                        <h3 className="mb-0">
                            <Link to={`/workspaces/${workspace.id}`}>
                                {workspace.name}
                            </Link>
                        </h3>
                        <small className="text-muted">
                            {(() => {
                                var prettyType = "Unknown type";
                                var workspaceType = workspaceTypes.find(
                                    (wt: WorkspaceType) =>
                                        wt.id === workspace.type
                                );
                                if (workspaceType) {
                                    prettyType = workspaceType.name;
                                }
                                return prettyType;
                            })()}
                        </small>
                    </div>
                </div>
                <div className="d-flex flex-column align-items-end">
                    <div className="d-flex align-items-center gap-2 mb-2">
                        {!workspace.healthy && (
                            <OverlayTrigger
                                trigger={["hover", "focus", "click"]}
                                placement="left"
                                overlay={
                                    <Popover
                                        id={`workspace-health-${workspace.id}`}
                                    >
                                        <Popover.Header as="h3">
                                            Workspace is unhealthy
                                        </Popover.Header>
                                        <Popover.Body>
                                            One or more containers are offline. Some services may be unavailable.
                                        </Popover.Body>
                                    </Popover>
                                }
                            >
                                <button
                                    type="button"
                                    className="btn btn-link p-0 text-warning text-decoration-none"
                                    aria-label="Workspace is unhealthy"
                                >
                                    <FontAwesomeIcon icon={faTriangleExclamation} />
                                </button>
                            </OverlayTrigger>
                        )}
                        <span
                            className={`badge rounded-pill px-2 py-1 border border-${GetWorkspaceStatusColor(workspace.status)} bg-${GetWorkspaceStatusColor(workspace.status)}-lt text-${GetWorkspaceStatusColor(workspace.status)}`}
                        >
                            {GetBeautyNameForStatus(workspace.status)}
                        </span>
                    </div>
                    <p
                        className="mb-0 text-muted"
                        style={{ fontSize: 12 }}
                    >
                        <small>
                            Last activity{" "}
                            {new Date(
                                workspace.updated_at
                            ).toLocaleString()}
                        </small>
                    </p>
                </div>
            </div>
            <hr className="my-3" />
        </React.Fragment>
    )
}