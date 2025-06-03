import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Col, Row } from "reactstrap";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import { RequestStatus } from "../../api/types";
import { Http } from "../../api/http";
import ReactMarkdown from 'react-markdown';

interface TemplateDetailsSummaryProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsSummary({ template }: TemplateDetailsSummaryProps) {
    const [latestVersion, setLatestVersion] = useState<WorkspaceTemplateVersion>();
    const [readmeContent, setReadmeContent] = useState<string>();

    const FetchTemplateLatestVersion = useCallback(async () => {
        // fetch template versions
        let [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${template.id}/latest-version`,
            "GET",
            null
        );

        if (status === RequestStatus.OK && statusCode === 200) {
            setLatestVersion(responseData as WorkspaceTemplateVersion);
        }
    }, [template.id]);

    useEffect(() => {
        FetchTemplateLatestVersion();
    }, [FetchTemplateLatestVersion]);

    const FetchReadme = useCallback(async () => {
        if (latestVersion) {
            let [status, statusCode, responseData] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${latestVersion.id}/entries/README.md`,
                "GET",
                null
            );

            if (status === RequestStatus.OK && statusCode === 200) {
                setReadmeContent(atob(responseData.content));
            }
        }
    }, [latestVersion, template.id]);

    useEffect(() => {
        FetchReadme();
    }, [FetchReadme]);

    return (
        <React.Fragment>
            {latestVersion ?
                (<React.Fragment>
                    <Row>
                        <Col md={12}>
                            <Card>
                                <CardBody>
                                    <Row>
                                        <Col md={6}>
                                            <span className="text-muted">Latest version:</span> {latestVersion.name}
                                        </Col>
                                        <Col md={6}>
                                            <span className="text-muted">Released on:</span> {new Date(latestVersion.published_on).toLocaleString()}
                                        </Col>
                                    </Row>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                    <Row className="mt-3">
                        <Col md={12}>
                            <Card>
                                <CardHeader className="border-0 pb-0">
                                    <h3>README.md</h3>
                                </CardHeader>
                                <CardBody className="pt-0">
                                    {readmeContent ? (
                                        <ReactMarkdown children={readmeContent} />
                                    ) : <h4>Not available</h4>}
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </React.Fragment>) : <h4 className="ms-3">No active version available</h4>
            }
        </React.Fragment>
    );
}
