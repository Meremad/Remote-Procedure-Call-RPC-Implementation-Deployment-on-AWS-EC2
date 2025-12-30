package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
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

func add(params map[string]interface{}) int {
	a := int(params["a"].(float64))
	b := int(params["b"].(float64))

	time.Sleep(5 * time.Second)

	return a + b
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func reverseString(params map[string]interface{}) string {
	s := params["s"].(string)
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var req RPCRequest
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	fmt.Println("[RPC REQUEST]", req)

	var res RPCResponse
	res.RequestID = req.RequestID

	switch req.Method {
	case "add":
		res.Status = "OK"
		res.Result = add(req.Params)
	case "get_time":
		res.Status = "OK"
		res.Result = getTime()
	case "reverse_string":
		res.Status = "OK"
		res.Result = reverseString(req.Params)
	default:
		res.Status = "ERROR"
		res.Error = "Unknown method"
	}

	encoder.Encode(res)
}

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("[*] RPC Server listening on port 5000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}
