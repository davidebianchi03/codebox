import "../assets/css/base.css"
import { Component, ReactNode } from "react";
import BasePage from "./base/Base";
import "../theme/theme.css"
import Card from "../theme/components/card/Card";


export default class HomePage extends Component{
    render(): ReactNode {
        return (
            <div>
                <Card>
                    Pippo
                </Card>
            </div>
            // <BasePage>
            //     <h1>Hello world</h1>
            // </BasePage>
        )
    }
}