"use client";

import { useState, useEffect, useCallback } from "react";
import { getEvents } from "@/utils/api";
import Link from "next/link";
import { useSearchParams, useRouter } from "next/navigation";
import CreateEventModal from "./create/page";
import EventDetailModal from "./[id]/page";

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

            <ul>
                {events.map((event) => (
                    <li key={event.id}>
                        <Link
                            href={`/events?modal=detail&id=${event.id}`}
                            style={{ color: "#0070f3", textDecoration: "underline" }}
                        >
                            {event.title}
                        </Link>
                    </li>
                ))}
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
            `}</style>
        </div>
    );
}
