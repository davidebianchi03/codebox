import React from "react";
import { ThirdPartyPackages } from "./thirdPartyPackages";
import { Badge, Card, Col, Row } from "react-bootstrap";

export function CreditsPage() {
    return (
        <React.Fragment>
            <div className="col mt-4">
                <h2 className="page-title">Credits</h2>
            </div>
            <Row className="mt-3">
                <p>
                    This software is distributed under an open source license, the source code is available &nbsp;
                    <a href="https://github.com/davidebianchi03/codebox" target="_blank">here</a>.
                </p>
            </Row>
            <Row className="mt-3">
                <Col>
                    <Card>
                        <Card.Header className="border-0 pb-0">
                            <h2>License (MIT)</h2>
                        </Card.Header>
                        <Card.Body className="pt-0">
                            <p style={{ maxWidth: 500, background: "var(--tblr-tertiary-bg)" }} className="p-3 rounded">
                                Copyright (c) {new Date().getFullYear()} Davide Bianchi
                                <br />
                                <br />
                                Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
                                <br />
                                <br />
                                The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
                                <br />
                                <br />
                                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
                            </p>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
            <Row className="mt-3">
                <Col>
                    <Card>
                        <Card.Header className="border-0 pb-0">
                            <h2>Third party packages</h2>
                        </Card.Header>
                        <Card.Body className="pt-0">
                            <p>
                                A special thanks to the open source community
                            </p>
                            <ul className="list-unstyled">
                                {ThirdPartyPackages.map((p, i) => (
                                    <React.Fragment key={i}>
                                        <li className="border-bottom p-2">
                                            <h3 className="mb-1">{p.name}</h3>
                                            <p className="mb-1">{p.author}</p>
                                            <p className="mb-2">
                                                <a href={p.url}>
                                                    {p.url}
                                                </a>
                                            </p>
                                            {p.description && (
                                                <p className="mb-2">
                                                    {p.description}
                                                </p>
                                            )}
                                            <Badge bg="light" className="text-dark">
                                                {p.license}
                                            </Badge>
                                        </li>
                                    </React.Fragment>
                                ))}
                            </ul>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </React.Fragment>
    );
}
