import axios from "axios";

export async function APILogin(email: string, password: string, rememberMe: boolean) {
    try {
        const r = await axios.post<{ token: string }>(`/api/v1/auth/login`, {
            email: email,
            password: password,
            remember_me: rememberMe
        });
        return r.data.token;
    } catch {
        return undefined;
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