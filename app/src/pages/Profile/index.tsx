import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Col,
  Container,
  FormFeedback,
  FormGroup,
  Input,
  Label,
  Row,
} from "reactstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";
import { useFormik } from "formik";
import * as Yup from "yup";
import { User } from "../../types/user";
import { ChangePasswordModal } from "./ChangePasswordModal";
import { APIRetrieveSshPublicKey, APIUpdateCurrentUserDetails, RetrieveCurrentUserDetails } from "../../api/common";

export default function Profile() {
  const [sshPublicKey, setSshPublicKey] = useState<string>("");
  const [currentUser, setCurrentUser] = useState<User>();
  const [showChangePasswordModal, setShowChangePasswordModal] =
    useState<boolean>(false);

  const validation = useFormik({
    initialValues: {
      firstName: currentUser?.first_name || "",
      lastName: currentUser?.last_name || "",
    },
    validationSchema: Yup.object({
      firstName: Yup.string().required("First name is required"),
      lastName: Yup.string().required("Last name is required"),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      if (await APIUpdateCurrentUserDetails(values.firstName, values.lastName)) {
        toast.info(`Profile has been updated`);
      } else {
        toast.error(`Failed to update profile, try again later`);
      }
    },
  });

  const ResetForm = useCallback(async () => {
    const user = await RetrieveCurrentUserDetails();
    if (user) {
      setCurrentUser(user);
      validation.resetForm();
      validation.setValues({
        firstName: user.first_name,
        lastName: user.last_name,
      });
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const FetchPublicKey = useCallback(async () => {
    const pk = await APIRetrieveSshPublicKey();
    if (pk) {
      setSshPublicKey(pk);
    }
  }, []);

  useEffect(() => {
    FetchPublicKey();
    ResetForm();
  }, [FetchPublicKey, ResetForm]);

  return (
    <Container className="mt-4 mb-4">
      <div className="row g-2 align-items-center mb-4">
        <div className="col">
          <div className="page-pretitle">Profile</div>
        </div>
      </div>
      <Row>
        <Col md={12}>
          <Card>
            <CardHeader className="border-0 pb-0">
              <div>
                <h3 className="mb-0 w-100">SSH Public Key</h3>
                <small className="text-muted">
                  Add this key to your Git server to enable authentication.
                </small>
              </div>
            </CardHeader>
            <CardBody>
              <div
                style={{
                  backgroundColor: "var(--tblr-dark-bg-subtle)",
                  borderRadius: 3,
                  fontFamily: "Consolas",
                  position: "relative",
                  cursor: "pointer",
                }}
                onClick={() => {
                  navigator.clipboard.writeText(sshPublicKey);
                  toast.info("Copied to clipboard");
                }}
                className="px-3 py-2"
              >
                <span
                  style={{
                    position: "absolute",
                    top: "5pt",
                    right: "5pt",
                  }}
                >
                  <FontAwesomeIcon icon={faCopy} />
                </span>
                {sshPublicKey}
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>
      <Row className="mt-4">
        <Col md={12}>
          <Card>
            <CardHeader className="border-0 pb-0">
              <h3 className="mb-0">Personal Information</h3>
            </CardHeader>
            <CardBody>
              <form
                onSubmit={(e) => {
                  validation.handleSubmit();
                  e.preventDefault();
                  return false;
                }}
              >
                <Row>
                  <Col md={6}>
                    <FormGroup>
                      <Label>First name</Label>
                      <Input
                        type="text"
                        placeholder="John"
                        name="firstName"
                        onChange={validation.handleChange}
                        value={validation.values.firstName}
                        invalid={validation.errors.firstName ? true : false}
                      />
                      <FormFeedback>{validation.errors.firstName}</FormFeedback>
                    </FormGroup></Col>
                  <Col md={6}>
                    <FormGroup>
                      <Label>First name</Label>
                      <Input
                        type="text"
                        placeholder="Doe"
                        name="lastName"
                        onChange={validation.handleChange}
                        value={validation.values.lastName}
                        invalid={validation.errors.lastName ? true : false}
                      />
                      <FormFeedback>{validation.errors.lastName}</FormFeedback>
                    </FormGroup>
                  </Col>
                </Row>
                <div className="d-flex justify-content-end">
                  <Button
                    color="accent"
                    onClick={(e) => {
                      e.preventDefault();
                      ResetForm();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button color="primary" className="ms-1" type="submit">
                    Save Changes
                  </Button>
                </div>
              </form>
            </CardBody>
          </Card>
        </Col>
      </Row>
      <Row className="mt-4">
        <Col md={12}>
          <Card>
            <CardBody>
              <div className="w-100 d-flex justify-content-between align-items-center">
                <h3 className="mb-0">Change Password</h3>
                <Button
                  color="primary"
                  onClick={() => setShowChangePasswordModal(true)}
                >
                  Change Password
                </Button>
              </div>
            </CardBody>
          </Card>
        </Col>
      </Row>
      <ChangePasswordModal
        isOpen={showChangePasswordModal}
        onClose={() => setShowChangePasswordModal(false)}
      />
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </Container>
  );
}
