import React, { useCallback, useEffect, useRef, useState } from "react";
import { Col, Input, Row } from "reactstrap";
import Editor from '@monaco-editor/react';
import { TemplateVersionEditorSidebar } from "./Sidebar";
import { useNavigate, useParams } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionEntry } from "../../types/templates";
import { toast } from "react-toastify";

export function TemplateVersionEditor() {
    const { templateId, versionId } = useParams();
    const navigate = useNavigate();
    const [selectedItemPath, setSelectedItemPath] = useState<string | null>(null);
    const [openFilePath, setOpenFilePath] = useState<string>("");
    const [fileContent, setFileContent] = useState<string>("");
    const [template, setTemplate] = useState<WorkspaceTemplate>();
    const [templateVersion, setTemplateVersion] = useState<WorkspaceTemplateVersion>();
    const timer = useRef<NodeJS.Timeout | null>(null);

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
    }, [navigate, templateId]);

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
    }, [navigate, templateId, versionId]);

    const fetchFileContent = useCallback(async () => {
        if (selectedItemPath) {
            let [status, statusCode, responseBody] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${templateId}/versions/${versionId}/entries/${encodeURIComponent(selectedItemPath)}`,
                "GET",
                null
            );
            if (status === RequestStatus.OK && statusCode === 200) {
                var entry = responseBody as WorkspaceTemplateVersionEntry;
                if (entry) {
                    if (entry.type === "file") {
                        // save previous content
                        if(timer.current) {
                            clearInterval(timer.current);
                        }
                        UpdateFileContent();
                        setFileContent(atob(entry.content));
                        setOpenFilePath(entry.name);
                    }
                } else {
                    toast.error("Failed to fetch file content");
                    setFileContent("");
                }
            } else {
                setFileContent("");
            }
        }
    }, [selectedItemPath, templateId, versionId]);


    const UpdateFileContent = useCallback(async () => {
        if (openFilePath) {
            let [status] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${templateId}/versions/${versionId}/entries/${encodeURIComponent(openFilePath)}`,
                "PUT",
                JSON.stringify({
                    path: openFilePath,
                    content: btoa(fileContent),
                })
            );

            if (status !== RequestStatus.OK) {
                toast.error("Failed to update file content");
            }
        }
    }, [fileContent, openFilePath, templateId, versionId])

    const EditorHandleChange = useCallback(async (value: string) => {
        setFileContent(value);
        if (timer.current) {
            clearTimeout(timer.current);
        }
        timer.current = setTimeout(UpdateFileContent, 800);
    }, [UpdateFileContent]);

    useEffect(() => {
        fetchTemplate();
        fetchTemplateVersion();
    }, [fetchTemplateVersion, fetchTemplate]);

    useEffect(() => {
        fetchFileContent();
    }, [fetchFileContent]);


    return (
        <React.Fragment>
            {template && templateVersion && (
                <React.Fragment>
                    <Row style={{
                        background: "#1f1f1f",
                        height: "90vh"
                    }}>
                        <Col md="2" style={{ background: "#181818" }}>
                            <TemplateVersionEditorSidebar
                                template={template}
                                templateVersion={templateVersion}
                                onSelectionChange={(si) => setSelectedItemPath(si)}
                            />
                        </Col>
                        <Col md={10} className="ps-0">
                            {
                                openFilePath && (
                                    <React.Fragment>
                                        <div
                                            style={{ fontFamily: "Consolas", background: "#181818", height: "45px" }}
                                            className="d-flex"
                                        >
                                            <Input
                                                value={openFilePath}
                                                disabled
                                                style={{ maxWidth: 250 }}
                                            />
                                        </div>
                                        <Editor
                                            height="90vh"
                                            language="dockerfile"
                                            value={fileContent}
                                            onChange={(value) => EditorHandleChange(value || "")}
                                            theme="vs-dark"
                                        />
                                    </React.Fragment>
                                )
                            }
                        </Col>
                    </Row>
                </React.Fragment >
            )}
        </React.Fragment >
    )
}