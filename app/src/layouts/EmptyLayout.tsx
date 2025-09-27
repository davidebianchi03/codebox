import React from "react";

interface EmptyLayoutProps {
    children: React.ReactNode;
}

export function EmptyLayout({ children }: EmptyLayoutProps) {
    return (
        <React.Fragment>
            {children}
        </React.Fragment>
    );
}
