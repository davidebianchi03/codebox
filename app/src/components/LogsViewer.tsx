import React from "react";
import { SystemLog } from "../types/logs";
import DataTable from "./DataTable";

export function LogsViewer({ logs }: { logs: SystemLog[] }) {
    return (
        <React.Fragment>
            <DataTable
                columns={[
                    {
                        key: "level",
                        label: "Level",
                        render: (value: string) => (
                            <span className={`badge bg-${value === "error" ? "danger" : value === "warn" ? "warning" : "info"} text-light`}>
                                {value.toUpperCase()}
                            </span>
                        )
                    },
                    {
                        key: "timestamp",
                        label: "Timestamp",
                        render: (value: string) => (
                            <small>
                                {new Date(value).toLocaleString()}
                            </small>
                        )
                    },
                    {
                        key: "log",
                        label: "Log",
                        render: (value: string) => (
                            <small>
                                {value}
                            </small>
                        )
                    },
                    {
                        key: "function",
                        label: "Function",
                        render: (value: string, row: SystemLog) => (
                            <small>
                                {row.module}.{value}
                            </small>
                        )
                    },
                ]}
                data={logs.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())}
                initialPageSize={20}
            />
        </React.Fragment>
    )
}
