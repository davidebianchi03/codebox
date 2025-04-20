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
import { Runner, RunnerType } from "../../types/runner";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { toast, ToastContainer } from "react-toastify";
import { useNavigate, useParams } from "react-router-dom";
import { useFormik } from "formik";
import * as Yup from "yup";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong } from "@fortawesome/free-solid-svg-icons";

export function AdminRunnerDetails() {
  const [runner, setRunner] = useState<Runner>();
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);

  const { id } = useParams();
  const navigate = useNavigate();

  const FetchRunnerTypes = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/runner-types`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setRunnerTypes(responseData as RunnerType[]);
    }
  }, []);

  const validation = useFormik({
    initialValues: {
      runnerName: "",
      runnerType: "",
      usePublicUrl: false,
      publicUrl: "",
    },
    validationSchema: Yup.object({
      runnerName: Yup.string()
        .required("Runner name is required")
        .test(
          "Another runner with the same name already exists",
          async (value) => {
            let [status, statusCode, responseData] = await Http.Request(
              `${Http.GetServerURL()}/api/v1/admin/runners`,
              "GET",
              null
            );
            if (status !== RequestStatus.OK && statusCode !== 200) {
              return false;
            }
            return (
              (responseData as Runner[]).find(
                (r) => r.name === value && r.id !== runner?.id
              ) === undefined
            );
          }
        ),
      runnerType: Yup.string().required("Runner type is required"),
      publicUrl: Yup.string().when("usePublicUrl", {
        is: true,
        then: (schema) =>
          schema
            .required("Public url is required")
            .test(
              "Another runner with the same public url already exists",
              async (value) => {
                let [status, statusCode, responseData] = await Http.Request(
                  `${Http.GetServerURL()}/api/v1/admin/runners`,
                  "GET",
                  null
                );
                if (status !== RequestStatus.OK && statusCode !== 200) {
                  toast.error(`Failed to update runner - ${statusCode}`);
                  return false;
                }
                return (
                  (responseData as Runner[]).find(
                    (r) =>
                      r.public_url === value &&
                      r.use_public_url &&
                      r.id !== runner?.id
                  ) === undefined
                );
              }
            ),
      }),
    }),
    validateOnBlur: false,
    validateOnChange: false,
    onSubmit: async (values) => {
      let [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/admin/runners/${id}`,
        "PUT",
        JSON.stringify({
          name: values.runnerName,
          type: values.runnerType,
          use_public_url: values.usePublicUrl,
          public_url: values.publicUrl,
        })
      );
      if (status !== RequestStatus.OK || statusCode !== 200) {
        setRunnerTypes(responseData as RunnerType[]);
      } else {
        validation.resetForm();
        FetchRunner();
      }
    },
  });

  const FetchRunner = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/admin/runners/${id}`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      var runner = responseData as Runner;
      validation.setValues({
        runnerName: runner.name,
        runnerType: runner.type,
        usePublicUrl: runner.use_public_url,
        publicUrl: runner.public_url,
      });
      setRunner(runner);
    } else {
      navigate("/");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id, navigate]);

  useEffect(() => {
    FetchRunner();
    FetchRunnerTypes();
  }, [FetchRunner, FetchRunnerTypes]);

  return (
    <>
      <Button
        color="accent"
        className="me-2 mb-4"
        onClick={(e) => {
          e.preventDefault();
          navigate("/admin/runners");
        }}
      >
        <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
        Back
      </Button>
      <Row>
        <Col md={12}>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              validation.handleSubmit();
              return false;
            }}
          >
            <div className="mb-3">
              <Label>Runner name</Label>
              <Input
                name="runnerName"
                value={validation.values.runnerName}
                onChange={validation.handleChange}
                invalid={!!validation.errors.runnerName}
              />
              <FormFeedback>{validation.errors.runnerName}</FormFeedback>
            </div>
            <div className="mb-3">
              <Label>Runner type</Label>
              <select
                name="runnerType"
                className={`form-control ${validation.errors.runnerType ? "is-invalid" : ""
                  }`}
                onChange={validation.handleChange}
                value={validation.values.runnerType}
              >
                <option value={""}>Select runner type</option>
                {runnerTypes.map((t) => (
                  <option value={t.id} key={t.id}>
                    {t.name}
                  </option>
                ))}
              </select>
              <small className="text-muted">
                Supported workspaces:{" "}
                {runnerTypes
                  .find((r) => r.id === validation.values.runnerType)
                  ?.supported_types.map((st) => st.name)
                  .join(", ")}
              </small>
              <FormFeedback>{validation.errors.runnerType}</FormFeedback>
            </div>
            <div className="mb-3">
              <label className="form-check">
                <input
                  className="form-check-input"
                  type="checkbox"
                  name="usePublicUrl"
                  onClick={(e) => {
                    validation.handleChange(e);
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
                invalid={validation.errors.publicUrl ? true : false}
                disabled={!validation.values.usePublicUrl}
              />
              <FormFeedback>{validation.errors.publicUrl}</FormFeedback>
            </div>
            <div className="d-flex justify-content-end mt-4">
              <Button
                color="accent"
                className="me-2"
                onClick={(e) => {
                  e.preventDefault();
                  FetchRunner();
                }}
              >
                Cancel
              </Button>
              <Button type="submit" color="primary">
                Save
              </Button>
            </div>
          </form>
        </Col>
      </Row>
      <ToastContainer />
    </>
  );
}
