import { Component, ReactNode, CSSProperties, MouseEventHandler } from "react"
import "./Button.css"

interface ButtonProps {
    children: any
    style?: CSSProperties
    onClick?: MouseEventHandler<HTMLButtonElement>
    type?: "button" | "submit" | "reset" | undefined
}

export default class Button extends Component<ButtonProps> {
    render(): ReactNode {
        return (
            <button className="button" style={this.props.style} onClick={this.props.onClick} type={this.props.type}>
                {this.props.children}
            </button>
        )
    }
}