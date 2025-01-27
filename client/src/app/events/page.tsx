"use client";

import { useState, useEffect, useCallback } from "react";
import { getEvents } from "@/utils/api";
import Link from "next/link";
import { useSearchParams, useRouter } from "next/navigation";
import CreateEventModal from "./create/page";
import EventDetailModal from "./[id]/page";

const statusOptions = ["To Do", "In Progress", "Pending", "In Review", "Done"];
const tagOptions = ["Urgent", "Medium", "Low"];

// Define color mappings
const statusColors: Record<string, string> = {
  "To Do": "#D3D3D3", // Gray
  "In Progress": "#87CEEB", // Light Blue
  "Pending": "#FFD700", // Gold
  "In Review": "#FFA500", // Orange
  "Done": "#32CD32", // Lime Green
};

const tagColors: Record<string, string> = {
  Urgent: "#FF6347", // Tomato Red
  Medium: "#FFD700", // Gold
  Low: "#90EE90", // Light Green
};

export default function EventsPage() {
    const [events, setEvents] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const searchParams = useSearchParams();
    const router = useRouter();

    // `modal` parameter from URL
    const modal = searchParams.get("modal");
    const eventId = searchParams.get("id"); // Event ID

    // Fetch events
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

    // Handle modal close and refresh events list
    const handleModalClose = async () => {
        await fetchEvents(); // refresh events list
        router.push("/events"); // Close modal
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

            <ul className="events-list">
                {events.map((event) => {
                    const isDone = event.status === "Done";
                    const isUrgent = event.tag === "Urgent" && event.status !== "Done";

                    return (
                        <li
                            key={event.id}
                            className="event-item"
                            style={{
                                border: isDone ? "2px solid #32CD32" : "1px solid #ccc",
                                backgroundColor: isUrgent ? "#FFDDDD" : "#f9f9f9", // Light Red
                            }}
                        >
                            <Link
                                href={`/events?modal=detail&id=${event.id}`}
                                className="event-title"
                            >
                                {event.title}
                                <span className="tooltip">{event.description}</span>
                            </Link>
                            <div className="event-info">
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
                                    <strong>Status:</strong>{" "}
                                    <span
                                        style={{
                                            padding: "0.2rem 0.5rem",
                                            backgroundColor: statusColors[event.status] || "#E0E0E0",
                                            borderRadius: "4px",
                                            color: "white",
                                            fontSize: "0.875rem",
                                        }}
                                    >
                                        {event.status}
                                    </span>
                                </p>
                                <p>
                                    <strong>Tag:</strong>{" "}
                                    <span
                                        style={{
                                            padding: "0.2rem 0.5rem",
                                            backgroundColor: tagColors[event.tag] || "#E0E0E0",
                                            borderRadius: "4px",
                                            color: "white",
                                            fontSize: "0.875rem",
                                        }}
                                    >
                                        {event.tag}
                                    </span>
                                </p>
                            </div>
                        </li>
                    );
                })}
            </ul>

            {/* Modal */}
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
                    background: #749AC7;
                }

                .events-list {
                    list-style: none;
                    padding: 0;
                }

                .event-item {
                    margin-bottom: 1rem;
                    padding: 1rem;
                    border: 1px solid #ccc;
                    border-radius: 4px;
                    background-color: #f9f9f9;
                    transition: box-shadow 0.3s ease, transform 0.3s ease;
                }

                .event-item:hover {
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); /* Shadow effect */
                    transform: translateY(-2px); /* Slight lift effect */
                    background-color: #ffffff; /* Optional: change background slightly */
                }

                .event-title {
                    color: #749ac7;
                    text-decoration: underline;
                    position: relative;
                    cursor: pointer;
                }

                /* Tooltip for description */
                .tooltip {
                    visibility: hidden;
                    position: absolute;
                    top: 100%;
                    left: 0;
                    margin-top: 0.5rem;
                    background: #333;
                    color: #fff;
                    padding: 0.5rem;
                    border-radius: 4px;
                    white-space: nowrap;
                    font-size: 0.875rem;
                    z-index: 10;
                    box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
                }

                .event-title:hover .tooltip {
                    visibility: visible;
                }

                .event-info {
                    margin-top: 0.5rem;
                }

                .event-info p {
                    margin: 0.25rem 0;
                }
            `}</style>
        </div>
    );
}
