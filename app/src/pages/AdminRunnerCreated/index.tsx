import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Alert, Button, Card, CardBody, CardHeader, Col, Container, Input, Label, Row, Spinner } from "reactstrap";
import { toast, ToastContainer } from "react-toastify";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faArrowLeftLong, faCopy, faInfoCircle, faKey, faLock, faRocket, faServer } from "@fortawesome/free-solid-svg-icons";
import { AdminPageHeader } from "../../components/AdminPageHeader";
import { AdminRetrieveRecommendedRunnerVersion, AdminRetrieveRunnerById } from "../../api/runner";
import { RunnerAdmin } from "../../types/runner";
import { AdminRunnerStatusBadge, IsRunnerOnline } from "../../components/AdminRunnerStatusBadge";

const RUNNER_STATUS_POLL_INTERVAL = 5000;

interface RunnerCredentials {
  runnerId: string;
  token: string;
  publicUrl?: string;
}

function isRunnerCredentials(value: unknown): value is RunnerCredentials {
  if (!value || typeof value !== "object") {
    return false;
  }

  return (
    "runnerId" in value &&
    typeof value.runnerId === "string" &&
    "token" in value &&
    typeof value.token === "string" &&
    (!("publicUrl" in value) || typeof value.publicUrl === "string" || value.publicUrl === undefined)
  );
}

function shellQuote(value: string) {
  return `'${value.replace(/'/g, "'\\''")}'`;
}

