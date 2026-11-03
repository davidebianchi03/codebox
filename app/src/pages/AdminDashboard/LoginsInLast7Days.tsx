import React from "react";
import { Card } from "reactstrap";
import { SvgAreaChart } from "../../components/SvgAreaChart";
import { StatsIcon } from "../../icons/Tabler";

interface LoginsInLast7DaysProps {
    data: number[];
}

export function LoginsInLast7Days({ data }: LoginsInLast7DaysProps) {
    const dayNames = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
    const today = new Date(Date.now()).getDay();
    const rotatedDayNames = [...dayNames.slice(today + 1), ...dayNames.slice(0, today + 1)];
    const computedData = data.map((d, i) => ({
        value: d,
        label: rotatedDayNames[i],
    }));

    return (
        <Card body className="dashboard-panel dashboard-activity h-100">
            <div className="dashboard-panel-heading">
                <div className="dashboard-panel-title">
                    <h2>Recent activity</h2>
                </div>
                <span className="dashboard-panel-caption">Sign-ins - last 7 days</span>
            </div>
            <SvgAreaChart
                data={computedData}
                paddingTop={24}
                height={180}
                strokeColor="var(--tblr-primary)"
                strokeWidth={3}
                showXAxis
            />
        </Card>
    );
}