import { faEllipsisVertical } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Dropdown, DropdownItem, DropdownMenu, DropdownToggle } from "reactstrap";
import { WorkspaceTemplate } from "../../types/templates";
import { WorkspaceType } from "../../types/workspace";
import Swal from "sweetalert2";
import { toast } from "react-toastify";
import { TemplateSettingsModal } from "./TemplateSettingsModal";
import { User } from "../../types/user";
import { APIListWorkspacesTypes } from "../../api/workspace";
import { APIDeleteTemplate, APIListTemplateVersionsByTemplate, APIListWorkspacesByTemplate, APIRetrieveTemplateById } from "../../api/templates";

export interface TemplateDetailsProps {
    template: WorkspaceTemplate
    user: User
}

export function TemplateDetailsHeader({ template: initialTemplate, user }: TemplateDetailsProps) {
    const [template, setTemplate] = useState<WorkspaceTemplate>(initialTemplate);
    const [showActionsDropdown, setShowActionsDropdown] = useState<boolean>(false);
    const [showSettingsModal, setShowSettingsModal] = useState<boolean>(false);
    const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
    const navigate = useNavigate();

    const fetchWorkspaceTypes = useCallback(async () => {
        const wt = await APIListWorkspacesTypes();
        if (wt) {
            setWorkspaceTypes(wt);
        }
    }, []);

    useEffect(() => {
        fetchWorkspaceTypes();
    }, [fetchWorkspaceTypes]);

    const handleDeleteTemplate = useCallback(async () => {
        // check if there are workspacea that are based on this template
        const workspaces = await APIListWorkspacesByTemplate(initialTemplate.id);

        if (workspaces) {
            if (workspaces.length === 0) {
                if ((await Swal.fire({
                    title: "Delete template",
                    text: `
                    This action cannot be undone.
                    Are you sure you want to proceed?
                    `,
                    icon: "warning",
                    showCancelButton: true,
                    reverseButtons: true,
                    cancelButtonText: "Cancel",
                    confirmButtonText: "Delete",
                    customClass: {
                        popup: "bg-dark text-light",
                        cancelButton: "btn btn-accent",
                        confirmButton: "btn btn-primary",
                    },
                })).isConfirmed) {
                    if (await APIDeleteTemplate(initialTemplate.id)) {
                        navigate("/templates");
                    } else {
                        toast.error("Unknown error");
                    }
                }
            } else {
                Swal.fire({
                    title: "Cannot delete template",
                    html: `
                        <p>
                            The template cannot be deleted because there are some workspaces that are using it.
                            Remove them before.
                        </p>
                        <p>
                            <small>Here is a list of latest created workspace using this template:</small>
                        </p>
                        <ul>
                            ${workspaces.reverse().slice(0, 5).map(
                        (w) => `<p class="mb-0"><small>${w.name} (owner: ${w.user.first_name} ${w.user.last_name})</small></p>`
                    ).join("")
                        }
                        </ul>
                    `,
                    icon: "warning",
                    reverseButtons: true,
                    cancelButtonText: "Cancel",
                    confirmButtonText: "Ok",
                    customClass: {
                        popup: "bg-dark text-light",
                        cancelButton: "btn btn-accent",
                        confirmButton: "btn btn-primary",
                    },
                });
            }
        } else {
            toast.error("unknown error");
        }
    }, [initialTemplate, navigate]);

    const handleEditFiles = useCallback(async () => {
        let versions = await APIListTemplateVersionsByTemplate(template.id);
        if (versions) {
            versions = versions.filter((v) => !v.published).reverse();
            if (versions.length > 0) {
                navigate(`/templates/${initialTemplate.id}/versions/${versions[0].id}/editor`)
            } else {
                toast.error("Unknown error");
            }
        } else {
            toast.error("Unknown error");
        }
    }, [initialTemplate.id, navigate, template.id]);

    const fetchTemplate = useCallback(async () => {
        const t = await APIRetrieveTemplateById(initialTemplate.id);

        if (t) {
            setTemplate(t);
        } else {
            navigate("/templates");
        }
    }, [initialTemplate, navigate]);

    return (
        <React.Fragment>
            <div className="row g-2 align-items-center justify-content-between">
                <div className="col">
                    <div className="page-pretitle">Templates</div>
                    <div className="d-flex mt-5 align-items-center">
                        {
                            template.icon ? (
                                <img
                                    src={template.icon}
                                    style={{
                                        width: 50,
                                        height: 50,
                                        fontSize: 20,
                                        padding: 3,
                                        opacity: 0.5,
                                        borderRadius: 4,
                                    }}
                                    alt="custom template icon"
                                />
                            ) : (
                                <div
                                    style={{
                                        width: 50,
                                        height: 50,
                                        fontSize: 20,
                                        opacity: 0.5,
                                        borderRadius: 4,
                                    }}
                                    className="bg-primary text-white d-flex align-items-center justify-content-center"
                                >
                                    {template.name[0].toUpperCase()}
                                </div>
                            )
                        }
                        <div className="ms-3">
                            <h2 className="mb-1">{template.name}</h2>
                            <h4 className="text-muted mb-0">{workspaceTypes.find(wt => wt.id === template.type)?.name || template.type}</h4>
                        </div>
                    </div>
                    <p className="mt-2">
                        {template.description}
                    </p>
                </div>
                <div className="col d-flex justify-content-end">
                    {(user.is_template_manager || user.is_superuser) && (
                        <Dropdown isOpen={showActionsDropdown} toggle={() => setShowActionsDropdown(!showActionsDropdown)}>
                            <DropdownToggle color="accent">
                                <FontAwesomeIcon icon={faEllipsisVertical} />
                            </DropdownToggle>
                            <DropdownMenu>
                                <DropdownItem onClick={() => setShowSettingsModal(true)}>
                                    Settings
                                </DropdownItem>
                                <DropdownItem onClick={handleEditFiles}>
                                    Edit files
                                </DropdownItem>
                                <DropdownItem onClick={handleDeleteTemplate}>
                                    <span className="text-warning">Delete template</span>
                                </DropdownItem>
                            </DropdownMenu>
                        </Dropdown>
                    )}
                </div>
            </div>
            <TemplateSettingsModal
                isOpen={showSettingsModal}
                onClose={() => {
                    fetchTemplate();
                    setShowSettingsModal(false);
                }}
                template={template}
            />
        </React.Fragment >
    );
}
