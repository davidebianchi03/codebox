import axios from "axios";
import { CLIBuild } from "../types/cli";

export async function ListCLIBuilds(): Promise<CLIBuild[] | undefined> {
    try {
        const response = await axios.get<CLIBuild[]>(`/api/v1/cli`);
        return response.data;
    } catch {
        return undefined;
    }
}
