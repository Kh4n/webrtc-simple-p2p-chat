Somewhat based off of https://github.com/mdn/samples-server/tree/master/s/webrtc-from-chat

Differences:
- I used Typescript because I prefer it (gives me intellisense, etc)
- Comes with a signaling server
- Simplified heavily. Only for text/data communication. 

I made this mainly because simple examples of P2P via WebRTC using datachannels only was limited.

To use:
- First build with `go build` and then start up the server
- Then run `tsc -p tsconfig.json` to build the Typescript
- Serve the directory locally and navigate to / (index.html)
- Click open and then connect. Then use send to send messages to peer