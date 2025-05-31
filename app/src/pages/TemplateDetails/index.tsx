import React, { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { WorkspaceTemplate } from "../../types/templates";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { Card, CardBody, Col, Container, Dropdown, DropdownItem, DropdownMenu, DropdownToggle, Row, Table } from "reactstrap";
import { TemplateDetailsVersions } from "./TemplateDetailsVersions";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEllipsisVertical } from "@fortawesome/free-solid-svg-icons";
import { Workspace } from "../../types/workspace";
import Swal from "sweetalert2";

export function TemplateDetailsPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [template, setTemplate] = useState<WorkspaceTemplate>();
    const [showActionsDropdown, setShowActionsDropdown] = useState<boolean>(false);

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
    }, [id, navigate]);

    const handleDeleteTemplate = useCallback(async () => {
        // check if there are workspacea that are based on this template
        var [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${id}/workspaces`,
            "GET",
            null
        );

        if (status === RequestStatus.OK && statusCode === 200) {
            var workspaces = responseData as Workspace[];
            if (workspaces.length === 0) {
                if((await Swal.fire({
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
                    [status] = await Http.Request(
                        `${Http.GetServerURL()}/api/v1/templates/${id}`,
                        "DELETE",
                        null
                    );

                    if (status === RequestStatus.OK) {
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
                            ${
                                workspaces.reverse().slice(0, 5).map(
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
    }, [id, navigate]);

    useEffect(() => {
        fetchTemplate();
    }, [fetchTemplate]);

    return (
        <React.Fragment>
            {template && (
                <Container className="mt-4 mb-4">
                    <div className="row g-2 align-items-center justify-content-between">
                        <div className="col">
                            <div className="page-pretitle">Templates</div>
                            <h2 className="page-title">{template.name}</h2>
                        </div>
                        <div className="col d-flex justify-content-end">
                            <Dropdown isOpen={showActionsDropdown} toggle={() => setShowActionsDropdown(!showActionsDropdown)}>
                                <DropdownToggle color="accent">
                                    <FontAwesomeIcon icon={faEllipsisVertical} />
                                </DropdownToggle>
                                <DropdownMenu>
                                    <DropdownItem onClick={handleDeleteTemplate}>
                                        <span className="text-warning">Delete template</span>
                                    </DropdownItem>
                                </DropdownMenu>
                            </Dropdown>
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