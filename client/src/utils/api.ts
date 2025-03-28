import axios from "axios";

// Axios Instance
const apiClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL, // Set in .env.local
    headers: {
        "Content-Type": "application/json",
    },
    // timeout: 5000,
});

// Fetch Data
export const fetchData = async (): Promise<{ message: string }> => {
    try {
        const response = await apiClient.get("/data");
        return response.data;
    } catch (error) {
        console.error("API Error:", error);
        throw new Error("Failed to fetch data");
    }
};

// Get Events
export const getEvents = async () => {
    const response = await apiClient.get("/events");
    return response.data;
};

// Get Event by ID
export const getEventById = async (id: string) => {
    const response = await apiClient.get(`/events/${id}`);
    return response.data;
};

// Post Event
export const createEvent = async (event: any) => {
    const response = await apiClient.post("/events", event);
    return response.data;
};

// Update Event
export const updateEvent = async (id: string, event: any) => {
    const response = await apiClient.put(`/events/${id}`, event);
    return response.data;
};

// Delete Event
export const deleteEvent = async (id: string) => {
    const response = await apiClient.delete(`/events/${id}`);
    return response.data;
};

// Stocks...
// Fetch Stock Data
export const fetchStockData = async (code: string) => {
    try {
        const response = await apiClient.get(`/stocks/${code}`);
        return response.data;
    } catch (error) {
        console.error("API Error:", error);
        throw new Error("Failed to fetch stock data");
    }
};

// Fetch Stock News
export const fetchStockNews = async (code: string) => {
    const response = await apiClient.get(`/stocks/${code}/news`);
    return response.data;
};

// Fetch Bloomberg News
export const fetchNewsArticle = async () => {
    try {
        const response = await apiClient.get("/bloomberg");
        return response.data;
    } catch (error) {
        console.error("API Error:", error);
        throw new Error("Failed to fetch news data");
    }
};

// Search Bloomberg News
export const searchNewsArticles = async (query: string) => {
    const response = await apiClient.get(`/bloomberg/search?q=${encodeURIComponent(query)}`);
    return response.data;
};

// Milestone Functions
// Add Milestone
export const addMilestone = async (article: any) => {
    const response = await apiClient.post("/milestones", article);
    return response.data;
};

// Fetch Milestones
export const getMilestones = async () => {
    const response = await apiClient.get("/milestones");
    return response.data;
};

// Remove Milestone
export const removeMilestone = async (id: string) => {
    const response = await apiClient.delete(`/milestones/${id}`);
    return response.data;
};
