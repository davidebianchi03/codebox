import { useFormik } from "formik";
import React, { useCallback, useEffect, useState } from "react";
import { Button, Col, FormFeedback, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Row } from "reactstrap";
import * as Yup from "yup";
import { WorkspaceType } from "../../types/workspace";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { WorkspaceTemplate } from "../../types/templates";
import { useNavigate } from "react-router-dom";

interface CreateTemplateModalProps {
    isOpen: boolean
    onClose: () => void
}

export function CreateTemplateModal({ isOpen, onClose }: CreateTemplateModalProps) {

    const navigate = useNavigate();
    const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
    const validation = useFormik({
        initialValues: {
            name: "",
            type: "",
            icon: "",
            description: ""
        },
        validationSchema: Yup.object({
            name: Yup.string().required("This field is required"),
            type: Yup.string().required("This field is required"),
        }),
        validateOnBlur: false,
        validateOnChange: false,
        onSubmit: async (values) => {
            let [status, statusCode, responseData] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates`,
                "POST",
                JSON.stringify({
                    name: values.name,
                    type: values.type,
                    description: values.description,
                    icon: values.icon
                })
            );

            if(status !== RequestStatus.OK || statusCode !== 201) {
                toast.error("Failed to create template");
            } else {
                var template = responseData as WorkspaceTemplate;
                navigate(`/templates/${template.id}`);
            }
        }
    });

    const fetchWorkspaceTypes = useCallback(async () => {
        let [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/workspace-types`,
            "GET",
            null
        );
        if (status === RequestStatus.OK) {
            setWorkspaceTypes((responseData as WorkspaceType[]).filter((wt) => (
                wt.supported_config_sources.find((scs) => scs === "template") !== undefined))
            );
        } else {
            console.log(`Error: received ${statusCode} from server`);
        }
    }, []);

    const handleCloseModal = useCallback(() => {
        validation.resetForm();
        onClose();
    }, [onClose, validation]);

    useEffect(() => {
        fetchWorkspaceTypes();
    }, [fetchWorkspaceTypes]);

    return (
        <React.Fragment>
            <Modal
                isOpen={isOpen}
                toggle={handleCloseModal}
                centered
                size="lg"
                modalClassName="modal-blur"
                fade
            >
                <ModalHeader toggle={handleCloseModal}>
                    Create template
                </ModalHeader>
                <ModalBody>
                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            validation.handleSubmit();
                            return false;
                        }}
                    >
                        <Row>
                            <Col md={6}>
                                <FormGroup>
                                    <Label>Name *</Label>
                                    <Input
                                        placeholder="my awesome template"
                                        name="name"
                                        value={validation.values.name}
                                        onChange={validation.handleChange}
                                        invalid={!!validation.errors.name}
                                    />
                                    <FormFeedback>
                                        {validation.errors.name}
                                    </FormFeedback>
                                </FormGroup>
                            </Col>
                            <Col md={6}>
                                <FormGroup>
                                    <Label>Type *</Label>
                                    <Input
                                        type="select"
                                        name="type"
                                        value={validation.values.type}
                                        onChange={validation.handleChange}
                                        invalid={!!validation.errors.type}
                                    >
                                        <option value="">Select a workspace type</option>
                                        {workspaceTypes.map((wt) => (
                                            <option value={wt.id}>{wt.name}</option>
                                        ))}
                                    </Input>
                                    <FormFeedback>
                                        {validation.errors.type}
                                    </FormFeedback>
                                </FormGroup>
                            </Col>
                        </Row>
                        <Row>
                            <Col md={12}>
                                <FormGroup>
                                    <Label>Icon</Label>
                                    <Input
                                        placeholder="https://www.youtube.com/watch?v=dQw4w9WgXcQ"
                                        name="icon"
                                        value={validation.values.icon}
                                        onChange={validation.handleChange}
                                        invalid={!!validation.errors.icon}
                                    />
                                    <FormFeedback>
                                        {validation.errors.icon}
                                    </FormFeedback>
                                </FormGroup>
                            </Col>
                        </Row>
                        <Row>
                            <Col md={12}>
                                <FormGroup>
                                    <Label>Description</Label>
                                    <Input
                                        type="textarea"
                                        name="description"
                                        value={validation.values.description}
                                        onChange={validation.handleChange}
                                        invalid={!!validation.errors.description}
                                    />
                                    <FormFeedback>
                                        {validation.errors.description}
                                    </FormFeedback>
                                </FormGroup>
                            </Col>
                        </Row>
                        <div className="d-flex justify-content-end">
                            <Button
                                color="light"
                                className="me-1"
                                onClick={(e) => {
                                    e.preventDefault();
                                    handleCloseModal();
                                }}
                            >
                                Cancel
                            </Button>
                            <Button color="primary" type="submit">
                                Create
                            </Button>
                        </div>
                    </form>
                </ModalBody>
            </Modal>
        </React.Fragment>
    );
}