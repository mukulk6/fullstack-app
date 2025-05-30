"use client";

import React, { useEffect, useState } from 'react';
import NavigationBar from "../../Components/NavigationBar";
import EditableComponent from "../../Components/Editable/EditableComponent";
import { useParams } from 'next/navigation';

export default function NewProductPage() {
    const { id } = useParams();
    const [data, setData] = useState({});

    useEffect(() => {
        const token = localStorage.getItem("token")
        if (id) {
            fetch(`http://localhost:8080/products/${id}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization':`Bearer ${token}`
                },
            })
                .then((res) => res.json())
                .then((data) => {
                    setData(data.product || data); // depends on your backend response
                })
                .catch((err) => console.log(err));
        }
    }, [id]);

    return (
        <>
            <NavigationBar />
            <EditableComponent edit={true} data={data} id={id} />
        </>
    );
}
