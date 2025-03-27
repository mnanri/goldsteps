"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import apiClient from "@/utils/api";

export default function SignupPage() {
    const [alias, setAlias] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const router = useRouter();

    const handleSignup = async () => {
        if (password.length < 8 || !/[a-zA-Z]/.test(password) || !/[0-9]/.test(password)) {
            setError("Password must be at least 8 characters long and contain letters and numbers");
            return;
        }

        try {
            await apiClient.post("/user/create", { alias, password });
            router.push("/auth/login");
        } catch (err: any) {
            setError(err.response?.data?.error || "Failed to sign up");
        }
    };

    return (
        <div className="max-w-md mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">Right, get start!</h1>
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
                onClick={handleSignup}
                className="bg-blue-500 text-white px-4 py-2 rounded w-full"
            >
                Sign up
            </button>
        </div>
    );
}
