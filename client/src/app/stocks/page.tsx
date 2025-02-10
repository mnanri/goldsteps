"use client";

import { useState } from "react";
import { fetchStockData } from "@/utils/api";

export default function StocksPage() {
    const [code, setCode] = useState("");
    const [data, setData] = useState<any>(null);
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

        try {
            const stockData = await fetchStockData(stockCode);
            setData(stockData);
            setHistory((prev) =>
                prev.includes(stockCode) ? prev : [...prev, stockCode]
            );
        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-xl mx-auto p-6">
            <h1 className="text-2xl font-bold mb-4">銘柄検索</h1>

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
                    検索
                </button>
            </div>

            {/* Display search history */}
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

            {loading && <p>データ取得中...</p>}
            {error && <p className="text-red-500">{error}</p>}

            {data && (
                <div className="border rounded p-4 mt-4">
                    <h2 className="text-xl font-bold">銘柄コード: {data.code}</h2>
                    <p>株価: {data.stock_price} 円</p>
                    <p>前日終値: {data.prev_close} 円</p>
                    <p>変動額: {data.price_change}</p>
                    <p>STOP高: {data.stop_high ? "あり" : "なし"}</p>
                    <p>時価総額: {data.market_cap}</p>
                    <p>発行済株数: {data.issued_shares}</p>
                    <p>平均PER: {data.average_per}</p>
                    <p>平均PBR: {data.average_pbr}</p>
                </div>
            )}
        </div>
    );
}
