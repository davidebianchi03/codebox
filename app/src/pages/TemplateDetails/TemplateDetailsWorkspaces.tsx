import React, { useCallback, useEffect, useState } from "react";
import { Badge, Card, CardBody, CardHeader, Input, Table } from "reactstrap";
import { WorkspaceTemplate } from "../../types/templates";
import { Workspace } from "../../types/workspace";
import { APIListWorkspacesByTemplate } from "../../api/templates";
import { GetBeautyNameForStatus, GetWorkspaceStatusColor } from "../../common/workspace";

interface TemplateDetailsWorkspacesProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsWorkspaces({ template }: TemplateDetailsWorkspacesProps) {
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
    const [searchText, setSearchText] = useState<string>("");

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
                    <Input
                        placeholder="Search workspaces..."
                        className="mb-3"
                        value={searchText}
                        onChange={(e) => setSearchText(e.target.value)}
                    />
                    <Table striped>
                        <thead>
                            <tr>
                                <th>
                                    Name
                                </th>
                                <th>
                                    Owner
                                </th>
                                <th>
                                    Status
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <React.Fragment>
                                {workspaces.filter(
                                    (w) => w.name.toLowerCase().includes(searchText.toLowerCase()) ||
                                        w.user.first_name.toLowerCase().includes(searchText.toLowerCase()) ||
                                        w.user.last_name.toLowerCase().includes(searchText.toLowerCase())
                                ).length > 0 ? (
                                    workspaces.filter(
                                        (w) => w.name.toLowerCase().includes(searchText.toLowerCase()) ||
                                            w.user.first_name.toLowerCase().includes(searchText.toLowerCase()) ||
                                            w.user.last_name.toLowerCase().includes(searchText.toLowerCase())
                                    ).map((workspace, index) => (
                                        <tr key={index}>
                                            <td>
                                                <div className="mt-1">
                                                    {workspace.name}
                                                </div>
                                            </td>
                                            <td>
                                                <div className="mt-1">
                                                    {workspace.user.first_name} {workspace.user.last_name}
                                                </div>
                                            </td>
                                            <td>
                                                <div className="mt-1">
                                                    <Badge
                                                        color={GetWorkspaceStatusColor(workspace.status)}
                                                        className="text-white mb-2"
                                                        style={{ fontSize: 11 }}
                                                    >
                                                        {GetBeautyNameForStatus(workspace.status)}
                                                    </Badge>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={3} className="text-center">
                                            {workspaces.length === 0 ? "No workspaces are using this template yet." : "No workspaces found."}
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
