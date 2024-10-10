import "./Navbar.css";
import CodeboxLogoWhite from "../../assets/images/logo-white.png";
import DefaultAvatar from "../../assets/images/default-avatar.png";
import { Component, ReactNode } from "react";
import { Link } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

interface NavbarProps {

}

interface NavbarState {
    displayedUserName: string
}

export class Navbar extends Component<NavbarProps, NavbarState> {

    constructor(props: any) {
        super(props);
        this.state = {
            displayedUserName: "",
        }
    }

    componentDidMount(): void {
        this.showUserName();
    }

    showUserName = async () => {
        let [status, statusCode, responseData] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/whoami`, "GET", null);
        if (status === RequestStatus.OK && statusCode === 200) {
            if (responseData.first_name && responseData.last_name) {
                this.setState({ displayedUserName: responseData.first_name + " " + responseData.last_name });
            } else {
                this.setState({ displayedUserName: responseData.email });
            }
        }
    }

    render(): ReactNode {
        return (
            <div className="navbar">
                <div></div>
                {/* <img src={CodeboxLogoWhite} alt="Codebox logo" width={"150px"}/> */}
                <div className="navbar-right">
                    <ul className="navbar-links">
                        <li>
                            <Link to={"/"}>Workspaces</Link>
                        </li>
                        <li>
                            <Link to={"/"}>Users</Link>
                        </li>
                    </ul>
                    <span className="divider"/>
                    <div className="navbar-user">
                        <span>
                            {this.state.displayedUserName}
                        </span>
                        <img src={DefaultAvatar} alt="User avatar" width={"40px"} height={"40px"} />
                    </div>
                </div>
            </div>
        );
    }
}