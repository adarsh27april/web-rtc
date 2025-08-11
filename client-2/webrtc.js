let dataChannel = null;
let latestSDP = null;
const sharedKey = 42; // Very basic XOR key, replace with proper crypto for real-world use

const DomainPort = "localhost:1337"
const API_BASE = `http://${DomainPort}`;
const WS_BASE = `ws://${DomainPort}`;


let ws;
function encrypt(msg, key) {
   return btoa([...msg].map(c => String.fromCharCode(c.charCodeAt(0) ^ key)).join(''));
}

function decrypt(enc, key) {
   const decoded = atob(enc);
   return [...decoded].map(c => String.fromCharCode(c.charCodeAt(0) ^ key)).join('');
}

function log(...args) {
   document.getElementById('log').textContent += args.join(' ') + '\n';
}

document.getElementById("create-room").onclick = async () => {
   log("create-room start")
   const res = await fetch(`${API_BASE}/api/rooms/create`, { method: "POST" });
   const data = await res.json();
   console.log(data);

   roomId = data.roomId;
   clientId = data.clientId;
   isOfferer = true;

   log(`Room created: ${roomId}, Client: ${clientId}`);
   connectWebSocket(roomId, clientId, isOfferer);
};

document.getElementById("join-room").onclick = async () => {
   const inputRoomId = document.getElementById("roomId").value.trim();
   if (!inputRoomId) {
      alert("Please enter a Room ID");
      return;
   }

   const res = await fetch(`${API_BASE}/api/rooms/join?roomId=${inputRoomId}`, { method: "POST" });
   const data = await res.json();

   roomId = data.roomId;
   clientId = data.clientId;
   isOfferer = false;

   log(`Joined room: ${roomId}, Client: ${clientId}`);
   connectWebSocket(roomId, clientId, isOfferer);
};

function connectWebSocket(roomId, clientId, isOfferer) {
   ws = new WebSocket(`${WS_BASE}/ws?roomId=${roomId}&clientId=${clientId}`);

   ws.onopen = () => {
      log("✅ WebSocket connected");
      initWebRTC(isOfferer, sendSignalToWS);
   };

   ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      handleSignalMessage(message);
   };

   ws.onclose = () => log("🔌 WebSocket disconnected");
   ws.onerror = (err) => log("❌ WebSocket error:", err);
}

function sendSignalToWS(type, data) {
   ws.send(JSON.stringify({ type, data }));
}

let pc;
let sendSignalCb;

function initWebRTC(isOfferer, sendSignalFn) {
   sendSignalCb = sendSignalFn;

   pc = new RTCPeerConnection();

   pc.onicecandidate = (event) => {
      if (event.candidate) {
         sendSignalCb("candidate", event.candidate);
      }
   };

   pc.ondatachannel = (event) => {
      dataChannel = event.channel;
      setupDataChannel();
   };

   if (isOfferer) {
      dataChannel = pc.createDataChannel("chat");
      setupDataChannel();

      pc.createOffer().then((offer) => {
         pc.setLocalDescription(offer);
         sendSignalCb("offer", offer);
      });
   }
}

function handleSignalMessage(msg) {
   if (msg.type === "offer") {
      pc.setRemoteDescription(new RTCSessionDescription(msg.data));
      pc.createAnswer().then((answer) => {
         pc.setLocalDescription(answer);
         sendSignalCb("answer", answer);
      });
   } else if (msg.type === "answer") {
      pc.setRemoteDescription(new RTCSessionDescription(msg.data));
   } else if (msg.type === "candidate") {
      pc.addIceCandidate(new RTCIceCandidate(msg.data));
   }
}

function setupDataChannel() {
   dataChannel.onopen = () => console.log("📡 Data channel open");
   dataChannel.onmessage = (e) => console.log("📩 Received:", e.data);
}

function sendMessage(msg) {
   if (dataChannel && dataChannel.readyState === "open") {
      dataChannel.send(msg);
   }
}

document.getElementById("elementId").onclick =()=>{
   const msg = document.getElementById("messageInput").value.trim();
   if (!msg) {
      alert("Please enter a message");
      return;
   }
   sendMessage(msg)
   appendMessage(`You: ${msg}`)
   document.getElementById("messageInput").value = ""
}

function appendMessage(msg) {
    messagesArea.value += msg + "\n";
    messagesArea.scrollTop = messagesArea.scrollHeight;
}