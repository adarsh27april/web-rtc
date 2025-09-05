import { initUI } from './ui';
import { WebRtcConnection } from './webrtc'

document.addEventListener('DOMContentLoaded', () => {
  initUI()
});

const apiBase = import.meta.env.VITE_API_BASE_URL;
const wsBase = import.meta.env.VITE_WS_BASE_URL;

const rtc = new WebRtcConnection(apiBase, wsBase);

const createRoomBtn = document.getElementById("create-room");
createRoomBtn?.addEventListener("click", async () => {
  const data = await rtc.createRoom()
  rtc.log("room created", JSON.stringify(data));
})

const joinRoomBtn = document.getElementById("join-room");
joinRoomBtn?.addEventListener("click", async () => {
  const roomIdInput = document.getElementById("roomId") as HTMLInputElement | null;
  const roomId = roomIdInput?.value.trim(); // Use optional chaining

  if (roomId) {
    const data = await rtc.joinRoom(roomId);
    rtc.log("room joined", JSON.stringify(data));
  } else alert("empty room id, cannot join");
});

const sendMsgBtn = document.getElementById("sendMsg")
sendMsgBtn?.addEventListener("click", async () => {
  const messageInput = document.getElementById("messageInput") as HTMLInputElement | null
  const msg = messageInput?.value.trim()
  if (msg) rtc.sendWebRTCmessage(msg)
  else alert("empty message, cannot send")
})
