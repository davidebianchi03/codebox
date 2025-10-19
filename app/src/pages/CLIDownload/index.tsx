import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faLinux, faWindows } from '@fortawesome/free-brands-svg-icons'
import React, { useCallback, useEffect, useState } from "react";
import { Badge, Card, CardBody, CardHeader, Col, Row, Table } from "reactstrap";
import { CLIBuild } from "../../types/cli";
import { ListCLIBuilds } from "../../api/cli";
import { toast } from "react-toastify";


export function CLIDownloadPage() {

    const [cliBuilds, setCliBuilds] = useState<CLIBuild[]>([]);

    const fetchCliBuilds = useCallback(async () => {
        const builds = await ListCLIBuilds();
        if (builds) {
            setCliBuilds(builds);
        } else {
            toast.error("Failed to fetch CLI builds");
        }
    }, []);

    useEffect(() => {
        fetchCliBuilds();
    }, [fetchCliBuilds]);

    return (
        <React.Fragment>
            <div className="col mt-4">
                <h2 className="page-title">CLI</h2>
            </div>
            <Row className="mt-3">
                <Col>
                    <p>
                        Download the codebox CLI from here
                    </p>
                </Col>
            </Row>
            <Row className="mt-3">
                <Col md={6}>
                    <Card>
                        <CardHeader>
                            <h2 className="mb-0">
                                <span className="pe-2">Windows</span>
                                <FontAwesomeIcon icon={faWindows as any} />
                            </h2>
                        </CardHeader>
                        <CardBody>
                            <Table responsive className="mb-0">
                                <tbody>
                                    {cliBuilds.filter(b => b.os === "windows").map(build => (
                                        <tr key={build.id}>
                                            <td>
                                                <p className="mb-2">
                                                    <a href={`${import.meta.env.VITE_SERVER_URL}/api/v1/cli/${build.id}/download`} target="_blank" rel="noopener noreferrer">
                                                        {build.name}
                                                    </a>
                                                </p>
                                                <p className="d-flex gap-2 mb-0">
                                                    <Badge color="primary" className="text-white">
                                                        {build.architecture}
                                                    </Badge>
                                                    <Badge color="primary" className="text-white">
                                                        {build.type}
                                                    </Badge>
                                                </p>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </CardBody>
                    </Card>
                </Col>
                <Col md={6}>
                    <Card>
                        <CardHeader>
                            <h2 className="mb-0">
                                <span className="pe-2">Linux</span>
                                <FontAwesomeIcon icon={faLinux as any} />
                            </h2>
                        </CardHeader>
                        <CardBody>
                            <Table responsive className="mb-0">
                                <tbody>
                                    {cliBuilds.filter(b => b.os === "linux").map(build => (
                                        <tr key={build.id}>
                                            <td>
                                                <p className="mb-2">
                                                    <a href={`${import.meta.env.VITE_SERVER_URL}/api/v1/cli/${build.id}/download`} target="_blank" rel="noopener noreferrer">
                                                        {build.name}
                                                    </a>
                                                </p>
                                                <p className="d-flex gap-2 mb-0">
                                                    <Badge color="primary" className="text-white">
                                                        {build.architecture}
                                                    </Badge>
                                                    <Badge color="primary" className="text-white">
                                                        {build.type}
                                                    </Badge>
                                                </p>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </CardBody>
                    </Card>
                </Col>
                {/* TODO: macos */}
            </Row>
        </React.Fragment>
    );
}
