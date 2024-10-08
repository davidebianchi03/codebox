import { Component,CSSProperties, ReactNode } from "react"
import "./Card.css"

interface CardProps {
    children: any
    style? : CSSProperties
}

export default class Card extends Component<CardProps> {
    render(): ReactNode {
        return(
            <div className="card" style={this.props.style}>
                {this.props.children}
            </div>
        )
    }
}