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
    value?: string
    helpText?: string
    errorMessage?: string
}

export default class TextInput extends Component<TextInputProps> {
    render(): ReactNode {
        let showError = this.props.errorMessage !== undefined && this.props.errorMessage !== "";
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
                    value={this.props.value}
                    className={showError ? "error" : ""}
                />
                {
                    showError ?
                        (
                            <p style={{ margin: 0 }}>
                                <small style={{ color: "var(--red)", fontSize: "11.5px" }}>{this.props.errorMessage}</small>
                            </p>
                        )
                        :
                        null
                }
                <small style={{ color: "var(--grey-400)", fontSize: "11.5px" }}>{this.props.helpText}</small>
            </div>
        )
    }
}