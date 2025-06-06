import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Col, Row } from "reactstrap";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import ReactMarkdown from 'react-markdown';
import { APIRetrieveTemplateLatestVersion, APIRetrieveTemplateVersionEntry } from "../../api/templates";

interface TemplateDetailsSummaryProps {
    template: WorkspaceTemplate
}

export function TemplateDetailsSummary({ template }: TemplateDetailsSummaryProps) {
    const [latestVersion, setLatestVersion] = useState<WorkspaceTemplateVersion>();
    const [readmeContent, setReadmeContent] = useState<string>();

    const FetchTemplateLatestVersion = useCallback(async () => {
        const lv = await APIRetrieveTemplateLatestVersion(template.id);

        if (lv) {
            setLatestVersion(lv);
        }
    }, [template.id]);

    useEffect(() => {
        FetchTemplateLatestVersion();
    }, [FetchTemplateLatestVersion]);

    const FetchReadme = useCallback(async () => {
        if (latestVersion) {
            const r = await APIRetrieveTemplateVersionEntry(template.id, latestVersion.id, "README.md");
            if (r) {
                setReadmeContent(atob(r.content));
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
