import React, { useCallback, useEffect, useRef, useState } from "react";
import { Button } from "reactstrap";
import Editor from '@monaco-editor/react';
import { TemplateVersionEditorSidebar } from "./Sidebar";
import { useNavigate, useParams } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionEntry } from "../../types/templates";
import { toast, ToastContainer } from "react-toastify";
import { FileMap, GetTypeForFile } from "./FileType";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear } from "@fortawesome/free-solid-svg-icons";
import Swal from "sweetalert2";
import { TemplateVersionSettingsModal } from "./TemplateVersionSettingsModal";

export function TemplateVersionEditor() {
    const { templateId, versionId } = useParams();
    const navigate = useNavigate();
    const [selectedItemPath, setSelectedItemPath] = useState<string | null>(null);
    const [openFilePath, setOpenFilePath] = useState<string>("");
    const [fileContent, setFileContent] = useState<string>("");
    const [template, setTemplate] = useState<WorkspaceTemplate>();
    const [templateVersion, setTemplateVersion] = useState<WorkspaceTemplateVersion>();
    const [openFileType, setOpenFileType] = useState<FileMap>();
    const [editing, setEditing] = useState<boolean>(false);
    const [showEditTemplateVersionModal, setShowEditTemplateVersionModal] = useState<boolean>(false);
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
            } else {
                setEditing(false);
            }
        }
    }, [fileContent, openFilePath, templateId, versionId]);

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
                        if (timer.current) {
                            clearInterval(timer.current);
                        }
                        setFileContent(atob(entry.content));
                        setOpenFilePath(entry.name);
                        setOpenFileType(GetTypeForFile(entry.name));
                    }
                } else {
                    toast.error("Failed to fetch file content");
                    setFileContent("");
                }
            } else {
                setFileContent("");
            }
        }
        setEditing(false);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedItemPath, templateId, versionId]);

    const EditorHandleChange = useCallback(async (value: string) => {
        setFileContent(value);
        setEditing(true);
        if (timer.current) {
            clearTimeout(timer.current);
            timer.current = null;
        }
        timer.current = setTimeout(UpdateFileContent, 800);
    }, [UpdateFileContent]);

    const handlePublishVersion = useCallback(async () => {
        var r = await Swal.fire({
            title: `Publish version`,
            text: "Publish new version of the template",
            input: 'text',
            inputLabel: `Name`,
            inputPlaceholder: `First commit`,
            inputValue: templateVersion?.name,
            showCancelButton: true,
            reverseButtons: true,
            inputValidator: (value) => {
                if (!value) {
                    return 'You need to write something!'
                }
            },
            customClass: {
                confirmButton: "btn btn-primary",
                cancelButton: "btn btn-accent me-1",
                popup: "bg-dark text-light",
            },
            buttonsStyling: false,
            confirmButtonText: "Publish",
        });

        if (r.isConfirmed && r.value) {
            let [status, statusCode] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${templateId}/versions/${versionId}`,
                "PUT",
                JSON.stringify({
                    name: r.value,
                    published: true,
                    config_file_path: templateVersion?.config_file_relative_path,
                })
            );
            if (status === RequestStatus.OK && statusCode === 200) {
                navigate(`/templates/${template?.id}`);
            } else {
                toast.error("Failed to publish version");
            }
        }

    }, [navigate, template?.id, templateId, templateVersion, versionId]);

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
                    <div style={{
                        background: "#1f1f1f",
                        height: "calc(100vh -  62.3125px)"
                    }}
                        className="d-flex align-items-start w-100"
                    >
                        <div style={{ background: "#181818", width: 250, height: "100%" }}>
                            <TemplateVersionEditorSidebar
                                template={template}
                                templateVersion={templateVersion}
                                onSelectionChange={(si) => setSelectedItemPath(si)}
                            />
                        </div>
                        <div className="ps-0 w-100 h-100">
                            <div
                                style={{ fontFamily: "Consolas", background: "#181818", height: "45px" }}
                                className="d-flex align-items-center justify-content-between"
                            >
                                {openFilePath ? (
                                    <React.Fragment>
                                        <div className="d-flex align-items-center ms-2">
                                            <img src={openFileType?.icon} alt="" width={15} height={15} />
                                            <span className="ms-1">{openFilePath}</span>
                                            <span style={{
                                                width: 7,
                                                height: 7,
                                                borderRadius: "100%",
                                                background: editing ? `var(--tblr-yellow)` : `var(--tblr-success)`,
                                                marginLeft: 8
                                            }}></span>
                                        </div>
                                    </React.Fragment>) :
                                    <span></span>
                                }
                                <div>
                                    <Button
                                        color="transparent"
                                        style={{ background: "none", border: "none" }}
                                        className="py-1 px-2 me-2"
                                        onClick={() => setShowEditTemplateVersionModal(true)}
                                    >
                                        <FontAwesomeIcon icon={faGear} />
                                    </Button>
                                    <Button
                                        color="success"
                                        size="sm"
                                        className="py-1 px-2 me-2"
                                        onClick={handlePublishVersion}
                                    >
                                        Publish
                                    </Button>
                                </div>
                            </div>
                            {openFilePath && (
                                <Editor
                                    height={"calc(100% - 45px)"}
                                    language={openFileType?.language}
                                    value={fileContent}
                                    onChange={(value) => EditorHandleChange(value || "")}
                                    theme="vs-dark"
                                />
                            )}
                        </div>
                    </div>
                </React.Fragment >
            )}
            <ToastContainer />
            {template && templateVersion && (
                <TemplateVersionSettingsModal
                    isOpen={showEditTemplateVersionModal}
                    onClose={() => {
                        fetchTemplateVersion();
                        setShowEditTemplateVersionModal(false);
                    }}
                    template={template}
                    templateVersion={templateVersion}
                />
            )}

        </React.Fragment >
    )
}