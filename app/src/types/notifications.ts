import { Workspace } from "./workspace";

export interface CodeboxNotification {
    type: string;
    event: string;
    workspace?: Workspace;
}