import { useFormik } from "formik";
import React, { useCallback, useEffect, useState } from "react";
import { Button, Card, Form } from "react-bootstrap";
import { APIAdminGetAnalyticsConfig, APIAdminUpdateAnalyticsConfig } from "../../api/admin";
import { toast } from "react-toastify";

export function AnalyticsConfig() {
    const [loading, setLoading] = useState<boolean>(true);

    const validation = useFormik({
        initialValues: {
            sendData: true,
        },
        onSubmit: async (values) => {
            if (!await APIAdminUpdateAnalyticsConfig(values.sendData)) {
                toast.error("Failed to update analytics config, try again later");
            } else {
                toast.success("Analytics config has been updated successfully");
                fetchConfig();
            }
        },
    });

    const fetchConfig = useCallback(async () => {
        setLoading(true);
        const config = await APIAdminGetAnalyticsConfig();
        if (config) {
            validation.setValues({
                sendData: config.send_analytics_data,
            })
            setLoading(false);
        } else {
            toast.error("Failed to fetch analytics config, try again later");
        }
    }, []);

    useEffect(() => {
        fetchConfig();
    }, [fetchConfig]);

    return (
        <React.Fragment>
            <Card body>
                <React.Fragment>
                    <Form onSubmit={validation.handleSubmit}>
                        <Form.Group className="d-flex gap-2">
                            <input
                                type="checkbox"
                                checked={validation.values.sendData}
                                onChange={validation.handleChange}
                                id="sendData"
                                name="sendData"
                                className={`form-check-input form-check-input-light ${loading && "placeholder"}`}
                            />
                            <Form.Label for="sendData" style={{ userSelect: "none" }} className={`${loading && "placeholder"}`}>
                                Send Analytics Data
                            </Form.Label>
                        </Form.Group>
                        <Button type="submit" className={`mt-2 ${loading && "placeholder"}`} variant="light">
                            Save
                        </Button>
                    </Form>
                </React.Fragment>
            </Card>
        </React.Fragment>
    )
}
