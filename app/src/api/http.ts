import axios, { AxiosRequestConfig, Method } from "axios";
import { LoginStatus, RequestStatus } from "./types";
import { GetCookie } from "./cookies";

export class Http {

    public static GetServerURL() {
        if (process.env.NODE_ENV === "development") {
            return "http://127.0.0.1:8080"
        } else {
            return "";
        }
    }

    private static GetJWTTokenFromCookies(): string | null {
        return GetCookie("jwtToken");
    }

    /**
     * Users authentication using /api/v1/auth/login api
     * @param email user's email address
     * @param password user's password
     * @returns result, token
     */
    public static async Login(email: string, password: string): Promise<[LoginStatus, string, Date]> {
        // prepare
        let requestUrl = `${this.GetServerURL()}/api/v1/auth/login`;
        let requestBody = JSON.stringify({
            email: email,
            password: password
        });

        // send request
        try {
            let response = await axios.post(
                requestUrl, requestBody,
                {
                    withCredentials: true,
                }
            );
            if (response.status >= 200 && response.status <= 299) {
                return [LoginStatus.OK, response.data.token, new Date(response.data.expiration)];
            }
        } catch (error) {
            if (axios.isAxiosError(error)) {
                if (error.status === 401) {
                    return [LoginStatus.INVALID_CREDENTIALS, "", new Date(Date.now())];
                }
            }
        }
        return [LoginStatus.UNKNOWN_ERROR, "", new Date(Date.now())];
    }

    /**
     * Users authentication using /api/v1/auth/login api
     * @param email user's email address
     * @param password user's password
     * @returns result, token
     */
    public static async SignUp(email: string, password: string, firstName: string, lastName: string): Promise<[LoginStatus, number]> {
        // prepare
        let requestUrl = `${this.GetServerURL()}/api/v1/auth/signup`;
        let requestBody = JSON.stringify({
            email: email,
            password: password,
            first_name: firstName,
            last_name: lastName
        });

        // send request
        try {
            let response = await axios.post(
                requestUrl, requestBody
            );
            if (response.status >= 200 && response.status <= 299) {
                return [LoginStatus.OK, response.status];
            }
        } catch (error) {
            if (axios.isAxiosError(error)) {
                return [LoginStatus.UNKNOWN_ERROR, error.response?.status || -1];
            }
        }
        return [LoginStatus.UNKNOWN_ERROR, -1];
    }

    public static async Request(url: string, method: Method, requestBody: any, contentType: string = "application/json"): Promise<[status: RequestStatus, statusCode: number | undefined, responseData: any, description: string]> {
        let errorDescription = "";

        let requestConfig: AxiosRequestConfig = {
            url: url,
            headers: {
                "Content-Type": contentType
            },
            method: method,
            data: requestBody,
        }

        try {
            let response = await axios.request(requestConfig);
            return [RequestStatus.OK, response.status, response.data, errorDescription];
        } catch (error) {
            if (axios.isAxiosError(error)) {
                if (error.response?.status === 401) {
                    return [RequestStatus.NOT_AUTHENTICATED, error.response?.status, error.response?.data, errorDescription];
                } else {
                    return [RequestStatus.UNKNOWN_ERROR, error.response?.status, error.response?.data, errorDescription];
                }
            }
            return [RequestStatus.UNKNOWN_ERROR, -1, null, errorDescription];
        }
    }
}