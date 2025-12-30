# LAB 1 — Simple RPC Implementation on AWS EC2 (Go)

## 📌 Overview
This project implements a **minimal Remote Procedure Call (RPC) system** using the Go programming language and raw TCP sockets.

The goal of this lab is to understand:
- How RPC works internally (without using gRPC or net/rpc)
- Client–server communication over the network
- Timeouts and retry logic
- RPC semantics such as **at-least-once execution**

The system is deployed on **two AWS EC2 instances**:
- One instance runs the RPC server
- Another instance runs the RPC client

---

## 🏗 Architecture

Client EC2 ─── TCP (port 5000) ───► Server EC2


### RPC Components
- **Client stub** – sends requests, handles retries and timeouts
- **Server stub** – receives requests, executes methods, returns results
- **Network transport** – TCP sockets
- **Serialization** – JSON

---

## ⚙️ Implemented RPC Method

### `add(a, b)`
Adds two integers and returns the result.

Example request:
```json
{
  "request_id": "uuid",
  "method": "add",
  "params": { "a": 10, "b": 20 },
  "timestamp": 1766997922
}

Example response:
{
  "request_id": "uuid",
  "result": 30,
  "status": "OK"
}

🚀 How to Run
1️⃣ Server Setup (EC2 Server Node)
sudo apt update
sudo apt install golang -y


Run the server:

go run server.go


Expected output:

[*] RPC Server listening on port 5000

2️⃣ Client Setup (EC2 Client Node)
sudo apt update
sudo apt install golang -y


Initialize Go module:

go mod init rpc-client
go get github.com/google/uuid


Run the client:

go run client.go


Expected output:

[Attempt 1] Sending RPC request
[OK] Result: 30
