import React, { useCallback, useEffect, useState } from "react";
import { Card, Table } from "reactstrap";
import { Runner } from "../../types/runner";
import { AdminListRunners } from "../../api/admin";
import { toast } from "react-toastify";
import { Link } from "react-router-dom";

export function RunnersListTable() {
    const [runners, setRunners] = useState<Runner[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    const fetchRunners = useCallback(async () => {
        const r = await AdminListRunners(5);
        if (r) {
            setRunners(r);
        } else {
            toast.error("Failed to fetch runners");
        }
        setLoading(false);
    }, []);

    useEffect(() => {
        fetchRunners();
    }, [fetchRunners]);

    return (
        <React.Fragment>
            <Card body>
                <h3 className="mb-2">Runners</h3>
                <Table>
                    <thead>
                        <th className="p-2">Name</th>
                        <th className="p-2">Type</th>
                        <th className="p-2">Last Contact</th>
                        <th className="p-2">Status</th>
                    </thead>
                    <tbody>
                        {!loading && (
                            <React.Fragment>
                                {runners.length === 0 && (
                                    <tr>
                                        <td colSpan={4} className="text-center p-2">
                                            No runners found.
                                        </td>
                                    </tr>
                                )}
                                {runners.map((runner) => (
                                    <tr key={runner.id}>
                                        <td className="p-2">
                                            <Link to={`/admin/runners/${runner.id}`}>
                                                <b>{runner.name}</b>
                                            </Link>
                                        </td>
                                        <td className="p-2">{runner.type}</td>
                                        <td className="p-2">{runner.last_contact ? new Date(runner.last_contact).toLocaleString() : "Never"}</td>
                                        <td className="p-2">
                                            {
                                                new Date().getTime() - new Date(runner.last_contact || "").getTime() < 5 * 60 * 1000 ?
                                                    (
                                                        <React.Fragment>
                                                            <span className="text-success">●</span> Online
                                                        </React.Fragment>
                                                    ) : (
                                                        <React.Fragment>
                                                            <span className="text-danger">●</span> Offline
                                                        </React.Fragment>
                                                    )
                                            }
                                        </td>
                                    </tr>
                                ))}
                                <tr>
                                    <td colSpan={4} className="text-center p-2">
                                        <Link to="/admin/runners">View All Runners</Link>
                                    </td>
                                </tr>
                            </React.Fragment>
                        )}
                    </tbody>
                </Table>
            </Card>
        </React.Fragment>
    )
}
