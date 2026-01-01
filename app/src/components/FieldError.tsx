import React from "react";
import { Form } from "react-bootstrap";

export interface FieldErrorProps {
    error?: string;
}

export function FieldError({ error }: FieldErrorProps) {
    return (
        <React.Fragment>
            <Form.Control.Feedback type="invalid" style={{ display: error ? 'block' : 'none' }}>
                {error}
            </Form.Control.Feedback>
        </React.Fragment>
    )
}
