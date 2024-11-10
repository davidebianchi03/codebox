import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom"
import BasePage from "./base/Base";
import { Http } from "../api/http";
import { RequestStatus } from "../api/types";
import Card from "../theme/components/card/Card";
import Button from "../theme/components/button/Button";
import TextInput from "../theme/components/textinput/TextInput";

interface CreateWorkspaceProps {

}

export default function CreateWorkspace(props: CreateWorkspaceProps) {

    const navigate = useNavigate();
    const [workspaceName, setWorkspaceName] = useState<string>("");
    const [workspaceNameError, setWorkspaceNameError] = useState<string>("");
    const [workspaceGitRepoURL, setWorkspaceGitRepoURL] = useState<string>("");
    const [workspaceGitRepoURLError, setWorkspaceGitRepoURLError] = useState<string>("");
    const [workspaceGitRepoConfigurationFolder, setWorkspaceGitRepoConfigurationFolder] = useState<string>(".devcontainer");
    const [workspaceGitRepoConfigurationFolderError, setWorkspaceGitRepoConfigurationFolderError] = useState<string>("");

    // check that users are authenticated
    useEffect(() => {
        (async () => {
            let [status, statusCode] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
            if (status === RequestStatus.NOT_AUTHENTICATED && statusCode === 401) {
                navigate("/login")
            }
        })();
    }, []);

    // form validation
    const validateWorkspaceName = (value: string) => {
        if (value === "") {
            setWorkspaceNameError("Cannot be empty");
            return false;
        }
        setWorkspaceNameError("");
        return true;
    }

    const validateWorkspaceGitRepositoryURL = (value: string) => {
        if (value === "") {
            setWorkspaceGitRepoURLError("Cannot be empty");
            return false;
        }
        setWorkspaceGitRepoURLError("");
        return true;
    }

    const validateWorkspaceGitRepoConfigurationFolder = (value: string) => {
        return true;
    }

    const validateForm = () => {
        return (
            validateWorkspaceName(workspaceName) &&
            validateWorkspaceGitRepositoryURL(workspaceGitRepoURL) &&
            validateWorkspaceGitRepoConfigurationFolder(workspaceGitRepoConfigurationFolder)
        );
    };

    // handle submit form event
    const handleSubmitForm = async (event: any) => {
        event.preventDefault();
        let formValid = validateForm();
        if (formValid) {
            let [status, statusCode, data, errorDescription] = await Http.Request(
                `${Http.GetServerURL()}/api/v1/workspace`,
                "POST",
                JSON.stringify({
                    name: workspaceName,
                    type: "devcontainer",
                    git_repo_url: workspaceGitRepoURL,
                    git_repo_configuration_folder: workspaceGitRepoConfigurationFolder,
                })
            );
            if (status === RequestStatus.OK) {
                navigate(`/workspaces/${data.id}`)
            } else {
                // TODO: show error message
            }
        }
    };

    return (
        <BasePage>
            <Card style={{
                width: "90%",
                margin: "auto",
                marginTop: "40pt",
                marginBottom: "30pt",
                paddingTop: "10pt",
            }}>
                <h3>Create workspace</h3>
                <form onSubmit={handleSubmitForm}>
                    <TextInput
                        label={"Workspace name"}
                        placeholder={"my-awesome-workspace"}
                        style={{ width: "calc(100% - 15pt)" }}
                        name="workspace-name"
                        onTextChanged={(event) => {
                            setWorkspaceName(event.target.value);
                            validateWorkspaceName(event.target.value);
                        }}
                        errorMessage={workspaceNameError}
                    />
                    <TextInput
                        label={"Git repository URL"}
                        placeholder={"git@github.com:codebox/example.git"}
                        style={{ width: "calc(100% - 15pt)", marginTop: "10pt" }}
                        name="workspace-git-repo"
                        onTextChanged={(event) => {
                            setWorkspaceGitRepoURL(event.target.value);
                            validateWorkspaceGitRepositoryURL(event.target.value);
                        }}
                        errorMessage={workspaceGitRepoURLError}
                    />
                    <TextInput
                        label={"Workspace configuration folder location (relative to git repository root)"}
                        placeholder={".devcontainer"}
                        style={{ width: "calc(100% - 15pt)", marginTop: "10pt" }}
                        name="workspace-git-repo"
                        value={workspaceGitRepoConfigurationFolder}
                        helpText="Location of workspace configuration files in git repository"
                        onTextChanged={(event) => {
                            setWorkspaceGitRepoConfigurationFolder(event.target.value);
                            validateWorkspaceGitRepoConfigurationFolder(event.target.value);
                        }}
                        errorMessage={workspaceGitRepoConfigurationFolderError}
                    />
                    <div style={{ marginTop: "20pt", display: "flex", justifyContent: "end" }}>
                        <Button type="submit">
                            Create workspace
                        </Button>
                    </div>
                </form>
            </Card>
        </BasePage>
    );
}