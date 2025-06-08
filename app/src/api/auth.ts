import axios from "axios";

export async function APILogin(email: string, password: string) {
    try {
        const r = await axios.post<{ token: string }>(`/api/v1/auth/login`, {
            email: email,
            password: password
        });
        return r.data.token;
    } catch {
        return undefined;
    }
}

export async function APISignUp(email: string, password: string, firstName: string, lastName: string) {
    try {
        const r = await axios.post<{ token: string }>(`/api/v1/auth/signup`, {
            email: email,
            password: password,
            first_name: firstName,
            last_name: lastName
        });
        return r.data.token;
    } catch {
        return undefined;
    }
}