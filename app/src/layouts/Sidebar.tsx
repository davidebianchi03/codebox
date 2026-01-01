import { Link, matchPath, useLocation } from "react-router-dom";
import CodeboxLogo from "../assets/images/codebox-logo-white.png";
import { UserDropdown } from "../components/UserDropdown";
import { useEffect, useState } from "react";

export interface SidebarItem {
    title: string;
    icon?: JSX.Element;
    type: "link" | "header";
    link?: string;
    activeOnLinks?: string[];
}

interface SidebarProps {
    sidebarItems: SidebarItem[];
}

export const Sidebar = ({ sidebarItems }: SidebarProps) => {
    const location = useLocation();
    const [isOpen, setIsOpen] = useState<boolean>(false);

    useEffect(() => {
        setIsOpen(false);
    }, [location]);

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
                    onClick={() => setIsOpen(!isOpen)}
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
                <div className={`navbar-collapse collapse ${isOpen && "show"}`} id="sidebar-menu">
                    <ul className="navbar-nav pt-lg-3">
                        {
                            sidebarItems.map((item, index) => {
                                if (item.type === "header") {
                                    return (
                                        <li className="menu-title" key={index}>
                                            <span data-key="t-menu">{item.title}</span>
                                        </li>
                                    );
                                } else if (item.type === "link" && item.link) {
                                    let active = matchPath(item.link, location.pathname);
                                    item.activeOnLinks?.forEach(v => {
                                        active = active || matchPath(v, location.pathname);
                                    })
                                    return (
                                        <li
                                            className={`nav-item ${active ? "active" : ""}`}
                                            key={index}
                                        >
                                            <Link
                                                to={item.link}
                                                className="nav-link">
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
