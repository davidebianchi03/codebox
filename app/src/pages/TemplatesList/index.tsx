import React, { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  CardBody,
  CardHeader,
  Col,
  Container,
  Input,
  Row,
} from "reactstrap";
import { toast, ToastContainer } from "react-toastify";
import { User } from "../../types/user";
import { WorkspaceTemplate } from "../../types/templates";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { WorkspaceType } from "../../types/workspace";
import { Link } from "react-router-dom";
import { CreateTemplateModal } from "./createTemplateModal";

export default function TemplatesList() {
  // TODO: check if current user can manage templates
  const [currentUser, setCurrentUser] = useState<User>();
  const [templates, setTemplates] = useState<WorkspaceTemplate[]>();
  const [workspaceTypes, setWorkspaceTypes] = useState<WorkspaceType[]>([]);
  const [showCreateTemplateModal, setShowCreateTemplateModal] = useState<boolean>(false);

  // TODO: debounce
  const [searchText, setSearchText] = useState<string>("");

  const fetchTemplates = useCallback(async () => {
    var [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/templates`,
      "GET",
      null
    );

    if (status === RequestStatus.OK && statusCode === 200) {
      setTemplates(responseData);
    } else {
      toast.error("Failed to fetch templates");
      setTemplates([]);
    }
  }, []);

  const fetchWorkspaceTypes = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/workspace-types`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setWorkspaceTypes(responseData as WorkspaceType[]);
    }
  }, []);

  useEffect(() => {
    fetchTemplates();
    fetchWorkspaceTypes();
  }, [fetchTemplates, fetchWorkspaceTypes]);

  return (
    <Container className="mt-4 mb-4">
      <div className="row g-2 align-items-center">
        <div className="col">
          <div className="page-pretitle">Overview</div>
          <h2 className="page-title">Templates</h2>
        </div>
        <div className="col-auto ms-auto d-print-none">
          <div className="btn-list">
            <Button
              color="primary"
              onClick={() => setShowCreateTemplateModal(true)}
            >
              Create template
            </Button>
          </div>
        </div>
      </div>
      <Row className="mt-4">
        <Col md={12}>
          <Card>
            <CardHeader>
              <Input
                placeholder="Search template"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
              />
            </CardHeader>
            <CardBody className="pt-1">
              {
                templates?.filter(template => template.name.indexOf(searchText) >= 0).map((template) => (
                  <React.Fragment key={template.id}>
                    <div className="d-flex pb-2 my-2 border-bottom align-items-center">
                      {
                        template.icon ? (
                          <img
                            src={template.icon}
                            style={{
                              width: 50,
                              height: 50,
                              fontSize: 20,
                              padding: 3,
                              opacity: 0.5,
                              borderRadius: 4,
                            }}
                            alt="custom template icon"
                          />
                        ) : (
                          <div
                            style={{
                              width: 50,
                              height: 50,
                              fontSize: 20,
                              opacity: 0.5,
                              borderRadius: 4,
                            }}
                            className="bg-primary text-white d-flex align-items-center justify-content-center"
                          >
                            {template.name[0].toUpperCase()}
                          </div>
                        )
                      }
                      <div className="w-100 d-flex justify-content-between">
                        <div className="ms-3">
                          <h4 className="mb-0">
                            <Link to={`/templates/${template.id}`}>{template.name}</Link>
                          </h4>
                          <small className="text-muted">{template.description}</small>
                        </div>
                        <div>
                          <Badge color="primary" className="text-white">
                            {(() => {
                              var templateType = workspaceTypes.find((wt) => wt.id === template.type);
                              return templateType ? templateType.name : template.type;
                            })()}
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </React.Fragment>
                ))
              }
            </CardBody>
          </Card>
        </Col>
      </Row>
      <CreateTemplateModal
        isOpen={showCreateTemplateModal}
        onClose={() => setShowCreateTemplateModal(false)}
      />
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </Container>
  );
}
