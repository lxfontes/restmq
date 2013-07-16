package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type VerbHandler struct {
	Context *MQContext
}

type VerbGet struct {
	Count uint64 `json:"count"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VerbDel struct {
	Queue  string `json:"queue"`
	Status int    `json:"status"`
}

func NewVerbHandler(context *MQContext) *VerbHandler {
	s := new(VerbHandler)
	s.Context = context
	return s
}

func (v VerbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	qname := vars["queue"]
	switch r.Method {
	case "GET":
		v.GET(w, r, qname)
	case "POST":
		v.POST(w, r, qname)
	case "DELETE":
		v.DELETE(w, r, qname)
	}
}

func (v VerbHandler) GET(w http.ResponseWriter, r *http.Request, qname string) {
	reply, err := v.Context.Get(qname, true)

	if err != nil {
		return
	}

	resp := VerbGet{
		Count: reply.Count,
		Key:   reply.Key,
		Value: reply.Value,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (v VerbHandler) DELETE(w http.ResponseWriter, r *http.Request, qname string) {
	status := 0
	err := v.Context.Flush(qname)

	if err != nil {
		status = 1
	}

	resp := VerbDel{
		Queue:  qname,
		Status: status,
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func (v VerbHandler) POST(w http.ResponseWriter, r *http.Request, qname string) {
	val := r.FormValue("value")
	reply, _ := v.Context.Push(qname, val)
	w.Write([]byte(reply.Value))
}
