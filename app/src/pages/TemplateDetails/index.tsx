import React, { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { WorkspaceTemplate } from "../../types/templates";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { Card, CardBody, CardHeader, Col, Container, Row, Table } from "reactstrap";
import { TemplateDetailsVersions } from "./TemplateDetailsVersions";

export function TemplateDetailsPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [template, setTemplate] = useState<WorkspaceTemplate>();

    const fetchTemplate = useCallback(async () => {
        var [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${id}`,
            "GET",
            null
        );

        if (status === RequestStatus.OK && statusCode === 200) {
            setTemplate(responseData);
        } else if (statusCode === 404) {
            navigate("/templates");
        } else {
            toast.error("Failed to fetch template details");
            setTemplate(undefined);
        }
    }, []);

    useEffect(() => {
        fetchTemplate();
    }, [fetchTemplate]);

    return (
        <React.Fragment>
            {template && (
                <Container className="mt-4 mb-4">
                    <div className="row g-2 align-items-center">
                        <div className="col">
                            <div className="page-pretitle">Templates</div>
                            <h2 className="page-title">{template.name}</h2>
                        </div>
                    </div>
                    <Row className="mt-4">
                        <Col md={12}>
                            <Card>
                                <CardBody>
                                    <Table striped>
                                        <tbody>
                                            <tr>
                                                <th>
                                                    Name
                                                </th>
                                                <td>
                                                    {template.name}
                                                </td>
                                            </tr>
                                            <tr>
                                                <th>
                                                    Description
                                                </th>
                                                <td>
                                                    {template.description}
                                                </td>
                                            </tr>
                                            <tr>
                                                <th>
                                                    Type
                                                </th>
                                                <td>
                                                    {template.type}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </Table>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                    <Row className="mt-4">
                        <Col md={12}>
                            <TemplateDetailsVersions template={template} />
                        </Col>
                    </Row>
                </Container>
            )}
        </React.Fragment>
    );
}