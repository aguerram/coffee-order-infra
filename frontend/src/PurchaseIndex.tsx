import React, { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";

export default function PurchaseIndex() {
    const { coffeeId } = useParams();
    const navigate = useNavigate();
    const [quantity, setQuantity] = useState(1);
    const [error, setError] = useState("");

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");

        try {
            const response = await axios.post(`${import.meta.env.VITE_ORDER_SERVICE_URL}/orders`, {
                coffee_id: parseInt(coffeeId as string),
                quantity: quantity
            });

            alert("Order created successfully!");
            navigate("/");
        } catch (err: any) {
            if (err.response && err.response.data && err.response.data.error) {
                setError(err.response.data.error);
            } else {
                setError("An unexpected error occurred");
            }
        }
    };

    return (
        <div>
            <h1>Purchase Coffee {coffeeId}</h1>
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="quantity">Quantity:</label>
                    <input
                        type="number"
                        id="quantity"
                        value={quantity}
                        onChange={(e) => setQuantity(parseInt(e.target.value))}
                        min="1"
                        required
                    />
                </div>
                <button type="submit">Place Order</button>
            </form>
            {error && <div style={{ color: "red" }}>{error}</div>}
        </div>
    );
}