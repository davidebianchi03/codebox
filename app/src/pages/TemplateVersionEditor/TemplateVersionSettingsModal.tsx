import { useFormik } from "formik";
import React, { useCallback, useEffect } from "react";
import { Button, Col, FormFeedback, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Row } from "reactstrap";
import { WorkspaceTemplate, WorkspaceTemplateVersion } from "../../types/templates";
import * as Yup from "yup";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";

interface TemplateVersionSettingsModalProps {
    isOpen: boolean;
    onClose: () => void;
    template: WorkspaceTemplate;
    templateVersion: WorkspaceTemplateVersion;
}

export function TemplateVersionSettingsModal({
    isOpen,
    onClose,
    template,
    templateVersion,
}: TemplateVersionSettingsModalProps) {

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
            let [status, statusCode] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}/versions/${templateVersion.id}`,
                "PUT",
                JSON.stringify({
                    name: values.name,
                    published: templateVersion.published,
                    config_file_path: values.configPath,
                })
            );
            if (status === RequestStatus.OK && statusCode === 200) {
                HandleCloseModal();
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
                    Edit template version
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
                            <Button className="ms-1" color="primary" type="submit">
                                Save
                            </Button>
                        </div>
                    </form>
                </ModalBody>
            </Modal>
        </React.Fragment>
    );
}
