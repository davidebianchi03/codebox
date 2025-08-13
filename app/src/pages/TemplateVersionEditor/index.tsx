import React, { useCallback, useEffect, useRef, useState } from "react";
import { Button } from "reactstrap";
import Editor from '@monaco-editor/react';
import { TemplateVersionEditorSidebar } from "./Sidebar";
import { useNavigate, useParams } from "react-router-dom";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import { toast, ToastContainer } from "react-toastify";
import { FileMap, GetTypeForFile } from "./FileType";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGear } from "@fortawesome/free-solid-svg-icons";
import { TemplateVersionSettingsModal } from "./TemplateVersionSettingsModal";
import { APIRetrieveTemplateById, APIRetrieveTemplateVersion, APIRetrieveTemplateVersionEntry, APIUpdateTemplateVersionEntry } from "../../api/templates";

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
    const [publishTemplateVersion, setPublishTemplateVersion] = useState<boolean>(false);
    const timer = useRef<NodeJS.Timeout | null>(null);

    const fetchTemplate = useCallback(async () => {
        if (templateId) {
            const t = await APIRetrieveTemplateById(parseInt(templateId));
            if (t) {
                setTemplate(t);
            } else {
                navigate(`/templates`);
            }
        }
    }, [navigate, templateId]);

    const fetchTemplateVersion = useCallback(async () => {
        if (templateId && versionId) {
            const tv = await APIRetrieveTemplateVersion(parseInt(templateId), parseInt(versionId));
            if (tv) {
                setTemplateVersion(tv);
            } else {
                navigate(`/templates`);
            }
        }
    }, [navigate, templateId, versionId]);


    const UpdateFileContent = useCallback(async (value: string) => {
        if (openFilePath && templateId && versionId) {
            if (await APIUpdateTemplateVersionEntry(
                parseInt(templateId),
                parseInt(versionId),
                openFilePath,
                openFilePath,
                btoa(value)
            )) {
                setEditing(false);
            } else {
                toast.error("Failed to update file content");
            }
        }
    }, [openFilePath, templateId, versionId]);

    const fetchFileContent = useCallback(async () => {
        if (selectedItemPath && templateId && versionId) {
            const entry = await APIRetrieveTemplateVersionEntry(parseInt(templateId), parseInt(versionId), selectedItemPath);
            if (entry) {
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
        timer.current = setTimeout(() => UpdateFileContent(value), 800);
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
                                onDelete={(item) => {
                                    if(selectedItemPath?.startsWith(item) || selectedItemPath === item) {
                                        setSelectedItemPath(null);
                                        setOpenFilePath("");
                                        setFileContent("");
                                    }
                                }}
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
                                        onClick={() => {
                                            setPublishTemplateVersion(false);
                                            setShowEditTemplateVersionModal(true)
                                        }}
                                    >
                                        <FontAwesomeIcon icon={faGear} />
                                    </Button>
                                    <Button
                                        color="success"
                                        size="sm"
                                        className="py-1 px-2 me-2"
                                        onClick={() => {
                                            setPublishTemplateVersion(true);
                                            setShowEditTemplateVersionModal(true);
                                        }}
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
                        setPublishTemplateVersion(false);
                        setShowEditTemplateVersionModal(false);
                    }}
                    template={template}
                    templateVersion={templateVersion}
                    publish={publishTemplateVersion}
                />
            )}

        </React.Fragment >
    )
}