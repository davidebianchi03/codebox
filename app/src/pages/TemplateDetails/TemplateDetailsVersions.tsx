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
                    <Table striped>
                        <thead>
                            <tr>
                                <th>
                                    Name
                                </th>
                                <th>
                                    Updated on
                                </th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody>
                            {versions && (
                                <React.Fragment>
                                    {versions.length > 0 ? (
                                        versions.map((version, index) => (
                                            <tr key={index}>
                                                <td>
                                                    <div className="mt-2">
                                                        {version.name}
                                                    </div>
                                                </td>
                                                <td>
                                                    <div className="mt-2">
                                                        {new Date(version.edited_on).toLocaleString()}
                                                    </div>
                                                </td>
                                                <td style={{ width: 150 }}>
                                                    {version.published ? (
                                                        <span className="btn border-success text-success w-100" style={{ cursor: "default" }}>
                                                            Released
                                                        </span>
                                                    ) : (
                                                        <span className="btn border-primary text-primary w-100" style={{ cursor: "default" }}>
                                                            Editing
                                                        </span>
                                                    )}
                                                </td>
                                            </tr>
                                        ))
                                    ) : (
                                        <tr>
                                            <td colSpan={3} className="text-center">
                                                No versions available
                                            </td>
                                        </tr>
                                    )}
                                </React.Fragment>
                            )}
                        </tbody>
                    </Table>
                </CardBody>
            </Card>
        </React.Fragment>
    );
}
