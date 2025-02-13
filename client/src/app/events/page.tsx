"use client";

import { useState, useEffect, useCallback } from "react";
import { getEvents, deleteEvent, fetchNewsArticle } from "@/utils/api";
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

    // For Bloomberg news
    const [news, setNews] = useState<any[]>([]);
    const [newsLoading, setNewsLoading] = useState(false);
    const [newsError, setNewsError] = useState<string | null>(null);

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

    const fetchNews = useCallback(async () => {
        setNewsLoading(true);
        setNewsError(null);
        try {
            const data = await fetchNewsArticle();
            setNews(data);
        } catch (err) {
            setNewsError("Failed to fetch news");
        } finally {
            setNewsLoading(false);
        }
    }, []);

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

    // Delete overdue events at once
    const deleteOverdueEvents = async () => {
        const overdueEvents = events.filter(event => countdowns[event.id] === "Overdue");
    
        if (overdueEvents.length === 0) {
            alert("No overdue tasks to delete.");
            return;
        }
    
        if (!confirm(`Are you sure you want to delete ${overdueEvents.length} overdue tasks?`)) {
            return;
        }
    
        setLoading(true);
    
        try {
            await Promise.all(overdueEvents.map(event => deleteEvent(event.id)));
            fetchEvents(); // Re-fetch
        } catch (error) {
            console.error("Failed to delete overdue tasks", error);
            alert("Failed to delete overdue tasks.");
        } finally {
            setLoading(false);
        }
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
                className={`toggle-overdue-button ${showOverdue ? "hide" : "show"}`}
            >
                {showOverdue ? "Hide Overdue" : "Show Overdue"}
            </button>

            <button 
                onClick={deleteOverdueEvents} 
                className="delete-overdue-button"
            >
                Delete Overdue
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
                                                <strong>In: </strong>
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

            {/* Bloomberg News */}
            <button onClick={fetchNews} className="news-button">View Headline</button>

            {newsLoading && <p>News Loading...</p>}
            {newsError && <p style={{ color: "red" }}>{newsError}</p>}
            {news.length > 0 && (
                <div className="news-container">
                    <h2 className="news-title">Bloomberg Latest News</h2>
                    <div className="news-grid">
                        {news.map((article, index) => (
                            <div key={index} className="news-card">
                                <a 
                                    href={article.link} 
                                    target="_blank" 
                                    rel="noopener noreferrer"
                                    className="news-link"
                                >
                                    {article.title}
                                </a>
                                {/* <p className="news-description">{article.description}</p> */}
                                <p>{article.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            )}

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

                .toggle-overdue-button {
                    // margin-top: 1rem;
                    padding: 0.5rem 1rem;
                    border: 2px solid #808080;
                    border-radius: 4px;
                    cursor: pointer;
                    transition: background 0.3s, color 0.3s;
                }

                /* Show Overdue */
                .toggle-overdue-button.show {
                    background: transparent;
                    color: #808080;
                }

                .toggle-overdue-button.show:hover {
                    background: #e0e0e0;
                }

                /* Hide Overdue */
                .toggle-overdue-button.hide {
                    background: #808080;
                    color: white;
                }

                .toggle-overdue-button.hide:hover {
                    background: #666666;
                }

                .delete-overdue-button {
                    margin-top: 1rem;
                    padding: 0.5rem 1rem;
                    background: transparent;
                    color: #e72121;
                    border: 2px solid #e72121;
                    border-radius: 4px;
                    cursor: pointer;
                    transition: background 0.3s;
                }

                .delete-overdue-button:hover {
                    background: #c6000c;
                    color: white;
                }

                .news-button {
                    margin-top: 1rem;
                    padding: 0.5rem 1rem;
                    background: #55beee;
                    color: white;
                    border: none;
                    border-radius: 4px;
                    cursor: pointer;
                    transition: background 0.3s;
                }

                .news-button:hover {
                    background: #749ac7;
                }

                .news-container {
                    margin-top: 1.5rem;
                    padding: 1rem;
                    border: 1px solid #ddd;
                    border-radius: 8px;
                    background-color: #f9f9f9;
                }

                .news-title {
                    font-size: 1.5rem;
                    font-weight: bold;
                    margin-bottom: 1rem;
                }

                .news-grid {
                    display: grid;
                    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
                    gap: 1rem;
                }

                .news-card {
                    padding: 1rem;
                    background: white;
                    border-radius: 8px;
                    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                    transition: transform 0.2s ease, box-shadow 0.2s ease;
                }

                .news-card:hover {
                    transform: translateY(-3px);
                    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
                }

                .news-link {
                    display: block;
                    font-weight: bold;
                    font-size: 1rem;
                    color: #0073e6;
                    text-decoration: none;
                    margin-bottom: 0.5rem;
                }

                .news-link:hover {
                    text-decoration: underline;
                }

                .news-description {
                    font-size: 0.875rem;
                    color: #444;
                    line-height: 1.5;
                    max-height: 0;
                    opacity: 0;
                    overflow: hidden;
                    transition: max-height 0.3s ease, opacity 0.3s ease;
                }

                .news-card:hover .news-description {
                    max-height: 500px;
                    opacity: 1;
                }
            `}</style>
        </div>
    );
}
