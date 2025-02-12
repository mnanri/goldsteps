"use client";

import { useState } from "react";
import { fetchStockData, fetchStockNews } from "@/utils/api";

export default function StocksPage() {
    const [code, setCode] = useState("");
    const [data, setData] = useState<any>(null);
    const [news, setNews] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [history, setHistory] = useState<string[]>([]); // Save search history

    const handleFetchStockData = async (stockCode: string) => {
        if (!stockCode) {
            setError("銘柄コードを入力してください");
            return;
        }

        setLoading(true);
        setError("");
        setData(null);
        setNews([]); // Initialize news data

        try {
            const stockData = await fetchStockData(stockCode);
            const stockNews = await fetchStockNews(stockCode); // Fetch news data

            setData(stockData);
            setNews(stockNews);

            setHistory((prev) =>
                prev.includes(stockCode) ? prev : [...prev, stockCode]
            );
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Highlight keyword in text
    const highlightLink = (title: string) => {
        return title.includes("決算") ? "bg-gray-600 text-white font-bold px-2 py-1 rounded inline-block" : "";
    };

    return (
        <div className="max-w-xl mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">銘柄検索</h1>

            {/* Search */}
            <div className="flex space-x-2 mb-4">
                <input
                    type="text"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="銘柄コードを入力 (例: 7203)"
                    className="border rounded p-2 w-full focus:outline-none"
                    style={{ borderColor: "#55beee" }}
                />
                <button
                    onClick={() => handleFetchStockData(code)}
                    className="text-white px-4 py-2 rounded transition-colors"
                    style={{ backgroundColor: "#55beee" }}
                    onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#749ac7")}
                    onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#55beee")}
                    disabled={loading}
                >
                    Run
                </button>
            </div>

            {/* Search History */}
            {history.length > 0 && (
                <div className="mb-4">
                    <h2 className="text-lg font-semibold mb-2">検索履歴</h2>
                    <div className="flex flex-wrap gap-2">
                        {history.map((item) => (
                            <button
                                key={item}
                                onClick={() => handleFetchStockData(item)}
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

            {/* Loading or error */}
            {loading && <p>データ取得中...</p>}
            {error && <p className="text-red-500">{error}</p>}

            {/* Stock data */}
            {data && (
                <div className="border rounded p-4 mt-4">
                    <h2 className="text-xl font-bold">銘柄コード: {data.code}</h2>
                    <p>株価: {data.stock_price} 円</p>
                    <p>前日終値: {data.prev_close} 円</p>
                    <p>変動額: {data.price_change}</p>
                    <p>STOP高: {data.stop_high ? "あり" : "なし"}</p>
                    <p>時価総額: {data.market_cap} 円</p>
                    <p>発行済株数: {data.issued_shares} 株</p>
                    <p>平均PER: {data.average_per}</p>
                    <p>平均PBR: {data.average_pbr}</p>
                </div>
            )}

            {/* AI Settle */}
            {data && (
                <div>
                    <p>
                        <a
                            href={`https://minkabu.jp/stock/${code}/settlement_summary`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-500 hover:underline"
                        >
                            決算短信AI Summary
                        </a>    
                    </p>
                </div>
            )}
            
            {/* {code && (
                <div className="mt-6 border rounded-lg shadow-md overflow-hidden">
                    <h2 className="text-xl font-bold mb-2 p-4 bg-gray-100">決算サマリー</h2>
                    <iframe
                        src={`https://minkabu.jp/stock/${code}/settlement_summary`}
                        width="100%"
                        height="800px"
                        className="border-none"
                    />
                </div>
            )} */}

            {/* News */}
            {news.length > 0 && (
                <div className="border rounded p-4 mt-4">
                    <h2 className="text-xl font-bold">【過去1年間】対象銘柄の施策 & 適時開示</h2>
                    <ul className="mt-2">
                        {news.map((article, index) => (
                            <li key={index} className="border-b py-2">
                                <a
                                    href={article.link}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    // className="text-blue-500 hover:underline"
                                    className={`text-blue-500 hover:underline ${highlightLink(article.title)}`}
                                >
                                    {article.title}
                                </a>
                                <p className="text-sm text-gray-600">
                                    {article.date} - {article.source}
                                </p>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
}
