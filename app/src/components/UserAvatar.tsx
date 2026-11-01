
import React from "react";

const AVATAR_COLORS = [
    "#7048e8", // Violet
    "#206bc4", // Blue
    "#4299e1", // Azure
    "#ae3ec9", // Purple
    "#d6336c", // Pink
    "#e8590c", // Orange
    "#2f9e44", // Green
    "#0ca678", // Teal
    "#e03131", // Red
    "#1098ad", // Cyan
];

function hashEmail(email: string): number {
    let hash = 0;

    for (let i = 0; i < email.length; i++) {
        hash = (hash * 31 + email.charCodeAt(i)) | 0;
    }

    return Math.abs(hash);
}

export function UserAvatar({
    firstName,
    lastName,
    email,
}: {
    firstName: string;
    lastName: string;
    email: string;
}) {
    const colorIndex = hashEmail(email.trim().toLowerCase()) % AVATAR_COLORS.length;
    const backgroundColor = AVATAR_COLORS[colorIndex];

    return (
        <div
            className="avatar avatar-sm d-flex align-items-center justify-content-center"
            style={{
                width: "35px",
                height: "35px",
                backgroundColor: `${backgroundColor}28`,
                color: backgroundColor,
            }}
        >
            {firstName.charAt(0).toUpperCase()}
            {lastName.charAt(0).toUpperCase()}
        </div>
    );
}
