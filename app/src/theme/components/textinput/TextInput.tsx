import "./TextInput.css"
import { Component, ChangeEventHandler, CSSProperties, ReactNode, KeyboardEventHandler } from "react"

interface TextInputProps {
    label?: string
    placeholder?: string
    secure?: boolean
    style?: CSSProperties
    autocomplete?: string
    name?: string
    onTextChanged?: ChangeEventHandler<HTMLInputElement>
    onKeyDown?: KeyboardEventHandler<HTMLInputElement>
}

export default class TextInput extends Component<TextInputProps> {
    render(): ReactNode {
        return (
            <div className="text-input" style={this.props.style}>
                <label>{this.props.label}</label>
                <input
                    type={this.props.secure ? "password" : "text"}
                    placeholder={this.props.placeholder}
                    onChange={this.props.onTextChanged}
                    onKeyDown={this.props.onKeyDown}
                    autoComplete={this.props.autocomplete}
                    name={this.props.name}
                />
            </div>
        )
    }
}