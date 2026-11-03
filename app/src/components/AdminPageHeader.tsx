import React from "react";

interface AdminPageHeaderProps {
  title: string;
  description?: React.ReactNode;
  actions?: React.ReactNode;
}

export function AdminPageHeader({
  title,
  description,
  actions,
}: AdminPageHeaderProps) {
  return (
    <div className="admin-page-heading">
      <div className="admin-page-heading-content">
        <h1>{title}</h1>
        {description && (
          <div className="admin-page-heading-description">{description}</div>
        )}
      </div>
      {actions && (
        <div className="admin-page-heading-actions">{actions}</div>
      )}
    </div>
  );
}
