import React, { useCallback, useEffect, useState } from "react";
import { ContainerPort, Workspace, WorkspaceContainer } from "../../types/workspace";
import { APIDeleteWorkspaceContainerPort, APIListWorkspaceContainerPorts, APICreateWorkspaceContainerPort } from "../../api/workspace";
import { toast, ToastContainer } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faXmark, faArrowUpRightFromSquare, faLock, faNetworkWired } from "@fortawesome/free-solid-svg-icons";
import { useFormik } from "formik";
import * as Yup from "yup";
import { Button, Dropdown, Form } from "react-bootstrap";

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
    const [containerExposedPorts, setContainerExposedPorts] = useState<ContainerPort[]>([]);

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

    const portFormik = useFormik({
        initialValues: {
            portNumber: "",
            serviceName: "",
            visibility: "public",
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
        onSubmit: async (values, helpers) => {
            const finalName = values.serviceName || `Port ${values.portNumber}`;
            const p = await APICreateWorkspaceContainerPort(
                workspace.id,
                container.container_name,
                parseInt(values.portNumber),
                finalName,
                values.visibility === "public"
            );
            if (p) {
                await FetchSelectedContainerPorts();
                helpers.resetForm();
                onChange();
            } else {
                toast.error("Failed to add port");
            }
        },
    });

    useEffect(() => {
        FetchSelectedContainerPorts();
    }, [FetchSelectedContainerPorts]);

    const publicPorts = containerExposedPorts.filter((port) => port.public);
    const privatePorts = containerExposedPorts.filter((port) => !port.public);

    return (
        <React.Fragment>
            <Dropdown align="end">
                <Dropdown.Toggle
                    variant="accent"
                    className="d-flex align-items-center gap-2 text-decoration-none"
                >
                    Open ports
                    <span className="badge rounded-pill bg-accent">{containerExposedPorts.length}</span>
                </Dropdown.Toggle>

                <Dropdown.Menu className="p-0 shadow overflow-auto" style={{ minWidth: 370 }}>
                    <div className="p-3">
                        <div className="pb-3 border-bottom">
                            <h5 className="mb-2">Exposed Ports</h5>
                            <p className="text-muted mb-0">
                                Manage ports exposed by this container and how they can be accessed.
                            </p>
                        </div>

                        <section className="py-3 border-bottom">
                            <div className="mb-2">
                                <h6 className="mb-1">Public ports</h6>
                                <p className="text-muted mb-0">Accessible externally</p>
                            </div>

                            {publicPorts.length === 0 ? (
                                <p className="text-muted small py-2 mb-0">No public ports exposed</p>
                            ) : (
                                <div>
                                    {publicPorts.map((port, index) => (
                                        <div
                                            key={port.port_number}
                                            className={`d-flex align-items-center gap-2 py-2`}
                                        >
                                            <FontAwesomeIcon icon={faNetworkWired} className="text-secondary" />
                                            <span className="fw-semibold text-nowrap" style={{ width: 60 }}>{port.port_number}</span>
                                            <span className="fw-semibold text-truncate">{port.service_name}</span>
                                            <div className="d-flex gap-1 ms-auto">
                                                {port.port_url && (
                                                    <Button
                                                        variant="link"
                                                        size="sm"
                                                        className="text-secondary p-1"
                                                        as="a"
                                                        href={port.port_url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        title="Open in new tab"
                                                        aria-label={`Open ${port.service_name} in a new tab`}
                                                    >
                                                        <FontAwesomeIcon icon={faArrowUpRightFromSquare} />
                                                    </Button>
                                                )}
                                                <Button
                                                    variant="link"
                                                    size="sm"
                                                    className="text-secondary p-1"
                                                    onClick={() => handleDeletePort(port)}
                                                    title="Remove port"
                                                    aria-label={`Remove ${port.service_name}`}
                                                >
                                                    <FontAwesomeIcon icon={faXmark} />
                                                </Button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                        </section>

                        <section className="py-3 border-bottom">
                            <div className="mb-2">
                                <h6 className="mb-1">Private ports</h6>
                                <p className="text-muted mb-0">Available inside the workspace</p>
                            </div>

                            {privatePorts.length === 0 ? (
                                <p className="text-muted small py-2 mb-0">No private ports exposed</p>
                            ) : (
                                <div>
                                    {privatePorts.map((port, index) => (
                                        <div
                                            key={port.port_number}
                                            className={`d-flex align-items-center gap-2 py-2`}
                                        >
                                            <FontAwesomeIcon icon={faLock} className="text-secondary" />
                                            <span className="fw-semibold text-nowrap" style={{ width: 60 }}>{port.port_number}</span>
                                            <span className="fw-semibold text-truncate">{port.service_name}</span>
                                            <div className="d-flex gap-1 ms-auto">
                                                {port.port_url && (
                                                    <Button
                                                        variant="link"
                                                        size="sm"
                                                        className="text-secondary p-1"
                                                        as="a"
                                                        href={port.port_url}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        title="Open in new tab"
                                                        aria-label={`Open ${port.service_name} in a new tab`}
                                                    >
                                                        <FontAwesomeIcon icon={faArrowUpRightFromSquare} />
                                                    </Button>
                                                )}
                                                <Button
                                                    variant="link"
                                                    size="sm"
                                                    className="text-secondary p-1"
                                                    onClick={() => handleDeletePort(port)}
                                                    title="Remove port"
                                                    aria-label={`Remove ${port.service_name}`}
                                                >
                                                    <FontAwesomeIcon icon={faXmark} />
                                                </Button>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}

                        </section>

                        <form
                            onSubmit={(e) => {
                                e.preventDefault();
                                portFormik.handleSubmit();
                            }}
                            className="mt-3 p-3 border rounded"
                        >
                            <div className="d-grid gap-2 mb-2">
                                <Form.Control
                                    type="number"
                                    placeholder="Port"
                                    name="portNumber"
                                    min={1}
                                    max={65535}
                                    value={portFormik.values.portNumber}
                                    onChange={portFormik.handleChange}
                                    isInvalid={!!portFormik.errors.portNumber}
                                />
                                <Form.Control
                                    type="text"
                                    placeholder="Service name (optional)"
                                    name="serviceName"
                                    value={portFormik.values.serviceName}
                                    onChange={portFormik.handleChange}
                                    isInvalid={!!portFormik.errors.serviceName}
                                />
                                <Form.Select
                                    name="visibility"
                                    value={portFormik.values.visibility}
                                    onChange={portFormik.handleChange}
                                    aria-label="Port visibility"
                                >
                                    <option value="public">Public - accessible externally</option>
                                    <option value="private">Private - workspace only</option>
                                </Form.Select>
                            </div>
                            {portFormik.errors.portNumber && (
                                <small className="text-danger d-block mb-2">{portFormik.errors.portNumber}</small>
                            )}
                            {portFormik.errors.serviceName && (
                                <small className="text-danger d-block mb-2">{portFormik.errors.serviceName}</small>
                            )}
                            <Button
                                type="submit"
                                variant="light"
                                className="w-100"
                            >
                                Add port
                            </Button>
                        </form>
                    </div>
                </Dropdown.Menu>
            </Dropdown>
            <ToastContainer toastClassName={"bg-dark"} />
        </React.Fragment>
    );
}