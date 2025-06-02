import React, { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { WorkspaceTemplate } from "../../types/templates";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast, ToastContainer } from "react-toastify";
import { Col, Container, Row } from "reactstrap";
import { TemplateDetailsVersions } from "./TemplateDetailsVersions";
import { TemplateDetailsHeader } from "./TemplateDetailsHeader";
import { TemplateDetailsSummary } from "./TemplateDetailsSummary";
import { TemplateDetailsWorkspaces } from "./TemplateDetailsWorkspaces";

export function TemplateDetailsPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [template, setTemplate] = useState<WorkspaceTemplate>();
    const [selectedTab, setSelectedTab] = useState<number>(0);

    const fetchTemplate = useCallback(async () => {
        var [status, statusCode, responseData] = await Http.Request(
            `${Http.GetServerURL()}/api/v1/templates/${id}`,
            "GET",
            null
        );

        if (status === RequestStatus.OK && statusCode === 200) {
            setTemplate(responseData);
        } else if (statusCode === 404) {
            navigate("/templates");
        } else {
            toast.error("Failed to fetch template details");
            setTemplate(undefined);
        }
    }, [id, navigate]);


    useEffect(() => {
        fetchTemplate();
    }, [fetchTemplate]);

    return (
        <React.Fragment>
            {template && (
                <Container className="mt-4 mb-4">
                    {template && (
                        <TemplateDetailsHeader
                            template={template}
                        />
                    )}
                    <Row className="mt-4">
                        <Col md={12}>
                            <header className="navbar-expand-md">
                                <div className="collapse navbar-collapse" id="navbar-menu">
                                    <div className="navbar border" style={{ borderRadius: 7 }}>
                                        <div className="container-xl">
                                            <div className="row flex-column flex-md-row flex-fill align-items-center">
                                                <div className="col">
                                                    <ul className="navbar-nav">
                                                        <li
                                                            className={`nav-item ${selectedTab === 0 && "active"}`}
                                                            onClick={() => setSelectedTab(0)}
                                                        >
                                                            <span className="nav-link pb-0 pt-0">
                                                                {/* <span className="nav-link-icon d-md-none d-lg-inline-block"></span> */}
                                                                <span className="nav-link-title">Summary </span>
                                                            </span>
                                                        </li>
                                                        <li
                                                            className={`nav-item ${selectedTab === 1 && "active"}`}
                                                            onClick={() => setSelectedTab(1)}
                                                        >
                                                            <span className="nav-link pb-0 pt-0">
                                                                {/* <span className="nav-link-icon d-md-none d-lg-inline-block"></span> */}
                                                                <span className="nav-link-title">Versions </span>
                                                            </span>
                                                        </li>
                                                        <li
                                                            className={`nav-item ${selectedTab === 2 && "active"}`}
                                                            onClick={() => setSelectedTab(2)}
                                                        >
                                                            <span className="nav-link pb-0 pt-0">
                                                                {/* <span className="nav-link-icon d-md-none d-lg-inline-block"></span> */}
                                                                <span className="nav-link-title">Workspaces that use this template</span>
                                                            </span>
                                                        </li>
                                                    </ul>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </header>
                        </Col>
                    </Row>
                    <Row className="mt-4 pb-5">
                        {selectedTab === 0 && (
                            <Col md={12}>
                                <TemplateDetailsSummary template={template} />
                            </Col>
                        )}
                        {selectedTab === 1 && (
                            <Col md={12}>
                                <TemplateDetailsVersions template={template} />
                            </Col>
                        )}
                        {selectedTab === 2 && (
                            <Col md={12}>
                                <TemplateDetailsWorkspaces template={template} />
                            </Col>
                        )}
                    </Row>
                </Container>
            )}
            <ToastContainer />
        </React.Fragment>
    );
}