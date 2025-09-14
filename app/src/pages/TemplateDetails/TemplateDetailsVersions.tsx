import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Table } from "reactstrap";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import { toast } from "react-toastify";
import { APIListTemplateVersionsByTemplate } from "../../api/templates";

interface TemplateDetailsVersionsProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsVersions({ template }: TemplateDetailsVersionsProps) {
    const [versions, setVersions] = useState<WorkspaceTemplateVersion[]>();

    const fetchVersions = useCallback(async () => {
        const v = await APIListTemplateVersionsByTemplate(template.id);
        if (v) {
            setVersions(v.reverse());
        } else {
            toast.error("Failed to fetch template versions");
            setVersions(undefined);
        }
    }, [template.id]);

    useEffect(() => {
        fetchVersions();
    }, [fetchVersions]);

    return (
        <React.Fragment>
            <Card>
                <CardHeader className="border-0 pb-0">
                    <h3>Versions</h3>
                </CardHeader>
                <CardBody className="pt-0">
                    <ul className="timeline timeline-simple">
                        {versions?.map((version, index) => (
                            <React.Fragment key={index}>
                                <li className="timeline-event">
                                    <div className="card timeline-event-card">
                                        <div className="card-body">
                                            <div className="text-secondary float-end">
                                                {new Date(version.edited_on).toLocaleString()}
                                            </div>
                                            <h4>{version.name}</h4>
                                            <p className="text-secondary">
                                                {version.published ? "Released" : "Editing"}
                                            </p>
                                        </div>
                                    </div>
                                </li>
                            </React.Fragment>
                        ))}
                    </ul>
                </CardBody>
            </Card>
        </React.Fragment>
    );
}
