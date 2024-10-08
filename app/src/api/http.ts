import axios from "axios";
import { LoginStatus } from "./types";

export class Http {

    private static GetServerURL() {
        return "http://127.0.0.1:8080"
    }

    /**
     * Users authentication using /api/v1/auth/login api
     * @param email user's email address
     * @param password user's password
     * @returns result, token
     */
    public static async Login(email: string, password: string): Promise<[LoginStatus, string]> {
        // prepare
        let requestUrl = `${this.GetServerURL()}/api/v1/auth/login`;
        let requestBody = JSON.stringify({
            email: email,
            password: password
        });

        // send request
        try {
            let response = await axios.post(
                requestUrl, requestBody
            );
            if (response.status >= 200 && response.status <= 299) {
                return [LoginStatus.OK, response.data.token];
            }
        } catch (error) {
            if (axios.isAxiosError(error)) {
                if(error.status === 401){
                    return [LoginStatus.INVALID_CREDENTIALS, ""];
                }
            }
        }
        return [LoginStatus.UNKNOWN_ERROR, ""];
    }
}