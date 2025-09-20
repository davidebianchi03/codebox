import React from "react";
import { Card, Table } from "reactstrap";

export function RunnersListTable() {
    return (
        <React.Fragment>
            <Card body>
                <h3 className="mb-2">Runners</h3>
                <Table>
                    <thead>
                        <th className="p-2">Name</th>
                        <th className="p-2">Email</th>
                        <th className="p-2">Last Login</th>
                        <th className="p-2">Status</th>
                    </thead>
                    <tbody>

                    </tbody>
                </Table>
            </Card>
        </React.Fragment>
    )
}
