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
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";
import { useFormik } from "formik";
import * as Yup from "yup";
import { User } from "../../types/user";

export default function Profile() {
  const [sshPublicKey, setSshPublicKey] = useState<string>("");
  const [currentUser, setCurrentUser] = useState<User>();

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
      var data = {
        first_name: values.firstName,
        last_name: values.lastName,
      };

      var [status, statusCode] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/auth/user-details`,
        "PATCH",
        JSON.stringify(data),
        "application/json"
      );
      if (status === RequestStatus.OK && statusCode === 200) {
        toast.info(`Profile has been updated`);
      } else {
        toast.error(`Failed to update profile, received status ${statusCode}`);
      }
    },
  });

  const ResetForm = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/user-details`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setCurrentUser(responseBody);
      validation.resetForm();
      validation.setValues({
        firstName: responseBody.first_name,
        lastName: responseBody.last_name,
      })
    }

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const FetchPublicKey = useCallback(async () => {
    let [status, statusCode, responseBody] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/auth/user-ssh-public-key`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setSshPublicKey(responseBody.public_key);
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
            <CardHeader>
              <div>
                <h4 className="mb-0 w-100">SSH public key</h4>
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
            <CardHeader>
              <h4 className="mb-0">Profile</h4>
            </CardHeader>
            <CardBody>
              <form
                onSubmit={(e) => {
                  validation.handleSubmit();
                  e.preventDefault();
                  return false;
                }}
              >
                <div className="mb-3">
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
                </div>
                <div className="mb-3">
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
                </div>
                <div className="d-flex justify-content-end">
                  <Button
                    color="outline-light"
                    onClick={(e) => {
                      e.preventDefault();
                      ResetForm();
                    }}
                  >
                    Cancel
                  </Button>
                  <Button color="primary" className="ms-1" type="submit">
                    Save changes
                  </Button>
                </div>
              </form>
            </CardBody>
          </Card>
        </Col>
      </Row>
      <ToastContainer />
    </Container>
  );
}
