import axios from "axios";

// Axios Instance
const apiClient = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL, // 
    timeout: 5000,
});

// データ取得の関数
export const fetchData = async (): Promise<{ message: string }> => {
  try {
    const response = await apiClient.get("/data");
    return response.data;
  } catch (error) {
    console.error("API Error:", error);
    throw new Error("Failed to fetch data");
  }
};
