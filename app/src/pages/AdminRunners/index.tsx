import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Col,
  Input,
  Label,
  Row,
  Table,
} from "reactstrap";
import { Runner, RunnerType } from "../../types/runner";
import { CreateRunnerModal } from "./CreateRunnerModal";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy, faXmark } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";
import { Link } from "react-router-dom";
import { ListRunnerTypes } from "../../api/runner";
import { AdminListRunners } from "../../api/admin";

export function AdminRunners() {
  const [runners, setRunners] = useState<Runner[]>([]);
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const [searchText, setSearchText] = useState<string>("");
  const [showCreateRunnerModal, setCreateRunnerModal] =
    useState<boolean>(false);
  const [runnerToken, setRunnerToken] = useState<string>("");
  const [runnerId, setRunnerID] = useState<string>("");

  const FetchRunners = useCallback(async () => {
    const r = await AdminListRunners();
    if (r) {
      setRunners(r);
    }
  }, []);

  const FetchRunnerTypes = useCallback(async () => {
    const rt = await ListRunnerTypes();
    if (rt) {
      setRunnerTypes(rt);
    }
  }, []);

  useEffect(() => {
    FetchRunners();
    FetchRunnerTypes();
  }, [FetchRunners, FetchRunnerTypes]);

  return (
    <>
      <div className="row g-2 align-items-center mb-4">
        <div className="col-auto ms-auto d-print-none">
          <Button color="primary" onClick={() => setCreateRunnerModal(true)}>
            Add new runner
          </Button>
        </div>
      </div>
      {runnerToken.length > 0 && runnerId.length > 0 && (
        <Row>
          <Col md={12}>
            <Card style={{ borderRadius: 5 }}>
              <CardBody className="bg-primary" style={{ borderRadius: 5 }}>
                <div className="d-flex justify-content-between">
                  <h3>Runner has been created</h3>
                  <Button
                    className="bg-transparent p-0 m-0"
                    style={{ height: 25 }}
                    onClick={() => {
                      setRunnerToken("");
                      setRunnerID("");
                    }}
                  >
                    <FontAwesomeIcon icon={faXmark} />
                  </Button>
                </div>
                <p>
                  Use the following ID and token to register the runner, the token will not
                  longer be visible
                </p>
                <div className="d-flex align-items-center">
                  <Label className="mt-2 me-6">
                    ID
                  </Label>
                  <Input
                    value={runnerId}
                    style={{ background: "var(--tblr-primary-darken)" }}
                    className="text-white"
                    disabled
                  />
                  <Button
                    className="bg-transparent"
                    onClick={() => {
                      navigator.clipboard.writeText(runnerId);
                      toast.info("Copied to clipboard");
                    }}
                  >
                    <FontAwesomeIcon icon={faCopy} />
                  </Button>
                </div>
                <div className="d-flex align-items-center mt-2">
                  <Label className="mt-2 me-3">
                    Token
                  </Label>
                  <Input
                    value={runnerToken}
                    style={{ background: "var(--tblr-primary-darken)" }}
                    className="text-white"
                    disabled
                  />
                  <Button
                    className="bg-transparent"
                    onClick={() => {
                      navigator.clipboard.writeText(runnerToken);
                      toast.info("Copied to clipboard");
                    }}
                  >
                    <FontAwesomeIcon icon={faCopy} />
                  </Button>
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>
      )}
      <Row className="mt-4">
        <Col md={12}>
          <Input
            placeholder="Filter runners"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <Table striped className="mt-4" responsive>
            <thead>
              <tr>
                <th>#</th>
                <th>Name</th>
                <th>Type</th>
                <th>Supported workspace types</th>
                <th>Last contact</th>
                <th>Version</th>
              </tr>
            </thead>
            <tbody>
              {runners.length === 0 ? (
                <tr>
                  <td colSpan={5}>There are no registered runners</td>
                </tr>
              ) : (
                runners.map((runner) => {
                  var runnerType = runnerTypes.find(
                    (type) => type.id === runner.type
                  );

                  if (runner.name.indexOf(searchText) >= 0) {
                    return (
                      <tr key={runner.id}>
                        <td>
                          <Link to={`/admin/runners/${runner.id}`}>{runner.id}</Link>
                        </td>
                        <td>{runner.name}</td>
                        <td
                          data-bs-toggle="tooltip"
                          data-bs-placement="top"
                          title={runnerType?.description}
                        >
                          {runnerType?.name || "N/A"}
                        </td>
                        <td>
                          {runnerType
                            ? runnerType.supported_types?.map((t, i) => {
                              if (runnerType) {
                                if (
                                  i <
                                  runnerType?.supported_types.length - 1
                                ) {
                                  return t.name + ", ";
                                }
                              }
                              return t.name;
                            })
                            : ""}
                        </td>
                        <td className={`${(new Date().getTime() - new Date(runner.last_contact).getTime()) > (5 * 60 * 1000) ? "text-warning" : ""
                          }`}>
                          {new Date(runner.last_contact).getFullYear() <
                            2000
                            ? "N/A"
                            : new Date(
                              runner.last_contact
                            ).toLocaleString()}
                        </td>
                        <td>
                          {runner.version.length > 0 ? runner.version : "N/A"}
                        </td>
                      </tr>
                    );
                  } else {
                    return null;
                  }
                })
              )}
            </tbody>
          </Table>
        </Col>
      </Row>
      <CreateRunnerModal
        isOpen={showCreateRunnerModal}
        onClose={(id, token) => {
          setCreateRunnerModal(false);
          FetchRunners();
          if (token && id) {
            setRunnerToken(token);
            setRunnerID(id);
          } else {
            setRunnerToken("");
            setRunnerID("");
          }
        }}
      />
      <ToastContainer
        toastClassName={"bg-dark"}
      />
    </>
  );
}
