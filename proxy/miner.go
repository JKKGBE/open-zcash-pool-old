package proxy

import (
	"fmt"
	"log"

	"github.com/jkkgbe/open-zcash-pool/equihash"
	"github.com/jkkgbe/open-zcash-pool/util"
)

func (s *ProxyServer) processShare(cs *Session, id string, params []string) (bool, *ErrorReply) {
	// workerId := params[0]
	// jobId := params[1]
	// nTime := params[2]
	extraNonce2 := params[3]
	solution := params[4]

	work := s.currentWork()
	header := work.BuildHeader(cs.extraNonce1, extraNonce2)
	ok, err := equihash.Verify(200, 9, header, util.HexToBytes(solution)[3:])
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		header = append(header, util.HexToBytes(solution)...)
		fmt.Println(util.BytesToHex(header))
		ok, err := s.rpc().SubmitBlock(util.BytesToHex(header))
		if err != nil {
			fmt.Println(err)
			log.Printf("Block submission failure")
			return false, &ErrorReply{Code: 23, Message: "Invalid share"}
		} else if !ok {
			log.Printf("Block rejected")
			return false, &ErrorReply{Code: 23, Message: "Invalid share"}
		} else {
			s.fetchWork()
			// exist, err := s.backend.WriteBlock(login, id, params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
			// if exist {
			// 	return true, false
			// }
			// if err != nil {
			// 	log.Println("Failed to insert block candidate into backend:", err)
			// } else {
			// 	log.Printf("Inserted block %v to backend", h.height)
			// }
			log.Printf("Block found by miner %v@%v at height", cs.login, cs.ip)
			return true, nil
		}
	} else {
		// exist, err := s.backend.WriteShare(login, id, params, shareDiff, h.height, s.hashrateExpiration)
		// if exist {
		// 	return true, false
		// }
		// if err != nil {
		// 	log.Println("Failed to insert share data into backend:", err)
		// }
		return false, &ErrorReply{Code: 23, Message: "Invalid share"}
	}
	// shareExists, validShare, errorReply := s.processShare(cs, id, t, params)
	// ok := s.policy.ApplySharePolicy(cs.ip, !shareExists && validShare)

	// if !validShare {
	// 	log.Printf("Invalid share from %s@%s", cs.login, cs.ip)
	// 	// Bad shares limit reached, return error and close
	// 	if !ok {
	// 		return false, false, errorReply
	// 	}
	// 	return false, false, nil
	// }
	// log.Printf("Valid share from %s@%s", cs.login, cs.ip)

	// if shareExists {
	// 	log.Printf("Duplicate share from %s@%s %v", cs.login, cs.ip, params)
	// 	return false, false, &ErrorReply{Code: 22, Message: "Duplicate share"}
	// }

	// if !ok {
	// 	return false, true, &ErrorReply{Code: -1, Message: "High rate of invalid shares"}
	// }
	// return false, true, nil
}

// func HashHeader(w *Work, header []byte) (ShareStatus, string) {
// 	round1 := sha256.Sum256(header)
// 	round2 := sha256.Sum256(round1[:])

// 	round2 = util.ReverseBuffer(round2[:])

// 	// Check against the global target
// 	if TargetCompare(round2, w.Template.Target) <= 0 {
// 		return ShareBlock, hex.EncodeToString(round2[:])
// 	}

// 	if TargetCompare(round2, shareTarget) > 1 {
// 		return ShareInvalid, ""
// 	}

// 	return ShareOK, ""
// }

// func TargetCompare(a []byte, b []byte) {
// 	x := big.NewInt(0)
// 	x.SetBytes(a[:])
// 	y := big.NewInt(0)
// 	y.SetBytes(b[:])
// 	return a.Cmp(b)
// }