export default function AdminRunnerCreatedPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const credentials = isRunnerCredentials(location.state) ? location.state : null;
  const [runner, setRunner] = useState<RunnerAdmin>();
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const [statusError, setStatusError] = useState(false);
  const [recommendedVersion, setRecommendedVersion] = useState<string>();
  const [versionError, setVersionError] = useState(false);

  useEffect(() => {
    if (!credentials) {
      navigate("/admin/runners", { replace: true });
      return;
    }

    const runnerId = Number(credentials.runnerId);
    if (!Number.isInteger(runnerId) || runnerId <= 0) {
      navigate("/admin/runners", { replace: true });
      return;
    }

    let active = true;
    let timeoutId: ReturnType<typeof setTimeout>;

    const pollRunnerStatus = async () => {
      const result = await AdminRetrieveRunnerById(runnerId);
      if (!active) {
        return;
      }

      if (result) {
        setRunner(result);
        setStatusError(false);
      } else {
        setStatusError(true);
      }
      setIsCheckingStatus(false);
      timeoutId = setTimeout(pollRunnerStatus, RUNNER_STATUS_POLL_INTERVAL);
    };

    void pollRunnerStatus();

    return () => {
      active = false;
      clearTimeout(timeoutId);
    };
  }, [credentials, navigate]);

  useEffect(() => {
    let active = true;

    const fetchRecommendedVersion = async () => {
      const version = await AdminRetrieveRecommendedRunnerVersion();
      if (!active) {
        return;
      }

      if (version) {
        setRecommendedVersion(version);
      } else {
        setVersionError(true);
      }
    };

    void fetchRecommendedVersion();

    return () => {
      active = false;
    };
  }, []);

  const copyToClipboard = async (value: string, label: string) => {
    try {
      await navigator.clipboard.writeText(value);
      toast.success(`${label} copied to clipboard`);
    } catch {
      toast.error(`Failed to copy ${label.toLowerCase()} to clipboard`);
    }
  };

  const dockerCommand = recommendedVersion
    ? [
      "docker run -d --name codebox-runner --restart unless-stopped",
      "-p 8100:8100",
      "-v /var/run/docker.sock:/var/run/docker.sock",
      "-v codebox-runner-config:/etc/codebox-runner",
      "-v codebox-runner-logs:/var/log/codebox-runner",
      `-e CODEBOX_SERVER_URL=${shellQuote(window.location.origin)}`,
      `-e CODEBOX_RUNNER_ID=${shellQuote(credentials?.runnerId ?? "")}`,
      `-e CODEBOX_TOKEN=${shellQuote(credentials?.token ?? "")}`,
      `-e CODEBOX_RUNNER_EXTERNAL_URL=${shellQuote(credentials?.publicUrl || "http://172.17.0.1:8100")}`,
      "-e CODEBOX_OBJECTS_PREFIX='codebox'",
      shellQuote(`dadebia/codebox-docker-runner:${recommendedVersion}`),
    ].join(" \\\n  ")
    : undefined;

  if (!credentials) {
    return null;
  }

  const isOnline = runner ? IsRunnerOnline(runner) : false;

  return (
    <>
      <Container>
        <AdminPageHeader
          title="Runner created"
          description="Your runner is registered. Save its credentials before leaving this page."
          actions={
            <Button color="accent" onClick={() => navigate("/admin/runners")}>
              <FontAwesomeIcon icon={faArrowLeftLong} className="me-2" />
              Back to runners
            </Button>
          }
        />
        <Row className="mt-4">
          <Col lg={7} className="mb-4">
            <Card>
              <CardHeader className="d-flex align-items-center gap-2">
                <h3 className="mb-0">Runner credentials</h3>
              </CardHeader>
              <CardBody>
                <Alert color="warning" className="d-flex align-items-center gap-2">
                  <FontAwesomeIcon icon={faLock} />
                  <span>
                    This token is shown only once. Copy it and store it securely.
                  </span>
                </Alert>
                <div className="mb-3">
                  <Label htmlFor="runner-id">Runner ID</Label>
                  <div className="d-flex gap-2">
                    <Input id="runner-id" value={credentials.runnerId} readOnly />
                    <Button color="accent" onClick={() => copyToClipboard(credentials.runnerId, "ID")} aria-label="Copy runner ID">
                      <FontAwesomeIcon icon={faCopy} />
                    </Button>
                  </div>
                </div>
                <div className="mb-3">
                  <Label htmlFor="runner-token">Token</Label>
                  <div className="d-flex gap-2">
                    <Input id="runner-token" value={credentials.token} readOnly />
                    <Button color="accent" onClick={() => copyToClipboard(credentials.token, "Token")} aria-label="Copy runner token">
                      <FontAwesomeIcon icon={faCopy} />
                    </Button>
                  </div>
                </div>
                <div className="d-flex justify-content-end gap-2">
                  <Button color="primary" onClick={() => copyToClipboard(credentials.token, "Token")}>
                    <FontAwesomeIcon icon={faCopy} className="me-2" />
                    Copy token
                  </Button>
                  <Button color="accent" onClick={() => navigate(`/admin/runners/${credentials.runnerId}`)}>
                    Done
                  </Button>
                </div>
              </CardBody>
            </Card>
          </Col>
          <Col lg={5} className="mb-4">
            <Card>
              <CardHeader className="d-flex align-items-center justify-content-between">
                <div className="d-flex align-items-center gap-2">
                  <h3 className="mb-0">Connection status</h3>
                </div>
                {runner && <AdminRunnerStatusBadge runner={runner} />}
              </CardHeader>
              <CardBody className="text-center">
                {runner ? (
                  <div className="mb-3">
                    <div className={`runner-connection-icon ${isOnline ? "is-online" : "is-offline"}`}>
                      <FontAwesomeIcon icon={faServer} />
                    </div>
                    <h3 className="mt-3 mb-2">
                      {isOnline
                        ? "Runner connected"
                        : runner.last_contact
                          ? "Runner is offline"
                          : "Waiting for first connection"}
                    </h3>
                    <p className="text-muted mb-1">
                      {isOnline
                        ? "Your runner is online and ready to use."
                        : "Start the runner with the installation command below. This page will update when it connects."}
                    </p>
                    {statusError && (
                      <small className="text-warning">Unable to refresh runner status. Retrying...</small>
                    )}
                  </div>
                ) : isCheckingStatus ? (
                  <div className="py-3">
                    <Spinner size="sm" />
                    <p className="text-muted mt-2 mb-0">Checking runner status...</p>
                  </div>
                ) : (
                  <p className="text-warning mb-0">
                    {statusError
                      ? "Unable to retrieve runner status. Retrying..."
                      : "Runner status is unavailable."}
                  </p>
                )}
                <div className="runner-setup-steps text-start mt-3">
                  {["Install the runner", "Configure the token", "Start the service"].map((step, index) => (
                    <div
                      className={`runner-setup-step ${isOnline ? "is-complete" : runner && index === 0 ? "is-active" : ""
                        }`}
                      key={step}
                    >
                      <span className="runner-setup-step-number">{index + 1}</span>
                      <span>{step}</span>
                    </div>
                  ))}
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>
        <Row>
          <Col>
            <Card>
              <CardHeader className="d-flex align-items-center gap-2">
                <h3 className="mb-0">Install with Docker</h3>
              </CardHeader>
              <CardBody>
                {dockerCommand ? (
                  <div className="d-flex gap-2">
                    <Input
                      type="textarea"
                      rows={3}
                      aria-label="Docker runner installation command"
                      className="font-monospace"
                      value={dockerCommand}
                      readOnly
                      style={{ minHeight: "250px" }}
                    />
                    <div>
                      <Button color="accent" onClick={() => copyToClipboard(dockerCommand, "Docker command")} aria-label="Copy Docker command">
                        <FontAwesomeIcon icon={faCopy} />
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="d-flex align-items-center gap-2 text-muted">
                    {versionError ? (
                      <span className="text-warning">Unable to load the recommended runner version. Reload the page to try again.</span>
                    ) : (
                      <>
                        <Spinner size="sm" />
                        <span>Loading recommended runner version...</span>
                      </>
                    )}
                  </div>
                )}
                <div className="d-flex align-items-center gap-2 text-muted mt-2">
                  <FontAwesomeIcon icon={faInfoCircle} className="text-info" />
                  <small>
                    For details, see the{" "}
                    <a
                      href="https://gitlab.com/codebox4073715/codebox/-/blob/master/docs/source/guide/quickstart/runner.md"
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      runner installation guide
                    </a>.
                  </small>
                </div>
              </CardBody>
            </Card>
          </Col>
        </Row>
      </Container>
      <ToastContainer toastClassName="bg-dark" />
    </>
  );
}
