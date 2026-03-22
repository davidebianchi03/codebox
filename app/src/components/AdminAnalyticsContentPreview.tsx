import React, { useCallback, useEffect, useState } from "react";
import { Card } from "react-bootstrap";
import { APIAdminGetAnalyticsPreviewContent } from "../api/admin";
import { toast } from "react-toastify";

export function AdminAnalyticsContentPreview() {
    const [loading, setLoading] = useState<boolean>(true);
    const [content, setContent] = useState<string>("");

    const fetchConfig = useCallback(async () => {
        setLoading(true);
        const c = await APIAdminGetAnalyticsPreviewContent();
        if (c) {
            setContent(c);
            setLoading(false);
        } else {
            toast.error("Failed to fetch analytics preview, try again later");
        }
    }, []);

    useEffect(() => {
        fetchConfig();
    }, [fetchConfig]);

    return (
        <React.Fragment>
            <Card body>
                {loading ? (
                    <React.Fragment>
                        <div className="placeholder mt-3">
                            "api_key": "preview",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "event": "codebox-analytics",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "properties",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "distinct_id": "012345678901234567890123456789012345678901234567890123456789",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "license_type": "community",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "server_version": "version",
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "total_users": 1,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "approved_users": 1,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "total_runners": 1,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "online_runners": 0,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "total_workspaces": 1,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "running_workspaces": 0,
                        </div>
                        <br/>
                        <div className="placeholder mt-3">
                            "total_templates": 0
                        </div>
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <pre>
                            {content}
                        </pre>
                    </React.Fragment>
                )}
            </Card>
        </React.Fragment>
    )
}
