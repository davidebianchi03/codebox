import { useFormik } from "formik";
import React, { useCallback, useEffect } from "react";
import { Alert, Button, Col, FormFeedback, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Row } from "reactstrap";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import * as Yup from "yup";
import { toast } from "react-toastify";
import { useNavigate } from "react-router-dom";
import { APIUpdateTemplateVersion } from "../../api/templates";

interface TemplateVersionSettingsModalProps {
    isOpen: boolean;
    onClose: () => void;
    template: WorkspaceTemplate;
    templateVersion: WorkspaceTemplateVersion;
    publish: boolean;
}

export function TemplateVersionSettingsModal({
    isOpen,
    onClose,
    template,
    templateVersion,
    publish,
}: TemplateVersionSettingsModalProps) {

    const navigate = useNavigate();

    const validation = useFormik({
        initialValues: {
            name: templateVersion.name,
            configPath: templateVersion.config_file_relative_path,
        },
        validationSchema: Yup.object({
            name: Yup.string().required("This field is required"),
            configPath: Yup.string().required("This field is required"),
        }),
        validateOnBlur: false,
        validateOnChange: false,
        onSubmit: async (values) => {
            if (
                await APIUpdateTemplateVersion(
                    template.id,
                    templateVersion.id,
                    values.name,
                    values.configPath,
                    templateVersion.published || publish,
                )
            ) {
                if (templateVersion.published || publish) {
                    navigate(`/templates/${template.id}`);
                } else {
                    HandleCloseModal();
                }
            } else {
                toast.error("Failed to update template version details");
            }
        }
    });

    useEffect(() => {
        validation.setValues({
            name: templateVersion.name,
            configPath: templateVersion.config_file_relative_path,
        });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [templateVersion])

    const HandleCloseModal = useCallback(() => {
        validation.setValues({
            name: templateVersion.name,
            configPath: templateVersion.config_file_relative_path,
        });
        onClose();
    }, [onClose, templateVersion.config_file_relative_path, templateVersion.name, validation]);

    return (
        <React.Fragment>
            <Modal
                isOpen={isOpen}
                toggle={HandleCloseModal}
                centered
                modalClassName="modal-blur"
                fade
            >
                <ModalHeader toggle={HandleCloseModal}>
                    {publish ? "Publish" : "Edit"} template version
                </ModalHeader>
                <ModalBody>
                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            validation.handleSubmit();
                            return false;
                        }}
                    >
                        {publish && (
                            <Row>
                                <Col md={12}>
                                    <Alert className="border-accent">
                                        <h4 className="mb-0">
                                            You're ready to release an updated version of this template.
                                        </h4>
                                    </Alert>
                                </Col>
                            </Row>
                        )}
                        <Row>
                            <Col md={12}>
                                <FormGroup>
                                    <Label>Name*</Label>
                                    <Input
                                        placeholder="Name"
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
                        </Row>
                        <Row>
                            <Col md={12}>
                                <FormGroup>
                                    <Label>Config path*</Label>
                                    <Input
                                        placeholder="e.g. docker-compose.yml, main.tf"
                                        name="configPath"
                                        value={validation.values.configPath}
                                        onChange={validation.handleChange}
                                        invalid={!!validation.errors.configPath}
                                    />
                                    <FormFeedback>
                                        {validation.errors.configPath}
                                    </FormFeedback>
                                </FormGroup>
                            </Col>
                        </Row>
                        <div className="d-flex justify-content-end">
                            <Button
                                color="accent"
                                onClick={(e) => {
                                    e.preventDefault();
                                    HandleCloseModal();
                                    return false;
                                }}
                            >
                                Cancel
                            </Button>
                            <Button className="ms-1" color="light" type="submit">
                                {publish ? "Publish" : "Save"}
                            </Button>
                        </div>
                    </form>
                </ModalBody>
            </Modal>
        </React.Fragment>
    );
}
