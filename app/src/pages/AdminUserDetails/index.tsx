import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Col,
  Container,
  FormFeedback,
  Input,
  Label,
  Row,
} from "reactstrap";
import { toast, ToastContainer } from "react-toastify";
import { useNavigate, useParams } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { User } from "../../types/user";
import { AdminChangePasswordModal } from "./AdminChangePasswordModal";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong } from "@fortawesome/free-solid-svg-icons";
import { RetrieveCurrentUserDetails } from "../../api/common";
import { AdminDeleteUser, AdminRetrieveUserByEmail, AdminUpdateUser } from "../../api/admin";
import React from "react";
import Swal from "sweetalert2";

export function AdminUserDetails() {
  const [user, setUser] = useState<User>();
  const [currentUser, setCurrentUser] = useState<User>();
  const [showChangePasswordModal, setShowChangePasswordModal] =
    useState<boolean>(false);

  const { email } = useParams();
  const navigate = useNavigate();

  const validation = useFormik({
    initialValues: {
      firstName: "",
      lastName: "",
      isAdmin: false,
      isTemplateManager: false,
    },
    validationSchema: Yup.object({
      firstName: Yup.string().required("This field is required"),
      lastName: Yup.string().required("This field is required"),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      if (user) {
        if ((
          await AdminUpdateUser(
            user.email,
            values.firstName,
            values.lastName,
            values.isAdmin,
            values.isTemplateManager
          )) !== undefined) {
          FetchUser();
          toast.success("User has been updated");
        } else {
          toast.error(`Failed to update user, try again later`);
        }
      }

    },
  });

  const FetchUser = useCallback(async () => {
    if (email) {
      const user = await AdminRetrieveUserByEmail(email);
      if (user) {
        validation.setValues({
          firstName: user.first_name || "",
          lastName: user.last_name || "",
          isAdmin: user.is_superuser,
          isTemplateManager: user.is_template_manager,
        });
        setUser(user);
      } else {
        navigate("/");
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [email, navigate]);

  const WhoAmI = useCallback(async () => {
    const user = await RetrieveCurrentUserDetails();
    if (user) {
      setCurrentUser(user);
    }
  }, []);

  useEffect(() => {
    FetchUser();
    WhoAmI();
  }, [FetchUser, WhoAmI]);

  const HandleDeleteUser = useCallback(async () => {
    if (user) {
      if (
        (await Swal.fire({
          title: "Delete User",
          text: `
          Are you sure that you want to delete this user?
          All workspaces managed by this user will be deleted
        `,
          icon: "warning",
          showCancelButton: true,
          reverseButtons: true,
          cancelButtonText: "Cancel",
          confirmButtonText: "Delete",
          customClass: {
            popup: "bg-dark text-light",
            cancelButton: "btn btn-accent",
            confirmButton: "btn btn-danger",
          },
        })).isConfirmed
      ) {
        if (await AdminDeleteUser(user.email)) {
          await Swal.fire({
            title: "User will be deleted shortly",
            text: `
              The deletion of ${user.email} has been scheduled,
              the user and all his workspaces will be deleted shortly.
            `,
            icon: "info",
            reverseButtons: true,
            cancelButtonText: "Cancel",
            confirmButtonText: "Ok",
            customClass: {
              popup: "bg-dark text-light",
              confirmButton: "btn btn-light",
            },
          });
          navigate("/admin/users");
        } else {
          toast.error("Failed to delete the user");
        }
      }
    }
  }, [user]);

  return (
    <React.Fragment>
      <Container>
        <Button
          color="accent"
          className="me-2 mb-4"
          onClick={(e) => {
            e.preventDefault();
            navigate("/admin/users");
          }}
        >
          <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
          Back
        </Button>
        <Row>
          <Col md={12}>
            <Card body className="pt-4">
              <div className="d-flex align-items-center">
                <div>
                  <img src="" />
                </div>
                <div>
                  <h2 className="mb-2">{user?.first_name} {user?.last_name}</h2>
                  <p className="text-muted">{user?.email}</p>
                </div>
              </div>
            </Card>
          </Col>
        </Row>
        <Row>
          <Col md={6} className="mt-4">
            <Card>
              <CardHeader className="pb-0 border-0">
                <h3>Basic Information</h3>
              </CardHeader>
              <CardBody className="pt-0">
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    validation.handleSubmit();
                    return false;
                  }}
                >
                  <div className="mb-3">
                    <Label>First Name</Label>
                    <Input
                      placeholder="John"
                      name="firstName"
                      onChange={validation.handleChange}
                      value={validation.values.firstName}
                      invalid={validation.errors.firstName ? true : false}
                    />
                    <FormFeedback>{validation.errors.firstName}</FormFeedback>
                  </div>
                  <div className="mb-3">
                    <Label>Last Name</Label>
                    <Input
                      placeholder="Doe"
                      name="lastName"
                      onChange={validation.handleChange}
                      value={validation.values.lastName}
                      invalid={validation.errors.lastName ? true : false}
                    />
                    <FormFeedback>{validation.errors.lastName}</FormFeedback>
                  </div>
                  <p>
                    <b>Roles</b>
                  </p>
                  <div className="d-flex align-items-center">
                    <div className="mb-3">
                      <label className="form-check">
                        <input
                          className="form-check-input"
                          type="checkbox"
                          name="isAdmin"
                          onClick={(e) => {
                            validation.setFieldValue("publicUrl", "");
                            validation.handleChange(e);
                          }}
                          checked={validation.values.isAdmin}
                          disabled={user?.email === currentUser?.email}
                        />
                        <span className="form-check-label">Admin</span>
                      </label>
                    </div>
                    <div className="mb-3 ms-5">
                      <label className="form-check">
                        <input
                          className="form-check-input"
                          type="checkbox"
                          name="isTemplateManager"
                          onClick={validation.handleChange}
                          checked={validation.values.isTemplateManager || validation.values.isAdmin}
                          disabled={validation.values.isAdmin}
                        />
                        <span className="form-check-label">Template Manager</span>
                      </label>
                    </div>
                  </div>
                  <div className="d-flex justify-content-end">
                    <Button
                      color="accent"
                      className="me-2"
                      onClick={(e) => {
                        e.preventDefault();
                        FetchUser();
                      }}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" color="primary">
                      Save
                    </Button>
                  </div>
                </form>
              </CardBody>
            </Card>
          </Col>
          <Col md={6} className="mt-4">
            <Row>
              <Col md={12}>
                <Card>
                  <CardHeader className="pb-0 border-0" r>
                    <h3>Security</h3>
                  </CardHeader>
                  <CardBody className="pt-0">
                    <Row>
                      <Col md={6}>
                        <p className="mb-2">
                          <b>Last Login:</b>
                        </p>
                        <p>
                          {new Date(user?.last_login || "").toLocaleString()}
                        </p>
                      </Col>
                      <Col md={6}>
                        <p className="mb-2">
                          <b>Created At:</b>
                        </p>
                        <p>
                          {new Date(user?.created_at || "").toLocaleString()}
                        </p>
                      </Col>
                    </Row>
                    <Button
                      type="submit"
                      color="orange"
                      className="me-2 mt-4"
                      onClick={(e) => {
                        e.preventDefault();
                        setShowChangePasswordModal(true);
                        return false;
                      }}
                    >
                      Change password
                    </Button>
                  </CardBody>
                </Card>
              </Col>
            </Row>
            <Row className="mt-4">
              <Col md={12}>
                <Card>
                  <CardHeader className="pb-0 border-0">
                    <h3>Actions</h3>
                  </CardHeader>
                  <CardBody className="pt-0">
                    {/* <Button color="yellow" className="me-2">
                      Impersonate (Coming soon)
                    </Button> */}
                    {/* {user.email} */}
                    <Button color="danger" onClick={HandleDeleteUser}>
                      Delete
                    </Button>
                  </CardBody>
                </Card>
              </Col>
            </Row>
          </Col>
        </Row>
        <ToastContainer
          toastClassName={"bg-dark"}
        />
        {user && (
          <AdminChangePasswordModal
            isOpen={showChangePasswordModal}
            onClose={() => setShowChangePasswordModal(false)}
            user={user}
          />
        )}
      </Container>
    </React.Fragment>
  );
}
