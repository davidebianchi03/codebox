import { Link } from "react-router-dom";
import "./WorkspaceListItem.css";
import { Component, ReactNode } from "react";
import { RetrieveBeautyNameForStatus, RetrieveColorForWorkspaceStatus } from "../../utils/workspaceStatus";

interface WorkspaceListItemProps {
    workspace: {
        id: number,
        name: string,
        status: string,
        last_activity_on: string,
    }
}

interface WorkspaceListItemState {

}

export class WorkspaceListItem extends Component<WorkspaceListItemProps, WorkspaceListItemState> {
    render(): ReactNode {
        return (
            <Link className="workspaces-list-item" to={`/workspaces/${this.props.workspace.id}`}>
                <div style={{
                    width: "60px",
                    height: "60px",
                    borderRadius: "7px",
                    margin: "5px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "40px",
                    fontFamily: "Consolas, monaco, monospace",
                    fontWeight: "bold",
                    backgroundColor: "#066fd188",
                    border: "solid #066fd1 1.5px"
                }}>
                    {this.props.workspace.name.toUpperCase()[0]}
                </div>
                <div style={{ marginLeft: "20pt", display: "flex", justifyContent: "space-between", width: "100%", alignItems: "center" }}>
                    <h4 style={{ marginBottom: 0, marginTop: 0 }}>{this.props.workspace.name}</h4>
                    <div style={{ textAlign: "right", marginRight: "20px" }}>
                        <span style={{
                            padding: "2px 7px",
                            fontSize: "11px",
                            background: `var(${RetrieveColorForWorkspaceStatus(this.props.workspace.status)})`,
                            borderRadius: "20px",
                            height: "15px",
                        }}>
                            {RetrieveBeautyNameForStatus(this.props.workspace.status)}
                        </span>
                        <div style={{ height: "5px" }}></div>
                        <span style={{ fontSize: "11px", color: "var(--grey-500)" }}>
                            Last activity on: {new Date(this.props.workspace.last_activity_on).toLocaleDateString()}
                        </span>
                    </div>
                </div>
            </Link>
        );
    }
}