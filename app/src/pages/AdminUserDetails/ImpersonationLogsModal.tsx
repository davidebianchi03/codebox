import React, { useCallback, useEffect, useState } from "react";
import {
  Modal,
  ModalBody,
  ModalHeader,
} from "reactstrap";
import { ImpersonationLogs } from "../../types/impersonationLogs";
import DataTable from "../../components/DataTable";
import { AdminListImpersonationLogs } from "../../api/admin";
import { toast } from "react-toastify";
import { AdminUser } from "../../types/user";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCircleCheck, faCircleXmark } from "@fortawesome/free-solid-svg-icons";

interface Props {
  isOpen: boolean;
  onClose: () => void;
  user: AdminUser;
}

export function ImpersonationLogsModal({ isOpen, onClose, user }: Props) {
  const [impersonationLogs, setImpersonationLogs] = useState<ImpersonationLogs[]>([]);

  const handleCloseModal = useCallback(() => {
    onClose();
  }, [onClose]);

  const FetchImpersonationLogs = useCallback(async () => {
    const l = await AdminListImpersonationLogs(user.email);
    if (l) {
      setImpersonationLogs(l);
    } else {
      toast.error("Failed to retrieve impersonation logs");
    }
  }, [user.email]);

  useEffect(() => {
    if (isOpen) {
      FetchImpersonationLogs();
    }
  }, [FetchImpersonationLogs, isOpen]);

  return (
    <React.Fragment>
      <Modal
        isOpen={isOpen}
        toggle={handleCloseModal}
        centered
        size="xl"
        modalClassName="modal-blur"
        fade
      >
        <ModalHeader toggle={handleCloseModal} className="border-0">
          Impersonation Logs
        </ModalHeader>
        <ModalBody className="pt-1">
          <DataTable
            columns={[
              {
                label: "Impersonator",
                key: "impersonator",
                render: (value) => value.email
              },
              {
                label: "Started On",
                key: "impersonation_started_at",
                render: (value) => new Date(value).toLocaleString()
              },
              {
                label: "Finished On",
                key: "impersonation_finished_at",
                render: (value) => value ? new Date(value).toLocaleString() : "N/A"
              },
              {
                label: "Session Expired",
                key: "session_expired",
                render: (value) => (
                  value ? (
                    <span className="text-success" style={{ marginLeft: 45 }}>
                      <FontAwesomeIcon icon={faCircleCheck} />
                    </span>
                  ) : (
                    <span className="text-danger" style={{ marginLeft: 45 }}>
                      <FontAwesomeIcon icon={faCircleXmark} />
                    </span>
                  )
                )
              },
            ]}
            data={impersonationLogs}
            initialPageSize={10}
          />
        </ModalBody>
      </Modal>
    </React.Fragment>
  );
}
