import React, { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Card, CardBody, Col, Row, Table } from "reactstrap";
import { User } from "../../types/user";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { Workspace } from "../../types/workspace";
import { toast } from "react-toastify";

export function AdminDashboard() {
    const [users, setUsers] = useState<User[]>([]);
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);

    const FetchUsers = useCallback(async () => {
        let [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/admin/users`,
            "GET",
            null
        );
        if (status === RequestStatus.OK && statusCode === 200) {
            setUsers(responseData as User[]);
        } else {
            toast.error("Failed to fetch users");
        }
    }, []);

    const FetchWorkspaces = useCallback(async () => {
        let [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/admin/workspaces`,
            "GET",
            null
        );
        if (status === RequestStatus.OK && statusCode === 200) {
            setWorkspaces(responseData as Workspace[]);
        } else {
            toast.error("Failed to fetch workspaces");
        }
    }, []);

    useEffect(() => {
        FetchUsers();
        FetchWorkspaces();
    }, [FetchUsers, FetchWorkspaces]);

    return (
        <React.Fragment>
            <Row>
                <Col md={12}>
                    <Row>
                        <Col md={4}>
                            <Card>
                                <CardBody>
                                    <h2>Users</h2>
                                    <div>
                                        <h1>
                                            {users.length} <small style={{ fontSize: 12 }}>users</small>
                                        </h1>
                                    </div>
                                    <Table>
                                        <tbody>
                                            {users.slice(-4).reverse().map(user => (
                                                <tr key={user.email}>
                                                    <td>
                                                        <Link to={`/admin/users/${user.email}`}>{user.email}</Link>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </Table>
                                    <div className="text-center">
                                        <Link to={"/admin/users"}>View All</Link>
                                    </div>
                                </CardBody>
                            </Card>
                        </Col>
                        <Col md={4}>
                            <Card>
                                <CardBody>
                                    <h2>Workspaces</h2>
                                    <div>
                                        <h1>
                                            {workspaces.length} <small style={{ fontSize: 12 }}>workspaces</small>
                                        </h1>
                                    </div>
                                    <Table>
                                        <tbody>
                                            {workspaces.reverse().slice(-4).map(workspace => (
                                                <tr key={workspace.id}>
                                                    <td>
                                                        <p className="mb-0">{workspace.name}</p>
                                                        <small className="text-muted">({workspace.user.email})</small>
                                                    </td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </Table>
                                </CardBody>
                            </Card>
                        </Col>
                        <Col md={4}>
                            <Card>
                                <CardBody>
                                    <h2>System info</h2>
                                    <Table>
                                        <tbody>
                                            <tr>
                                                <th>Version</th>
                                                <td>
                                                    {import.meta.env.VITE_APP_VERSION}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </Table>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </React.Fragment>
    )
}
