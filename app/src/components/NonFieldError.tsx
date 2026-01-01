import React from "react";

export interface NonFieldErrorProps {
    error: string;
}

export function NonFieldError({ error }: NonFieldErrorProps) {
    return (
        <React.Fragment>
            <div className="mt-3">
                <p
                    className="border-0 d-flex justify-content-start px-3 py-2 mb-3 mt-2 text-danger"
                    style={{
                        background: "rgba(var(--tblr-danger-rgb), 0.1)",
                        borderLeft: "4px solid var(--tblr-danger) !important",
                        borderRadius: "4px",
                    }}
                    role="alert"
                >
                    {error}
                </p>
            </div>
        </React.Fragment>
    )
}