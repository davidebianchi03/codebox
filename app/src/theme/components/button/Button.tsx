import { CSSProperties, MouseEventHandler } from "react"
import "./Button.css"

interface ButtonProps {
    children: any
    style?: CSSProperties
    onClick?: MouseEventHandler<HTMLButtonElement>
    type?: "button" | "submit" | "reset" | "link" | undefined
    extraClass?: "primary" | "warning" | "outline-primary" | "outline-white" | undefined
    linkHref?: string | undefined
}

export default function Button(props: ButtonProps) {

    let btnExtraClass = "primary";

    if (props.extraClass !== undefined) {
        btnExtraClass = props.extraClass;
    }

    if (props.type === "link") {
        return (
            <a className={`button ${btnExtraClass}`} style={props.style} href={props.linkHref}>
                {props.children}
            </a>
        )
    } else {
        return (
            <button className={`button ${btnExtraClass}`} style={props.style} onClick={props.onClick} type={props.type}>
                {props.children}
            </button>
        )
    }
}