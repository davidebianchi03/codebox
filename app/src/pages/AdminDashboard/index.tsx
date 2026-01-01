import React, { useCallback, useEffect, useState } from "react";
import { Card, Col, Container, Row } from "reactstrap";
import { toast } from "react-toastify";
import { AdminRetrieveStats } from "../../api/admin";
import { AdminStats } from "../../types/admin";
import { LoginsInLast7Days } from "./LoginsInLast7Days";
import { UsersListTable } from "./UsersListTable";
import { RunnersListTable } from "./RunnersListTable";
import { useSelector } from "react-redux";
import { RootState } from "../../redux/store";


export function AdminDashboard() {
    const [stats, setStats] = useState<AdminStats>();
    const user = useSelector((state:RootState) => state.user);

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

    useEffect(() => {
        FetchAdminStats();
    }, [FetchAdminStats]);

    return (
        <React.Fragment>
            <Container>
                <h1>Admin Dashboard</h1>
                <Row>
                    <Col md={12} className="mt-5">
                        <Row>
                            <Col md={4} className="mb-3">
                                <Card body>
                                    <h3>Total Users</h3>
                                    <h1>{stats?.total_users}</h1>
                                </Card>
                            </Col>
                            <Col md={4} className="mb-3">
                                <Card body>
                                    <h3>Active Workspaces</h3>
                                    <h1>{stats?.online_workspaces}</h1>
                                </Card>
                            </Col>
                            <Col md={4} className="mb-3">
                                <Card body>
                                    <h3>Online Runners</h3>
                                    <h1>{stats?.online_runners}</h1>
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
