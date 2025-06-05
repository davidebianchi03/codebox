import { useCallback, useEffect, useState } from "react";
import {
  Button,
  FormFeedback,
  Input,
  Label,
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
import { AdminRetrieveUserByEmail, AdminUpdateUser } from "../../api/admin";

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

  return (
    <>
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
              disabled={user?.email === currentUser?.email}
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
        <div className="d-flex justify-content-end mt-5">
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
    </>
  );
}
