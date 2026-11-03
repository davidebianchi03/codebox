import { faDownload } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useEffect, useState } from "react";
import { Button, Modal, Spinner } from "react-bootstrap";
import { toast } from "react-toastify";
import { APIAdminListRunnerLogs } from "../../api/admin";
import { LogsViewer } from "../../components/LogsViewer";
import { SystemLog } from "../../types/logs";

interface RunnerLogsModalProps {
  isOpen: boolean;
  runnerId: number;
  runnerName: string;
  onClose: () => void;
}

export function RunnerLogsModal({
  isOpen,
  runnerId,
  runnerName,
  onClose,
}: RunnerLogsModalProps) {
  const [logs, setLogs] = useState<SystemLog[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!isOpen) {
      setLogs([]);
      setLoading(false);
      return;
    }

    let active = true;
    const fetchLogs = async () => {
      setLoading(true);
      const result = await APIAdminListRunnerLogs(runnerId);
      if (!active) {
        return;
      }
      if (result === undefined) {
        toast.error("Failed to fetch runner logs");
        setLogs([]);
      } else {
        setLogs(result);
      }
      setLoading(false);
    };

    fetchLogs();
    return () => {
      active = false;
    };
  }, [isOpen, runnerId]);

  const downloadLogs = () => {
    const blob = new Blob([JSON.stringify(logs, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `runner-${runnerId}-logs.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.setTimeout(() => URL.revokeObjectURL(url), 0);
  };

  return (
    <Modal
      show={isOpen}
      onHide={onClose}
      centered
      size="xl"
      scrollable
      className="modal-blur"
    >
      <Modal.Header closeButton>
        <Modal.Title>{runnerName} logs</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {loading ? (
          <div className="d-flex justify-content-center py-5">
            <Spinner animation="border" role="status">
              <span className="visually-hidden">Loading logs...</span>
            </Spinner>
          </div>
        ) : (
          <LogsViewer logs={logs.slice()} />
        )}
      </Modal.Body>
      <Modal.Footer>
        <Button
          variant="outline-light"
          onClick={downloadLogs}
          disabled={loading || logs.length === 0}
        >
          <FontAwesomeIcon icon={faDownload} className="me-2" />
          Download logs
        </Button>
        <Button variant="accent" onClick={onClose}>
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
