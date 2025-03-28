"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

interface Milestone {
    id: number;
    title: string;
    link: string;
}

interface MilestoneModalProps {
    milestones: Milestone[];
    onClose: () => void;
}

export default function MilestoneModal({ milestones, onClose }: MilestoneModalProps) {
    const router = useRouter();

    const handleClose = () => {
        onClose();
        router.push("/events");
    };

    return (
        <div className="modal-overlay">
            <div className="modal-content">
                <h2>Milestones</h2>
                <div className="milestone-list">
                    {milestones.map((milestone) => (
                        <div key={milestone.id} className="milestone-card">
                            <h3>{milestone.title}</h3>
                            <a href={milestone.link} target="_blank" rel="noopener noreferrer">
                                Read More →
                            </a>
                        </div>
                    ))}
                </div>
                <button onClick={handleClose} className="close-button">Close</button>
            </div>
            <style jsx>{`
                .modal-overlay {
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background: rgba(0, 0, 0, 0.5);
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    padding: 20px;
                }
                .modal-content {
                    background: white;
                    padding: 24px;
                    border-radius: 12px;
                    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
                    text-align: center;
                    width: 90%;
                    max-width: 700px; /* 横幅を広く */
                }
                .milestone-list {
                    display: flex;
                    flex-direction: column;
                    gap: 16px;
                    margin: 20px 0;
                }
                .milestone-card {
                    background: #f9f9f9;
                    padding: 16px;
                    border-radius: 8px;
                    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
                    text-align: left;
                    width: 100%;
                }
                .milestone-card h3 {
                    margin: 0 0 8px;
                    font-size: 18px;
                }
                .milestone-card a {
                    color: #0073e6;
                    text-decoration: none;
                    font-weight: bold;
                }
                .milestone-card a:hover {
                    text-decoration: underline;
                }
                .close-button {
                    margin-top: 10px;
                    padding: 12px 24px;
                    background: #0073e6;
                    color: white;
                    border: none;
                    border-radius: 6px;
                    cursor: pointer;
                    font-size: 16px;
                }
                .close-button:hover {
                    background: #005bb5;
                }
            `}</style>
        </div>
    );
}
