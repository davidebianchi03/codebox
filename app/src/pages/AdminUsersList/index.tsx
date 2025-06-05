import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Col,
  Input,
  Row,
  Table,
} from "reactstrap";
import { CreateUserModal } from "./CreateUserModal";
import { ToastContainer } from "react-toastify";
import { User } from "../../types/user";
import { Link, useNavigate } from "react-router-dom";
import { AdminListUsers } from "../../api/admin";

export function AdminUsersList() {
  const [users, setUsers] = useState<User[]>([]);
  const [searchText, setSearchText] = useState<string>("");
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
    <>
      <div className="row g-2 align-items-center mb-4">
        <div className="col-auto ms-auto d-print-none">
          <Button
            color="primary"
            onClick={() => setShowCreateUserModal(true)}
          >
            Create new user
          </Button>
        </div>
      </div>
      <Row className="mt-4">
        <Col md={12}>
          <Input
            placeholder="Filter users"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <Table striped className="mt-4">
            <thead>
              <tr>
                <th>Email</th>
                <th>First Name</th>
                <th>Last Name</th>
                <th className="text-center">Admin</th>
                <th className="text-center">Template manager</th>
              </tr>
            </thead>
            <tbody>
              {users.length === 0 ? (
                <tr>
                  <td colSpan={6}>There are no users</td>
                </tr>
              ) : (
                users
                  .filter((user) => user.email.indexOf(searchText) >= 0)
                  .map((user, index) => (
                    <tr key={index}>
                      <td>
                        <Link to={`/admin/users/${user.email}`}>
                          {user.email}
                        </Link>
                      </td>
                      <td>{user.first_name}</td>
                      <td>{user.last_name}</td>
                      <td className="text-center">
                        {user.is_superuser ? (
                          <span className="badge bg-success text-white">
                            Yes
                          </span>
                        ) : (
                          <span className="badge bg-danger text-white">
                            No
                          </span>
                        )}
                      </td>
                      <td className="text-center">
                        {user.is_template_manager ? (
                          <span className="badge bg-success text-white">
                            Yes
                          </span>
                        ) : (
                          <span className="badge bg-danger text-white">
                            No
                          </span>
                        )}
                      </td>
                    </tr>
                  ))
              )}
            </tbody>
          </Table>
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
    </>
  );
}
