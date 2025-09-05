import { initUI } from './ui';
import { WebRtcConnection } from './webrtc'

document.addEventListener('DOMContentLoaded', () => {
  initUI()
});

const apiBase = import.meta.env.VITE_API_BASE_URL;
const wsBase = import.meta.env.VITE_WS_BASE_URL;

const rtc = new WebRtcConnection(apiBase, wsBase);

document.getElementById("create-room")?.addEventListener("click", async () => {
  const data = await rtc.createRoom()
  rtc.log("room created", JSON.stringify(data));
})

document.getElementById("join-room")?.addEventListener("click", async () => {
  const roomId = (document.getElementById("roomId") as HTMLInputElement).value.trim()

  if (roomId) {
    const data = await rtc.joinRoom(roomId)
    rtc.log("room joined", JSON.stringify(data));
  } else alert("empty room id, cannot join")
})

document.getElementById("sendMsg")?.addEventListener("click", async () => {
  const msg = (document.getElementById("messageInput") as HTMLInputElement).value.trim()
  if (msg) rtc.sendWebRTCmessage(msg)
  else alert("empty message, cannot send")
})