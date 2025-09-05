interface RoleMessage {
   type: 'role';
   data: { role: 'offerer' | 'answerer' };
}

interface OfferMessage {
   type: 'offer';
   data: RTCSessionDescriptionInit;
}

interface AnswerMessage {
   type: 'answer';
   data: RTCSessionDescriptionInit;
}

interface CandidateMessage {
   type: 'candidate';
   data: RTCIceCandidateInit;
}

interface TimeoutMessage {
   type: 'timeout';
   message: string;
}

type SignalingMessage = RoleMessage | OfferMessage | AnswerMessage | CandidateMessage | TimeoutMessage;


export class WebRtcConnection {
   private apiBase: string
   private wsBase: string
   private peerConn: RTCPeerConnection;
   private dataChannel: RTCDataChannel | null = null
   private ws: WebSocket | null = null
   // private isOfferer: boolean = false

   constructor(apiBase: string, wsBase: string) {
      this.peerConn = new RTCPeerConnection({
         iceServers: [
            { urls: "stun:stun.l.google.com:19302" } // free public STUN
         ]
      })
      this.apiBase = apiBase
      this.wsBase = wsBase
   }

   public async createRoom() {
      const res = await fetch(`${this.apiBase}/api/rooms/create`, { method: "POST" })
      const data: { roomId: string, clientId: string } = await res.json()
      this.connectWebSocket(data.roomId, data.clientId)
      return data
   }

   public async joinRoom(roomId: string) {
      const res = await fetch(`${this.apiBase}/api/rooms/join?roomId=${roomId}`, { method: "POST" })
      const data: { roomId: string, clientId: string } = await res.json()
      this.connectWebSocket(data.roomId, data.clientId)
      return data
   }

   private connectWebSocket(roomId: string, clientId: string) {
      this.ws = new WebSocket(`${this.wsBase}/ws?roomId=${roomId}&clientId=${clientId}`)

      this.ws.onopen = () => {
         this.log("websocket connected ✅");
         this.initWebRtc()
      }

      this.ws.onmessage = (e) => {
         const msg = JSON.parse(e.data)
         this.handleSignalingMessage(msg);
      }
   }

   private sendSignalToWS(message:SignalingMessage) {
      this.ws?.send(JSON.stringify(message))
   }

   private initWebRtc() {
      this.peerConn.onicecandidate = (e) => {
         if (e.candidate) {
            this.sendSignalToWS({ type: "candidate", data: e.candidate })
         }
      }

      this.peerConn.ondatachannel = (e) => {
         this.dataChannel = e.channel
         this.setupDataChannel()
      }

   }

   private handleSignalingMessage(msg: SignalingMessage) {
      switch (msg.type) {
         case "role":
            {
               if (msg.data?.role === "offerer") {
                  this.dataChannel = this.peerConn.createDataChannel("chat");
                  this.setupDataChannel()
                  this.peerConn.createOffer().then((offer) => {
                     this.peerConn.setLocalDescription(offer)
                     this.sendSignalToWS({ type: "offer", data: offer })
                  })
               }
            }
            break;

         case "offer":
            {
               this.peerConn.setRemoteDescription(new RTCSessionDescription(msg.data))

               this.peerConn.createAnswer().then((answer) => {
                  this.peerConn.setLocalDescription(answer)
                  this.sendSignalToWS({ type: "answer", data: answer })
               })
            }
            break;

         case "answer":
            {
               this.peerConn.setRemoteDescription(msg.data)
            }
            break;

         case "candidate":
            {
               this.peerConn.addIceCandidate(new RTCIceCandidate(msg.data))
            }
            break;

         case "timeout":
            {
               this.log("❌", msg.message);
            }
            break;

         default:
            break;
      }
   }

   private setupDataChannel() {
      this.dataChannel!.onopen = () => this.log("📡 Data channel open")
      this.dataChannel!.onmessage = (e) => this.log("📩 Received:", e.data)
   }

   public sendWebRTCmessage(msg: string) {
      if (this.dataChannel?.readyState === "open") {
         this.dataChannel.send(msg)
      }
   }

   public log(...args: any[]) {
      const logEl = document.getElementById("log") as HTMLDivElement;
      if (!logEl) return;
      // if (logEl) logEl.textContent += args.join(" ") + "\n";

      const newLogEntry = document.createElement("div");
      newLogEntry.className = "log-entry"; // Optional: for styling
      newLogEntry.textContent = args.join(" ");

      logEl.appendChild(newLogEntry);

      // Optional: Auto-scroll to the bottom
      logEl.scrollTop = logEl.scrollHeight;
   }
}

