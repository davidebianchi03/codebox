import "../theme/theme.css"
import { Component, ReactNode } from "react";
import Card from "../theme/components/card/Card";
import TextInput from "../theme/components/textinput/TextInput";
import Button from "../theme/components/button/Button";
import CodeboxLogoWhite from "../assets/images/logo-white.png";
import { Http } from "../api/http";
import { LoginStatus, RequestStatus } from "../api/types";
import { Navigate } from "react-router-dom";


interface LoginPageProps {

}

interface LoginPageState {
    loginEmail: string
    loginPassword: string
    errorMessage: string
    redirect: boolean
    redirectUrl: string
}

export default class LoginPage extends Component<LoginPageProps, LoginPageState> {

    constructor(props: any) {
        super(props);
        this.state = {
            loginEmail: "",
            loginPassword: "",
            errorMessage: "",
            redirect: false,
            redirectUrl: "",
        }
    }

    componentDidMount(): void {
        this.IsAuthenticated();
    }

    private IsAuthenticated = async () => {
        // redirect to home if user is already authenticated
        let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
        if(status === RequestStatus.OK && statusCode === 200) {
            this.setState({ redirect: true, redirectUrl: "/" });
        }
    }

    private SubmitLoginForm = async (event: any) => {
        event.preventDefault();
        
        // validate fields
        if (this.state.loginEmail === "") {
            this.setState({ errorMessage: "Missing email" });
            return;
        }
        if (this.state.loginPassword === "") {
            this.setState({ errorMessage: "Missing password" });
            return;
        }

        // process login
        let [status, jwtToken, expirationDate] = await Http.Login(this.state.loginEmail, this.state.loginPassword);
        if (status === LoginStatus.OK) {
            this.setState({ errorMessage: "" });
            document.cookie = `jwtToken=${jwtToken};expires=${expirationDate.toUTCString()}`;
            this.setState({ redirect: true, redirectUrl: "/" });
        } else {
            document.cookie = `jwtToken=invalidtoken;expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
            if (status === LoginStatus.INVALID_CREDENTIALS) {
                this.setState({ errorMessage: "Invalid credentials" });
            } else {
                this.setState({ errorMessage: "Unknown error, check that server is reachable" });
            }
        }
    }

    render(): ReactNode {
        if (this.state.redirect) {
            return <Navigate to={this.state.redirectUrl} />
        }

        return (
            <div style={{
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                width: "100%",
                height: "100%",
            }}>
                <Card style={{ width: "350px", display: "flex", flexDirection: "column" }}>
                    <div style={{ display: "flex", justifyContent: "center", marginTop: "10pt", marginBottom: "20pt" }}>
                        <img src={CodeboxLogoWhite} style={{ maxWidth: "250px" }} alt="Codebox logo" />
                    </div>
                    <div style={{ textAlign: "center", marginBottom: "10pt", color: "var(--red)" }}>
                        {this.state.errorMessage}
                    </div>
                    <form onSubmit={this.SubmitLoginForm}>
                        <TextInput
                            label={"Email"}
                            placeholder={"john@doe.com"}
                            style={{ width: "calc(100% - 15pt)" }}
                            onTextChanged={(event) => { this.setState({ loginEmail: event.target.value }) }}
                            autocomplete="email"
                            name="email"
                        />
                        <TextInput
                            label={"Password"}
                            placeholder={"password"}
                            secure={true}
                            style={{ width: "calc(100% - 15pt)", marginTop: "10pt" }}
                            onTextChanged={(event) => { this.setState({ loginPassword: event.target.value }) }}
                            autocomplete="password"
                            name="password"
                        />
                        <Button
                            style={{
                                display: "flex",
                                justifyContent: "center",
                                width: "200px",
                                margin: "auto",
                                marginTop: "30pt"
                            }}
                            type="submit"
                        >
                            Login
                        </Button>
                    </form>
                </Card>
            </div>
        )
    }
}