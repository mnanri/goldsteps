"use client";

import { useState } from "react";
import { createEvent } from "@/utils/api";
import { useRouter } from "next/navigation";

export default function CreateEventModal({ onClose }: { onClose: () => void }) {
    const router = useRouter();

    // Define form state
    const [form, setForm] = useState({
        title: "",
        description: "",
        start_time: "",
        end_time: "",
        deadline: "",
        status: "",
        tag: "",
    });

    // Define options for `status` and `tag` select fields
    const statusOptions = ["To Do", "In Progress", "Pending", "In Review", "Done"];
    const tagOptions = ["Urgent", "Medium", "Low"];

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
    ) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    // const handleSubmit = async (e: React.FormEvent) => {
    //     e.preventDefault();
    //     try {
    //         await createEvent(form);
    //         router.push("/events"); // Back to events page
    //     } catch (error) {
    //         console.error("Failed to create event:", error);
    //     }
    // };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const formattedData = {
            ...form,
            start_time: form.start_time ? new Date(form.start_time).toISOString() : new Date().toISOString(),
            end_time: form.end_time ? new Date(form.end_time).toISOString() : new Date().toISOString(),
            deadline: form.deadline ? new Date(form.deadline).toISOString() : new Date().toISOString(),
        };

        console.log("Formatted event data:", formattedData); // Debug

        try {
            await createEvent(formattedData);
            // router.push("/events");
            if (onClose) {
                await onClose(); // Refresh events list
            } else {
                console.warn("onClose is not defined, events list will not refresh.");
            }
        } catch (error) {
            console.error("Failed to create an event:", error);
        }
    };

    const handleClose = () => {
        router.push("/events"); // Close modal
    };

    return (
        <div className="modal-overlay">
            <div className="modal">
                <button className="close-button" onClick={handleClose}>
                    &times;
                </button>
                <h1>Create New Event</h1>
                <form onSubmit={handleSubmit}>
                    <label>
                        Title:
                        <input
                            name="title"
                            type="text"
                            value={form.title}
                            onChange={handleChange}
                            required
                        />
                    </label>

                    <label>
                        Description:
                        <textarea
                            name="description"
                            value={form.description}
                            onChange={handleChange}
                            // required
                        />
                    </label>

                    <label>
                        Start Time:
                        <input
                            name="start_time"
                            type="datetime-local"
                            value={form.start_time}
                            onChange={handleChange}
                            // required
                        />
                    </label>

                    <label>
                        End Time:
                        <input
                            name="end_time"
                            type="datetime-local"
                            value={form.end_time}
                            onChange={handleChange}
                            // required
                        />
                    </label>

                    <label>
                        Deadline:
                        <input
                            name="deadline"
                            type="datetime-local"
                            value={form.deadline}
                            onChange={handleChange}
                            required
                        />
                    </label>

                    <label>
                        Status:
                        <select
                            name="status"
                            value={form.status}
                            onChange={handleChange}
                            // required
                        >
                            <option value="" disabled>
                                Select a status
                            </option>
                            {statusOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>
                    </label>

                    <label>
                        Tag:
                        <select
                            name="tag"
                            value={form.tag}
                            onChange={handleChange}
                            // required
                        >
                            <option value="" disabled>
                                Select a tag
                            </option>
                            {tagOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>
                    </label>

                    <button type="submit">Create</button>
                </form>
            </div>

            <style jsx>{`
                .modal-overlay {
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background: rgba(0, 0, 0, 0.5);
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    z-index: 1000;
                }
                .modal {
                    background: white;
                    padding: 2rem;
                    border-radius: 8px;
                    width: 500px;
                    max-width: 90%;
                    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
                    position: relative;
                }
                .close-button {
                    position: absolute;
                    top: 10px;
                    right: 10px;
                    background: transparent;
                    border: none;
                    font-size: 1.5rem;
                    padding: 0 0.5rem;
                    cursor: pointer;
                    transition: background-color 0.3s ease;
                }
                .close-button:hover {
                    background-color: #ccc;
                }
                label {
                display: block;
                margin-bottom: 1rem;
                }
                input,
                textarea,
                select {
                width: 100%;
                padding: 0.5rem;
                margin-top: 0.5rem;
                border: 1px solid #ccc;
                border-radius: 4px;
                }
                button[type="submit"] {
                    background-color: #55BEEE;
                    color: white;
                    border: none;
                    padding: 0.5rem 1rem;
                    border-radius: 4px;
                    cursor: pointer;
                    margin-top: 1rem;
                    transition: background-color 0.3s ease;
                }
                button[type="submit"]:hover {
                    background-color: #749AC7;
                }
            `}</style>
        </div>
    );
}
