# Broadcast Server
Simple broadcast server that allows clients to connect and send messages, which will be broadcast to all connected clients.

### Commands

````bash
#This command will start the server.
broadcast-server --operation start --port 8080

#This command will connect the client to the server.
broadcast-server --operation connect --port 8080 --username client01
````

### Inspiration idea
Project idea basis [roadmap.sh](https://roadmap.sh/projects/broadcast-server)