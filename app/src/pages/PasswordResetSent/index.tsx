import React from "react";
import { PageWithCardLayout } from "../../layouts/PageWithCardLayout";
import { Link } from "react-router-dom";

export function PasswordResetSentPage() {
    return (
        <React.Fragment>
            <PageWithCardLayout
                title="Password Reset Requested"
            >
                <p>
                    If an account with the email address you provided exists, you will receive
                    an email with instructions on how to reset your password.
                </p>
                <Link to="/login" className="btn btn-light w-100">
                    Back to Login
                </Link>
            </PageWithCardLayout>
        </React.Fragment>
    );
}
