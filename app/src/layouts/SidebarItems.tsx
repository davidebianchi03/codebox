import { SidebarItem } from "./Sidebar";
import { BackhoeIcon, HomeIcon, ShieldIcon, UserIcon } from "../icons/Tabler";

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
        title: "Instance Settings",
        type: "header"
    },
    {
        title: "Authentication",
        icon: <ShieldIcon />,
        type: "link",
        link: "/admin/authentication-settings"
    },
];
