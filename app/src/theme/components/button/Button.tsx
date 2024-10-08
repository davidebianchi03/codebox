import { Component, ReactNode, CSSProperties, MouseEventHandler } from "react"
import "./Button.css"

interface ButtonProps {
    children: any
    style?:CSSProperties
    onClick?:MouseEventHandler<HTMLButtonElement>
}

export default class Button extends Component<ButtonProps> {
    render(): ReactNode {
        return(
            <button className="button" style={this.props.style} onClick={this.props.onClick}>
                {this.props.children}
            </button>
        )
    }
}