import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Table } from "reactstrap";
import { WorkspaceTemplate } from "../../types/templates";
import { RequestStatus } from "../../api/types";
import { Http } from "../../api/http";
import { Workspace } from "../../types/workspace";

interface TemplateDetailsWorkspacesProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsWorkspaces({ template }: TemplateDetailsWorkspacesProps) {
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);

    const FetchWorkspacesUsingTemplate = useCallback(async () => {
        // fetch template versions
        let [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${template.id}/workspaces`,
            "GET",
            null
        );

        if (status === RequestStatus.OK && statusCode === 200) {
            setWorkspaces(responseData as Workspace[]);
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
