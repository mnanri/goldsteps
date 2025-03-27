"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import apiClient from "@/utils/api";

export default function LoginPage() {
    const [alias, setAlias] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const router = useRouter();

    const handleLogin = async () => {
        try {
            const response = await apiClient.post("/auth/login", { alias, password });
            localStorage.setItem("token", response.data.token);
            router.push("/");
        } catch (err: any) {
            setError(err.response?.data?.error || "Failed to Login");
        }
    };

    return (
        <div className="max-w-md mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">Hello!</h1>
            {error && <p className="text-red-500">{error}</p>}
            <input
                type="text"
                placeholder="Alias"
                value={alias}
                onChange={(e) => setAlias(e.target.value)}
                className="border rounded p-2 w-full mb-2"
            />
            <input
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="border rounded p-2 w-full mb-4"
            />
            <button
                onClick={handleLogin}
                className="bg-blue-500 text-white px-4 py-2 rounded w-full"
            >
                Login
            </button>
        </div>
    );
}
