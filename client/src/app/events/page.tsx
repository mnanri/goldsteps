"use client";

import { useState, useEffect, useCallback } from "react";
import { getEvents } from "@/utils/api";
import { useSearchParams, useRouter } from "next/navigation";
import CreateEventModal from "./create/page";
import EventDetailModal from "./[id]/page";

const statusOptions = ["To Do", "In Progress", "Pending", "In Review", "Done"];
const tagOptions = ["Urgent", "Medium", "Low"];

const statusColors: Record<string, string> = {
    "To Do": "#00008B",
    "In Progress": "#333333",
    "Pending": "#666666",
    "In Review": "#999999",
    "Done": "#BBBBBB",
};

const tagColors: Record<string, string> = {
    Urgent: "#E72121",
    Medium: "#E6B422",
    Low: "#AA89BD",
};

export default function EventsPage() {
    const [events, setEvents] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const searchParams = useSearchParams();
    const router = useRouter();

    const modal = searchParams.get("modal");
    const eventId = searchParams.get("id");

    const fetchEvents = useCallback(async () => {
        setLoading(true);
        try {
            const data = await getEvents();
            setEvents(data);
        } catch (err) {
            setError("Failed to fetch events");
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchEvents();
    }, [fetchEvents]);

    const handleModalClose = async () => {
        await fetchEvents();
        router.push("/events");
    };

    if (loading) return <p>Loading...</p>;
    if (error) return <p style={{ color: "red" }}>{error}</p>;

    return (
        <div>
            <h1>Events</h1>
            <button
                onClick={() => router.push("/events?modal=create")}
                className="create-button"
            >
                Create Event
            </button>

            <div className="events-container">
                {statusOptions.map((status) => (
                    <div key={status} className="status-column">
                        <h2
                            style={{
                                backgroundColor: statusColors[status] || "#E0E0E0",
                                color: "white",
                                padding: "0.5rem",
                                textAlign: "center",
                                borderRadius: "4px",
                            }}
                        >
                            {status}
                        </h2>
                        <div className="events-frame">
                            {events
                                .filter((event) => event.status === status)
                                .sort(
                                    (a, b) =>
                                        new Date(a.deadline).getTime() -
                                        new Date(b.deadline).getTime()
                                ) // Sort by deadline
                                .map((event) => {
                                    const isDone = event.status === "Done";
                                    const isUrgent =
                                        event.tag === "Urgent" && event.status !== "Done";

                                    return (
                                        <div
                                            key={event.id}
                                            className="event-item"
                                            style={{
                                                border: isDone
                                                    ? "2px solid #32CD32"
                                                    : "1px solid #ccc",
                                                backgroundColor: isUrgent ? "#F8CACA" : "#f9f9f9",
                                            }}
                                        >
                                            <p>
                                                <strong>Title:</strong> {event.title}
                                            </p>
                                            <p>
                                                <strong>Deadline:</strong>{" "}
                                                {new Intl.DateTimeFormat("ja-JP", {
                                                    year: "numeric",
                                                    month: "2-digit",
                                                    day: "2-digit",
                                                    hour: "2-digit",
                                                    minute: "2-digit",
                                                }).format(new Date(event.deadline))}
                                            </p>
                                            <p>
                                                <strong>Tag:</strong>{" "}
                                                <span
                                                    style={{
                                                        padding: "0.2rem 0.5rem",
                                                        backgroundColor:
                                                            tagColors[event.tag] || "#E0E0E0",
                                                        borderRadius: "4px",
                                                        color: "white",
                                                    }}
                                                >
                                                    {event.tag}
                                                </span>
                                            </p>
                                            <div className="actions">
                                                <button
                                                    onClick={() =>
                                                        router.push(`/events?modal=detail&id=${event.id}`)
                                                    }
                                                    className="edit-button"
                                                >
                                                    Edit
                                                </button>
                                            </div>
                                        </div>
                                    );
                                })}
                        </div>
                    </div>
                ))}
            </div>

            {modal === "create" && <CreateEventModal onClose={handleModalClose} />}
            {modal === "detail" && eventId && (
                <EventDetailModal /* id={eventId} */ onClose={handleModalClose} />
            )}

            <style jsx>{`
                .create-button {
                margin-bottom: 1rem;
                padding: 0.5rem 1rem;
                background: #55beee;
                color: white;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                transition: background 0.3s;
                }
                .create-button:hover {
                background: #749ac7;
                }

                .events-container {
                display: flex;
                gap: 1rem;
                margin-top: 1rem;
                }

                .status-column {
                flex: 1;
                display: flex;
                flex-direction: column;
                }

                .events-frame {
                margin-top: 0.5rem;
                padding: 0.5rem;
                border: 1px solid #ccc;
                border-radius: 4px;
                background-color: #f1f1f1;
                max-height: 400px;
                overflow-y: auto;
                }

                .event-item {
                margin-bottom: 0.5rem;
                padding: 0.5rem;
                border-radius: 4px;
                transition: box-shadow 0.3s ease, transform 0.3s ease;
                }

                .event-item:hover {
                box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
                transform: translateY(-2px);
                background-color: #ffffff;
                }

                .actions {
                margin-top: 0.5rem;
                display: flex;
                gap: 0.5rem;
                }

                .edit-button,
                .delete-button {
                padding: 0.1rem 0.6rem;
                border: none;
                border-radius: 4px;
                cursor: pointer;
                }

                .edit-button {
                background-color: #999999;
                color: white;
                }

                .delete-button {
                background-color: #e72121;
                color: white;
                }

                .edit-button:hover {
                background-color: #749AC7;
                }
            `}</style>
        </div>
    );
}
