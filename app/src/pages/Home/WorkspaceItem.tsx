import React from "react";
import { Workspace, WorkspaceType } from "../../types/workspace";
import { GetBeautyNameForStatus, GetWorkspaceStatusColor } from "../../common/workspace";
import { Badge } from "react-bootstrap";
import { Link } from "react-router-dom";

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
                    <Badge
                        bg={GetWorkspaceStatusColor(workspace.status)}
                        className="text-white mb-2"
                        style={{ fontSize: 11 }}
                    >
                        {GetBeautyNameForStatus(workspace.status)}
                    </Badge>
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