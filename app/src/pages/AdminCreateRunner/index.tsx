import { useCallback, useEffect, useState } from "react";
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
import { useNavigate } from "react-router-dom";
import { AdminCreateRunner, ListRunnerTypes } from "../../api/runner";
import { AdminPageHeader } from "../../components/AdminPageHeader";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong } from "@fortawesome/free-solid-svg-icons";
import { RunnerType } from "../../types/runner";

export default function AdminCreateRunnerPage() {
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const navigate = useNavigate();

  const FetchRunnerTypes = useCallback(async () => {
    const types = await ListRunnerTypes();
    if (types) {
      setRunnerTypes(types);
    } else {
      toast.error("Failed to fetch runner types, try again later");
    }
  }, []);

  useEffect(() => {
    FetchRunnerTypes();
  }, [FetchRunnerTypes]);

  const validation = useFormik({
    initialValues: {
      runnerName: "",
      description: "",
      runnerType: "",
      usePublicUrl: false,
      publicUrl: "",
    },
    validationSchema: Yup.object({
      runnerName: Yup.string().required("Runner name is required"),
      runnerType: Yup.string().required("Runner type is required"),
      publicUrl: Yup.string().when("usePublicUrl", {
        is: true,
        then: (schema) =>
          schema
            .required("Public url is required"),
      }),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      const runner = await AdminCreateRunner(
        values.runnerName,
        values.runnerType,
        values.description,
        values.usePublicUrl,
        values.publicUrl,
      );

      if (runner) {
        navigate("/admin/runners/created", {
          state: {
            runnerId: runner.id.toString(),
            token: runner.token,
            publicUrl: values.usePublicUrl ? values.publicUrl : undefined,
          },
        });
      } else {
        toast.error("Failed to create runner, try again later");
      }
    },
  });

  return (
    <>
      <Container>
        <AdminPageHeader
          title="Register a new runner"
          description="Configure a runner instance."
          actions={
            <Button color="accent" onClick={() => navigate("/admin/runners")}>
              <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
              Back to runners
            </Button>
          }
        />
        <Row className="mt-4">
          <Col md={12}>
            <Card>
              <CardHeader className="pb-0 border-0 flex-column align-items-start mb-4">
                <h3 className="mb-0">Runner details</h3>
              </CardHeader>
              <CardBody className="pt-0">
                <form onSubmit={validation.handleSubmit}>
                  <div className="mb-3">
                    <Label>Runner name</Label>
                    <Input
                      placeholder="my runner"
                      name="runnerName"
                      onChange={validation.handleChange}
                      value={validation.values.runnerName}
                      invalid={Boolean(validation.errors.runnerName)}
                    />
                    <FormFeedback>{validation.errors.runnerName}</FormFeedback>
                  </div>
                  <div className="mb-3">
                    <Label>Description</Label>
                    <Input
                      type="textarea"
                      rows={3}
                      placeholder="Describe this runner"
                      name="description"
                      onChange={validation.handleChange}
                      value={validation.values.description}
                    />
                  </div>
                  <div className="mb-3">
                    <Label>Runner type</Label>
                    <select
                      name="runnerType"
                      className={`form-control ${validation.errors.runnerType ? "is-invalid" : ""}`}
                      onChange={validation.handleChange}
                      value={validation.values.runnerType}
                    >
                      <option value="">Select runner type</option>
                      {runnerTypes.map((type) => (
                        <option value={type.id} key={type.id}>
                          {type.name}
                        </option>
                      ))}
                    </select>
                    <FormFeedback>{validation.errors.runnerType}</FormFeedback>
                  </div>
                  {/* <div className="mb-3">
                    <label className="form-check">
                      <input
                        className="form-check-input form-check-input-light"
                        type="checkbox"
                        name="usePublicUrl"
                        onChange={(event) => {
                          validation.setFieldValue("publicUrl", "");
                          validation.handleChange(event);
                        }}
                        checked={validation.values.usePublicUrl}
                      />
                      <span className="form-check-label">Use public url</span>
                    </label>
                  </div>
                  <div className="mb-3">
                    <Label>Public url</Label>
                    <Input
                      placeholder="http://my-host.example.com:12345"
                      name="publicUrl"
                      onChange={validation.handleChange}
                      value={validation.values.publicUrl}
                      invalid={Boolean(validation.errors.publicUrl)}
                      disabled={!validation.values.usePublicUrl}
                    />
                    <FormFeedback>{validation.errors.publicUrl}</FormFeedback>
                  </div> */}
                  <hr className="my-3" />
                  <div className="d-flex justify-content-end">
                    <Button
                      type="button"
                      color="accent"
                      className="me-2"
                      onClick={() => navigate("/admin/runners")}
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
