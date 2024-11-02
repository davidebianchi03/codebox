import "./StatusBar.css"
import { Component, ReactNode } from "react";

interface StatusBarProps {
    serverPing: number
}

interface StatusBarState {
    downloadSpeed: number
    uploadSpeed: number
    pingTime: number
}

export default class StatusBar extends Component<StatusBarProps, StatusBarState> {

    constructor(props: any) {
        super(props);
        this.state = {
            downloadSpeed: 5,
            uploadSpeed: 5,
            pingTime: 5,
        }
    }

    render(): ReactNode {
        return (
            <div className="statusbar">
                <span style={{ width: "200px", marginLeft: "5px" }}>&copy; Codebox 2024</span>
                <span style={{ width: "100%" }}></span>
                <div style={{ display: "flex", alignItems: "center" }}>
                    <span style={{ width: "75px" }}>Ping: {this.props.serverPing}ms</span>
                </div>
            </div>
        )
    }
}