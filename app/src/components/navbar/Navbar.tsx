import "./Navbar.css";
import CodeboxLogoSquare from "../../assets/images/logo-square.png";
import DefaultAvatar from "../../assets/images/default-avatar.png";
import MenuIcon from "../../assets/images/menu.png";
import { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import CryptoJS from "crypto-js";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUser, faRightFromBracket } from '@fortawesome/free-solid-svg-icons'
import { Http } from "../../api/http";


interface NavbarProps {
    firstName: string
    lastName: string
    email: string
    useGravatar: boolean,
}

export function Navbar(props: NavbarProps) {

    const [showActionsDropdown, setShowActionsDropdown] = useState<boolean>(false);
    const [showUserDropdown, setShowUserDropdown] = useState<boolean>(false);

    const actionsDropdown = useRef<any>(null);
    const userDropdown = useRef<any>(null);

    const navigate = useNavigate();

    const handleClickOutsideMenuDropDown = (e: MouseEvent) => {
        if (!actionsDropdown.current.contains(e.target)) {
            setShowActionsDropdown(false);
        }
        if (!userDropdown.current.contains(e.target)) {
            setShowUserDropdown(false);
        }
    };

    const logoutUser = async() => {
        Http.Request(`${Http.GetServerURL()}/api/v1/auth/logout`, "POST", "");
        document.cookie = `jwtToken=loggedout;expires=${new Date(1970, 1, 1).toUTCString()}`;
        navigate("/login");
    }

    useEffect(() => {
        document.addEventListener("mousedown", handleClickOutsideMenuDropDown);
        return () => {
            document.removeEventListener("mousedown", handleClickOutsideMenuDropDown);
        }
    }, []);

    return (
        <div className="navbar">
            {/* Menu */}
            <div style={{ display: "flex", alignItems: "center" }}>
                <Link to={"/"}>
                    <img src={CodeboxLogoSquare} alt="Codebox logo" width={"40px"} />
                </Link>
                <img src={MenuIcon}
                    className="dropdown-menu-hamburger"
                    alt="Menu"
                    width={"30px"}
                    onClick={() => setShowActionsDropdown(!showActionsDropdown)}
                    style={{
                        marginLeft: "10pt"
                    }}
                />
                <ul className="navbar-links" style={showActionsDropdown ? { display: "block" } : undefined} ref={actionsDropdown}>
                    {/* <li>
                        <Link to={"/"}>Workspaces</Link>
                    </li>
                    <li>
                        <Link to={"/"}>Users</Link>
                    </li> */}
                </ul>
            </div>
            {/* User */}
            <div className="navbar-right">
                <div className="navbar-user" onClick={() => setShowUserDropdown(!showUserDropdown)}>
                    {/* User details */}
                    <span className="user-details">
                        <span>{props.firstName} {props.lastName}</span>
                        <small>{props.email}</small>
                    </span>
                    {
                        !props.useGravatar ?
                            <img src={DefaultAvatar} alt="User avatar" width={"35px"} height={"35px"} />
                            :
                            <img src={`https://www.gravatar.com/avatar/${CryptoJS.SHA256(props.email)}`} alt="User avatar" width={"35px"} height={"35px"} />
                    }
                    {/* Dropdown */}
                    <ul className="user-dropdown"
                        style={showUserDropdown ? { display: "block" } : { display: "none" }}
                        ref={userDropdown}
                    >
                        <li onClick={() => navigate("/profile")}>
                            <FontAwesomeIcon icon={faUser} />
                            <span style={{ width: "100%" }}>Profile</span>
                        </li>
                        <li onClick={() => logoutUser()}>
                            <FontAwesomeIcon icon={faRightFromBracket} />
                            <span>Logout</span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    );
}
