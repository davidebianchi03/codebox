import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faApple, faLinux } from '@fortawesome/free-brands-svg-icons'
import React, { useCallback, useEffect, useState } from "react";
import { CLIBuild } from "../../types/cli";
import { ListCLIBuilds } from "../../api/cli";
import { toast } from "react-toastify";
import { WindowsIcon } from "../../icons/Tabler";
import { Badge, Card, Col, Row, Table } from "react-bootstrap";


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
                <Col md={12}>
                    <Card>
                        <Card.Header className="d-flex flex-column  align-items-start">
                            <h2 className="mb-0 d-flex align-items-center">
                                <span className="pe-2 mb-1">Windows</span>
                                <WindowsIcon />
                            </h2>
                            <p className="mb-0">
                                <small>
                                    Warning: Windows Defender may flag the CLI as a false positive,
                                    this is because the CLI is not signed with a valid certificate,
                                    you can safely ignore this warning and allow the CLI to run on your system.
                                </small>
                            </p>
                        </Card.Header>
                        <Card.Body>
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
                                                    <Badge bg="light" className="text-dark">
                                                        {build.architecture}
                                                    </Badge>
                                                    <Badge bg="light" className="text-dark">
                                                        {build.type}
                                                    </Badge>
                                                </p>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
            <Row className="mt-3">
                <Col md={12}>
                    <Card>
                        <Card.Header>
                            <h2 className="mb-0">
                                <span className="pe-2">Linux</span>
                                <FontAwesomeIcon icon={faLinux as any} />
                            </h2>
                        </Card.Header>
                        <Card.Body>
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
                                                    <Badge bg="light" className="text-dark">
                                                        {build.architecture}
                                                    </Badge>
                                                    <Badge bg="light" className="text-dark">
                                                        {build.type}
                                                    </Badge>
                                                </p>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
            <Row className="mt-3">
                <Col md={12}>
                    <Card>
                        <Card.Header>
                            <h2 className="mb-0">
                                <span className="pe-2">MacOS</span>
                                <FontAwesomeIcon icon={faApple as any} />
                            </h2>
                        </Card.Header>
                        <Card.Body>
                            <p>
                                <h4>Install codebox-cli using brew:</h4>
                                <div className="w-100 p-2 rounded" style={{ fontFamily: "Consolas", background: "var(--tblr-dark-bg-subtle)" }}>
                                    <p className="mb-0">brew tap codebox/codebox-cli https://gitlab.com/codebox4073715/codebox-homebrew-tap.git</p>
                                    <p className="mb-0">brew update</p>
                                    <p className="mb-0 text-success"># if codebox-cli is already installed on your mac you've to uninstall it first</p>
                                    <p className="mb-0">brew uninstall codebox-cli</p>
                                    <p className="mb-0">brew install codebox/codebox-cli/codebox-cli</p>
                                </div>
                            </p>
                            <p>
                                <small>
                                    Brew is required; you must install it first by viewing the official guide, you can find it at &nbsp;
                                    <a href="https://brew.sh/">https://brew.sh/</a>
                                </small>
                            </p>
                            <Table responsive className="mb-0">
                                <tbody>
                                    {cliBuilds.filter(b => b.os === "darwin").map(build => (
                                        <tr key={build.id}>
                                            <td>
                                                <p className="mb-2">
                                                    <a href={`${import.meta.env.VITE_SERVER_URL}/api/v1/cli/${build.id}/download`} target="_blank" rel="noopener noreferrer">
                                                        {build.name}
                                                    </a>
                                                </p>
                                                <p className="d-flex gap-2 mb-0">
                                                    <Badge bg="light" className="text-dark">
                                                        {build.architecture}
                                                    </Badge>
                                                    <Badge bg="light" className="text-dark">
                                                        {build.type}
                                                    </Badge>
                                                </p>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </Table>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </React.Fragment>
    );
}
