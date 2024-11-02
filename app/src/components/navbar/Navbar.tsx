import "./Navbar.css";
import CodeboxLogoSquare from "../../assets/images/logo-square.png";
import DefaultAvatar from "../../assets/images/default-avatar.png";
import MenuIcon from "../../assets/images/menu.png";
import { Component, ReactNode, createRef } from "react";
import { Link } from "react-router-dom";
import CryptoJS from "crypto-js";

interface NavbarProps {
    firstName: string
    lastName: string
    email: string
    useGravatar: boolean,
}

interface NavbarState {
    showActionDropdows: boolean
}

export class Navbar extends Component<NavbarProps, NavbarState> {

    mouseDownEventListener: any
    menuDropdownRef: any

    constructor(props: any) {
        super(props);
        this.state = {
            showActionDropdows: false
        }
        this.menuDropdownRef = createRef();
    }

    componentDidMount(): void {
        document.addEventListener("mousedown", this.handleClickOutsideMenuDropDown);
    }

    componentWillUnmount(): void {
        document.removeEventListener("mousedown", this.handleClickOutsideMenuDropDown);
    }

    handleClickOutsideMenuDropDown = (e: MouseEvent) => {
        if (!this.menuDropdownRef.current.contains(e.target)) {
            if (this.state.showActionDropdows) {
                this.setState({ showActionDropdows: false })
            }
        }
    }

    render(): ReactNode {
        return (
            <div className="navbar">
                <div style={{ display: "flex", alignItems: "center" }}>
                    <Link to={"/"}>
                        <img src={CodeboxLogoSquare} alt="Codebox logo" width={"40px"} />
                    </Link>
                    <img src={MenuIcon}
                        className="dropdown-menu-hamburger"
                        alt="Menu"
                        width={"30px"}
                        onClick={() => this.setState({ showActionDropdows: !this.state.showActionDropdows })}
                        style={{
                            marginLeft: "10pt"
                        }}
                    />
                    <ul className="navbar-links" style={this.state.showActionDropdows ? { display: "block" } : undefined} ref={this.menuDropdownRef}>
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