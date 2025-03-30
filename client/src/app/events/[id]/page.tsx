"use client";

import { useEffect, useState } from "react";
import { getEventById, updateEvent, deleteEvent } from "@/utils/api";
import { useRouter, useSearchParams } from "next/navigation";

export default function EventDetailModal({ onClose }: { onClose: () => void }) {
    const searchParams = useSearchParams(); // Query parameters
    const id = searchParams.get("id"); // `id` parameter
    console.log("ID from useSearchParams:", id); // Debug

    const router = useRouter();
    const [event, setEvent] = useState<any>(null);
    const [error, setError] = useState<string | null>(null);

    // Define options for `status` and `tag` select fields
    const statusOptions = ["To Do", "In Progress", "Pending", "In Review", "Done"];
    const tagOptions = ["Urgent", "Medium", "Low"];

    useEffect(() => {
        const fetchData = async () => {
            try {
                const data = await getEventById(id as string);
                setEvent(data);
            } catch (err) {
                setError("Failed to fetch an event");
            }
        };

        fetchData();
    }, [id]);

    const handleUpdate = async () => {
        try {
            console.log("Updating event:", event); // Debug
            await updateEvent(id as string, event);
            // alert("Event updated!");
            // router.push("/events"); // Back to events page
            if (onClose) {
                await onClose(); // Refresh events list
            } else {
                console.warn("onClose is not defined, events list will not refresh.");
            }
        } catch (err) {
            console.error("Failed to update the event:", err);
        }
    };

    const handleDelete = async () => {
        try {
            await deleteEvent(id as string);
            // router.push("/events"); // Back to events page
            if (onClose) {
                await onClose(); // Refresh events list
            } else {
                console.warn("onClose is not defined, events list will not refresh.");
            }
        } catch (err) {
            console.error("Failed to delete the event:", err);
        }
    };

    const handleClose = () => {
        router.push("/events"); // Close modal
    };

    if (error) return <p style={{ color: "#E72121" }}>{error}</p>;
    if (!event) return <p>Loading...</p>;

    return (
        <div className="modal-overlay">
            <div className="modal">
                <button className="close-button" onClick={handleClose}>
                    &times;
                </button>
                <h1>Edit Event</h1>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        handleUpdate();
                    }}
                >

                    <label>
                        Title:
                        <input
                            name="title"
                            type="text"
                            value={event.title}
                            onChange={(e) => setEvent({ ...event, title: e.target.value })}
                            required
                        />
                    </label>

                    <label>
                        Description:
                        <textarea
                            name="description"
                            value={event.description}
                            onChange={(e) =>
                                setEvent({ ...event, description: e.target.value })
                            }
                            // required
                        />
                    </label>

                    <label>
                        Start Time:
                        <input
                            name="start_time"
                            type="datetime-local"
                            value={new Date(event.start_time).toISOString().slice(0, 16)}
                            onChange={(e) =>
                                setEvent({ ...event, start_time: e.target.value })
                            }
                            // required
                        />
                    </label>

                    <label>
                        End Time:
                        <input
                            name="end_time"
                            type="datetime-local"
                            value={new Date(event.end_time).toISOString().slice(0, 16)}
                            onChange={(e) =>
                                setEvent({ ...event, end_time: e.target.value })
                            }
                            // required
                        />
                    </label>

                    <label>
                        Deadline:
                        <input
                            name="deadline"
                            type="datetime-local"
                            value={new Date(event.deadline).toISOString().slice(0, 16)}
                            onChange={(e) =>
                                setEvent({ ...event, deadline: e.target.value })
                            }
                            required
                        />
                    </label>

                    <label>
                        Status:
                        <select
                            name="status"
                            value={event.status}
                            onChange={(e) => setEvent({ ...event, status: e.target.value })}
                            // required
                        >
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
                            value={event.tag}
                            onChange={(e) => setEvent({ ...event, tag: e.target.value })}
                            // required
                        >
                            {tagOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>
                    </label>

                    <button type="submit">Update</button>

                    {/* <button type="button" onClick={handleDelete} style={{
                        color: "#E72121",
                        border: "1px solid #E72121",
                        padding: "0.5rem 1rem",
                        borderRadius: "4px",
                        cursor: "pointer",
                        margin: "1rem",
                        // alignSelf: "right",
                    }}>
                        Delete
                    </button> */}
                    <button
                        type="button"
                        className="delete-button"
                        onClick={handleDelete}
                    >
                        Delete
                    </button>
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
                    margin: 1rem 0;
                    transition: background-color 0.3s ease;
                }
                button[type="submit"]:hover {
                    background-color: #749AC7;
                }
                .delete-button {
                    color: #e72121;
                    background-color: transparent;
                    border: 1px solid #e72121;
                    padding: 0.5rem 1rem;
                    border-radius: 4px;
                    cursor: pointer;
                    margin: 1rem;
                    transition: all 0.3s ease;
                }
                .delete-button:hover {
                    color: white;
                    background-color: #e72121;
                }
            `}</style>
        </div>
    );
}
