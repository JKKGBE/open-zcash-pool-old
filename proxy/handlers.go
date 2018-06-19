package proxy

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/jkkgbe/open-zcash-pool/util"
)

// Allow only lowercase hexadecimal with 0x prefix
var noncePattern = regexp.MustCompile("^0x[0-9a-f]{16}$")
var hashPattern = regexp.MustCompile("^0x[0-9a-f]{64}$")
var workerPattern = regexp.MustCompile("^[0-9a-zA-Z-_]{1,8}$")

func (s *ProxyServer) handleSubscribeRPC(cs *Session, extraNonce1 string) []byte {
	cs.extraNonce1 = extraNonce1
	return json.RawMessage(`[null, "` + extraNonce1 + `"]`)
}

func (s *ProxyServer) handleAuthorizeRPC(cs *Session, params []string) (bool, *ErrorReply) {
	if len(params) == 0 {
		return false, &ErrorReply{Code: -1, Message: "Invalid params"}
	}

	login := strings.ToLower(params[0])
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
	if !noncePattern.MatchString(cs.extraNonce1) {
		return false, &ErrorReply{Code: 25, Message: "Not subscribed"}
	}
	return s.handleSubmitRPC(cs, params, id)
}

func (s *ProxyServer) handleSubmitRPC(cs *Session, params []string, id string) (bool, *ErrorReply) {
	if !workerPattern.MatchString(id) {
		id = "0"
	}
	return false, &ErrorReply{Code: -1, Message: "Invalid params"} //temp
	// if len(params) != 5 {
	// 	s.policy.ApplyMalformedPolicy(cs.ip)
	// 	log.Printf("Malformed params from %s@%s %v", cs.login, cs.ip, params)
	// 	return false, &ErrorReply{Code: -1, Message: "Invalid params"}
	// }

	// if !hashPattern.MatchString(params[1]) || !hashPattern.MatchString(params[2]) {
	// 	s.policy.ApplyMalformedPolicy(cs.ip)
	// 	log.Printf("Malformed PoW result from %s@%s %v", cs.login, cs.ip, params)
	// 	return false, &ErrorReply{Code: -1, Message: "Malformed PoW result"}
	// }

	// t := s.currentWork()
	// shareExists, validShare, errorReply := s.processShare(cs, id, t, params)
	// ok := s.policy.ApplySharePolicy(cs.ip, !shareExists && validShare)

	// if !validShare {
	// 	log.Printf("Invalid share from %s@%s", cs.login, cs.ip)
	// 	// Bad shares limit reached, return error and close
	// 	if !ok {
	// 		return false, errorReply
	// 	}
	// 	return false, nil
	// }
	// log.Printf("Valid share from %s@%s", cs.login, cs.ip)

	// if shareExists {
	// 	log.Printf("Duplicate share from %s@%s %v", cs.login, cs.ip, params)
	// 	return false, &ErrorReply{Code: 22, Message: "Duplicate share"}
	// }

	// if !ok {
	// 	return true, &ErrorReply{Code: -1, Message: "High rate of invalid shares"}
	// }
	// return true, nil
}

func (s *ProxyServer) handleUnknownRPC(cs *Session, m string) *ErrorReply {
	log.Printf("Unknown request method %s from %s", m, cs.ip)
	s.policy.ApplyMalformedPolicy(cs.ip)
	return &ErrorReply{Code: -3, Message: "Method not found"}
}
