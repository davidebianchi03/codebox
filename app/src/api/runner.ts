import axios from "axios";
import { Runner, RunnerType } from "../types/runner";

export async function ListRunners(): Promise<Runner[] | undefined> {
    try {
        const r = await axios.get<Runner[]>(`/api/v1/runners`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function ListRunnerTypes(): Promise<RunnerType[] | undefined> {
    try {
        const r = await axios.get<RunnerType[]>(`/api/v1/runner-types`);
        return r.data;
    } catch {
        return undefined;
    }
}
