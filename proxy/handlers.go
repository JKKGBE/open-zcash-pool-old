package proxy

import (
	"log"
	"regexp"

	"github.com/jkkgbe/open-zcash-pool/util"
)

// Allow only lowercase hexadecimal with 0x prefix
var nTimePattern = regexp.MustCompile("^[0-9a-f]{8}$")
var noncePattern = regexp.MustCompile("^[0-9a-f]{64}$")

// var solutionPattern = regexp.MustCompile("^[0-9a-f]{2694}$")
var workerPattern = regexp.MustCompile("^[0-9a-zA-Z-_]{1,8}$")

func (s *ProxyServer) handleSubscribeRPC(cs *Session, extraNonce1 string) []string {
	cs.extraNonce1 = extraNonce1
	array := []string{"0", extraNonce1}
	return array
	// return json.RawMessage(`[null, "` + extraNonce1 + `"]`)
}

func (s *ProxyServer) handleAuthorizeRPC(cs *Session, params []string) (bool, *ErrorReply) {
	if len(params) == 0 {
		return false, &ErrorReply{Code: -1, Message: "Invalid params"}
	}

	login := params[0]
	if !util.IsValidtAddress(login) {
		return false, &ErrorReply{Code: -1, Message: "Invalid login"}
	}
	if !s.policy.ApplyLoginPolicy(login, cs.ip) {
		return false, &ErrorReply{Code: -1, Message: "You are blacklisted"}
	}
	cs.login = login
	s.registerSession(cs)
	log.Printf("Stratum miner connected %v@%v", login, cs.ip)
	return true, nil
}

func (s *ProxyServer) handleTCPSubmitRPC(cs *Session, params []string, id string) (bool, *ErrorReply) {
	s.sessionsMu.RLock()
	_, ok := s.sessions[cs]
	s.sessionsMu.RUnlock()

	if !ok {
		return false, &ErrorReply{Code: 24, Message: "Not authorized"}
	}
	if cs.extraNonce1 == "" {
		return false, &ErrorReply{Code: 25, Message: "Not subscribed"}
	}
	return s.handleSubmitRPC(cs, params, id)
}

func (s *ProxyServer) handleSubmitRPC(cs *Session, params []string, id string) (bool, *ErrorReply) {
	if !workerPattern.MatchString(id) {
		id = "0"
	}

	if len(params) != 5 {
		s.policy.ApplyMalformedPolicy(cs.ip)
		log.Printf("Malformed params from %s@%s %v", cs.login, cs.ip, params)
		return false, &ErrorReply{Code: -1, Message: "Invalid params"}
	}

	if !nTimePattern.MatchString(params[2]) {
		s.policy.ApplyMalformedPolicy(cs.ip)
		log.Printf("Malformed nTime result from %s@%s %v", cs.login, cs.ip, params)
		return false, &ErrorReply{Code: -1, Message: "Malformed nTime result"}
	}

	if !noncePattern.MatchString(cs.extraNonce1 + params[3]) {
		s.policy.ApplyMalformedPolicy(cs.ip)
		log.Printf("Malformed nonce result from %s@%s %v", cs.login, cs.ip, params)
		return false, &ErrorReply{Code: -1, Message: "Malformed nonce result"}
	}

	if len(params[4]) != 2694 {
		s.policy.ApplyMalformedPolicy(cs.ip)
		log.Printf("Malformed solution result from %s@%s %v", cs.login, cs.ip, params)
		return false, &ErrorReply{Code: -1, Message: "Malformed solution result, != 2694 length"}
	}

	return s.processShare(cs, id, params)
}

func (s *ProxyServer) handleUnknownRPC(cs *Session, m string) *ErrorReply {
	log.Printf("Unknown request method %s from %s", m, cs.ip)
	s.policy.ApplyMalformedPolicy(cs.ip)
	return &ErrorReply{Code: -3, Message: "Method not found"}
}
