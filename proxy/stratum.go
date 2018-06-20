package proxy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/jkkgbe/open-zcash-pool/util"
)

const (
	MaxReqSize = 10240
)

func (s *ProxyServer) ListenTCP() {
	timeout := util.MustParseDuration(s.config.Proxy.Stratum.Timeout)
	s.timeout = timeout

	addr, err := net.ResolveTCPAddr("tcp", s.config.Proxy.Stratum.Listen)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer server.Close()

	log.Printf("Stratum listening on %s", s.config.Proxy.Stratum.Listen)
	var accept = make(chan int, s.config.Proxy.Stratum.MaxConn)
	n := 0

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}
		conn.SetKeepAlive(true)

		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())

		if s.policy.IsBanned(ip) || !s.policy.ApplyLimitPolicy(ip) {
			conn.Close()
			continue
		}
		n += 1
		cs := &Session{conn: conn, ip: ip}

		accept <- n
		go func(cs *Session) {
			err = s.handleTCPClient(cs)
			if err != nil {
				s.removeSession(cs)
				conn.Close()
			}
			<-accept
		}(cs)
	}
}

func (s *ProxyServer) handleTCPClient(cs *Session) error {
	cs.enc = json.NewEncoder(cs.conn)
	connbuff := bufio.NewReaderSize(cs.conn, MaxReqSize)
	s.setDeadline(cs.conn)

	for {
		data, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			fmt.Println(string(data))
			log.Printf("Socket flood detected from %s", cs.ip)
			s.policy.BanClient(cs.ip)
			return err
		} else if err == io.EOF {
			log.Printf("Client %s disconnected", cs.ip)
			s.removeSession(cs)
			break
		} else if err != nil {
			log.Printf("Error reading from socket: %v", err)
			return err
		}

		if len(data) > 1 {
			var req StratumReq
			err = json.Unmarshal(data, &req)
			if err != nil {
				s.policy.ApplyMalformedPolicy(cs.ip)
				log.Printf("Malformed stratum request from %s: %v", cs.ip, err)
				return err
			}
			s.setDeadline(cs.conn)
			err = cs.handleTCPMessage(s, &req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cs *Session) handleTCPMessage(s *ProxyServer, req *StratumReq) error {
	var params []string
	err := json.Unmarshal(req.Params, &params)
	if err != nil {
		log.Println("Malformed stratum request params from", cs.ip)
		return err
	}

	var reply interface{}
	var errReply *ErrorReply
	// Handle RPC methods
	switch req.Method {
	case "mining.subscribe":
		extraNonce1 := s.extraNonceCounter.getNextExtraNonce1()
		reply = s.handleSubscribeRPC(cs, extraNonce1)
		fmt.Println("mining.subscribe", params, reply, errReply)
	case "mining.authorize":
		reply, errReply = s.handleAuthorizeRPC(cs, params)
		fmt.Println("mining.authorize", params, reply, errReply)
		if errReply != nil {
			return cs.sendTCPError(req.Id, errReply)
		}
		cs.sendTCPResult(req.Id, reply)
		// TODO set_target
		var d = []interface{}{s.diff}
		cs.setTarget(&d)
		t := s.currentWork()
		if t == nil || s.isSick() {
			return nil
		}
		reply := t.CreateJob()
		return cs.pushNewJob(&reply)
	case "mining.submit":
		reply, errReply = s.handleTCPSubmitRPC(cs, params, req.Worker)
		fmt.Println("mining.submit", params, reply, errReply)
	case "mining.extranonce.subscribe":
		errReply = &ErrorReply{Code: 20, Message: "Not supported."}
		fmt.Println("mining.extranonce.subscribe", params, reply, errReply)
	default:
		errReply = s.handleUnknownRPC(cs, req.Method)
	}

	if errReply != nil {
		return cs.sendTCPError(req.Id, errReply)
	}
	return cs.sendTCPResult(req.Id, reply)
}

func (cs *Session) sendTCPResult(id json.RawMessage, result interface{}) error {
	cs.Lock()
	defer cs.Unlock()

	message := JSONRpcResp{Id: id, Version: "2.0", Error: nil, Result: result}
	return cs.enc.Encode(&message)
}

func (cs *Session) setTarget(params *[]interface{}) error {
	cs.Lock()
	defer cs.Unlock()
	message := JSONPushMessage{Version: "2.0", Method: "mining.set_target", Params: *params, Id: 0}
	fmt.Println("setTarget", &message)
	return cs.enc.Encode(&message)
}

func (cs *Session) pushNewJob(params *[]interface{}) error {
	cs.Lock()
	defer cs.Unlock()
	message := JSONPushMessage{Version: "2.0", Method: "mining.notify", Params: *params, Id: 0}
	fmt.Println("pushNewJob", message)
	return cs.enc.Encode(&message)
}

func (cs *Session) sendTCPError(id json.RawMessage, reply *ErrorReply) error {
	cs.Lock()
	defer cs.Unlock()

	message := JSONRpcResp{Id: id, Version: "2.0", Error: reply}
	return cs.enc.Encode(&message)
	// err := cs.enc.Encode(&message)
	// if err != nil {
	// return err
	// }
	// return errors.New(reply.Message)
}

func (self *ProxyServer) setDeadline(conn *net.TCPConn) {
	conn.SetDeadline(time.Now().Add(self.timeout))
}

func (s *ProxyServer) registerSession(cs *Session) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	s.sessions[cs] = struct{}{}
}

func (s *ProxyServer) removeSession(cs *Session) {
	s.sessionsMu.Lock()
	defer s.sessionsMu.Unlock()
	delete(s.sessions, cs)
}

func (s *ProxyServer) broadcastNewJobs() {
	t := s.currentWork()
	// if t == nil || len(t.Header) == 0 || s.isSick() {
	if t == nil || s.isSick() {
		return
	}
	reply := t.CreateJob()

	s.sessionsMu.RLock()
	defer s.sessionsMu.RUnlock()

	count := len(s.sessions)
	log.Printf("Broadcasting new job to %v stratum miners", count)

	start := time.Now()
	bcast := make(chan int, 1024)
	n := 0
	for m, _ := range s.sessions {
		n++
		bcast <- n

		go func(cs *Session) {
			err := cs.pushNewJob(&reply)
			<-bcast
			if err != nil {
				log.Printf("Job transmit error to %v@%v: %v", cs.login, cs.ip, err)
				s.removeSession(cs)
			} else {
				s.setDeadline(cs.conn)
			}
		}(m)
	}
	log.Printf("Jobs broadcast finished %s", time.Since(start))
}
