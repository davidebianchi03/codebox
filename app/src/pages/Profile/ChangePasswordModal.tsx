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
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import Swal from "sweetalert2";

interface Props {
  isOpen: boolean;
  onClose: () => void;
}

export function ChangePasswordModal({ isOpen, onClose }: Props) {
  var validation = useFormik({
    initialValues: {
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
    validationSchema: Yup.object({
      currentPassword: Yup.string().required("This field is required"),
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
        })
        .test({
          name: "checkPrevious",
          exclusive: false,
          params: {},
          message: "The new password must be different from the previous one",
          test: (value, context) => value !== context.parent.currentPassword,
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
      var [status, statusCode] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/auth/change-password`,
        "POST",
        JSON.stringify({
          current_password: values.currentPassword,
          new_password: values.newPassword,
        })
      );

      if (status === RequestStatus.OK && statusCode === 200) {
        await Swal.fire(
          "Password changed",
          "Password has been changed successfully!",
          "success"
        );
        handleCloseModal();
        // TODO: logout
      } else if (statusCode === 417) {
        validation.setFieldError("currentPassword", "Wrong password");
      } else {
        await Swal.fire("Unknown error", "Password change failed!", "error");
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
          Change Password
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
              <Label>Current Password</Label>
              <Input
                name="currentPassword"
                value={validation.values.currentPassword}
                onChange={validation.handleChange}
                invalid={!!validation.errors.currentPassword}
                placeholder="Current password"
                type="password"
              />
              <FormFeedback>{validation.errors.currentPassword}</FormFeedback>
            </div>
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
            <div className="d-flex justify-content-end mt-3">
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
              <Button className="ms-1" color="primary">
                Submit
              </Button>
            </div>
          </form>
        </ModalBody>
      </Modal>
    </React.Fragment>
  );
}
