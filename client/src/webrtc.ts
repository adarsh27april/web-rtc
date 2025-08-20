export class WebRtcConnection {
   private apiBase: string
   private wsBase: string
   private peerConn: RTCPeerConnection;
   private dataChannel: RTCDataChannel | null = null
   private ws: WebSocket | null = null
   // private isOfferer: boolean = false

   constructor(apiBase: string, wsBase: string) {
      this.peerConn = new RTCPeerConnection()
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

   private sendSignalToWS(type: string, data: any) {
      this.ws?.send(JSON.stringify({ type, data }))
   }

   private initWebRtc() {
      this.peerConn.onicecandidate = (e) => {
         if (e.candidate) {
            this.sendSignalToWS("candidate", e.candidate)
         }
      }

      this.peerConn.ondatachannel = (e) => {
         this.dataChannel = e.channel
         this.setupDataChannel()
      }

   }

   private handleSignalingMessage(msg: any) {
      switch (msg.type) {
         case "role":
            {
               if (msg.data?.role === "offerer") {
                  // this.isOfferer = true;

                  // since from signaling server we are told that this is offerer, hence creating the offer.
                  this.dataChannel = this.peerConn.createDataChannel("chat");
                  this.setupDataChannel()
                  this.peerConn.createOffer().then((offer) => {
                     this.peerConn.setLocalDescription(offer)
                     this.sendSignalToWS("offer", offer)
                  })
               } //else if (msg.data?.role === "answerer") {
               //    this.isOfferer = false;
               // }
            }
            break;

         case "offer":
            {
               this.peerConn.setRemoteDescription(new RTCSessionDescription(msg.data))

               this.peerConn.createAnswer().then((answer) => {
                  this.peerConn.setLocalDescription(answer)
                  this.sendSignalToWS("answer", answer)
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
               this.log("⏳ No peer joined within", msg.data?.afterSec, "seconds");
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
      if (logEl) logEl.textContent += args.join(" ") + "\n";
   }
}

