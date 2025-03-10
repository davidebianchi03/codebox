import { useCallback, useEffect, useState } from "react";
import { Card, CardBody, Col, Container, Input, Row, Table } from "reactstrap";
import { Runner, RunnerType } from "../../types/runner";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

export function AdminRunners() {
  const [runners, setRunners] = useState<Runner[]>([]);
  const [runnerTypes, setRunnerTypes] = useState<RunnerType[]>([]);
  const [searchText, setSearchText] = useState<string>("");

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
          <div className="col-auto ms-auto d-print-none"></div>
        </div>
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
      </Container>
    </>
  );
}
