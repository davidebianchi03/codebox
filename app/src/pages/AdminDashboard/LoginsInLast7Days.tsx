import React from "react";
import { Card } from "reactstrap";
import { SvgAreaChart } from "../../components/SvgAreaChart";

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
    }))

    return (
        <React.Fragment>
            <Card body>
                <h3 className="mb-2">Recent Activity</h3>
                <p>Logins in the last 7 days</p>
                <SvgAreaChart
                    data={computedData}
                    paddingTop={60}
                    height={260}
                    strokeWidth={4}
                />
            </Card>
        </React.Fragment>
    );
}