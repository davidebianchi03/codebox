import React, { useCallback, useEffect } from "react";
import { Button, Col, FormFeedback, FormGroup, Input, Label, Modal, ModalBody, ModalHeader, Row } from "reactstrap";
import { WorkspaceTemplate } from "../../types/templates";
import { useFormik } from "formik";
import * as Yup from "yup";
import { RequestStatus } from "../../api/types";
import { toast } from "react-toastify";
import { Http } from "../../api/http";

interface TemplateSettingsModalProps {
    isOpen: boolean
    onClose: () => void
    template: WorkspaceTemplate
}

export function TemplateSettingsModal({ isOpen, onClose, template }: TemplateSettingsModalProps) {

    const validation = useFormik({
        initialValues: {
            name: template.name,
        },
        validationSchema: Yup.object({
            name: Yup.string().required("This field is required").test(
                "already_exists",
                "Another template with the same name already exists",
                async (name) => {
                    let [status, statusCode, responseData] = await Http.Request(
                        `${Http.GetServerURL()}/api/v1/templates-by-name/${encodeURIComponent(name)}`,
                        "GET",
                        null
                    );
                    if(status === RequestStatus.OK && statusCode === 200) {
                        var resp = responseData as WorkspaceTemplate;
                        return resp.id === template.id;
                    }

                    return false;
                }
            ),
        }),
        validateOnBlur: false,
        validateOnChange: false,
        onSubmit: async (values) => {
            let [status, statusCode] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/templates/${template.id}`,
                "PUT",
                JSON.stringify({
                    name: values.name,
                })
            );

            if (status === RequestStatus.OK && statusCode === 200) {
                validation.resetForm();
                onClose();
            } else {
                toast.error("Unknown error");
            }
        },
    });

    useEffect(() => {
        validation.setFieldValue("name", template.name);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isOpen]);

    const HandleCloseModal = useCallback(() => {
        validation.resetForm();
        onClose();
    }, [onClose, validation]);

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
                    Settings
                </ModalHeader>
                <ModalBody>
                    <form
                        onSubmit={validation.handleSubmit}
                    >
                        <Row>
                            <Col md={12}>
                                <FormGroup>
                                    <Label>Name</Label>
                                    <Input
                                        name="name"
                                        onChange={validation.handleChange}
                                        value={validation.values.name}
                                        invalid={!!validation.errors.name}
                                    />
                                    <FormFeedback>{validation.errors.name}</FormFeedback>
                                </FormGroup>
                            </Col>
                        </Row>
                        <div className="d-flex justify-content-end">
                            <Button color="accent" onClick={HandleCloseModal}>
                                Cancel
                            </Button>
                            <Button color="primary" className="ms-1" type="submit">
                                Save
                            </Button>
                        </div>
                    </form>
                </ModalBody>
            </Modal>
        </React.Fragment>
    );
}
