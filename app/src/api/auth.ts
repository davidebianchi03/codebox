import axios, { isAxiosError } from "axios";

export enum APILoginCode {
    SUCCESS,
    EMAIL_NOT_VERIFIED,
    ERROR,
}

export async function APILogin(
    email: string,
    password: string,
    rememberMe: boolean
): Promise<{ code: APILoginCode, token: string | null }> {
    try {
        const r = await axios.post<{ token: string }>(`/api/v1/auth/login`, {
            email: email,
            password: password,
            remember_me: rememberMe
        });
        return {
            code: APILoginCode.SUCCESS,
            token: r.data.token,
        };
    } catch (error) {
        if (isAxiosError(error)) {
            if (error.response?.status === 412) {
                return {
                    code: APILoginCode.EMAIL_NOT_VERIFIED,
                    token: null,
                };
            }
        }

        return {
            code: APILoginCode.ERROR,
            token: null,
        };
    }
}

export async function APISignUp(email: string, password: string, firstName: string, lastName: string) {
    try {
        await axios.post(`/api/v1/auth/signup`, {
            email: email,
            password: password,
            first_name: firstName,
            last_name: lastName
        });
        return true;
    } catch {
        return false;
    }
}