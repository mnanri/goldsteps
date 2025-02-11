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
    Low: "#C4C4C4",
};

export default function EventsPage() {
    const [events, setEvents] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [countdowns, setCountdowns] = useState<Record<string, string>>({});
    const [showOverdue, setShowOverdue] = useState(false);

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

    // Calculate countdown
    const calculateCountdown = (deadline: string) => {
        const now = new Date().getTime();
        const deadlineTime = new Date(deadline).getTime();
        const timeDiff = deadlineTime - now;

        if (timeDiff <= 0) return "Overdue";

        const days = Math.floor(timeDiff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((timeDiff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((timeDiff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((timeDiff % (1000 * 60)) / 1000);

        return `${days}day ${hours}h ${minutes}min ${seconds}sec`;
    };

    // Update countdowns every second
    useEffect(() => {
        const updateCountdowns = () => {
            const newCountdowns: Record<string, string> = {};
            events.forEach((event) => {
                newCountdowns[event.id] = calculateCountdown(event.deadline);
            });
            setCountdowns(newCountdowns);
        };

        updateCountdowns();
        const interval = setInterval(updateCountdowns, 1000); // Update every seconds

        return () => clearInterval(interval);
    }, [events]);

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

            <br />
            <button 
                onClick={() => setShowOverdue((prev) => !prev)} 
                className="toggle-overdue-button"
            >
                {showOverdue ? "Hide Overdue" : "Show Overdue"}
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
                                .filter((event) => {
                                    const isOverdue = countdowns[event.id] === "Overdue";
                                    return event.status === status && (showOverdue || !isOverdue);
                                })
                                .sort(
                                    (a, b) =>
                                        new Date(a.deadline).getTime() -
                                        new Date(b.deadline).getTime()
                                ) // Sort by deadline
                                .map((event) => {
                                    const isDone = event.status === "Done";
                                    const isUrgent =
                                        event.tag === "Urgent" && event.status !== "Done";
                                    const isOverdue = countdowns[event.id] === "Overdue";

                                    return (
                                        <div
                                            key={event.id}
                                            className="event-item"
                                            style={{
                                                border: isDone
                                                    ? "2px solid #55beee" 
                                                    : isUrgent ? "2px solid #c6000c" : "1px solid #ccc",
                                                backgroundColor: isUrgent ? "rgba(198, 0, 12, 0.1)" : "#f9f9f9",
                                            }}
                                        >
                                            <p>
                                                <strong>Title:</strong> {event.title}
                                            </p>
                                            <p>
                                                <strong>In:</strong>
                                                <span style={{ color: isOverdue ? "red" : "black" }}>
                                                    {countdowns[event.id] || "Counting..."}
                                                </span>
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
                                            <div className="tag-edit-container">
                                                <strong>Severity:</strong>{" "}
                                                <span
                                                    className="tag"
                                                    style={{
                                                        backgroundColor: tagColors[event.tag] || "#E0E0E0",
                                                    }}
                                                >
                                                    {event.tag}
                                                </span>
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
                box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
                transform: translateY(-3px);
                background-color: #ffffff;
                }

                .actions {
                margin-top: 0.5rem;
                display: flex;
                gap: 0.5rem;
                }

                .tag-edit-container {
                    display: flex;
                    
                    align-items: center;
                }
                .tag {
                    padding: 0.2rem 0.5rem;
                    color: white;
                    border-radius: 4px;
                    margin: 0 0.5rem;

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
