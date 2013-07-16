package main

import (
	"encoding/json"
	"net/http"
)

type CmdHandler struct {
	Context *MQContext
}

type CmdRequest struct {
	Command string `json:"cmd"`
	Queue   string `json:"queue"`
	Value   string `json:"value"`
	Key     string `json:"key"`
}

type CmdResponse struct {
	Queue   string `json:"queue"`
	Value   string `json:"value,omitempty"`
	Key     string `json:"key,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
	Count   uint64 `json:"count,omitempty"`
}

func NewCmdHandler(context *MQContext) *CmdHandler {
	s := new(CmdHandler)
	s.Context = context
	return s
}

func (v CmdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("body")
	req := new(CmdRequest)
	err := json.Unmarshal([]byte(data), &req)
	if err != nil {
		http.Error(w, "Boo boo", 500)
		return
	}

	switch req.Command {
	case "get":
		v.GET(w, r, req)
	case "add":
		v.ADD(w, r, req)
	case "del":
		v.DELETE(w, r, req)
	case "take":
		v.TAKE(w, r, req)
	}
}

func (v CmdHandler) GET(w http.ResponseWriter, r *http.Request, req *CmdRequest) {
	reply, err := v.Context.Get(req.Queue, false)

	if err != nil {
		return
	}

	resp := CmdResponse{
		Queue: req.Queue,
		Count: reply.Count,
		Key:   reply.Key,
		Value: reply.Value,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (v CmdHandler) TAKE(w http.ResponseWriter, r *http.Request, req *CmdRequest) {
	reply, err := v.Context.Get(req.Queue, true)

	if err != nil {
		return
	}

	resp := CmdResponse{
		Queue: req.Queue,
		Key:   reply.Key,
		Value: reply.Value,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (v CmdHandler) DELETE(w http.ResponseWriter, r *http.Request, req *CmdRequest) {
	deleted, _ := v.Context.Del(req.Queue, req.Key)

	resp := CmdResponse{
		Queue:   req.Queue,
		Key:     req.Key,
		Deleted: deleted,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (v CmdHandler) ADD(w http.ResponseWriter, r *http.Request, req *CmdRequest) {
	reply, _ := v.Context.Push(req.Queue, req.Value)

	resp := CmdResponse{
		Queue: req.Queue,
		Key:   reply.Value,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}
