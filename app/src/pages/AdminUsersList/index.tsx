import React, { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Card,
  Col,
  Container,
  Row,
  Spinner,
} from "reactstrap";
import { CreateUserModal } from "./CreateUserModal";
import { ToastContainer } from "react-toastify";
import { AdminUser } from "../../types/user";
import { Link, useNavigate } from "react-router-dom";
import { AdminListUsers } from "../../api/users";
import DataTable from "../../components/DataTable";
import { TimeSince } from "../../common/time";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCheckCircle, faCircleXmark } from "@fortawesome/free-solid-svg-icons";

export function AdminUsersList() {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [showCreateUserModal, setShowCreateUserModal] =
    useState<boolean>(false);

  const navigate = useNavigate();

  const FetchUsers = useCallback(async () => {
    const u = await AdminListUsers();
    if (u) {
      setUsers(u);
    }
  }, []);

  useEffect(() => {
    FetchUsers();
  }, [FetchUsers]);

  return (
    <React.Fragment>
      <Container>
        <div className="row g-2 align-items-center mb-2">
          <div className="col">
            <h2 className="mb-0 mt-2">Users</h2>
            <p className="text-muted">List of all users</p>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <Button
              color="light"
              onClick={() => setShowCreateUserModal(true)}
            >
              Create new user
            </Button>
          </div>
        </div>
        <Row className="mt-2">
          <Col md={12}>
            <Card body>
              <DataTable
                columns={[
                  {
                    label: "Email",
                    key: "email",
                    sortable: true,
                    render: (_, user: AdminUser) => (
                      <Link to={`/admin/users/${user.email}`} className="d-flex gap-2 align-items-center">
                        <b>{user.email}</b>
                        {user.deletion_in_progress && (
                          <React.Fragment>
                            <Badge color="orange" className="text-white">
                              Deletion in progress
                              <Spinner size="sm" />
                            </Badge>
                          </React.Fragment>
                        )}
                      </Link>
                    ),
                  },
                  {
                    label: "First Name",
                    key: "first_name",
                    sortable: true,
                  },
                  {
                    label: "Last Name",
                    key: "last_name",
                    sortable: true,
                  },
                  {
                    label: "Admin",
                    key: "is_superuser",
                    sortable: true,
                    render: (_, user: AdminUser) =>
                      <React.Fragment>
                        <span className={`ps-3 ${user.is_superuser ? "text-success" : "text-danger"}`}>
                          {user.is_superuser ?
                            <FontAwesomeIcon icon={faCheckCircle} /> :
                            <FontAwesomeIcon icon={faCircleXmark} />
                          }
                        </span>
                      </React.Fragment>,
                  },
                  {
                    label: "Template manager",
                    key: "is_template_manager",
                    sortable: true,

                    render: (_, user: AdminUser) =>
                      <React.Fragment>
                        <span
                          className={`${user.is_template_manager ? "text-success" : "text-danger"}`}
                          style={{ paddingLeft: "3.5rem" }}
                        >
                          {user.is_template_manager ?
                            <FontAwesomeIcon icon={faCheckCircle} /> :
                            <FontAwesomeIcon icon={faCircleXmark} />
                          }
                        </span>
                      </React.Fragment>,
                  },
                  {
                    label: "Email verified",
                    key: "email_verified",
                    sortable: true,

                    render: (_, user: AdminUser) =>
                      <React.Fragment>
                        <span
                          className={`${user.email_verified ? "text-success" : "text-danger"}`}
                          style={{ paddingLeft: "2.5rem" }}
                        >
                          {user.email_verified ?
                            <FontAwesomeIcon icon={faCheckCircle} /> :
                            <FontAwesomeIcon icon={faCircleXmark} />
                          }
                        </span>
                      </React.Fragment>,
                  },
                  {
                    label: "Last login",
                    key: "last_login",
                    sortable: true,
                    render: (_, user: AdminUser) =>
                      user.last_login
                        ? TimeSince(new Date(user.last_login))
                        : "Never",
                  }
                ]}
                data={users}
                pageSizeOptions={[1, 5, 10, 20, 50, 100]}
                initialPageSize={10}
              />
            </Card>
          </Col>
        </Row>
        <CreateUserModal
          isOpen={showCreateUserModal}
          onClose={(user) => {
            setShowCreateUserModal(false);
            if (user) {
              navigate(`/admin/users/${user.email}`);
            }
          }}
        />
        <ToastContainer
          toastClassName={"bg-dark"}
        />
      </Container>
    </React.Fragment>
  );
}
