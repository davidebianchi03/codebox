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

export enum APISignUpCode {
    SUCCESS,
    CANNOT_SIGNUP,
    ERROR,
}

export async function APISignUp(
    email: string,
    password: string,
    firstName: string,
    lastName: string
) : Promise<APISignUpCode> {
    try {
        await axios.post(`/api/v1/auth/signup`, {
            email: email,
            password: password,
            first_name: firstName,
            last_name: lastName
        });
        return APISignUpCode.SUCCESS;
    } catch (error) {
        if (isAxiosError(error)) {
            if (error.response?.status === 406) {
                return APISignUpCode.CANNOT_SIGNUP;
            }
        }
        return APISignUpCode.ERROR;
    }
}

export enum APIVerifyEmailCode {
    SUCCESS,
    INVALID_CODE,
    EMAIL_ALREADY_VERIFIED,
    USER_LOGGED_IN,
    UNKNOWN_ERROR,
}

export async function APIVerifyEmailAddress(code: string): Promise<APIVerifyEmailCode> {
    try {
        await axios.post(`/api/v1/auth/verify-email-address`, {
            code: code,
        });
        return APIVerifyEmailCode.SUCCESS;
    } catch (error) {
        if (isAxiosError(error)) {
            if (error.response?.status === 406) {
                return APIVerifyEmailCode.INVALID_CODE;
            } else if (error.response?.status === 409) {
                return APIVerifyEmailCode.EMAIL_ALREADY_VERIFIED;
            } else if (error.response?.status === 412) {
                return APIVerifyEmailCode.USER_LOGGED_IN;
            }
        }
        return APIVerifyEmailCode.UNKNOWN_ERROR;
    }
}
