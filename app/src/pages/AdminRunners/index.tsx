import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Card,
  CardBody,
  Col,
  Container,
  Input,
  Row,
  Table,
} from "reactstrap";
import { Runner, RunnerType } from "../../types/runner";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { CreateRunnerModal } from "./CreateRunnerModal";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy, faXmark } from "@fortawesome/free-solid-svg-icons";
import { toast, ToastContainer } from "react-toastify";

export function AdminRunners() {
  const [runners, setRunners] = useState<Runner[]>([]);
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const [searchText, setSearchText] = useState<string>("");
  const [showCreateRunnerModal, setCreateRunnerModal] =
    useState<boolean>(false);
  const [runnerToken, setRunnerToken] = useState<string>("");

  const FetchRunners = useCallback(async () => {
    let [status, statusCode, responseData] = await Http.Request(
      `${Http.GetServerURL()}/api/v1/admin/runners`,
      "GET",
      null
    );
    if (status === RequestStatus.OK && statusCode === 200) {
      setRunners(responseData as Runner[]);
    }
  }, []);

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
    FetchRunners();
    FetchRunnerTypes();
  }, [FetchRunners, FetchRunnerTypes]);

  return (
    <>
      <Container className="mt-4">
        <div className="row g-2 align-items-center mb-4">
          <div className="col">
            <div className="page-pretitle">Admin</div>
            <h2 className="page-title">Runners</h2>
          </div>
          <div className="col-auto ms-auto d-print-none">
            <Button color="primary" onClick={() => setCreateRunnerModal(true)}>
              Add new runner
            </Button>
          </div>
        </div>
        {runnerToken.length > 0 && (
          <Row>
            <Col md={12}>
              <Card style={{ borderRadius: 5 }}>
                <CardBody className="bg-success" style={{ borderRadius: 5 }}>
                  <div className="d-flex justify-content-between">
                    <h3>Runner has been created</h3>
                    <Button
                      className="bg-transparent p-0 m-0"
                      style={{ height: 25 }}
                      onClick={() => {
                        setRunnerToken("");
                      }}
                    >
                      <FontAwesomeIcon icon={faXmark} />
                    </Button>
                  </div>
                  <p>
                    Use the following token to register the runner, it will not
                    longer be visible
                  </p>
                  <div className="d-flex align-items-center">
                    <Input
                      value={runnerToken}
                      style={{
                        borderColor: "#248c35",
                        background: "#39d952",
                        color: "#fff",
                      }}
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
            <Card>
              <CardBody>
                <Input
                  placeholder="Filter runners"
                  value={searchText}
                  onChange={(e) => setSearchText(e.target.value)}
                />
                <Table striped className="mt-4">
                  <thead>
                    <tr>
                      <th>#</th>
                      <th>Name</th>
                      <th>Type</th>
                      <th>Supported workspace types</th>
                      <th>Last contact</th>
                    </tr>
                  </thead>
                  <tbody>
                    {runners.length === 0 ? (
                      <tr>
                        <td>There are no registered runners</td>
                      </tr>
                    ) : (
                      runners.map((runner) => {
                        var runnerType = runnerTypes.find(
                          (type) => type.id === runner.type
                        );

                        if (runner.name.indexOf(searchText) >= 0) {
                          return (
                            <tr key={runner.id}>
                              <td>{runner.id}</td>
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
                              <td>
                                {new Date(runner.last_contact).toLocaleString()}
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
              </CardBody>
            </Card>
          </Col>
        </Row>
        <CreateRunnerModal
          isOpen={showCreateRunnerModal}
          onClose={(token) => {
            setCreateRunnerModal(false);
            FetchRunners();
            if (token) {
              setRunnerToken(token);
            } else {
              setRunnerToken("");
            }
          }}
        />
        <ToastContainer />
      </Container>
    </>
  );
}
