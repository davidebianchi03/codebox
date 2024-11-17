import { useState } from "react";
import BasePage from "./base/Base";
import { Http } from "../api/http";
import Card from "../theme/components/card/Card";
import Button from "../theme/components/button/Button";
import TextInput from "../theme/components/textinput/TextInput";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCopy } from '@fortawesome/free-solid-svg-icons'
import { useNavigate } from "react-router-dom";

interface CreateWorkspaceProps {

}

export default function Profile(props: CreateWorkspaceProps) {

    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [sshPublicKey, setSshPublicKey] = useState("");

    const navigate = useNavigate();

    useState(async () => {
        let [, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
        if (statusCode === 200) {
            setFirstName(responseData.first_name);
            setLastName(responseData.last_name);
            setSshPublicKey(responseData.public_key);
        }
    });

    const handleSubmitForm = async (event: any) => {
        event.preventDefault();

        let requestBody = {
            first_name: firstName,
            last_name: lastName
        }

        let [, statusCode, responseData, description] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "PATCH", JSON.stringify(requestBody));
        console.log(description)
        if (statusCode === 200) {
            navigate(0);
        } else {
            console.log(`${statusCode} - ${responseData}`)
        }
    }

    return (
        <BasePage authRequired={true}>
            <Card style={{
                width: "90%",
                margin: "auto",
                marginTop: "40pt",
                marginBottom: "30pt",
                paddingTop: "10pt",
            }}>
                <h3 style={{ padding: 0, margin: 0 }}>Profile</h3>
                <form onSubmit={(event) => handleSubmitForm(event)}>
                    <TextInput
                        style={{ width: "calc(100% - 15pt)", marginTop: "15pt" }}
                        placeholder="John"
                        onTextChanged={(event) => setFirstName(event.target.value)}
                        value={firstName}
                        label="First Name"
                    />
                    <TextInput
                        style={{ width: "calc(100% - 15pt)", marginTop: "15pt" }}
                        placeholder="Doe"
                        onTextChanged={(event) => setFirstName(event.target.value)}
                        value={lastName}
                        label="Last Name"
                    />
                    <div style={{ marginTop: "10pt", display: "flex", justifyContent: "end" }}>
                        <Button type="link" linkHref="/profile" extraClass="outline-white">
                            Cancel
                        </Button>
                        <Button type="submit" style={{ marginLeft: "10pt" }}>
                            Update
                        </Button>
                    </div>
                </form>
            </Card>
            <Card style={{
                width: "90%",
                margin: "auto",
                marginTop: "40pt",
                marginBottom: "30pt",
                paddingTop: "10pt",
            }}>
                <h3 style={{ padding: 0, margin: 0 }}>SSH public key</h3>
                <p style={{ padding: 0, margin: 0 }}>
                    <small>Add this key to your Git server to enable authentication.</small>
                </p>
                <div style={{
                    border: "solid var(--background-divider) 1px",
                    margin: "10pt",
                    borderRadius: "4pt",
                    padding: "10pt",
                    fontSize: "11pt",
                    color: "var(--grey-500)",
                    wordWrap: "break-word",
                    position: "relative",
                    cursor: "pointer",
                }}
                    onClick={() => {
                        navigator.clipboard.writeText(sshPublicKey);
                    }}

                >
                    <span style={{
                        position: "absolute",
                        top: "5pt",
                        right: "5pt",
                    }}>
                        <FontAwesomeIcon icon={faCopy} />
                    </span>
                    <code>
                        {sshPublicKey}
                    </code>
                </div>
            </Card>
        </BasePage>
    );
}