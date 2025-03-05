import { withRouter } from "../../common/router"
import LogoSquare from "../../assets/images/logo-square.png"
import DefaultAvatar from "../../assets/images/default-avatar.png"
import { Http } from "../../api/http"
import { useCallback, useEffect, useState } from "react"
import { RequestStatus } from "../../api/types"
import { useNavigate } from "react-router-dom"
import { User } from "../../types/user"
import { InstanceSettings } from "../../types/settings"
import sha256 from 'crypto-js/sha256';

type Props = {
    children: string | JSX.Element | JSX.Element[] | (() => JSX.Element)
}

function AuthRequired({ children }: Props) {

    const navigate = useNavigate();
    const [user, setUser] = useState<User | null>(null);
    const [settings, setSettings] = useState<InstanceSettings | null>(null);

    const WhoAmI = useCallback(async () => {
        let [status, statusCode, responseBody] = await Http.Request(`${Http.GetServerURL()}/api/v1/auth/user-details`, "GET", null);
        if (status === RequestStatus.OK && statusCode === 200) {
            setUser(responseBody as User);
        } else {
            navigate("/login");
        }
    }, [navigate]);

    const FetchSettings = useCallback(async () => {
        let [status, statusCode, responseBody] = await Http.Request(`${Http.GetServerURL()}/api/v1/instance-settings`, "GET", null);
        if (status === RequestStatus.OK && statusCode === 200) {
            setSettings(responseBody as InstanceSettings);
        }
    }, []);

    useEffect(() => {
        WhoAmI();
        FetchSettings();
    }, [WhoAmI, FetchSettings]);

    return (
        <>
            <header className="navbar navbar-expand-md d-print-none">
                <div className="container-xl">
                    <a className="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3" href="/">
                        <span className="d-flex align-items-center">
                            <img src={LogoSquare} alt="logo" width={35} />
                            <h2 className="mb-0 ms-2">Codebox</h2>
                        </span>
                    </a>
                    <div className="nav-item dropdown">
                        <span className="nav-link d-flex lh-1 p-0 px-2" data-bs-toggle="dropdown" aria-label="Open user menu">
                            <img className="avatar avatar-sm" src={settings?.use_gravatar && user ? `https://www.gravatar.com/avatar/${sha256(user?.email)}` : DefaultAvatar} alt="avatar" />
                            <div className="d-none d-xl-block ps-2">
                                <div className="text-white">{user?.first_name} {user?.last_name}</div>
                                <div className="mt-1 small text-secondary">{user?.email}</div>
                            </div>
                        </span>
                        <div className="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
                            <a href="./profile.html" className="dropdown-item">Profile</a>
                            <a href="./profile.html" className="dropdown-item">Templates</a>
                            <a href="./profile.html" className="dropdown-item">Users</a>
                            <div className="dropdown-divider"></div>
                            <a href="./settings.html" className="dropdown-item">Settings</a>
                            <a href="./sign-in.html" className="dropdown-item">Logout</a>
                        </div>
                    </div>
                </div>
            </header>
            {children}
        </>
    )
}

export default withRouter(AuthRequired);