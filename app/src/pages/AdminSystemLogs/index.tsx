import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, Col, Container, Row } from "react-bootstrap";
import { ToastContainer } from "react-toastify";
import { SystemLog } from "../../types/logs";
import { APIAdminListSystemLogs } from "../../api/admin";
import { LogsViewer } from "../../components/LogsViewer";

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
                <div>
                    <h2 className="mb-1">System Logs</h2>
                    <small>View system logs</small>
                </div>
                <div className="d-flex justify-content-end">
                    <Button className="mt-3" variant="outline-light" onClick={downloadLogs}>
                        Download Logs
                    </Button>
                </div>
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
