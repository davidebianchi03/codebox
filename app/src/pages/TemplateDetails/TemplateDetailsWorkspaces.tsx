import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Table } from "reactstrap";
import { WorkspaceTemplate } from "../../types/templates";
import { Workspace } from "../../types/workspace";
import { APIListWorkspacesByTemplate } from "../../api/templates";

interface TemplateDetailsWorkspacesProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsWorkspaces({ template }: TemplateDetailsWorkspacesProps) {
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);

    const FetchWorkspacesUsingTemplate = useCallback(async () => {
        const w = await APIListWorkspacesByTemplate(template.id);
        if (w) {
            setWorkspaces(w);
        }
    }, [template.id]);

    useEffect(() => {
        FetchWorkspacesUsingTemplate();
    }, [FetchWorkspacesUsingTemplate]);

    return (
        <React.Fragment>
            <Card>
                <CardHeader className="border-0 pb-0">
                    <h3>Workspaces</h3>
                </CardHeader>
                <CardBody className="pt-0">
                    <Table striped>
                        <thead>
                            <tr>
                                <th>
                                    Name
                                </th>
                                <th>
                                    Owner
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                                <React.Fragment>
                                    {workspaces.length > 0 ? (
                                        workspaces.map((workspace, index) => (
                                            <tr key={index}>
                                                <td>
                                                    {workspace.name}
                                                </td>
                                                <td style={{ width: 150 }}>
                                                    {workspace.user.first_name} {workspace.user.last_name}
                                                </td>
                                            </tr>
                                        ))
                                    ) : (
                                        <tr>
                                            <td colSpan={2} className="text-center">
                                                There are't workspaces that use this template
                                            </td>
                                        </tr>
                                    )}
                                </React.Fragment>
                        </tbody>
                    </Table>
                </CardBody>
            </Card>
        </React.Fragment>
    );
}
