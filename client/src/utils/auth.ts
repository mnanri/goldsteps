import apiClient from "./api";

export async function signupUser(alias: string, password: string) {
    try {
        const res = await apiClient.post("/user/create", { alias, password });
        return res.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || "Fail to Sign up");
    }
}

export async function loginUser(alias: string, password: string) {
    try {
        const res = await apiClient.post("/auth/login", { alias, password });

        // Save token
        if (typeof window !== "undefined") {
            localStorage.setItem("token", res.data.token);
        }

        return res.data;
    } catch (error: any) {
        throw new Error(error.response?.data?.message || "Fail to Login");
    }
}

export const getCurrentUser = async () => {
    try {
        const response = await apiClient.get("/auth/ref");
        return response.data;
    } catch (error) {
        console.error("Failed to fetch current user", error);
        return null;
    }
};

export const logout = async () => {
    try {
        await apiClient.post("/auth/logout");
    } catch (error) {
        console.error("Failed to log out", error);
    }
};
