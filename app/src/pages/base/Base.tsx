import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";
import { Navbar } from "../../components/navbar/Navbar";
import StatusBar from "../../components/statusbar/StatusBar";
import "./Base.css"
import { Component, ReactNode } from "react";

interface BasePageProps {
    children: any
}

interface BasePageState {
    displayedUsername: string,
    serverPing: number
}

export default class BasePage extends Component<BasePageProps, BasePageState> {

    updateUIInterval: NodeJS.Timer|null

    constructor(props: BasePageProps) {
        super(props);
        this.state = {
            displayedUsername: "",
            serverPing: 0,
        }
        this.updateUIInterval = null;
    }

    componentDidMount(): void {
        if(!this.updateUIInterval){
            this.retrieveUsername();
            this.updateUIInterval = setInterval(this.retrieveUsername, 30000);
        }
    }

    retrieveUsername = async () => {
        let requestStart = new Date(Date.now());
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/whoami`, "GET", null);
        let requestEnd = new Date(Date.now());
        this.setState({ serverPing: requestEnd.getMilliseconds() - requestStart.getMilliseconds() });
        if (status === RequestStatus.OK && statusCode === 200) {
            if (responseData.first_name && responseData.last_name) {
                this.setState({ displayedUsername: responseData.first_name + " " + responseData.last_name });
            } else {
                this.setState({ displayedUsername: responseData.email });
            }
        }
    }

    render(): ReactNode {
        return (
            <div className="basepage-container">
                <Navbar username={this.state.displayedUsername} />
                <div className="basepage-content">
                    {this.props.children}
                </div>
                <StatusBar serverPing={this.state.serverPing} />
            </div>
        )
    }
}