"use client";

import { useState } from "react";
import { searchNewsArticles } from "@/utils/api";

export default function SearchPage() {
    const [query, setQuery] = useState("");
    const [results, setResults] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [history, setHistory] = useState<string[]>([]);

    const handleSearch = async () => {
        if (!query.trim()) {
            setError("Missing words");
            return;
        }

        setLoading(true);
        setError("");
        setResults([]);

        try {
            const data = await searchNewsArticles(query);
            setResults(data);

            setHistory((prev) =>
                prev.includes(query) ? prev : [...prev, query]
            );
        } catch (err: any) {
            setError("Search Failed");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-xl mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">記事検索</h1>

            {/* Search Box */}
            <div className="flex space-x-2 mb-4">
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="検索キーワード（日本語）"
                    className="border rounded p-2 w-4/5 focus:outline-none"
                    style={{ borderColor: "#55beee" }}
                />
                <button
                    onClick={handleSearch}
                    className="text-white px-4 py-2 rounded transition-colors"
                    style={{ backgroundColor: "#55beee" }}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#749ac7")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#55beee")}
                    disabled={loading}
                >
                    Search
                </button>
            </div>

            {/* Search History */}
            {history.length > 0 && (
                <div className="mb-4">
                    {/* <h2 className="text-lg font-semibold mb-2">検索履歴</h2> */}
                    <div className="flex flex-wrap gap-2">
                        {history.map((item) => (
                            <button
                                key={item}
                                onClick={() => {
                                    setQuery(item);
                                    handleSearch();
                                }}
                                className="px-3 py-1 rounded-full text-white"
                                style={{ backgroundColor: "#55beee" }}
                                onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#749ac7")}
                                onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#55beee")}
                            >
                                {item}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {/* Loading or Error */}
            {loading && <p>データ取得中...</p>}
            {error && <p className="text-red-500">{error}</p>}

            {/* Search Results */}
            {results.length > 0 && (
                <div className="mt-6">
                    {/* <h2 className="text-xl font-bold mb-4">検索結果</h2> */}
                    <div className="grid grid-cols-1 gap-4">
                        {results.map((article, index) => (
                            <div 
                                key={index} 
                                className="bg-white p-4 rounded-lg shadow-md transition-transform transform hover:scale-105 hover:shadow-lg"
                            >
                                <a 
                                    href={article.link} 
                                    target="_blank" 
                                    rel="noopener noreferrer"
                                    className="text-blue-600 font-bold text-md hover:underline"
                                >
                                    {article.title}
                                </a>
                                {/* <p className="text-gray-700 mt-2 text-sm line-clamp-3">
                                    {article.description || "説明がありません"}
                                </p> */}
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
