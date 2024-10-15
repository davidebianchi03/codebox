import { Link } from "react-router-dom";
import "./WorkspaceListItem.css";
import { Component, ReactNode } from "react";

interface WorkspaceListItemProps {
    workspaceName: string
    workspaceId: number
}

interface WorkspaceListItemState {

}

export class WorkspaceListItem extends Component<WorkspaceListItemProps, WorkspaceListItemState> {
    render(): ReactNode {
        return (
            <Link className="workspaces-list-item" to={`/workspaces/${this.props.workspaceId}`}>
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
                    backgroundColor:"#066fd188",
                    border: "solid #066fd1 1.5px"
                }}>
                    {this.props.workspaceName.toUpperCase()[0]}
                </div>
                <h4 style={{ marginLeft: "20pt" }}>{this.props.workspaceName}</h4>
            </Link>
        );
    }
}