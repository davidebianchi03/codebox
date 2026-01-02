import React from "react";
import { Card, CardBody, Col, Row } from "reactstrap";

export function SettingsFormPlaceholder() {
    return (
        <React.Fragment>
            <Row className="mt-4">
                <Col>
                    <Card body>
                        <h4 className="placeholder" style={{ height: 40, maxWidth: 300 }}>
                            Title
                        </h4>
                        <div className="placeholder">
                            Signup open
                        </div>
                        <div className="placeholder mt-3">
                            Restricted
                        </div>
                        <div className="placeholder mt-3" style={{ height: 100 }}>
                            Allowed regex
                        </div>
                        <div className="placeholder mt-3" style={{ height: 100 }}>
                            Blocked regex
                        </div>
                    </Card>
                </Col>
            </Row>
            <Row className="mt-4">
                <Col>
                    <Card className="mt-3">
                        <CardBody className="d-flex justify-content-end">
                            <div className="disabled placeholder me-2" style={{ width: 100, height: 40 }}>
                                &nbsp;
                            </div>
                            <div className="bg-primary disabled placeholder" style={{ width: 80, height: 40 }}>
                                &nbsp;
                            </div>
                        </CardBody>
                    </Card>
                </Col>
            </Row>
        </React.Fragment>
    )
}
