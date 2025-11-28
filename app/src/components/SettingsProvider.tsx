import React, { useCallback, useEffect } from "react";
import { useDispatch } from "react-redux";
import { FetchSettings } from "../redux/slices/settings";
import { AppDispatch } from "../redux/store";

export interface SettingsProviderProps {
    children?: React.ReactNode;
}

export function SettingsProvider({ children }: SettingsProviderProps) {

    const dispach: AppDispatch = useDispatch();

    const retrieveSettings = useCallback(() => {
        dispach(FetchSettings());
    }, [dispach]);

    useEffect(() => {
        retrieveSettings();
    }, [retrieveSettings]);

    return (
        <React.Fragment>
            {children}
        </React.Fragment>
    )
}