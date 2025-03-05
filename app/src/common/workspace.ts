export function GetWorkspaceStatusColor(status:string | undefined) {
    if (status === "running") {
        return "success";
    } else if (status === "stopping") {
        return "warning";
    } else if (status === "starting") {
        return "info";
    } else if (status === "error") {
        return "danger";
    } else if (status === "deleting") {
        return "warning";
    } else {
        return "secondary";
    }
}

export function GetBeautyNameForStatus(status: string | undefined) {
    if (status === "running") {
        return "Running";
    } else if (status === "stopping") {
        return "Stopping";
    } else if (status === "starting") {
        return "Starting";
    } else if (status === "error") {
        return "Error";
    } else if (status === "deleting") {
        return "Deleting";
    } else {
        return "Stopped";
    }
}