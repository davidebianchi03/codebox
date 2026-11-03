import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, Col, Container, Row } from "react-bootstrap";
import { ToastContainer } from "react-toastify";
import { SystemLog } from "../../types/logs";
import { APIAdminListSystemLogs } from "../../api/admin";
import { LogsViewer } from "../../components/LogsViewer";
import { AdminPageHeader } from "../../components/AdminPageHeader";

export function AdminSystemLogsPage() {

    const [logs, setLogs] = useState<SystemLog[]>([]);

    const fetchLogs = useCallback(async () => {
        const r = await APIAdminListSystemLogs();
        if (r) {
            setLogs(r);
        }
    }, []);

    const downloadLogs = () => {
        var element = document.createElement('a');
        element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(JSON.stringify(logs)));
        element.setAttribute('download', 'system-logs.json');

        element.style.display = 'none';
        document.body.appendChild(element);

        element.click();

        document.body.removeChild(element);
    }

    useEffect(() => {
        fetchLogs();
    }, [fetchLogs]);

    return (
        <React.Fragment>
            <Container>
                <AdminPageHeader
                    title="System Logs"
                    description="View system logs"
                    actions={
                      <Button variant="outline-light" onClick={downloadLogs}>
                        Download Logs
                      </Button>
                    }
                />
                <Row className="mt-4">
                    <Col>
                        <Card body>
                            <LogsViewer logs={logs} />
                        </Card>
                    </Col>
                </Row>
            </Container>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    )
}
