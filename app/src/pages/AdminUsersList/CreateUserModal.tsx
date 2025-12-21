import { useFormik } from "formik";
import {
  Button,
  FormFeedback,
  Input,
  Label,
  Modal,
  ModalBody,
  ModalHeader,
} from "reactstrap";
import * as Yup from "yup";
import { toast } from "react-toastify";
import { AdminUser } from "../../types/user";
import { AdminCreateUser, AdminRetrieveUserByEmail } from "../../api/users";
import { useCallback } from "react";

interface Props {
  isOpen: boolean;
  onClose: (user: AdminUser | null) => void;
}

export function CreateUserModal({ isOpen, onClose }: Props) {
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
          test: (value, context) => {
            if (value.length < 10) {
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
      const u = await AdminCreateUser(
        values.email,
        values.password,
        values.firstName,
        values.lastName,
        values.isAdmin,
        values.isTemplateManager,
      );
      if (u !== undefined) {
        HandleCloseModal(u);
      } else {
        toast.error(`Failed to create user, unkwnown error`);
      }
    },
  });

  const HandleCloseModal = useCallback((user: AdminUser | null) => {
    validation.resetForm();
    onClose(user);
  }, [onClose, validation]);

  return (
    <Modal
      isOpen={isOpen}
      toggle={() => {
        HandleCloseModal(null);
      }}
      centered
      size="lg"
      modalClassName="modal-blur"
      fade
    >
      <ModalHeader
        toggle={() => {
          HandleCloseModal(null);
        }}
      >
        Create new user
      </ModalHeader>
      <ModalBody>
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
              name="email"
              onChange={validation.handleChange}
              value={validation.values.email}
              invalid={validation.errors.email ? true : false}
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
                onClick={validation.handleChange}
                checked={validation.values.isAdmin}
              />
              <span className="form-check-label">Admin</span>
            </label>
          </div>
          {
            !validation.values.isAdmin && (
              <div className="mb-3">
                <label className="form-check">
                  <input
                    className="form-check-input"
                    type="checkbox"
                    name="isTemplateManager"
                    onClick={validation.handleChange}
                    checked={validation.values.isTemplateManager}
                  />
                  <span className="form-check-label">Template Manager</span>
                </label>
              </div>
            )
          }
          <div className="mb-3">
            <Label>Password</Label>
            <Input
              placeholder="Password"
              name="password"
              onChange={validation.handleChange}
              value={validation.values.password}
              invalid={validation.errors.password ? true : false}
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
              invalid={validation.errors.confirmPassword ? true : false}
              type="password"
            />
            <FormFeedback>{validation.errors.confirmPassword}</FormFeedback>
          </div>
          <hr className="my-3" />
          <div className="d-flex justify-content-end">
            <Button
              color="accent"
              className="me-2"
              onClick={(e) => {
                e.preventDefault();
                HandleCloseModal(null);
              }}
            >
              Cancel
            </Button>
            <Button type="submit" color="light">
              Create
            </Button>
          </div>
        </form>
      </ModalBody>
    </Modal>
  );
}
