package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

type RPCRequest struct {
	RequestID string                 `json:"request_id"`
	Method    string                 `json:"method"`
	Params    map[string]interface{} `json:"params"`
	Timestamp int64                  `json:"timestamp"`
}

type RPCResponse struct {
	RequestID string      `json:"request_id"`
	Status    string      `json:"status"`
	Result    interface{} `json:"result,omitempty"`
	Error     string      `json:"error,omitempty"`
}

const (
	SERVER_ADDR = "SERVER_PUBLIC_IP:5000"
	TIMEOUT     = 2 * time.Second
	MAX_RETRIES = 3
)

func rpcCall(method string, params map[string]interface{}) {
	request := RPCRequest{
		RequestID: uuid.New().String(),
		Method:    method,
		Params:    params,
		Timestamp: time.Now().Unix(),
	}

	for attempt := 1; attempt <= MAX_RETRIES; attempt++ {
		fmt.Printf("[Attempt %d] Sending RPC request\n", attempt)

		conn, err := net.DialTimeout("tcp", SERVER_ADDR, TIMEOUT)
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}

		conn.SetDeadline(time.Now().Add(TIMEOUT))

		encoder := json.NewEncoder(conn)
		decoder := json.NewDecoder(conn)

		if err := encoder.Encode(request); err != nil {
			fmt.Println("Send error:", err)
			conn.Close()
			continue
		}

		var response RPCResponse
		if err := decoder.Decode(&response); err != nil {
			fmt.Println("[TIMEOUT] No response, retrying...")
			conn.Close()
			continue
		}

		if response.Status == "OK" {
			fmt.Println("[SUCCESS] Result:", response.Result)
			conn.Close()
			return
		}

		fmt.Println("[RPC ERROR]", response.Error)
		conn.Close()
	}

	fmt.Println("[FAILED] RPC call failed after retries")
}

func main() {
	rpcCall("add", map[string]interface{}{
		"a": 10,
		"b": 20,
	})
}
