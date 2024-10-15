export function RetrieveColorForWorkspaceStatus(status: string | undefined) {
    if (status === "creating") {
        return "--cyan";
    } else if (status === "running") {
        return "--green";
    } else if (status === "stopping") {
        return "--orange";
    } else if (status === "starting") {
        return "--blue";
    } else if (status === "error") {
        return "--red";
    } else {
        return "--background-divider";
    }
}

export function RetrieveBeautyNameForStatus(status: string | undefined) {
    if (status === "creating") {
        return "Creating";
    } else if (status === "running") {
        return "Running";
    } else if (status === "stopping") {
        return "Stopping";
    } else if (status === "starting") {
        return "Starting";
    } else if (status === "error") {
        return "Error";
    } else {
        return "Stopped";
    }
}