import React, { useCallback, useEffect, useState } from "react";
import { Button, Dropdown, DropdownMenu, DropdownToggle, Input, InputGroup } from "reactstrap";
import { ContainerPort, Workspace, WorkspaceContainer } from "../../types/workspace";
import { APIDeleteWorkspaceContainerPort, APIListWorkspaceContainerPorts } from "../../api/workspace";
import { toast, ToastContainer } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faTrash } from "@fortawesome/free-solid-svg-icons";
import { EditExposedPortsAddPortModal } from "./EditExposedPortsAddPortModal";

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
    const [containerExposedPorts, setContainerExposedPorts] = useState<
        ContainerPort[]
    >([]);
    const [showAddPortModal, setShowAddPortModal] = useState<boolean>(false);

    const FetchSelectedContainerPorts = useCallback(async () => {
        const ports = await APIListWorkspaceContainerPorts(
            workspace.id,
            container.container_name
        )

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
                FetchSelectedContainerPorts();
                onChange();
            } else {
                toast.error("Failed to remove port");
            }
        },
        [workspace.id, container.container_name, FetchSelectedContainerPorts, onChange]
    );

    useEffect(() => {
        FetchSelectedContainerPorts();
    }, [FetchSelectedContainerPorts])


    return (
        <React.Fragment>
            <Dropdown
                toggle={() => setIsOpen(!isOpen)}
                isOpen={isOpen}
            >
                <DropdownToggle
                    color="accent"
                >
                    Edit exposed ports
                </DropdownToggle>
                <DropdownMenu
                    style={{ width: 350 }}
                    className="p-3"
                >
                    <h3 className="mb-3">Exposed ports</h3>
                    <p className="text-muted text-uppercase" style={{ fontSize: 12 }}>
                        Public ports
                    </p>
                    {containerExposedPorts.filter((port) => port.public).length === 0 ? (
                        <React.Fragment>
                            <small className="text-muted">
                                There are no public exposed ports
                            </small>
                        </React.Fragment>
                    ) : (
                        <React.Fragment>
                            {containerExposedPorts.filter((port) => port.public).map((port, index) => (
                                <React.Fragment key={index}>
                                    <InputGroup className="py-1">
                                        <Input
                                            value={port.service_name}
                                            disabled
                                        />
                                        <Input
                                            value={port.port_number}
                                            readOnly
                                        />
                                        <Button
                                            color="accent"
                                            onClick={() => handleDeletePort(port)}
                                        >
                                            <FontAwesomeIcon icon={faTrash} />
                                        </Button>
                                    </InputGroup>
                                </React.Fragment>
                            ))}
                        </React.Fragment>
                    )}
                    <p className="text-muted text-uppercase mt-3" style={{ fontSize: 12 }}>
                        Private ports
                    </p>
                    {containerExposedPorts.filter((port) => !port.public).length === 0 ? (
                        <React.Fragment>
                            <small className="text-muted">
                                There are no private exposed ports
                            </small>
                        </React.Fragment>
                    ) : (
                        <React.Fragment>
                            {containerExposedPorts.filter((port) => !port.public).map((port, index) => (
                                <React.Fragment key={index}>
                                    <InputGroup className="py-1">
                                        <Input
                                            value={port.service_name}
                                            disabled
                                        />
                                        <Input
                                            value={port.port_number}
                                            readOnly
                                        />
                                        <Button
                                            color="accent"
                                            onClick={() => handleDeletePort(port)}
                                        >
                                            <FontAwesomeIcon icon={faTrash} />
                                        </Button>
                                    </InputGroup>
                                </React.Fragment>
                            ))}
                        </React.Fragment>
                    )}
                    <Button
                        className="w-100 mt-5"
                        color="accent"
                        onClick={() => setShowAddPortModal(true)}
                    >
                        Expose a port
                    </Button>
                </DropdownMenu>
            </Dropdown>
            <ToastContainer
                toastClassName={"bg-dark"}
            />
            <EditExposedPortsAddPortModal
                container={container}
                workspace={workspace}
                isOpen={showAddPortModal}
                onClose={() => {
                    setShowAddPortModal(false);
                    FetchSelectedContainerPorts();
                    onChange();
                }}
            />
        </React.Fragment>
    )
}