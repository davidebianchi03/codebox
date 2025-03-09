import { useCallback, useEffect, useState } from "react";
import { Card, CardBody, Col, Container, Row, Table } from "reactstrap";
import { Runner } from "../../types/runner";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

export function AdminRunners() {
  const [runners, setRunners] = useState<Runner[]>([]);

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

  useEffect(() => {
    FetchRunners();
  }, [FetchRunners]);

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
                <Table striped>
                  <thead>
                    <th>#</th>
                    <th>Name</th>
                    <th>Type</th>
                    <th>Supported workspace types</th>
                    <th>Last contact</th>
                  </thead>
                  <tbody>
                    {runners.length === 0 ? (
                      <tr>
                        <td>There are no registered runners</td>
                      </tr>
                    ) : (
                      runners.map((runner) => (
                        <tr key={runner.id}>
                            <td>{runner.id}</td>
                            <td>{runner.name}</td>
                            <td>{runner.type}</td>
                            <td></td>
                            <td></td>
                        </tr>
                      ))
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
