import React from "react";
import { Card, Table } from "reactstrap";
import { RunnerAdmin } from "../../types/runner";
import { Link } from "react-router-dom";
import { AdminRunnerStatusBadge } from "../../components/AdminRunnerStatusBadge";

interface RunnersListTableProps {
    runners: RunnerAdmin[] | undefined;
    loading: boolean;
}

export function RunnersListTable({ runners, loading }: RunnersListTableProps) {
    return (
        <Card body className="dashboard-panel dashboard-table">
            <div className="dashboard-panel-heading">
                <div className="dashboard-panel-title">
                    <h2>Runners</h2>
                </div>
                <Link to="/admin/runners" className="dashboard-panel-link">View all runners <span aria-hidden="true">→</span></Link>
            </div>
            <div className="table-responsive dashboard-data-table-wrapper mt-3">
                <Table className="table table-hover dashboard-data-table">
                    <thead className="table-light border-0">
                        <tr>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Name</div></th>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Type</div></th>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Last contact</div></th>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Status</div></th>
                        </tr>
                    </thead>
                    <tbody>
                        {loading ? (
                            <tr><td colSpan={4} className="text-center p-2">Loading runners...</td></tr>
                        ) : runners === undefined ? (
                            <tr><td colSpan={4} className="text-center p-2">Runner data is unavailable.</td></tr>
                        ) : runners.length === 0 ? (
                            <tr><td colSpan={4} className="text-center p-2">No runners found.</td></tr>
                        ) : (
                            <React.Fragment>
                                {runners.slice(0, 5).map((runner) => {
                                    return (
                                        <tr key={runner.id}>
                                            <td className="p-2">
                                                <Link to={`/admin/runners/${runner.id}`}>
                                                    <b>{runner.name}</b>
                                                </Link>
                                            </td>
                                            <td className="p-2">{runner.type}</td>
                                            <td className="p-2">{runner.last_contact ? new Date(runner.last_contact).toLocaleString() : "Never"}</td>
                                            <td className="p-2">
                                                <AdminRunnerStatusBadge runner={runner} />
                                            </td>
                                        </tr>
                                    );
                                })}
                                <tr>
                                    <td colSpan={4} className="text-center p-2">
                                        <Link to="/admin/runners">View all runners</Link>
                                    </td>
                                </tr>
                            </React.Fragment>
                        )}
                    </tbody>
                </Table>
            </div>
        </Card>
    );
}
