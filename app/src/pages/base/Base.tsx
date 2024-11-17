import "./Base.css"
import { Http } from "../../api/http";
import { Component, ReactNode, useEffect, useState } from "react";
import { RequestStatus } from "../../api/types";
import { Navbar } from "../../components/navbar/Navbar";
import StatusBar from "../../components/statusbar/StatusBar";
import { useNavigate } from "react-router-dom";

interface BasePageProps {
    children: any,
    authRequired?: boolean
}

export default function BasePage(props: BasePageProps) {

    const [pingInterval, setPingInterval] = useState<number>(0);
    const [useGravatar, setUseGravatar] = useState<boolean>(true);
    const [firstName, setFirstName] = useState<string>("");
    const [lastName, setLastName] = useState<string>("");
    const [emailAddress, setEmailAddress] = useState<string>("");

    const navigate = useNavigate();

    var updateUIInterval: NodeJS.Timer | null = null;

    const retrieveUsername = async () => {
        let requestStart = new Date(Date.now());
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
        let requestEnd = new Date(Date.now());
        setPingInterval(requestEnd.getMilliseconds() - requestStart.getMilliseconds());
        if (status === RequestStatus.OK && statusCode === 200) {
            setFirstName(responseData.first_name);
            setLastName(responseData.last_name);
            setEmailAddress(responseData.email);
        } else if(status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
            if(props.authRequired === undefined || props.authRequired === true) {
                navigate("/login");
                return;
            }
        }
    }

    const retrieveSettings = async () => {
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/instance-settings`, "GET", null);
        if (status === RequestStatus.OK && statusCode === 200) {
            setUseGravatar(responseData.use_gravatar);
        }
    }

    useEffect(() => {
        if (updateUIInterval === null) {
            updateUIInterval = setInterval(retrieveUsername, 30000);
        }

        retrieveUsername();
        retrieveSettings();

        return () => {
            if (updateUIInterval !== null) {
                clearInterval(updateUIInterval);
            }
        };
    }, []);

    return (
        <div className="basepage-container">
            <Navbar firstName={firstName} lastName={lastName} email={emailAddress} useGravatar={useGravatar} />
            <div className="basepage-content">
                {props.children}
            </div>
            <StatusBar serverPing={pingInterval} />
        </div>
    )
}
