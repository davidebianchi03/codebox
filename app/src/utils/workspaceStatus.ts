export function RetrieveColorForWorkspaceStatus(status: string | undefined) {
    if (status === "creating") {
        return "--cyan";
    } else if (status === "running") {
        return "--green";
    } else if (status === "stopping") {
        return "--yellow";
    } else if (status === "starting") {
        return "--blue";
    } else if (status === "error") {
        return "--red";
    } else if (status === "deleting") {
        return "--orange";
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
    } else if (status === "deleting") {
        return "Deleting";
    } else {
        return "Stopped";
    }
}