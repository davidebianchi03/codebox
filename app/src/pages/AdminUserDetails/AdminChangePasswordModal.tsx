import { useFormik } from "formik";
import React from "react";
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
import { AdminUser } from "../../types/user";
import { toast } from "react-toastify";
import { AdminSetUserPassword } from "../../api/users";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  user: AdminUser;
}

export function AdminChangePasswordModal({ isOpen, onClose, user }: Props) {
  var validation = useFormik({
    initialValues: {
      newPassword: "",
      confirmPassword: "",
    },
    validationSchema: Yup.object({
      newPassword: Yup.string()
        .required("This field is required")
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
        .required("This field is required")
        .test({
          name: "checkMatch",
          exclusive: false,
          params: {},
          message: "Passwords do not match",
          test: (value, context) => value === context.parent.newPassword,
        }),
    }),
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      if (await AdminSetUserPassword(user.email, values.newPassword)) {
        toast.success("Password has been changed successfully");
        handleCloseModal();
      } else {
        toast.error("Failed to set new password");
      }
    },
  });

  const handleCloseModal = () => {
    validation.resetForm();
    onClose();
  };

  return (
    <React.Fragment>
      <Modal
        isOpen={isOpen}
        toggle={handleCloseModal}
        centered
        size="lg"
        modalClassName="modal-blur"
        fade
      >
        <ModalHeader toggle={handleCloseModal} className="border-0">
          Change password
        </ModalHeader>
        <ModalBody className="pt-1">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              validation.handleSubmit();
              return false;
            }}
          >
            <div className="mb-3">
              <Label>New Password</Label>
              <Input
                name="newPassword"
                value={validation.values.newPassword}
                onChange={validation.handleChange}
                invalid={!!validation.errors.newPassword}
                placeholder="New password"
                type="password"
              />
              <FormFeedback>{validation.errors.newPassword}</FormFeedback>
            </div>
            <div>
              <Label>Confirm Password</Label>
              <Input
                name="confirmPassword"
                value={validation.values.confirmPassword}
                onChange={validation.handleChange}
                invalid={!!validation.errors.confirmPassword}
                placeholder="Confirm password"
                type="password"
              />
              <FormFeedback>{validation.errors.confirmPassword}</FormFeedback>
            </div>
            <div className="d-flex justify-content-end mt-4">
              <Button
                color="accent"
                onClick={(e) => {
                  e.preventDefault();
                  handleCloseModal();
                  return false;
                }}
              >
                Cancel
              </Button>
              <Button className="ms-1" color="light">
                Submit
              </Button>
            </div>
          </form>
        </ModalBody>
      </Modal>
    </React.Fragment>
  );
}
