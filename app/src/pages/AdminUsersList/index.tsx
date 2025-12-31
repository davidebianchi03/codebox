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
import { toast, ToastContainer } from "react-toastify";
import { AdminUser } from "../../types/user";
import { Link, useNavigate } from "react-router-dom";
import { AdminListUsers } from "../../api/users";
import DataTable from "../../components/DataTable";
import { TimeSince } from "../../common/time";
import { useSelector } from "react-redux";
import { RootState } from "../../redux/store";
import { AuthenticationSettings } from "../../types/settings";
import { APIAdminRetrieveAuthenticationSettings } from "../../api/common";
import { UserListCardPlaceholder } from "./UserListCardPlaceholder";

export function AdminUsersList() {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [showCreateUserModal, setShowCreateUserModal] = useState<boolean>(false);
  const [authSettings, setAuthSettings] = useState<AuthenticationSettings>();
  const [loading, setLoading] = useState<boolean>(true);

  const currentUser = useSelector((state: RootState) => state.user);

  const navigate = useNavigate();

  const FetchInfo = useCallback(async () => {
    setLoading(true);
    const s = await APIAdminRetrieveAuthenticationSettings();
    if (s) {
      setAuthSettings(s);
    } else {
      toast.error("Failed to fetch settings, try again later");
    }

    const u = await AdminListUsers();
    if (u) {
      setUsers(u);
    } else if (u === undefined) {
      toast.error("Failed to fetch users, try again later");
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    FetchInfo();
  }, [FetchInfo]);

  return (
    <React.Fragment>
      <Container>
        {loading ? (
          <React.Fragment>
            <UserListCardPlaceholder />
          </React.Fragment>
        ) : (
          <React.Fragment>

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
                  New User
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
                            <div className="d-flex flex-column text-light">
                              <div>
                                <b>{user.first_name} {user.last_name}</b>
                                {user.is_superuser && (
                                  <React.Fragment>
                                    <Badge color="info" className="text-white ms-2">
                                      Admin
                                    </Badge>
                                  </React.Fragment>
                                )}
                                {user.is_template_manager && !user.is_superuser && (
                                  <React.Fragment>
                                    <Badge color="accent" className="text-white ms-2">
                                      Template Manager
                                    </Badge>
                                  </React.Fragment>
                                )}
                                {user.email === currentUser?.email && (
                                  <React.Fragment>
                                    <Badge color="success" className="text-white ms-2">
                                      It's you
                                    </Badge>
                                  </React.Fragment>
                                )}
                                {!user.approved && authSettings?.users_must_be_approved && (
                                  <React.Fragment>
                                    <Badge color="danger" className="text-white ms-2">
                                      Not Approved
                                    </Badge>
                                  </React.Fragment>
                                )}
                                {user.deletion_in_progress && (
                                  <React.Fragment>
                                    <Badge color="orange" className="text-white ms-2">
                                      Deletion in progress
                                      <Spinner size="sm" />
                                    </Badge>
                                  </React.Fragment>
                                )}
                              </div>
                              <small className="text-muted">{user.email}</small>
                            </div>
                          </Link>
                        ),
                      },
                      {
                        label: "Last login",
                        key: "last_login",
                        sortable: true,
                        render: (_, user: AdminUser) => (
                          <React.Fragment>
                            <div className="mt-2">
                              {user.last_login
                                ? TimeSince(new Date(user.last_login))
                                : "Never"}
                            </div>
                          </React.Fragment>
                        ),
                      }
                    ]}
                    data={users}
                    pageSizeOptions={[1, 5, 10, 20, 50, 100]}
                    initialPageSize={10}
                  />
                </Card>
              </Col>
            </Row>
          </React.Fragment>
        )}
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
