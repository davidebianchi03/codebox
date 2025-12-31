import React from "react";
import { Card, Col, Row } from "reactstrap";

export function UserDetailsPlaceholder() {
    return (
        <React.Fragment>
            <Row className="mt-4">
                <Col>
                    <Card body>
                        <div className="placeholder mt-3" style={{ width: 300, height: 25 }} />
                        <div className="placeholder mt-3" style={{ width: 300, height: 15 }} />
                    </Card>
                </Col>
            </Row>
            <Row className="mt-4">
                <Col md={6}>
                    <Card body>
                        <div className="placeholder mt-3" style={{ width: 300, height: 25 }} />
                        <div className="placeholder mt-3" style={{ width: 200, height: 20 }} />
                        <div className="placeholder mt-3" style={{ width: "100%", height: 35 }} />
                        <div className="placeholder mt-5" style={{ width: 200, height: 20 }} />
                        <div className="placeholder mt-3" style={{ width: "100%", height: 35 }} />
                        <div className="placeholder mt-5" style={{ width: 200, height: 20 }} />
                        <div className="placeholder mt-3" style={{ width: "100%", height: 35 }} />
                        <div className="mt-5">
                            <div className="placeholder mt-3" style={{ width: 20, height: 20 }} />
                            <div className="placeholder mt-3 ms-2" style={{ width: 150, height: 20 }} />
                        </div>
                        <div className="placeholder mt-5" style={{ width: 200, height: 20 }} />
                        <div>
                            <div className="placeholder mt-3" style={{ width: 20, height: 20 }} />
                            <div className="placeholder mt-3 ms-2" style={{ width: 150, height: 20 }} />
                        </div>
                        <div className="placeholder mt-5" style={{ width: 200, height: 20 }} />
                        <div>
                            <div className="placeholder mt-3" style={{ width: 20, height: 20 }} />
                            <div className="placeholder mt-3 ms-2" style={{ width: 150, height: 20 }} />
                            <div className="placeholder mt-3 ms-3" style={{ width: 20, height: 20 }} />
                            <div className="placeholder mt-3 ms-2" style={{ width: 150, height: 20 }} />
                        </div>
                        <div className="d-flex justify-content-end">
                            <div className="placeholder mt-3" style={{ width: 75, height: 35 }} />
                            <div className="placeholder mt-3 ms-2" style={{ width: 75, height: 35 }} />
                        </div>
                    </Card>
                </Col>
                <Col md={6}>
                    <Row>
                        <Col>
                            <Card body>
                                <div className="placeholder mt-3" style={{ width: 300, height: 25 }} />
                                <div className="placeholder mt-3" style={{ width: 200, height: 20 }} />
                                <div className="placeholder mt-3" style={{ width: "100%", height: 35 }} />
                                <div className="placeholder mt-5" style={{ width: 200, height: 20 }} />
                            </Card>
                        </Col>
                    </Row>
                    <Row className="mt-3">
                        <Col>
                            <Card body>
                                <div className="placeholder mt-3" style={{ width: 300, height: 25 }} />
                                <div>
                                    <div className="placeholder mt-3" style={{ width: 130, height: 35 }} />
                                    <div className="placeholder mt-3 ms-3" style={{ width: 130, height: 35 }} />
                                    <div className="placeholder mt-3 ms-3" style={{ width: 130, height: 35 }} />
                                </div>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </React.Fragment>
    )
}
