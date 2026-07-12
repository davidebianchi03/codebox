import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, Col, Container, Form, Row, Spinner } from "react-bootstrap";
import { toast, ToastContainer } from "react-toastify";
import { SystemLog } from "../../types/logs";
import { APIAdminListRunnerLogs, APIAdminListSystemLogs } from "../../api/admin";
import { LogsViewer } from "../../components/LogsViewer";
import { AdminListRunners } from "../../api/runner";
import { RunnerAdmin } from "../../types/runner";

export function AdminRunnerLogsPage() {
    const [runners, setRunners] = useState<RunnerAdmin[]>([]);
    const [selectedRunner, setSelectedRunner] = useState<RunnerAdmin | null>(null);
    const [logs, setLogs] = useState<SystemLog[]>([]);
    const [loadingRunners, setLoadingRunners] = useState<boolean>(true);
    const [loadingLogs, setLoadingLogs] = useState<boolean>(true);

    const FetchRunners = useCallback(async () => {
        setLoadingRunners(true);
        const r = await AdminListRunners();
        if (r) {
            setRunners(r);
            if (r.length > 0) {
                setSelectedRunner(r[0]);
            }
        } else {
            toast.error("Failed to fetch runners");
        }
        setLoadingRunners(false);
    }, []);


    useEffect(() => {
        FetchRunners();
    }, [FetchRunners]);

    const fetchLogs = useCallback(async () => {
        setLoadingLogs(true);
        if (selectedRunner) {
            const r = await APIAdminListRunnerLogs(selectedRunner.id);
            if (r) {
                setLogs(r);
            }
        } else {
            setLogs([]);
        }
        setLoadingLogs(false);
    }, [selectedRunner]);

    useEffect(() => {
        fetchLogs();
    }, [fetchLogs]);

    const downloadLogs = () => {
        var element = document.createElement('a');
        element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(JSON.stringify(logs)));
        element.setAttribute('download', `runner-${selectedRunner?.id || "unknown"}-logs.json`);

        element.style.display = 'none';
        document.body.appendChild(element);

        element.click();

        document.body.removeChild(element);
    }


    return (
        <React.Fragment>
            <Container>
                <div>
                    <h2 className="mb-1">Runner Logs</h2>
                    <small>View runner logs</small>
                </div>
                {loadingRunners ? (
                    <React.Fragment>
                        <Spinner animation="border" role="status">
                            <span className="visually-hidden">Loading...</span>
                        </Spinner>
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <div className="d-flex justify-content-end">
                            <Button className="mt-3" variant="outline-light" onClick={downloadLogs}>
                                Download Logs
                            </Button>
                        </div>
                        <Row className="mt-4">
                            <Col>
                                <Card body>
                                    <Form.Group>
                                        <Form.Label>Runner</Form.Label>
                                        <Form.Select
                                            onChange={(e) => {
                                                setSelectedRunner(
                                                    runners.find(r => String(r.id) === e.target.value) || null
                                                );
                                            }}
                                            value={selectedRunner?.id || ""}
                                        >
                                            <option value="">Select a runner</option>
                                            {runners.map((runner) => (
                                                <option key={runner.id} value={runner.id}>
                                                    {runner.name} (ID: {runner.id})
                                                </option>
                                            ))}
                                        </Form.Select>
                                    </Form.Group>
                                </Card>
                            </Col>
                        </Row>
                        {loadingLogs ? (
                            <React.Fragment>
                                <Spinner animation="border" role="status">
                                    <span className="visually-hidden">Loading...</span>
                                </Spinner>
                            </React.Fragment>
                        ) : (
                            <React.Fragment>
                                <Row className="mt-4">
                                    <Col>
                                        <Card body>
                                            <LogsViewer logs={logs} />
                                        </Card>
                                    </Col>
                                </Row>
                            </React.Fragment>
                        )}
                    </React.Fragment>
                )}
            </Container>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    )
}
