import React, { useCallback, useEffect, useState } from "react";
import { Card, Col, Container, Row } from "reactstrap";
import { toast } from "react-toastify";
import { AdminRetrieveStats } from "../../api/admin";
import { AdminListRunners } from "../../api/runner";
import { AdminStats } from "../../types/admin";
import { RunnerAdmin } from "../../types/runner";
import { LoginsInLast7Days } from "./LoginsInLast7Days";
import { UsersListTable } from "./UsersListTable";
import { RunnersListTable } from "./RunnersListTable";
import { AdminRunnerStatusBadge, IsRunnerOnline } from "../../components/AdminRunnerStatusBadge";
import { BackhoeIcon, PackagesIcon, StatsIcon, UserIcon } from "../../icons/Tabler";
import { Link } from "react-router-dom";
import { AdminPageHeader } from "../../components/AdminPageHeader";

export default function AdminDashboard() {
    const [stats, setStats] = useState<AdminStats>();
    const [runners, setRunners] = useState<RunnerAdmin[]>();
    const [runnersLoading, setRunnersLoading] = useState(true);

    const fetchAdminStats = useCallback(async () => {
        const data = await AdminRetrieveStats();
        if (data === undefined) {
            toast.error("Failed to fetch admin stats");
            return;
        }
        setStats(data);
    }, []);

    const fetchRunners = useCallback(async () => {
        const data = await AdminListRunners();
        if (data === undefined) {
            toast.error("Failed to fetch runners");
        } else {
            setRunners(data);
        }
        setRunnersLoading(false);
    }, []);

    useEffect(() => {
        fetchAdminStats();
        fetchRunners();
    }, [fetchAdminStats, fetchRunners]);

    const onlineRunnerCount = runners?.filter((runner) => IsRunnerOnline(runner)).length ?? 0;
    const offlineRunnerCount = (runners?.length ?? 0) - onlineRunnerCount;
    const onlineRunnerPercentage = runners?.length
        ? (onlineRunnerCount / runners.length) * 100
        : 0;

    const metrics = [
        {
            label: "Total users",
            value: stats?.total_users ?? "—",
            icon: <UserIcon />,
            tone: "blue",
        },
        {
            label: "Active workspaces",
            value: stats?.online_workspaces ?? "—",
            icon: <PackagesIcon />,
            tone: "indigo",
        },
        {
            label: "Online runners",
            value: stats?.online_runners ?? "—",
            icon: <BackhoeIcon />,
            tone: "green",
        },
        {
            label: "Codebox version",
            value: import.meta.env.VITE_APP_VERSION ?? "—",
            icon: <StatsIcon />,
            tone: "cyan",
        },
    ];

    return (
        <Container className="admin-dashboard">
            <AdminPageHeader
                title="Overview"
                description="A quick look at your Codebox instance"
            />

            <Row className="g-3 dashboard-metrics">
                {metrics.map((metric) => (
                    <Col key={metric.label} xs={12} sm={6} xl={3}>
                        <Card body className={`dashboard-stat dashboard-stat-${metric.tone} h-100`}>
                            <div className="dashboard-stat-top">
                                <span className="dashboard-stat-icon">{metric.icon}</span>
                                <span className="dashboard-stat-label">{metric.label}</span>
                            </div>
                            <div className="dashboard-stat-value">{metric.value}</div>
                        </Card>
                    </Col>
                ))}
            </Row>

            <Row className="g-3 dashboard-panels">
                <Col xl={6}>
                    <LoginsInLast7Days data={stats?.login_counts_last_7_days ?? []} />
                </Col>
                <Col xl={6}>
                    <Card body className="dashboard-panel dashboard-runner-status h-100">
                        <div className="dashboard-panel-heading">
                            <div className="dashboard-panel-title">
                                <h2>Runner status</h2>
                            </div>
                            <Link to="/admin/runners" className="dashboard-panel-link">View runners <span aria-hidden="true">→</span></Link>
                        </div>
                        <div className="dashboard-runner-content mt-3">
                            <div
                                className={`dashboard-runner-ring${runnersLoading || !runners?.length ? " is-empty" : ""}`}
                                style={{
                                    background: runners?.length
                                        ? `conic-gradient(var(--tblr-success) ${onlineRunnerPercentage}%, var(--tblr-danger) ${onlineRunnerPercentage}% 100%)`
                                        : "rgba(120, 146, 180, 0.25)",
                                }}
                                aria-label={
                                    runnersLoading
                                        ? "Loading runner status"
                                        : runners
                                            ? `${onlineRunnerCount} online, ${offlineRunnerCount} offline`
                                            : "Runner status unavailable"
                                }
                            >
                                <div className="dashboard-runner-ring-center">
                                    <strong>{runnersLoading || !runners ? "—" : onlineRunnerCount}</strong>
                                    <span>Online</span>
                                </div>
                            </div>
                            <div className="dashboard-runner-legend">
                                <div>
                                    <div>
                                        <AdminRunnerStatusBadge status="online" />
                                    </div>
                                    <strong>{runnersLoading || !runners ? "—" : onlineRunnerCount}</strong>
                                </div>
                                <div>
                                    <div>
                                        <AdminRunnerStatusBadge status="offline" />
                                    </div>
                                    <strong>{runnersLoading || !runners ? "—" : offlineRunnerCount}</strong>
                                </div>
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>

            <Row className="g-3 dashboard-tables">
                <Col xs={12}><UsersListTable /></Col>
                <Col xs={12}><RunnersListTable runners={runners} loading={runnersLoading} /></Col>
            </Row>
        </Container>
    );
}
