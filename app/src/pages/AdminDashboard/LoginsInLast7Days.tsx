import React from "react";
import ReactApexChart from "react-apexcharts";
import { Card } from "reactstrap";

interface LoginsInLast7DaysProps {
    data: number[];
}

export function LoginsInLast7Days({ data }: LoginsInLast7DaysProps) {
    const dayNames = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
    const today = new Date(Date.now()).getDay();
    const rotatedDayNames = [...dayNames.slice(today + 1), ...dayNames.slice(0, today + 1)];

    return (
        <React.Fragment>
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
                        labels: rotatedDayNames,
                        xaxis: {
                            labels: {
                                show: false
                            },
                            axisBorder: {
                                show: false
                            },
                            axisTicks: {
                                show: false
                            }
                        },
                        yaxis: {
                            show: false,
                        },
                        legend: {
                            show: false
                        },
                        grid: {
                            show: false
                        },
                        tooltip: {
                            // enabled: false,
                            custom(options) {
                                return `
                                <div class="p-1 border bg-light text-dark fs-5" style="border-radius: 0.15rem;">
                                    ${options.series[options.seriesIndex][options.dataPointIndex]} logins
                                </div>`;
                            },
                        }
                    }}
                    series={[{
                        name: "Logins in last 7 days",
                        data: data
                    }]}
                    type="area"
                    height={200}
                />
            </Card>
        </React.Fragment>
    );
}