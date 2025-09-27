import { SidebarItem } from "./Sidebar";
import { BackhoeIcon, HomeIcon, LicenseIcon, PackagesIcon, PasswordUserIcon, UserIcon } from "../icons/Tabler";

export const SuperUserSidebarItems: SidebarItem[] = [
    {
        title: "Overview",
        type: "header"
    },
    {
        title: "Dashboard",
        icon: <HomeIcon />,
        type: "link", link: "/admin"
    },
    {
        title: "Users",
        icon: <UserIcon />,
        type: "link",
        link: "/admin/users",
        activeOnLinks: ["/admin/users", "/admin/users/:userEmail"]
    },
    // { title: "Groups", icon: <GroupIcon />, type: "link", link: "/admin/groups" },
    {
        title: "Runners",
        icon: <BackhoeIcon />, type: "link", link: "/admin/runners",
        activeOnLinks: ["/admin/runners", "/admin/runners/:id"]
    },
    {
        title: "System",
        type: "header"
    },
    {
        title: "Authentication",
        icon: <PasswordUserIcon />,
        type: "link",
        link: "/admin/auth"
    },
    {
        title: "Credits",
        type: "header"
    },
    {
        title: "License",
        icon: <LicenseIcon />,
        type: "link",
        link: "/admin/license"
    },
    {
        title: "Third party packages",
        icon: <PackagesIcon />,
        type: "link",
        link: "/admin/3rd-packages"
    },
];
