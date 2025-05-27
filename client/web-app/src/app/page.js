'use client'
import Image from "next/image";
import styles from "./page.module.css";
import LoginPage from "./home/page";
import { useEffect, useState } from "react";

export default function Home() {
 

  // useEffect(()=>{
  //   fetch("http://localhost:8080/ping")
  //   .then((res) => res.json())
  //   .then((data) => setMsg(data.message))
  //   .catch((err) => console.error("Error fetching from backend:", err));
  // },[])

  return (
 <>
 <div style={{margin:`1%`}}>
 <LoginPage />
 </div>
 </>
  );
}
