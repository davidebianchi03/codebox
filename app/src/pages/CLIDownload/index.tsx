import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faLinux, faWindows } from '@fortawesome/free-brands-svg-icons'
import React from "react";
import { Card, CardBody, CardHeader, Col, Row } from "reactstrap";

export function CLIDownloadPage() {
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
                                <FontAwesomeIcon icon={faWindows} />
                            </h2>
                        </CardHeader>
                        <CardBody>

                        </CardBody>
                    </Card>
                </Col>
                <Col md={6}>
                    <Card>
                        <CardHeader>
                            <h2 className="mb-0">
                                <span className="pe-2">Linux</span>
                                <FontAwesomeIcon icon={faLinux} />
                            </h2>
                        </CardHeader>
                        <CardBody>

                        </CardBody>
                    </Card>
                </Col>
                {/* TODO: macos */}
            </Row>
        </React.Fragment>
    );
}
