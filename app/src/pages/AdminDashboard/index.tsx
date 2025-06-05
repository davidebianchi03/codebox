import React, { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Card, CardBody, Col, Row, Table } from "reactstrap";
import { User } from "../../types/user";
import { Workspace } from "../../types/workspace";
import { toast } from "react-toastify";
import { AdminListUsers, AdminListWorkspaces } from "../../api/admin";

export function AdminDashboard() {
    const [users, setUsers] = useState<User[]>([]);
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);

    const FetchUsers = useCallback(async () => {
        const u = await AdminListUsers();
        if (u) {
            setUsers(u);
        } else {
            toast.error("Failed to fetch users");
        }
    }, []);

    const FetchWorkspaces = useCallback(async () => {
        const w = await AdminListWorkspaces();
        if (w) {
            setWorkspaces(w);
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
