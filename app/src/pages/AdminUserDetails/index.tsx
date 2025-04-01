import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Col,
  Container,
  FormFeedback,
  Input,
  Label,
  Row,
} from "reactstrap";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast, ToastContainer } from "react-toastify";
import { useNavigate, useParams } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { User } from "../../types/user";
import { AdminChangePasswordModal } from "./AdminChangePasswordModal";

export function AdminUserDetails() {
  const [user, setUser] = useState<User>();
  const [showChangePasswordModal, setShowChangePasswordModal] =
    useState<boolean>(false);

  const { email } = useParams();
  const navigate = useNavigate();

  const validation = useFormik({
    initialValues: {
      firstName: "",
      lastName: "",
      isAdmin: false,
    },
    validationSchema: Yup.object({
      firstName: Yup.string().required("This field is required"),
      lastName: Yup.string().required("This field is required"),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      var requestBody = {
        first_name: values.firstName,
        last_name: values.lastName,
        is_superuser: values.isAdmin,
      };
      let [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/admin/users/${user?.email}`,
        "PATCH",
        JSON.stringify(requestBody),
        "application/json"
      );
      if (status === RequestStatus.OK && statusCode === 200) {
        FetchUser();
        toast.success("User has been updated");
      } else if (statusCode === 409) {
        toast.error(responseData.detail);
      } else {
        toast.error(`Failed to update user, received status ${statusCode}`);
      }
    },
  });

  const FetchUser = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/admin/users/${email}`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      var user = responseData as User;
      validation.setValues({
        firstName: user.first_name || "",
        lastName: user.last_name || "",
        isAdmin: user.is_superuser,
      });
      setUser(user);
    } else {
      navigate("/");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [email, navigate]);

  useEffect(() => {
    FetchUser();
  }, [FetchUser]);

  return (
    <>
      <Container className="mt-4">
        <div className="row g-2 align-items-center mb-4">
          <div className="col">
            <div className="page-pretitle">Runners</div>
            <h2 className="page-title">{user?.email}</h2>
          </div>
        </div>
        <Button
          color="outline-light"
          className="me-2"
          onClick={(e) => {
            e.preventDefault();
            navigate("/admin/users");
          }}
        >
          Go back
        </Button>
        <Row className="mt-4">
          <Col md={12}>
            <Card>
              <CardBody>
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    validation.handleSubmit();
                    return false;
                  }}
                >
                  <div className="mb-3">
                    <Label>Email</Label>
                    <Input
                      placeholder="johndoe@example.com"
                      value={user?.email}
                      disabled
                    />
                  </div>
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
                      />
                      <span className="form-check-label">Admin</span>
                    </label>
                  </div>
                  <hr className="my-3" />
                  <div className="d-flex justify-content-end">
                    <Button
                      type="submit"
                      color="orange"
                      className="me-2"
                      onClick={(e) => {
                        e.preventDefault();
                        setShowChangePasswordModal(true);
                        return false;
                      }}
                    >
                      Change password
                    </Button>
                    <Button
                      color="outline-light"
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
        </Row>
        <ToastContainer />
        {user && (
          <AdminChangePasswordModal
            isOpen={showChangePasswordModal}
            onClose={() => setShowChangePasswordModal(false)}
            user={user}
          />
        )}
      </Container>
    </>
  );
}
