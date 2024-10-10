import { Component, ReactNode } from "react";
import "../theme/theme.css"
import BasePage from "./base/Base";
import { Http } from "../api/http";
import { RequestStatus } from "../api/types";
import { Navigate } from "react-router-dom";
import Card from "../theme/components/card/Card";

interface HomePageProps {

}

interface HomePageState {
    redirect: boolean
    redirectUrl: string
}


export default class HomePage extends Component<HomePageProps, HomePageState> {

    constructor(props: any) {
        super(props);
        this.state = {
            redirect: false,
            redirectUrl: "",
        }
    }

    componentDidMount(): void {
        this.IsAuthenticated();
    }

    private IsAuthenticated = async () => {
        // redirect to home if user is already authenticated
        let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/whoami`, "GET", null);
        if (status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
            this.setState({ redirect: true, redirectUrl: "/login" });
        }
    }

    render(): ReactNode {
        if (this.state.redirect) {
            return <Navigate to={this.state.redirectUrl} />
        }

        return (
            <BasePage>
                <Card style={{ width: "90%", minWidth: "450px", margin: "auto", marginTop: "40pt" }}>
                    <h3>Workspaces</h3>
                </Card>
            </BasePage>
        )
    }
}