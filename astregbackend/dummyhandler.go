package main

import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strings"
	"time"
)

type DummyHandler struct {
	rediscon *redis.Client
	digits   int
	section  string
}

func NewDummyHandler(rediscon *redis.Client, digits int, section string) *DummyHandler {
	return &DummyHandler{
		rediscon: rediscon,
		digits:   digits,
		section:  section,
	}
}

func parseConditions(q url.Values) (whereClause string, values []interface{}) {
	conditions, values := make([]string, 0, len(q)), make([]interface{}, 0, len(q))
	var column, op string
	for field := range q {
		parts := strings.SplitN(field, " ", 2)
		if len(parts) > 1 {
			column = parts[0]
			op = parts[1]
		} else {
			column = field
			op = "="
		}

		conditions = append(conditions, fmt.Sprintf("`%s` %s ?", column, op))
		values = append(values, q.Get(field))
	}

	whereClause = strings.Join(conditions, " AND ")
	return
}

func (h *DummyHandler) Single(w http.ResponseWriter, r *http.Request) {
	var peer, reply, key string
	var buffer bytes.Buffer

	whereClause, values := parseConditions(r.PostForm)

	for _, v := range values {
		peer = v.(string)
	}

	if len(peer) != Config.Digits {
		if Config.Verbose {
			log.Printf("len(%s) != %d -> ignore\n", peer, Config.Digits)
		}
		return
	}

	if Config.Verbose {
		log.Printf("GET %s WHERE %s (%s)\n", h.section, whereClause, values)
	}

	switch h.section {
	case "endpoint":
		//reply = fmt.Sprintf("id=%s&type=endpoint&auth=%s&aors=%s&transport=transport-tcp&context=nodes&disallow=all&allow=g722,opus,alaw&send_pai=yes", peer, peer, peer)
		reply = fmt.Sprintf("id=%s&type=endpoint&aors=%s&transport=transport-tcp&context=nodes&disallow=all&allow=g722,opus,alaw&send_pai=yes", peer, peer)
		key = fmt.Sprintf("endpoint_%s", peer)
	case "auth":
		reply = fmt.Sprintf("id=%s&type=auth&authtype=userpass&username=%s&password=%s", peer, peer, peer)
		key = fmt.Sprintf("auth_%s", peer)
	case "aors":
		reply = fmt.Sprintf("id=%s&type=aor&max_contacts=1&qualify_frequency=60&default_expiration=60", peer)
		key = fmt.Sprintf("aors_%s", peer)
	default:
		return
	}

	err := h.rediscon.Set(ctx, key, reply, time.Duration(Config.RedisExpiration)*time.Second).Err()

	if err != nil {
		panic(err)
	} else {
		if Config.Verbose {
			log.Printf("> %s", reply)
		}
	}

	buffer.Reset()
	buffer.Write([]byte(reply))
	buffer.Write([]byte("\n"))

	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (h *DummyHandler) Multi(w http.ResponseWriter, r *http.Request) {
	var peer string
	var val []string
	var buffer bytes.Buffer
	var err error

	whereClause, values := parseConditions(r.PostForm)

	for _, v := range values {
		peer = v.(string)
	}

	if Config.Verbose {
		//log.Printf("get : %s - where : %s - str: %s - values : %s\n", h.section, whereClause, str, values)
		log.Printf("GET %s WHERE %s (%s)\n", h.section, whereClause, values)
	}

	switch h.section {
	case "endpoint":
		val, err = h.rediscon.Keys(ctx, "endpoint_*").Result()
	case "auth":
		val, err = h.rediscon.Keys(ctx, "auth_*").Result()
	case "aors":
		val, err = h.rediscon.Keys(ctx, "aors_*").Result()
	default:
		return
	}

	if err != nil {
		log.Println("ERROR: ", err)
		return
	}

	buffer.Reset()

	for _, v := range val {
		peer, err = h.rediscon.Get(ctx, v).Result()
		log.Printf("> %s", peer)
		buffer.Write([]byte(peer))
		buffer.Write([]byte("\n"))
	}

	_, err = buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
		return
	}

	return
}
