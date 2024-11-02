import "./Base.css"
import { Http } from "../../api/http";
import { Component, ReactNode } from "react";
import { RequestStatus } from "../../api/types";
import { Navbar } from "../../components/navbar/Navbar";
import StatusBar from "../../components/statusbar/StatusBar";

interface BasePageProps {
    children: any
}

interface BasePageState {
    firstName: string,
    lastName: string,
    email: string,
    serverPing: number
    useGravatar: boolean,
}

export default class BasePage extends Component<BasePageProps, BasePageState> {

    updateUIInterval: NodeJS.Timer|null

    constructor(props: BasePageProps) {
        super(props);
        this.state = {
            firstName: "",
            lastName: "",
            email: "",
            serverPing: 0,
            useGravatar: false,
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
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
        let requestEnd = new Date(Date.now());
        this.setState({ serverPing: requestEnd.getMilliseconds() - requestStart.getMilliseconds() });
        if (status === RequestStatus.OK && statusCode === 200) {
            this.setState({
                firstName: responseData.first_name,
                lastName: responseData.last_name,
                email: responseData.email
            })
        }
    }

    retrieveSettings = async()=>{
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/instance-settings`, "GET", null);
        if (status === RequestStatus.OK && statusCode === 200) {
            this.setState({
                useGravatar: responseData.use_gravatar
            })
        }
    }

    render(): ReactNode {
        return (
            <div className="basepage-container">
                <Navbar firstName={this.state.firstName} lastName={this.state.lastName} email={this.state.email} useGravatar={this.state.useGravatar}/>
                <div className="basepage-content">
                    {this.props.children}
                </div>
                <StatusBar serverPing={this.state.serverPing} />
            </div>
        )
    }
}