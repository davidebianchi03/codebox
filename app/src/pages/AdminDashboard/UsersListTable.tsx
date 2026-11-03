import React, { useCallback, useEffect, useState } from "react";
import { Card, Table } from "reactstrap";
import { AdminUser } from "../../types/user";
import { AdminListUsers } from "../../api/users";
import { Link } from "react-router-dom";

export function UsersListTable() {
    const [users, setUsers] = useState<AdminUser[]>([]);

    const FetchUsers = useCallback(async () => {
        const u = await AdminListUsers(5);
        if (u) {
            setUsers(u);
        }
    }, []);

    useEffect(() => {
        FetchUsers();
    }, [FetchUsers]);

    return (
        <Card body className="dashboard-panel dashboard-table">
            <div className="dashboard-panel-heading">
                <div className="dashboard-panel-title">
                    <h2>Users</h2>
                </div>
                <Link to="/admin/users" className="dashboard-panel-link">View all users <span aria-hidden="true">→</span></Link>
            </div>
            <div className="table-responsive dashboard-data-table-wrapper mt-3">
                <Table className="table table-hover dashboard-data-table">
                    <thead className="table-light border-0">
                        <tr className="border-0">
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Name</div></th>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Email</div></th>
                            <th scope="col" className="border-0"><div className="d-flex align-items-center py-2 text-muted">Last login</div></th>
                        </tr>
                    </thead>
                    <tbody>
                        {users.map((u, idx) => (
                            <tr key={idx}>
                                <td className="p-2">
                                    <Link to={`/admin/users/${u.email}`}>
                                        <b>{u.first_name} {u.last_name}</b>
                                    </Link>
                                </td>
                                <td className="p-2">{u.email}</td>
                                <td className="p-2">{u.last_login ? new Date(u.last_login).toLocaleString() : "Never logged in"}</td>
                            </tr>
                        ))}
                    </tbody>
                </Table>
            </div>
        </Card>
    );
}
