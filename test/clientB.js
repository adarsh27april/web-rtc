// clientB.js (room joiner)

import WebSocket from "ws";
import fetch from "node-fetch";

const DomainPort = "localhost:1337"
const API_BASE = `http://${DomainPort}`;
const WS_BASE = `ws://${DomainPort}`

// Provide the roomId from Client A output
const ROOM_ID = process.argv[2];
if (!ROOM_ID) {
   console.error("Usage: node clientB.js <roomId>");
   process.exit(1);
}

const joinRoom = async (roomId) => {
   const res = await fetch(`${API_BASE}/api/rooms/join?roomId=${roomId}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
   });
   const data = await res.json();
   console.log("Joined room:", data);
   return data;
};

const startClientB = async () => {
   const { clientId } = await joinRoom(ROOM_ID);

   // Connect via WebSocket
   const ws = new WebSocket(`${WS_BASE}/ws?roomId=${ROOM_ID}&clientId=${clientId}`);

   ws.on("open", () => {
      console.log("Client B connected to room:", ROOM_ID);
      setInterval(() => {
         ws.send(`Hi 👊 from Client B at ${new Date().toLocaleTimeString()}`);
      }, 5000);
   });

   ws.on("message", (msg) => {
      console.log("Client B received:", msg.toString());
   });

   ws.on("close", () => console.log("Client B disconnected"));
};

startClientB();
