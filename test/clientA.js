// clientA.js
import fetch from "node-fetch";
import WebSocket from "ws";

const DomainPort = "localhost:1337"
const API_BASE = `http://${DomainPort}`;
const WS_BASE = `ws://${DomainPort}`

async function main() {
   try {
      // 1. Create room to get valid roomId and clientId
      const res = await fetch(`${API_BASE}/api/rooms/create`, { method: "POST" });
      if (!res.ok) throw new Error(`Create room failed: ${res.status}`);
      const data = await res.json();
      console.log("Room created:", data);

      // 2. Connect WebSocket
      const wsUrl = `${WS_BASE}/ws?roomId=${data.roomId}&clientId=${data.clientId}`;
      const socket = new WebSocket(wsUrl);

      socket.on("open", () => {
         console.log("✅ WebSocket connection opened");
         setInterval(() => {
            socket.send(`Hello 👋 from Client A at ${new Date().toLocaleTimeString()}`);
         }, 5000);
      });

      socket.on("message", (data) => {
         console.log("📩 Received from server:", data.toString());
      });

      socket.on("close", () => {
         console.log("🔌 WebSocket closed");
      });

      socket.on("error", (err) => {
         console.error("❌ WebSocket error:", err);
      });

   } catch (err) {
      console.error("Error in clientA:", err);
   }
}

main();
