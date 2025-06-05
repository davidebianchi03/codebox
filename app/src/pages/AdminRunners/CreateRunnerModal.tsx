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
import { Runner, RunnerType } from "../../types/runner";
import { toast } from "react-toastify";
import { ListRunnerTypes } from "../../api/runner";
import { AdminCreateRunner, AdminListRunners } from "../../api/admin";

interface Props {
  isOpen: boolean;
  onClose: (token: string | null) => void;
}

export function CreateRunnerModal({ isOpen, onClose }: Props) {
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);

  const FetchRunnerTypes = useCallback(async () => {
    const rt = await ListRunnerTypes();
    if (rt) {
      setRunnerTypes(rt);
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
            const runners = await AdminListRunners();

            if (runners) {
              return (
                (runners as Runner[]).find((r) => r.name === value) ===
                undefined
              );
            }

            return false;
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
                const runners = await AdminListRunners();
                if (!runners) {
                  return false;
                }
                return (
                  runners.find(
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
      const runner = await AdminCreateRunner(
        values.runnerName,
        values.runnerType,
        values.usePublicUrl,
        values.publicUrl,
      );

      if (runner) {
        HandleCloseModal(runner.token);
      } else {
        toast.error(`Failed to create runner, try again later`);
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
              className={`form-control ${validation.errors.runnerType ? "is-invalid" : ""
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
              color="accent"
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
