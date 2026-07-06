import { SidebarItem } from "./Sidebar";
import { BackhoeIcon, EmailIcon, HomeIcon, LogsIcon, ShieldIcon, StatsIcon, UserIcon } from "../icons/Tabler";

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
    {
        title: "Email Sender",
        icon: <EmailIcon />,
        type: "link",
        link: "/admin/email-sender"
    },
    {
        title: "Analytics",
        icon: <StatsIcon />,
        type: "link",
        link: "/admin/analytics"
    },
    {
        title: "Monitoring",
        type: "header"
    },
    {
        title: "System Logs",
        icon: <LogsIcon />,
        type: "link",
        link: "/admin/system-logs"
    },
];
