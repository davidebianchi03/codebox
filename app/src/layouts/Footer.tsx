import React from "react";
import { Link } from "react-router-dom";

export function Footer() {
    return (
        <React.Fragment>
            <footer className="footer footer-transparent d-print-none">
                <div className="container-xl">
                    <div className="row text-center align-items-center flex-row-reverse">
                        <div className="col-lg-auto ms-lg-auto">
                            <ul className="list-inline list-inline-dots mb-0">
                                {/* <li className="list-inline-item">
                                    <a href="https://docs.tabler.io" target="_blank" className="link-secondary" rel="noopener">Documentation</a>
                                </li> */}
                                <li className="list-inline-item">
                                    <Link to={"/credits"} className="link-secondary">
                                        Credits
                                    </Link>
                                </li>
                            </ul>
                        </div>
                        <div className="col-12 col-lg-auto mt-3 mt-lg-0">
                            <ul className="list-inline list-inline-dots mb-0">
                                <li className="list-inline-item">
                                    Copyright &copy; {new Date().getFullYear()} &nbsp;
                                    <a href="https://github.com/davidebianchi03/codebox/" className="link-secondary">Codebox</a>. All rights reserved.
                                </li>
                                {/* <li className="list-inline-item">
                                    <a href="./changelog.html" className="link-secondary" rel="noopener"> v1.4.0 </a>
                                </li> */}
                            </ul>
                        </div>
                    </div>
                </div>
            </footer>
        </React.Fragment>
    )
}