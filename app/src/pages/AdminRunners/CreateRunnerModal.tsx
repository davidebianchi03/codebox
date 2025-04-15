import { useFormik } from "formik";
import { useCallback, useEffect, useState } from "react";
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
import { Runner, RunnerType } from "../../types/runner";
import { toast } from "react-toastify";

interface Props {
  isOpen: boolean;
  onClose: (token: string | null) => void;
}

export function CreateRunnerModal({ isOpen, onClose }: Props) {
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);

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

  useEffect(() => {
    FetchRunnerTypes();
  }, [FetchRunnerTypes]);

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
              (responseData as Runner[]).find((r) => r.name === value) ===
              undefined
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
                  return false;
                }
                return (
                  (responseData as Runner[]).find(
                    (r) => r.public_url === value && r.use_public_url
                  ) === undefined
                );
              }
            ),
      }),
    }),
    validateOnChange: false,
    validateOnBlur: false,
    onSubmit: async (values) => {
      var requestBody = {
        name: values.runnerName,
        type: values.runnerType,
        use_public_url: values.usePublicUrl,
        public_url: values.publicUrl,
      };

      let [status, statusCode, responseData] = await Http.Request(
        `${Http.GetServerURL()}/api/v1/admin/runners`,
        "POST",
        JSON.stringify(requestBody),
        "application/json"
      );
      if (status === RequestStatus.OK && statusCode === 201) {
        HandleCloseModal(responseData.token);
      } else if (statusCode === 409) {
        toast.error(responseData.detail);
      } else {
        toast.error(`Failed to create runner, received status ${statusCode}`);
      }
    },
  });

  const HandleCloseModal = (token: string | null) => {
    validation.resetForm();
    onClose(token);
  };

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
        Add new runner
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
            <Label>Runner name</Label>
            <Input
              placeholder="my runner"
              name="runnerName"
              onChange={validation.handleChange}
              value={validation.values.runnerName}
              invalid={validation.errors.runnerName ? true : false}
            />
            <FormFeedback>{validation.errors.runnerName}</FormFeedback>
          </div>
          <div className="mb-3">
            <Label>Runner type</Label>
            <select
              name="runnerType"
              className={`form-control ${
                validation.errors.runnerType ? "is-invalid" : ""
              }`}
              onChange={validation.handleChange}
              value={validation.values.runnerType}
            >
              <option value={""}>Select runner type</option>
              {runnerTypes.map((t) => (
                <option value={t.id} key={t.id}>{t.name}</option>
              ))}
            </select>
            <FormFeedback>{validation.errors.runnerType}</FormFeedback>
          </div>
          <div className="mb-3">
            <label className="form-check">
              <input
                className="form-check-input"
                type="checkbox"
                name="usePublicUrl"
                onChange={(e) => {
                  validation.setFieldValue("publicUrl", "");
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
          <hr className="my-3" />
          <div className="d-flex justify-content-end">
            <Button
              color="outline-light"
              className="me-2"
              onClick={(e) => {
                e.preventDefault();
                HandleCloseModal(null);
              }}
            >
              Cancel
            </Button>
            <Button type="submit" color="primary">
              Create
            </Button>
          </div>
        </form>
      </ModalBody>
    </Modal>
  );
}
