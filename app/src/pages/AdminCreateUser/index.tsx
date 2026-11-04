import { useFormik } from "formik";
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
import * as Yup from "yup";
import { toast, ToastContainer } from "react-toastify";
import { AdminCreateUser, AdminRetrieveUserByEmail } from "../../api/users";
import { useNavigate } from "react-router-dom";
import { AdminPageHeader } from "../../components/AdminPageHeader";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong } from "@fortawesome/free-solid-svg-icons";

export default function AdminCreateUserPage() {
  const navigate = useNavigate();

  const validation = useFormik({
    initialValues: {
      email: "",
      firstName: "",
      lastName: "",
      isAdmin: false,
      isTemplateManager: false,
      password: "",
      confirmPassword: "",
    },
    validationSchema: Yup.object({
      email: Yup.string()
        .required("This field is required")
        .email("This must be a valid email address")
        .test({
          name: "check_if_user_exists",
          exclusive: false,
          params: {},
          message: "Another user with the same email address already exists",
          test: async (value) => {
            const user = await AdminRetrieveUserByEmail(value);

            if (user === undefined) {
              toast.error(
                `Failed to check if user already exists, try again later`
              );
              return false;
            }

            return user === null;
          },
        }),
      firstName: Yup.string().required("This field is required"),
      lastName: Yup.string().required("This field is required"),
      password: Yup.string()
        .required("A password is required")
        .test({
          name: "password",
          exclusive: false,
          params: {},
          message:
            "The password must be at least 10 characters long and include at least one uppercase letter and one special symbol (!_-,.?!).",
          test: (value) => {
            if (!value || value.length < 10) {
              return false;
            }
            const hasUppercase = /[A-Z]/.test(value);
            const hasSpecialSymbol = /[!_\-,.?]/.test(value);
            return hasUppercase && hasSpecialSymbol;
          },
        }),
      confirmPassword: Yup.string()
        .required("Confirm the password")
        .test({
          name: "confirmPassword",
          exclusive: false,
          params: {},
          message: "Passwords do not match",
          test: (value, context) => value === context.parent.password,
        }),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      const user = await AdminCreateUser(
        values.email,
        values.password,
        values.firstName,
        values.lastName,
        values.isAdmin,
        values.isTemplateManager,
      );
      if (user !== undefined) {
        navigate(`/admin/users/${user.email}`);
      } else {
        toast.error("Failed to create user, unknown error");
      }
    },
  });

  return (
    <>
      <Container>
        <AdminPageHeader
          title="New user"
          description="Create a user account and set its permissions."
          actions={
            <Button color="accent" onClick={() => navigate("/admin/users")}>
              <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
              Back to users
            </Button>
          }
        />
        <Row className="mt-4">
          <Col md={12}>
            <Card>
              <CardHeader className="pb-0 border-0 flex-column align-items-start mb-4">
                <h3 className="mb-0">Account details</h3>
                <p className="text-muted mb-0">
                  Enter the user's profile, permissions and password.
                </p>
              </CardHeader>
              <CardBody className="pt-0">
                <form
                  onSubmit={(event) => {
                    event.preventDefault();
                    validation.handleSubmit();
                  }}
                >
                  <div className="mb-3">
                    <Label>Email</Label>
                    <Input
                      placeholder="johndoe@example.com"
                      name="email"
                      onChange={validation.handleChange}
                      value={validation.values.email}
                      invalid={Boolean(validation.errors.email)}
                    />
                    <FormFeedback>{validation.errors.email}</FormFeedback>
                  </div>
                  <div className="mb-3">
                    <Label>First Name</Label>
                    <Input
                      placeholder="John"
                      name="firstName"
                      onChange={validation.handleChange}
                      value={validation.values.firstName}
                      invalid={Boolean(validation.errors.firstName)}
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
                      invalid={Boolean(validation.errors.lastName)}
                    />
                    <FormFeedback>{validation.errors.lastName}</FormFeedback>
                  </div>
                  {/* <div className="mb-3">
                    <label className="form-check">
                      <input
                        className="form-check-input form-check-input-light"
                        type="checkbox"
                        name="isAdmin"
                        onChange={validation.handleChange}
                        checked={validation.values.isAdmin}
                      />
                      <span className="form-check-label">Admin</span>
                    </label>
                  </div>
                  <div className="mb-3">
                    <label className="form-check">
                      <input
                        className="form-check-input form-check-input-light"
                        type="checkbox"
                        name="isTemplateManager"
                        onChange={validation.handleChange}
                        checked={
                          validation.values.isTemplateManager ||
                          validation.values.isAdmin
                        }
                        disabled={validation.values.isAdmin}
                      />
                      <span className="form-check-label">Template Manager</span>
                    </label>
                  </div> */}
                  <div className="mb-3">
                    <Label>Password</Label>
                    <Input
                      placeholder="Password"
                      name="password"
                      onChange={validation.handleChange}
                      value={validation.values.password}
                      invalid={Boolean(validation.errors.password)}
                      type="password"
                    />
                    <FormFeedback>{validation.errors.password}</FormFeedback>
                  </div>
                  <div className="mb-3">
                    <Label>Confirm password</Label>
                    <Input
                      placeholder="Confirm password"
                      name="confirmPassword"
                      onChange={validation.handleChange}
                      value={validation.values.confirmPassword}
                      invalid={Boolean(validation.errors.confirmPassword)}
                      type="password"
                    />
                    <FormFeedback>
                      {validation.errors.confirmPassword}
                    </FormFeedback>
                  </div>
                  <hr className="my-3" />
                  <div className="d-flex justify-content-end">
                    <Button
                      type="button"
                      color="accent"
                      className="me-2"
                      onClick={() => navigate("/admin/users")}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" color="light">
                      Create
                    </Button>
                  </div>
                </form>
              </CardBody>
            </Card>
          </Col>
        </Row>
      </Container>
      <ToastContainer toastClassName="bg-dark" />
    </>
  );
}
