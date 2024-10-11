import "./Navbar.css";
import CodeboxLogoWhite from "../../assets/images/logo-white.png";
import DefaultAvatar from "../../assets/images/default-avatar.png";
import { Component, ReactNode } from "react";
import { Link } from "react-router-dom";
import { Http } from "../../api/http";
import { RequestStatus } from "../../api/types";

interface NavbarProps {
    username: string
}

interface NavbarState {

}

export class Navbar extends Component<NavbarProps, NavbarState> {

    constructor(props: any) {
        super(props);
    }

    render(): ReactNode {
        return (
            <div className="navbar">
                <div></div>
                {/* <img src={CodeboxLogoWhite} alt="Codebox logo" width={"200px"}/> */}
                <div className="navbar-right">
                    <ul className="navbar-links">
                        <li>
                            <Link to={"/"}>Workspaces</Link>
                        </li>
                        <li>
                            <Link to={"/"}>Users</Link>
                        </li>
                    </ul>
                    <span className="divider" />
                    <div className="navbar-user">
                        <span>
                            {this.props.username}
                        </span>
                        <img src={DefaultAvatar} alt="User avatar" width={"40px"} height={"40px"} />
                    </div>
                </div>
            </div>
        );
    }
}