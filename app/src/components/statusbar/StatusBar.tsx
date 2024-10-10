import "./StatusBar.css"
import { Component, ReactNode } from "react";

type StatusBarState = {
    downloadSpeed: number,
    uploadSpeed: number,
    pingTime: number,
}

export default class StatusBar extends Component{

    state: StatusBarState

    constructor(props:any) {
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
                {/* <span className="app-name">Codebox</span> */}
                <span style={{width: "100%"}}></span>
                <div style={{display: "flex", alignItems: "center"}}>
                    <b style={{marginRight: "10px"}}>Transmission: </b>
                    <span style={{width: "110px"}}>Upload: {this.state.uploadSpeed} MB/s</span>
                    <span style={{width: "130px"}}>Download: {this.state.downloadSpeed}  MB/s</span>
                    <span style={{width: "75px"}}>Ping: {this.state.downloadSpeed}ms</span>
                </div>
            </div>
        )
    }
}