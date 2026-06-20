import React, { useCallback, useEffect, useState } from "react";
import { ContainerPort, Workspace, WorkspaceContainer } from "../../types/workspace";
import { APIDeleteWorkspaceContainerPort, APIListWorkspaceContainerPorts, APICreateWorkspaceContainerPort } from "../../api/workspace";
import { toast, ToastContainer } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrash, faPlus, faArrowUpRightFromSquare } from "@fortawesome/free-solid-svg-icons";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Button, Dropdown, Form, InputGroup } from "react-bootstrap";

export interface ExposedPortsDropdownProps {
    workspace: Workspace;
    container: WorkspaceContainer;
    onChange: () => void;
}

export function ExposedPortsDropdown({
    workspace,
    container,
    onChange,
}: ExposedPortsDropdownProps) {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [containerExposedPorts, setContainerExposedPorts] = useState<ContainerPort[]>([]);
    const [addingPublic, setAddingPublic] = useState<boolean>(false);
    const [addingPrivate, setAddingPrivate] = useState<boolean>(false);

    const FetchSelectedContainerPorts = useCallback(async () => {
        const ports = await APIListWorkspaceContainerPorts(
            workspace.id,
            container.container_name
        );

        if (ports) {
            setContainerExposedPorts(ports);
        } else {
            toast.error("Failed to fetch workspace container ports");
            setContainerExposedPorts([]);
        }
    }, [container.container_name, workspace.id]);

    const handleDeletePort = useCallback(
        async (port: ContainerPort) => {
            if (await APIDeleteWorkspaceContainerPort(workspace.id, container.container_name, port.port_number)) {
                await FetchSelectedContainerPorts();
                onChange();
            } else {
                toast.error("Failed to remove port");
            }
        },
        [workspace.id, container.container_name, FetchSelectedContainerPorts, onChange]
    );

    const createFormikForSection = (isPublic: boolean) => {
        return useFormik({
            initialValues: {
                portNumber: "",
                serviceName: "",
            },
            validateOnChange: false,
            validateOnBlur: false,
            validationSchema: Yup.object({
                portNumber: Yup.number()
                    .required("Port number is required")
                    .min(1, "Min: 1")
                    .max(65535, "Max: 65535")
                    .test(
                        "port_number",
                        "Port already exposed",
                        (value) =>
                            containerExposedPorts.find((port) => port.port_number === value) ===
                            undefined
                    ),
                serviceName: Yup.string()
                    .test(
                        "name_unique",
                        "Name already exists",
                        (value) =>
                            !value ||
                            containerExposedPorts.find(
                                (port) => port.service_name === value
                            ) === undefined
                    ),
            }),
            onSubmit: async (values) => {
                const finalName = values.serviceName || `Port ${values.portNumber}`;
                const p = await APICreateWorkspaceContainerPort(
                    workspace.id,
                    container.container_name,
                    parseInt(values.portNumber),
                    finalName,
                    isPublic
                );
                if (p) {
                    await FetchSelectedContainerPorts();
                    if (isPublic) {
                        setAddingPublic(false);
                    } else {
                        setAddingPrivate(false);
                    }
                    onChange();
                } else {
                    toast.error("Failed to add port");
                }
            },
        });
    };

    const publicFormik = createFormikForSection(true);
    const privateFormik = createFormikForSection(false);

    useEffect(() => {
        if (!addingPublic) {
            publicFormik.resetForm();
        }
    }, [addingPublic, publicFormik]);

    useEffect(() => {
        if (!addingPrivate) {
            privateFormik.resetForm();
        }
    }, [addingPrivate, privateFormik]);

    useEffect(() => {
        FetchSelectedContainerPorts();
    }, [FetchSelectedContainerPorts]);

    const publicPorts = containerExposedPorts.filter((port) => port.public);
    const privatePorts = containerExposedPorts.filter((port) => !port.public);

    return (
        <React.Fragment>
            <Dropdown
                onToggle={() => setIsOpen(!isOpen)}
                show={isOpen}
            >
                <Dropdown.Toggle
                    variant="accent"
                    className="d-flex align-items-center gap-2"
                >
                    Exposed Ports
                    {containerExposedPorts.length > 0 && (
                        <span className="badge bg-light text-dark" style={{ fontSize: "0.7rem" }}>
                            {containerExposedPorts.length}
                        </span>
                    )}
                </Dropdown.Toggle>

                <Dropdown.Menu className="exposed-ports-dropdown" style={{ width: 380, maxHeight: 500, overflowY: "auto" }}>
                    <div className="p-3">
                        {/* Public Ports Section */}
                        <div className="mb-4">
                            <h6 className="text-uppercase text-muted mb-2" style={{ fontSize: "0.75rem", letterSpacing: "0.5px" }}>
                                Public Ports
                            </h6>

                            {publicPorts.length === 0 ? (
                                <p className="text-muted small mb-2">No public ports</p>
                            ) : (
                                <div className="mb-2">
                                    {publicPorts.map((port, index) => (
                                        <div key={index} className="port-row mb-2 d-flex align-items-center justify-content-between">
                                            <div className="flex-grow-1">
                                                <div className="port-name fw-bold" style={{ fontSize: "0.9rem" }}>
                                                    {port.service_name}
                                                </div>
                                                <div className="port-number text-muted small">
                                                    Port {port.port_number}
                                                </div>
                                            </div>
                                            <div className="d-flex gap-1">
                                                {port.port_url && (
                                                    <Button
                                                        variant="link"
                                                        size="sm"
                                                        className="p-0 text-decoration-none"
                                                        as="a"
                                                        href={port.port_url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        title="Open in new tab"
                                                    >
                                                        <FontAwesomeIcon icon={faArrowUpRightFromSquare} className="text-muted" />
                                                    </Button>
                                                )}
                                                <Button
                                                    variant="link"
                                                    size="sm"
                                                    className="p-0 text-danger text-decoration-none"
                                                    onClick={() => handleDeletePort(port)}
                                                    title="Remove port"
                                                >
                                                    <FontAwesomeIcon icon={faTrash} />
                                                </Button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                            {!addingPublic ? (
                                <Button
                                    variant="outline-light"
                                    className="w-100"
                                    onClick={() => setAddingPublic(true)}
                                >
                                    <FontAwesomeIcon icon={faPlus} className="me-1" /> Add Public Port
                                </Button>
                            ) : (
                                <form
                                    onSubmit={(e) => {
                                        e.preventDefault();
                                        publicFormik.handleSubmit();
                                    }}
                                    className="add-port-form"
                                >
                                    <InputGroup className="mb-2">
                                        <Form.Control
                                            type="number"
                                            placeholder="Port (required)"
                                            name="portNumber"
                                            min={1}
                                            max={65535}
                                            value={publicFormik.values.portNumber}
                                            onChange={publicFormik.handleChange}
                                            isInvalid={!!publicFormik.errors.portNumber}
                                            autoFocus
                                        />
                                        <Form.Control
                                            type="text"
                                            placeholder="Name (optional)"
                                            name="serviceName"
                                            value={publicFormik.values.serviceName}
                                            onChange={publicFormik.handleChange}
                                            isInvalid={!!publicFormik.errors.serviceName}
                                        />
                                    </InputGroup>
                                    {publicFormik.errors.portNumber && (
                                        <small className="text-danger d-block mb-2">{publicFormik.errors.portNumber}</small>
                                    )}
                                    {publicFormik.errors.serviceName && (
                                        <small className="text-danger d-block mb-2">{publicFormik.errors.serviceName}</small>
                                    )}
                                    <div className="d-flex gap-2">
                                        <Button
                                            type="submit"
                                            variant="light"
                                            className="flex-grow-1"
                                        >
                                            Add
                                        </Button>
                                        <Button
                                            type="button"
                                            variant="outline-secondary"
                                            className="flex-grow-1"
                                            onClick={() => {
                                                setAddingPublic(false);
                                                publicFormik.resetForm();
                                            }}
                                        >
                                            Cancel
                                        </Button>
                                    </div>
                                </form>
                            )}
                        </div>

                        <hr className="my-3" />

                        {/* Private Ports Section */}
                        <div>
                            <h6 className="text-uppercase text-muted mb-2" style={{ fontSize: "0.75rem", letterSpacing: "0.5px" }}>
                                Private Ports
                            </h6>

                            {privatePorts.length === 0 ? (
                                <p className="text-muted small mb-2">No private ports</p>
                            ) : (
                                <div className="mb-2">
                                    {privatePorts.map((port, index) => (
                                        <div key={index} className="port-row mb-2 d-flex align-items-center justify-content-between">
                                            <div className="flex-grow-1">
                                                <div className="port-name fw-bold" style={{ fontSize: "0.9rem" }}>
                                                    {port.service_name}
                                                </div>
                                                <div className="port-number text-muted small">
                                                    Port {port.port_number}
                                                </div>
                                            </div>
                                            <div className="d-flex gap-1">
                                                {port.port_url && (
                                                    <Button
                                                        variant="link"
                                                        size="sm"
                                                        className="p-0 text-decoration-none"
                                                        as="a"
                                                        href={port.port_url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        title="Open in new tab"
                                                    >
                                                        <FontAwesomeIcon icon={faArrowUpRightFromSquare} className="text-muted" />
                                                    </Button>
                                                )}
                                                <Button
                                                    variant="link"
                                                    size="sm"
                                                    className="p-0 text-danger text-decoration-none"
                                                    onClick={() => handleDeletePort(port)}
                                                    title="Remove port"
                                                >
                                                    <FontAwesomeIcon icon={faTrash} />
                                                </Button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                            {!addingPrivate ? (
                                <Button
                                    variant="outline-light"
                                    className="w-100"
                                    onClick={() => setAddingPrivate(true)}
                                >
                                    <FontAwesomeIcon icon={faPlus} className="me-1" /> Add Private Port
                                </Button>
                            ) : (
                                <form
                                    onSubmit={(e) => {
                                        e.preventDefault();
                                        privateFormik.handleSubmit();
                                    }}
                                >
                                    <InputGroup className="mb-2">
                                        <Form.Control
                                            type="number"
                                            placeholder="Port (required)"
                                            name="portNumber"
                                            min={1}
                                            max={65535}
                                            value={privateFormik.values.portNumber}
                                            onChange={privateFormik.handleChange}
                                            isInvalid={!!privateFormik.errors.portNumber}
                                            autoFocus
                                        />
                                        <Form.Control
                                            type="text"
                                            placeholder="Name (optional)"
                                            name="serviceName"
                                            value={privateFormik.values.serviceName}
                                            onChange={privateFormik.handleChange}
                                            isInvalid={!!privateFormik.errors.serviceName}
                                        />
                                    </InputGroup>
                                    {privateFormik.errors.portNumber && (
                                        <small className="text-danger d-block mb-2">{privateFormik.errors.portNumber}</small>
                                    )}
                                    {privateFormik.errors.serviceName && (
                                        <small className="text-danger d-block mb-2">{privateFormik.errors.serviceName}</small>
                                    )}
                                    <div className="d-flex gap-2">
                                        <Button
                                            type="submit"
                                            variant="light"
                                            className="flex-grow-1"
                                        >
                                            Add
                                        </Button>
                                        <Button
                                            type="button"
                                            variant="outline-secondary"
                                            className="flex-grow-1"
                                            onClick={() => {
                                                setAddingPrivate(false);
                                                privateFormik.resetForm();
                                            }}
                                        >
                                            Cancel
                                        </Button>
                                    </div>
                                </form>
                            )}
                        </div>
                    </div>
                </Dropdown.Menu>
            </Dropdown>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    );
}