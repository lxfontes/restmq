package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	redis "menteslibres.net/gosexy/redis"
	"net/http"
	"os"
	"strconv"
)

type MQContext struct {
	Redis *redis.Client
}

type MQReply struct {
	Value string
	Key   string
	Count uint64
}

func (mq MQContext) Flush(queue string) error {
	kname := fmt.Sprintf("%s:queue", queue)
	_, err := mq.Redis.Del(kname)
	return err
}

func (mq MQContext) Get(queue string, pop bool) (*MQReply, error) {
	var kname string
	var err error
	var ename string

	kname = fmt.Sprintf("%s:queue", queue)
	if pop {
		ename, err = mq.Redis.RPop(kname)
	} else {
		ename, err = mq.Redis.LIndex(kname, -1)
	}

	if err != nil {
		return nil, err
	}

	val, err := mq.Redis.Get(ename)

	if err != nil {
		return nil, err
	}

	refcount := uint64(0)

	if !pop {
		kname = fmt.Sprintf("%s:refcount", ename)
		refs, err := mq.Redis.Get(kname)
		if err == nil {
			refcount, _ = strconv.ParseUint(refs, 10, 64)
		}
	}

	reply := new(MQReply)
	reply.Value = val
	reply.Key = ename
	reply.Count = refcount

	return reply, nil
}

func (mq MQContext) Del(queue string, key string) (bool, error) {
	kname := fmt.Sprintf("%s:queue", queue)
	_, err := mq.Redis.LRem(kname, -1, key)
	if err != nil {
		return false, err
	}
	count, err := mq.Redis.Del(key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (mq MQContext) Push(queue string, value string) (*MQReply, error) {
	var kname string
	var err error

	kname = fmt.Sprintf("%s:UUID", queue)
	qid, err := mq.Redis.Incr(kname)
	if err != nil {
		return nil, err
	}
	ename := fmt.Sprintf("%s:%d", queue, qid)
	_, err = mq.Redis.Set(ename, value)
	if err != nil {
		return nil, err
	}

	kname = fmt.Sprintf("%s:queue", queue)
	_, err = mq.Redis.LPush(kname, ename)

	if err != nil {
		return nil, err
	}

	reply := new(MQReply)
	reply.Value = ename
	reply.Count = 0

	return reply, nil
}

var rHost = flag.String("redis_host", "localhost", "Redis IP Address")
var rPort = flag.Uint("redis_port", 6379, "Redis Port")

func main() {

	flag.Parse()

	mqContext := new(MQContext)
	mqContext.Redis = redis.New()

	err := mqContext.Redis.Connect(*rHost, *rPort)
	if err != nil {
		log.Print("Error finding redis server")
		os.Exit(1)
		return
	}

	r := mux.NewRouter()
	r.Handle("/q/{queue}", NewVerbHandler(mqContext))
	r.Handle("/queue", NewCmdHandler(mqContext))
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
