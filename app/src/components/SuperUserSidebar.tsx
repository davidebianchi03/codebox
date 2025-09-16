import { Link, matchPath, useLocation, useNavigate } from "react-router-dom";
import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import { BackhoeIcon, HomeIcon, LicenseIcon, PackagesIcon, PasswordUserIcon, UserIcon } from "../icons/Tabler";
import { Navbar } from "./Navbar";
import { UserDropdown } from "./UserDropdown";

interface SidebarItem {
    title: string;
    icon?: JSX.Element;
    type: "link" | "header";
    link?: string;
}

const SidebarItems: SidebarItem[] = [
    { title: "Overview", type: "header" },
    { title: "Dashboard", icon: <HomeIcon />, type: "link", link: "/admin" },
    { title: "Users", icon: <UserIcon />, type: "link", link: "/admin/users" },
    // { title: "Groups", icon: <GroupIcon />, type: "link", link: "/admin/groups" },
    { title: "Runners", icon: <BackhoeIcon />, type: "link", link: "/admin/runners" },
    { title: "System", type: "header" },
    { title: "Authentication", icon: <PasswordUserIcon />, type: "link", link: "/admin/runners" },
    { title: "Credits", type: "header" },
    { title: "License", icon: <LicenseIcon />, type: "link", link: "/admin/runners" },
    { title: "Third party packages", icon: <PackagesIcon />, type: "link", link: "/admin/runners" },
];


export const SuperUserSidebar = () => {
    const location = useLocation();    

    return (
        <aside className="navbar navbar-vertical navbar-expand-lg" data-bs-theme="dark">
            <div className="container-fluid">
                <button
                    className="navbar-toggler collapsed"
                    type="button"
                    data-bs-toggle="collapse"
                    data-bs-target="#sidebar-menu"
                    aria-controls="sidebar-menu"
                    aria-expanded="false"
                    aria-label="Toggle navigation"
                >
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="navbar-brand navbar-brand-autodark">
                    <Link to={"/"}>
                        <img src={CodeboxLogo} style={{ width: "130px" }} alt="logo" />
                    </Link>
                </div>
                <div className="navbar-nav flex-row d-lg-none">
                    <UserDropdown />
                </div>
                <div className="navbar-collapse collapse" id="sidebar-menu">
                    <ul className="navbar-nav pt-lg-3">
                        {
                            SidebarItems.map((item, index) => {
                                if (item.type === "header") {
                                    return (
                                        <li className="menu-title" key={index}>
                                            <span data-key="t-menu">{item.title}</span>
                                        </li>
                                    );
                                } else if (item.type === "link" && item.link) {
                                    return (
                                        <li className="nav-item" key={index}>
                                            <Link to={item.link} className={`nav-link ${matchPath(item.link, location.pathname) ? "active" : ""}`}>
                                                {item.icon && (
                                                    <span className="nav-link-icon d-md-none d-lg-inline-block">
                                                        {item.icon}
                                                    </span>
                                                )}
                                                <span className="nav-link-title">
                                                    {item.title}
                                                </span>
                                            </Link>
                                        </li>
                                    );
                                }
                                return null;
                            })
                        }
                    </ul>
                </div>
            </div>
        </aside>
    );
}
