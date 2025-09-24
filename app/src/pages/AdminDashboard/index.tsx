import React, { useCallback, useEffect, useState } from "react";
import { Card, Col, Container, Row, Table } from "reactstrap";
import { toast } from "react-toastify";
import { AdminRetrieveStats } from "../../api/admin";
import { AdminStats } from "../../types/admin";
import { LoginsInLast7Days } from "./LoginsInLast7Days";
import { User } from "../../types/user";
import { RetrieveCurrentUserDetails } from "../../api/common";
import { TimeSince } from "../../common/time";
import { UsersListTable } from "./UsersListTable";
import { RunnersListTable } from "./RunnersListTable";


export function AdminDashboard() {
    const [stats, setStats] = useState<AdminStats>();
    const [user, setUser] = useState<User>();

    const FetchAdminStats = useCallback(
        async () => {
            const data = await AdminRetrieveStats();
            if (data === undefined) {
                toast.error("Failed to fetch admin stats");
                return;
            }
            setStats(data);
        },
        [],
    );

    const FetchUserDetails = useCallback(
        async () => {
            const data = await RetrieveCurrentUserDetails();
            if (data === undefined) {
                toast.error("Failed to fetch user details");
                return;
            }
            setUser(data);
        },
        [],
    );

    useEffect(() => {
        FetchAdminStats();
        FetchUserDetails();
    }, [FetchAdminStats, FetchUserDetails]);

    return (
        <React.Fragment>
            <Container>
                <h1>Admin Dashboard</h1>
                <Row>
                    <Col md={12} className="mt-5">
                        <Row>
                            <Col md={3} className="mb-4">
                                <Card body>
                                    <h3>Total Users</h3>
                                    <h1>{stats?.total_users}</h1>
                                </Card>
                            </Col>
                            <Col md={3} className="mb-4">
                                <Card body>
                                    <h3>Active Workspaces</h3>
                                    <h1>{stats?.online_workspaces}</h1>
                                </Card>
                            </Col>
                            <Col md={3} className="mb-4">
                                <Card body>
                                    <h3>Online Runners</h3>
                                    <h1>{stats?.online_runners}</h1>
                                </Card>
                            </Col>
                            <Col md={3} className="mb-4">
                                <Card body>
                                    <h3>Last Login</h3>
                                    <h1>{user?.last_login ? TimeSince(new Date(user.last_login)) : "N/A"}</h1>
                                </Card>
                            </Col>
                        </Row>
                        <Row>
                            <Col md={4} className="mb-4">
                                <LoginsInLast7Days data={stats?.login_counts_last_7_days || []} />
                            </Col>
                            <Col md={8} className="mb-4">
                                <UsersListTable />
                            </Col>
                        </Row>
                        <Row>
                            <Col md={12} className="mb-4">
                                <RunnersListTable />
                            </Col>
                        </Row>
                    </Col>
                </Row>
            </Container>
        </React.Fragment>
    )
}
