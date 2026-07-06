import React, { useCallback, useEffect, useState } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import { AdminAnalyticsConfig } from "../../components/AdminAnalyticsConfig";
import { ToastContainer } from "react-toastify";
import { AdminAnalyticsContentPreview } from "../../components/AdminAnalyticsContentPreview";
import { SystemLog } from "../../types/logs";
import { APIAdminListSystemLogs } from "../../api/admin";
import DataTable from "../../components/DataTable";

export function AdminSystemLogsPage() {

    const [logs, setLogs] = useState<SystemLog[]>([]);

    const fetchLogs = useCallback(async () => {
        const r = await APIAdminListSystemLogs();
        if (r) {
            setLogs(r);
        }
    }, []);

    useEffect(() => {
        fetchLogs();
    }, [fetchLogs]);

    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">System Logs</h2>
                    <small>View system logs</small>
                </div>
                <Row className="mt-4">
                    <Col>
                        <Card body>
                            <DataTable
                                columns={[
                                    {
                                        key: "level",
                                        label: "Level",
                                        render: (value: string) => (
                                            <span className={`badge bg-${value === "error" ? "danger" : value === "warning" ? "warning" : "info"} text-dark`}>
                                                {value.toUpperCase()}
                                            </span>
                                        )
                                    },
                                    {
                                        key: "timestamp",
                                        label: "Timestamp",
                                        render: (value: string) => (
                                            <small>
                                                {new Date(value).toLocaleString()}
                                            </small>
                                        )
                                    },
                                    {
                                        key: "log",
                                        label: "Log",
                                        render: (value: string) => (
                                            <small>
                                                {value}
                                            </small>
                                        )
                                    },
                                    {
                                        key: "function",
                                        label: "Function",
                                        render: (value: string, row: SystemLog) => (
                                            <small>
                                                {row.module}.{value}
                                            </small>
                                        )
                                    },
                                ]}
                                data={logs.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())}
                                initialPageSize={20}
                            />
                        </Card>
                    </Col>
                </Row>
            </Container>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    )
}
