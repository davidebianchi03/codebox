import React from "react";
import { Card, CardBody, Col, Row } from "reactstrap";

export function UserListCardPlaceholder() {
    return (
        <React.Fragment>
            <div className="mt-4 d-flex flex-column">
                <h4 className="placeholder" style={{ height: 30, maxWidth: 150 }}>
                    Title
                </h4>
                <div className="placeholder" style={{ maxWidth: 300 }}>
                    Signup open
                </div>
            </div>
            <Row className="mt-4">
                <Col>
                    <Card body>
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
        </React.Fragment>
    )
}
