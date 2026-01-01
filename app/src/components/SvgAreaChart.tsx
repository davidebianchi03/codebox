import React, { useMemo, useState } from "react";

interface ChartDataPoint {
    value: number;
    label: string;
}

interface SvgAreaChartProps {
    data: ChartDataPoint[];
    width?: number;
    height?: number;
    strokeColor?: string;
    strokeWidth?: number;
    showXAxis?: boolean;
    paddingTop?: number;
}

interface Point extends ChartDataPoint {
    x: number;
    y: number;
}

export function SvgAreaChart({
    data,
    width = 600,
    height = 200,
    strokeColor = "#fff",
    strokeWidth = 2,
    showXAxis = false,
    paddingTop = 30,
}: SvgAreaChartProps) {
    const paddingBottom = 20;
    const paddingX = 20;

    const [hovered, setHovered] = useState<Point | null>(null);

    const values = data.map(d => d.value);
    const min = Math.min(...values);
    const max = Math.max(...values);

    const points: Point[] = useMemo(() => {
        return data.map((item, index) => {
            const x =
                paddingX +
                (index / (data.length - 1)) * (width - paddingX * 2);

            const y =
                height -
                paddingBottom -
                ((item.value - min) / (max - min || 1)) *
                (height - paddingTop - paddingBottom);

            return { x, y, value: item.value, label: item.label };
        });
    }, [data, width, height, min, max]);

    const linePath = useMemo(() => {
        if (!points.length) return "";

        let d = `M ${points[0].x} ${points[0].y}`;

        for (let i = 1; i < points.length; i++) {
            const prev = points[i - 1];
            const curr = points[i];
            const cx = (prev.x + curr.x) / 2;

            d += ` C ${cx} ${prev.y}, ${cx} ${curr.y}, ${curr.x} ${curr.y}`;
        }

        return d;
    }, [points]);

    const areaPath = `
        ${linePath}
        L ${points[points.length - 1]?.x ?? 0} ${height - paddingBottom}
        L ${points[0]?.x ?? 0} ${height - paddingBottom}
        Z
    `;

    return (
        <React.Fragment>
            <svg width="100%" viewBox={`0 0 ${width} ${height}`}>
                <defs>
                    <linearGradient id="areaGradient" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#ffffff66" />
                        <stop offset="100%" stopColor="#ffffff00" />
                    </linearGradient>
                </defs>

                <path d={areaPath} fill="url(#areaGradient)" />

                <path
                    d={linePath}
                    fill="none"
                    stroke={strokeColor}
                    strokeWidth={strokeWidth}
                />

                {points.map((point, index) => (
                    <circle
                        key={index}
                        cx={point.x}
                        cy={point.y}
                        r={8}
                        fill="transparent"
                        onMouseEnter={() => setHovered(point)}
                        onMouseLeave={() => setHovered(null)}
                    />
                ))}

                {showXAxis && points.map((point, index) => (
                    <text
                        key={`label-${index}`}
                        x={point.x}
                        y={paddingTop - 5}
                        textAnchor="middle"
                        fontSize={12}
                        fill="#fff"
                        opacity={0.9}
                    >
                        {point.label}
                    </text>
                ))}

                {hovered && (
                    <foreignObject
                        x={hovered.x - 50}
                        y={hovered.y - 50}
                        width={100}
                        height={45}
                    >
                        <div
                            style={{
                                background: "#fff",
                                color: "#000",
                                padding: "4px 6px",
                                fontSize: "14px",
                                borderRadius: "4px",
                                textAlign: "center",
                                pointerEvents: "none",
                            }}
                        >
                            <div>{hovered.label}</div>
                            <strong>{hovered.value} logins</strong>
                        </div>
                    </foreignObject>
                )}
            </svg>
        </React.Fragment>
    );
}
