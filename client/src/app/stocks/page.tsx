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
        setNews([]);

        try {
            const stockData = await fetchStockData(stockCode);
            const stockNews = await fetchStockNews(stockCode);

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

    const highlightLink = (title: string) => {
        return title.includes("決算") ? "bg-sky-500 text-white font-bold px-2 py-1 rounded inline-block" : "";
    };

    return (
        <div className="max-w-6xl mx-auto p-6">
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
                    Run
                </button>
            </div>

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

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    {data && (
                        <div className="border rounded p-4">
                            <h2 className="text-xl font-bold">銘柄コード: {data.stockData.code}</h2>
                            <p>株価: &yen; {data.stockData.stock_price} </p>
                            <p>前日終値: &yen; {data.stockData.prev_close}</p>
                            <p>変動額: {data.stockData.price_change}</p>
                            <p>STOP高: {data.stockData.stop_high ? "あり" : "なし"}</p>
                            <p>時価総額: &yen; {data.stockData.market_cap} </p>
                            <p>発行済株数: {data.stockData.issued_shares} 株</p>
                            <p>平均PER: {data.stockData.average_per}</p>
                            <p>平均PBR: {data.stockData.average_pbr}</p>
                        </div>
                    )}

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

                    {data && data.stock && (
                        <div className="border rounded p-4 mt-4">
                            <h2 className="text-xl font-bold">{data.stock.company_name}</h2>
                            <p><strong>英語名:</strong> {data.stock.english_company_name}</p>
                            <p><strong>市場:</strong> {data.stock.listing_market}</p>
                            <p><strong>業種:</strong> {data.stock.industry}</p>
                            <p><strong>所在地:</strong> {data.stock.address}</p>
                            <p><strong>電話番号:</strong> {data.stock.phone}</p>
                            <p><strong>代表者:</strong> {data.stock.representative}</p>
                            <p><strong>決算月:</strong> {data.stock.settlement_month}月</p>
                            <p><strong>資本金:</strong> &yen; {data.stock.capital.toLocaleString()}</p>
                            <p><strong>上場日:</strong> {data.stock.listing_date}</p>
                            <p><strong>単位株数:</strong> {data.stock.unit_shares} 株</p>
                        </div>
                    )}

                    {data && (
                        <div className="border rounded p-4 mt-4">
                            <h2 className="text-xl font-bold">経営</h2>
                            <p><strong>特徴:</strong> {data.stockDetail.feature || "情報なし"}</p>
                            <p><strong>事業内容:</strong> {data.stockDetail.business || "情報なし"}</p>
                            <p><strong>従業員数 (単体):</strong> {data.stockDetail.employees_solo ? data.stockDetail.employees_solo.toLocaleString() : "情報なし"} 人</p>
                            <p><strong>従業員数 (連結):</strong> {data.stockDetail.employees_consolidated ? data.stockDetail.employees_consolidated.toLocaleString() : "情報なし"} 人</p>
                            <p><strong>平均年齢:</strong> {data.stockDetail.average_age ? data.stockDetail.average_age.toFixed(1) + " 歳" : "情報なし"}</p>
                            <p><strong>平均年収:</strong> {data.stockDetail.average_salary ? `¥ ${data.stockDetail.average_salary.toLocaleString()}` : "情報なし"}</p>
                        </div>
                    )}
                </div>

                <div>
                    {news.length > 0 && (
                        <div className="border rounded p-4 h-[85vh] overflow-y-auto">
                            <h2 className="text-xl font-bold">【過去1年間】対象銘柄の施策 & 適時開示</h2>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
                                {news.map((article, index) => (
                                    <div key={index} className="border rounded-lg p-4 shadow-md hover:shadow-lg transition-transform transform hover:scale-105">
                                        <a href={article.link} target="_blank" rel="noopener noreferrer" className={`text-blue-500 hover:underline ${highlightLink(article.title)}`}>
                                            {article.title}
                                        </a>
                                        <p className="text-sm text-gray-600">{article.date} - {article.source}</p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

