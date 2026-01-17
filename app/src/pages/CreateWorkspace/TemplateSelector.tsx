import React, { ChangeEventHandler } from "react";
import { Form } from "react-bootstrap";
import { WorkspaceTemplate } from "../../types/templates";

export interface TemplateSelectorProps {
    value: string;
    onChange?: ChangeEventHandler<HTMLSelectElement> | undefined;
    isInvalid?: boolean;
    error?: string;
    templates: WorkspaceTemplate[];
}

export function TemplateSelector({
    value,
    onChange,
    isInvalid,
    error,
    templates,
}: TemplateSelectorProps) {
    return (
        <React.Fragment>
            <Form.Group className="mt-2">
                <Form.Label>Template</Form.Label>
                <select
                    className={`form-control ${isInvalid && "invalid"}`}
                    name="template"
                    value={value}
                    onChange={onChange}
                >
                    <option value={""}>Select a template</option>
                    {
                        templates.map((template, index) => (
                            <option key={index} value={template.id}>{template.name}</option>
                        ))
                    }
                </select>
                <Form.Control.Feedback>
                    {error}
                </Form.Control.Feedback>
            </Form.Group>
        </React.Fragment>
    )
}
