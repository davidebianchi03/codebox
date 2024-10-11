import { Component, ReactNode } from "react";
import "../theme/theme.css"
import BasePage from "./base/Base";
import { Http } from "../api/http";
import { RequestStatus } from "../api/types";
import { Navigate } from "react-router-dom";
import Card from "../theme/components/card/Card";
import { WorkspaceListItem } from "../components/workspaceListItem/WorkspaceListItem";
import Button from "../theme/components/button/Button";
import TextInput from "../theme/components/textinput/TextInput";

interface HomePageProps {

}

interface HomePageState {
    redirect: boolean
    redirectUrl: string
    workspaces: Array<any>,
    workspacesFilterText: string,
}


export default class HomePage extends Component<HomePageProps, HomePageState> {

    constructor(props: any) {
        super(props);
        this.state = {
            redirect: false,
            redirectUrl: "",
            workspaces: [],
            workspacesFilterText: "",
        }
    }

    componentDidMount(): void {
        this.IsAuthenticated();
        this.ListWorkspaces();
    }

    private IsAuthenticated = async () => {
        // redirect to home if user is already authenticated
        let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/whoami`, "GET", null);
        if (status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
            this.setState({ redirect: true, redirectUrl: "/login" });
        }
    }

    private ListWorkspaces = async () => {
        let [status, statusCode, responseData, errorDescription] = await Http.Request(`${Http.GetServerURL()}/api/v1/workspace`, "GET", null);
        if (status == RequestStatus.OK) {
            this.setState({ workspaces: responseData });
        } else if (status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
            this.setState({ redirect: true, redirectUrl: "/login" });
        } else {
            console.log(`Error: received ${statusCode} from server`);
        }
    }

    render(): ReactNode {
        if (this.state.redirect) {
            return <Navigate to={this.state.redirectUrl} />
        }

        let filteredWorkspaces: Array<any> = [];
        this.state.workspaces.forEach((workspace) => {
            if ((workspace.name as string).indexOf(this.state.workspacesFilterText) != -1) {

                filteredWorkspaces.push(workspace);
            }
        });

        return (
            <BasePage>
                <Card style={{ width: "90%", minWidth: "450px", margin: "auto", marginTop: "40pt", marginBottom: "30pt" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                        <h3>Workspaces</h3>
                        <Button style={{ height: "40px" }}>Create Workspace</Button>
                    </div>
                    <TextInput
                        style={{ width: "calc(100% - 15pt)" }}
                        placeholder="Filter workspaces"
                        onTextChanged={(event) => { this.setState({ workspacesFilterText: event.target.value }) }}
                    />
                    <div style={{ marginTop: "20px" }}>
                        {
                            filteredWorkspaces.length > 0 ?
                                (
                                    filteredWorkspaces.map((workspace) => {
                                        return (
                                            <WorkspaceListItem workspaceName={workspace.name} key={workspace.id} />
                                        )
                                    })
                                )
                                :
                                (
                                    <span style={{ display: "flex", alignItems: "center", justifyContent: "center", width: "100%" }}>
                                        {
                                            this.state.workspaces.length == 0 ?
                                                (
                                                    <span>
                                                        <a style={{ textDecoration: "underline" }}>Create your first workspace</a>
                                                    </span>
                                                ) :
                                                (
                                                    <span>
                                                        No workspace found matching '{this.state.workspacesFilterText}', <a style={{ textDecoration: "underline" }}>create it</a>
                                                    </span>
                                                )
                                        }
                                    </span>
                                )
                        }
                    </div>
                </Card>
            </BasePage>
        )
    }
}