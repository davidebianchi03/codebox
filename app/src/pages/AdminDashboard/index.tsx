import React, { useCallback, useEffect, useState } from "react";
import { Card, CardBody, CardHeader, Col, Row, Table } from "reactstrap";
import { Workspace } from "../../types/workspace";
import { toast } from "react-toastify";
import { AdminListWorkspaces } from "../../api/admin";
import Chart from "react-apexcharts";
import ReactApexChart from "react-apexcharts";


const ActivityChart = {
    series: [{
        name: "STOCK ABC",
        data: [30, 40, 45, 50, 49, 60, 70]
    }],
    options: {
        chart: {
            type: 'area',
            height: 350,
            zoom: {
                enabled: false
            }
        },
        dataLabels: {
            enabled: false
        },
        stroke: {
            curve: 'straight'
        },

        title: {
            text: 'Fundamental Analysis of Stocks',
            align: 'left'
        },
        subtitle: {
            text: 'Price Movements',
            align: 'left'
        },
        labels: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
        xaxis: {
            type: 'datetime',
        },
        yaxis: {
            opposite: true
        },
        legend: {
            horizontalAlign: 'left'
        }
    },


}

export function AdminDashboard() {
    const [workspaces, setWorkspaces] = useState<Workspace[]>([]);


    const FetchWorkspaces = useCallback(async () => {
        const w = await AdminListWorkspaces();
        if (w) {
            setWorkspaces(w);
        } else {
            toast.error("Failed to fetch workspaces");
        }
    }, []);

    useEffect(() => {
        FetchWorkspaces();
    }, [FetchWorkspaces]);

    return (
        <React.Fragment>
            <h1>Admin Dashboard</h1>
            <Row>
                <Col md={12} className="mt-5">
                    <Row>
                        <Col md={3} className="mb-4">
                            <Card body>
                                <h3>Total Users</h3>
                                <h1>12</h1>
                            </Card>
                        </Col>
                        <Col md={3} className="mb-4">
                            <Card body>
                                <h3>Active Workspaces</h3>
                                <h1>7</h1>
                            </Card>
                        </Col>
                        <Col md={3} className="mb-4">
                            <Card body>
                                <h3>Online Runners</h3>
                                <h1>5</h1>
                            </Card>
                        </Col>
                        <Col md={3} className="mb-4">
                            <Card body>
                                <h3>Last Login</h3>
                                <h1>5 minutes ago</h1>
                            </Card>
                        </Col>
                    </Row>
                    <Row>
                        <Col md={4} className="mb-4">
                            <Card body>
                                <h3 className="mb-2">Recent Activity</h3>
                                <p>Logins in the last 7 days</p>
                                <ReactApexChart
                                    options={{
                                        chart: {
                                            type: 'area',
                                            height: 350,
                                            zoom: {
                                                enabled: false
                                            },
                                            toolbar: {
                                                show: false
                                            }
                                        },
                                        dataLabels: {
                                            enabled: false
                                        },
                                        stroke: {
                                            curve: 'smooth'
                                        },
                                        labels: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
                                        xaxis: {
                                            // type: 'datetime',
                                            labels: {
                                                style: {
                                                    colors: '#fff'
                                                }
                                            }
                                        },
                                        yaxis: {
                                            opposite: true,
                                            labels: {
                                                style: {
                                                    colors: '#fff'
                                                }
                                            }
                                        },
                                        legend: {
                                            show: false
                                        },
                                        grid: {
                                            show: false
                                        },
                                        tooltip: {
                                            enabled: false,
                                        }
                                    }}
                                    series={ActivityChart.series}
                                    type="area"
                                    height={200}
                                />
                            </Card>
                        </Col>
                        <Col md={8} className="mb-4">
                            <Card body>
                                <h3 className="mb-2">Users</h3>
                                <Table className="table table-vcenter card-table">
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
                        </Col>
                    </Row>
                    <Row>
                        <Col md={12} className="mb-4">
                            <Card body>
                                <h3 className="mb-2">Runners</h3>
                                <Table>
                                    <thead>
                                        <th>Name</th>
                                        <th>Email</th>
                                        <th>Last Login</th>
                                        <th>Status</th>
                                    </thead>
                                    <tbody>

                                    </tbody>
                                </Table>
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </React.Fragment>
    )
}
