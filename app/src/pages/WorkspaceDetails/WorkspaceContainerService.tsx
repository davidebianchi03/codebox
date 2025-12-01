import { faChevronRight } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React from "react";

interface WorkspaceContainerServiceProps {
    icon: string;
    title: string;
    url: string;
    description: string;
}

export function WorkspaceContainerService({
    icon,
    title,
    url,
    description,
}: WorkspaceContainerServiceProps) {
    return (
        <React.Fragment>
            <div style={{ marginTop: 5 }} className="my-1">
                <a
                    className="d-flex alert rounded align-items-center px-2 text-light"
                    style={{ cursor: "pointer", height: 50 }}
                    href={url}
                    target="_blank"
                >
                    <img src={icon} alt="vscode" width={25} className="me-3" />
                    <div className="d-flex justify-content-between align-items-center w-100 me-2">
                        <div className="d-flex align-items-center">
                            <h4 className="mb-0">{title}</h4>
                            <span className="text-muted ms-5">
                                {description}
                            </span>
                        </div>
                        <span className="text-muted">
                            <FontAwesomeIcon icon={faChevronRight} />
                        </span>
                    </div>
                </a>
            </div>
        </React.Fragment>
    )
}
