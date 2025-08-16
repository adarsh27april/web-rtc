export class WebRtcConnection {
   private apiBase: string
   private wsBase: string
   private peerConn: RTCPeerConnection;
   private dataChannel: RTCDataChannel | null = null
   private ws: WebSocket | null = null
   private isOfferer: boolean = false

   constructor(apiBase: string, wsBase: string) {
      this.peerConn = new RTCPeerConnection()
      this.apiBase = apiBase
      this.wsBase = wsBase
   }

   public async createRoom() {
      const res = await fetch(`${this.apiBase}/api/rooms/create`, { method: "POST" })
      const data: { roomId: string, clientId: string } = await res.json()
      this.isOfferer = true
      this.connectWebSocket(data.roomId, data.clientId)
      return data
   }

   public async joinRoom(roomId: string) {
      const res = await fetch(`${this.apiBase}/api/rooms/join?roomId=${roomId}`, { method: "POST" })
      const data: { roomId: string, clientId: string } = await res.json()
      this.isOfferer = false
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

      if (this.isOfferer) {
         this.dataChannel = this.peerConn.createDataChannel("chat")
         this.setupDataChannel()

         this.peerConn.createOffer().then((offer) => {
            this.peerConn.setLocalDescription(offer)
            this.sendSignalToWS("offer", offer)
         })
      }
   }

   private handleSignalingMessage(msg: any) {
      if (msg.type === "offer") {
         this.peerConn.setRemoteDescription(new RTCSessionDescription(msg.data))

         this.peerConn.createAnswer().then((answer) => {
            this.peerConn.setLocalDescription(answer)
            this.sendSignalToWS("answer", answer)
         })
      } else if (msg.type === "answer") {
         this.peerConn.setRemoteDescription(msg.data)
      } else if (msg.type === "candidate") {
         this.peerConn.addIceCandidate(new RTCIceCandidate(msg.data))
      }
   }

   private setupDataChannel() {
      this.dataChannel!.onopen = () => this.log("📡 Data channel open")
      this.dataChannel!.onmessage = (e) => this.log("📩 Received:", e.data)
   }

   public sendMessage(msg: string) {
      if (this.dataChannel?.readyState === "open") {
         this.dataChannel.send(msg)
      }
   }

   public log(...args: any[]) {
      const logEl = document.getElementById("log") as HTMLDivElement;
      if (logEl) logEl.textContent += args.join(" ") + "\n";
   }
}

