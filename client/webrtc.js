let peerConnection = new RTCPeerConnection();
let dataChannel = null;
let latestSDP = null;
const sharedKey = 42; // Very basic XOR key, replace with proper crypto for real-world use

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

peerConnection.onicecandidate = () => {
   if (peerConnection.localDescription) {
      latestSDP = JSON.stringify(peerConnection.localDescription, null, 2);
      log("✅ SDP ready. Use 'Copy SDP' or 'Show SDP'");
   }
};

peerConnection.ondatachannel = (event) => {
   dataChannel = event.channel;
   dataChannel.onopen = () => {
      log("✅ Connection opened (Peer B)");
      dataChannel.send(encrypt("Hello from Peer B!", sharedKey));
   };
   dataChannel.onmessage = (e) => {
      log("🔓 Message from Peer A:", decrypt(e.data, sharedKey));
   };
};

function setRemote() {
   try {
      const remote = JSON.parse(document.getElementById('sdp').value);
      peerConnection.setRemoteDescription(remote).then(() => {
         log("📥 Remote SDP set");

         if (remote.type === "offer") {
            peerConnection.createAnswer().then(answer => {
               peerConnection.setLocalDescription(answer);
            });
         }
      });
   } catch (e) {
      alert("❌ Invalid SDP JSON");
   }
}

function copySDP() {
   if (!latestSDP) return alert("⚠️ SDP not ready");
   if (!navigator.clipboard) {
      alert("📋 Clipboard API not available. Use 'Show SDP' instead.");
      return;
   }
   navigator.clipboard.writeText(latestSDP).then(() => {
      alert("✅ SDP copied to clipboard");
   }).catch(err => {
      alert("❌ Failed to copy: " + err);
   });
}

function showSDP() {
   if (!latestSDP) return alert("⚠️ SDP not ready");
   document.getElementById('sdpOutput').value = latestSDP;
}

// Offerer (Peer A) logic
if (confirm("Are you the Offerer (Peer A)?")) {
   dataChannel = peerConnection.createDataChannel("secureChannel");

   document.getElementById("role").textContent = "🔵 You are Peer A (Offerer)"

   dataChannel.onopen = () => {
      log("✅ Connection opened (Peer A)");
      dataChannel.send(encrypt("Hello from Peer A!", sharedKey));
   };

   dataChannel.onmessage = (e) => {
      // log("🔓 Message from Peer B:", decrypt(e.data, sharedKey));
      try {
         const decrypted = decrypt(e.data, sharedKey);
         log("Success! Decrypted:", decrypted);
      } catch (err) {
         log("Error! Failed to decrypt:", e.data, err);
      }
   };

   peerConnection.createOffer().then(offer => {
      peerConnection.setLocalDescription(offer);
   });
} else {
   document.getElementById("role").textContent = "🟡 You are Peer B (Answerer)";
}