import "./Navbar.css";
import CodeboxLogoSquare from "../../assets/images/logo-square.png";
import DefaultAvatar from "../../assets/images/default-avatar.png";
import { Component, ReactNode } from "react";
import { Link } from "react-router-dom";
import CryptoJS from "crypto-js";

interface NavbarProps {
    firstName: string
    lastName: string
    email: string
    useGravatar: boolean,
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
                <div style={{display:"flex", alignItems:"center"}}>
                    <img src={CodeboxLogoSquare} alt="Codebox logo" width={"40px"}/>
                    <ul className="navbar-links">
                        <li>
                            <Link to={"/"}>Workspaces</Link>
                        </li>
                        <li>
                            <Link to={"/"}>Users</Link>
                        </li>
                    </ul>
                </div>
                <div className="navbar-right">
                    <div className="navbar-user">
                        <span className="user-details">
                            <span>{this.props.firstName} {this.props.lastName}</span>
                            <small>{this.props.email}</small>
                        </span>
                        { 
                            this.props.useGravatar ?
                            <img src={DefaultAvatar} alt="User avatar" width={"35px"} height={"35px"} /> 
                            :
                            <img src={`https://www.gravatar.com/avatar/${CryptoJS.SHA256(this.props.email)}`} alt="User avatar" width={"35px"} height={"35px"} /> 
                        }
                    </div>
                </div>
            </div>
        );
    }
}