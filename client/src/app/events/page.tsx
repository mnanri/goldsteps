"use client";

import { useState, useEffect } from "react";
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

    useEffect(() => {
        const fetchEvents = async () => {
            try {
                const data = await getEvents();
                setEvents(data);
            } catch (err) {
                setError("Failed to fetch events");
            } finally {
                setLoading(false);
            }
        };

        fetchEvents();
    }, []);

    if (loading) return <p>Loading...</p>;
    if (error) return <p style={{ color: "red" }}>{error}</p>;

    return (
        <div>
            <h1>Events</h1>
            <button
                onClick={() => router.push("/events?modal=create")}
                style={{
                    marginBottom: "1rem",
                    padding: "0.5rem 1rem",
                    background: "#0070f3",
                    color: "white",
                    border: "none",
                    borderRadius: "4px",
                    cursor: "pointer",
                }}
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

            {/* Render a modal */}
            {modal === "create" && <CreateEventModal />}
            {modal === "detail" && eventId && <EventDetailModal />}
        </div>
    );
}
