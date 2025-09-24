import React, { useCallback, useEffect, useState } from "react";
import { Card, Table } from "reactstrap";
import { User } from "../../types/user";
import { AdminListUsers } from "../../api/admin";
import { Link } from "react-router-dom";

export function UsersListTable() {
    const [users, setUsers] = useState<User[]>([]);

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
        <React.Fragment>
            <Card body>
                <h3 className="mb-2">Users</h3>
                <Table className="table table-vcenter card-table">
                    <thead>
                        <th className="p-2">Name</th>
                        <th className="p-2">Email</th>
                        <th className="p-2">Last Login</th>
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
                        <tr>
                            <td colSpan={3} className="p-2 text-center">
                                <Link to="/admin/users">View all users</Link>
                            </td>
                        </tr>
                    </tbody>
                </Table>
            </Card>
        </React.Fragment>
    )
}
