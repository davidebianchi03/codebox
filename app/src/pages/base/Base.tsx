import StatusBar from "../../components/statusbar/StatusBar";
import "./Base.css"
import { Component, ReactNode } from "react";

interface BasePageProps {
    children: any
}

interface BasePageState {
}

export default class BasePage extends Component<BasePageProps, BasePageState> {

    constructor(props: BasePageProps) {
        super(props);
    }

    render(): ReactNode {
        return (
            <div className="basepage-container">
                <div className="basepage-content">
                    {this.props.children}
                </div>
                <StatusBar />
            </div>
        )
    }
}