import React, { useCallback, useEffect, useState } from "react";
import { Form, Modal } from "react-bootstrap";
import { useSelector } from "react-redux";
import { RootState } from "../redux/store";
import { AdminAnalyticsConfig } from "./AdminAnalyticsConfig";
import { AdminAnalyticsContentPreview } from "./AdminAnalyticsContentPreview";
import { APIAdminGetAnalyticsDataBannerSent, APIAdminSetAnalyticsDataBannerSent } from "../api/admin";
import { toast } from "react-toastify";

export function AdminAnalyticsModal() {

    const [isOpen, setIsOpen] = useState<boolean>(false);
    const user = useSelector((state: RootState) => state.user);

    const fetchData = useCallback(async () => {
        if (user?.is_superuser) {
            const dataSent = await APIAdminGetAnalyticsDataBannerSent();
            if (!dataSent) {
                setIsOpen(true);
            }
        }
    }, [user]);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    const handleAnalyticsSave = async () => {
        if (!await APIAdminSetAnalyticsDataBannerSent()) {
            toast.error("Failed to set analytics preferences, try again later");
        } else {
            toast.success("Analytics preferences have been set");
            setIsOpen(false)
        }
    }

    return (
        <React.Fragment>
            <Modal
                show={isOpen}
                centered
                size="lg"
                className="modal-blur"
            >
                <Modal.Header>
                    <h3 className="mb-0">Help us to improve Codebox</h3>
                </Modal.Header>
                <Modal.Body>
                    <p className="mb-1">
                        We'd like to collect technical usage statistics (such as number of users, workspaces, and server version)
                        to improve the product.
                    </p>
                    <p className="mb-1">
                        We do not collect personal data or user content.
                    </p>
                    <p className="mb-1">
                        You can change this choice at any time from the administration interface.
                    </p>
                    <p className="mb-1">
                        Below is an example of the data sent to our analytics server:
                    </p>
                    <AdminAnalyticsContentPreview />
                    <div className="mt-4">
                        <AdminAnalyticsConfig
                            overrideInitialValue={true}
                            overriddenInitialValue={true}
                            buttonPosition="end"
                            onSave={handleAnalyticsSave}
                        />
                    </div>
                </Modal.Body>
            </Modal>
        </React.Fragment>
    )
}
