import React, { useCallback, useEffect, useState } from "react";
import { Col, Row } from "reactstrap";
import Editor, { DiffEditor, useMonaco, loader } from '@monaco-editor/react';
import { TemplateVersionEditorSidebar } from "./Sidebar";
import { useNavigate, useParams } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";



export function TemplateVersionEditor() {
    const { templateId, versionId } = useParams();
    const navigate = useNavigate();
    const [template, setTemplate] = useState<WorkspaceTemplate>();
    const [templateVersion, setTemplateVersion] = useState<WorkspaceTemplateVersion>();

    const fetchTemplate = useCallback(async () => {
        let [status, statusCode, responseBody] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${templateId}`,
            "GET",
            null
        );
        if (status === RequestStatus.OK && statusCode === 200) {
            setTemplate(responseBody as WorkspaceTemplate);
        } else {
            navigate(`/templates`);
        }
    }, []);

    const fetchTemplateVersion = useCallback(async () => {
        let [status, statusCode, responseBody] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${templateId}/versions/${versionId}`,
            "GET",
            null
        );
        if (status === RequestStatus.OK && statusCode === 200) {
            setTemplateVersion(responseBody as WorkspaceTemplateVersion);
        } else {
            navigate(`/templates`);
        }
    }, []);

    useEffect(() => {
        fetchTemplate();
        fetchTemplateVersion();
    }, [fetchTemplateVersion, fetchTemplate]);


    return (
        <React.Fragment>
            {template && templateVersion && (
                <React.Fragment>
                    <Row style={{
                        background: "#181818",
                        height: "90vh"
                    }}>
                        <Col md="2">
                            <TemplateVersionEditorSidebar
                                template={template}
                                templateVersion={templateVersion}
                            />
                        </Col>
                        <Col md={10}>
                            <Editor height="90vh" defaultLanguage="javascript" defaultValue="// some comment" theme="vs-dark" />
                        </Col>
                    </Row>
                </React.Fragment >
            )}
        </React.Fragment >
    )
}