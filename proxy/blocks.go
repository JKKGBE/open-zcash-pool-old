package proxy

import (
	"log"
	"math/big"
	"sync"

	"github.com/jkkgbe/open-zcash-pool/merkleTree"
	"github.com/jkkgbe/open-zcash-pool/util"
)

const maxBacklog = 3

type heightDiffPair struct {
	diff   *big.Int
	height uint64
}

type transaction struct {
	data string `json:"data"`
	hash string `json:"hash"`
	fee  int    `json:"fee"`
}

type coinbaseTxn struct {
	data           string `json:"data"`
	hash           string `json:"hash"`
	foundersReward int    `json:"foundersreward"`
}

type BlockTemplate struct {
	sync.RWMutex
	version       uint32        `json:"version"`
	prevBlockHash string        `json:"previousblockhash"`
	transactions  []transaction `json:"transactions"`
	coinbaseTxn   coinbaseTxn   `json:"coinbasetxn"`
	longpollId    string        `json:"longpollid"`
	target        string        `json:"target"`
	minTime       int           `json:"mintime"`
	nonceRange    string        `json:"noncerange"`
	sigOpLimit    int           `json:"sigoplimit"`
	sizeLimit     int           `json:"sizelimit"`
	curTime       uint32        `json:"curtime"`
	bits          string        `json:"bits"`
	height        int           `json:"height"`
}

type Work struct {
	jobId              string
	version            string
	prevHashReversed   string
	merkleRootReversed string
	reservedField      string
	time               string
	bits               string
	cleanJobs          bool
	nonce              string
	solutionSize       [3]byte
	solution           [1344]byte
	header             [4 + 32 + 32 + 32 + 4 + 4 + 32 + 3 + 1344]byte
}

// func (b Work) Difficulty() *big.Int     { return b.difficulty }
// func (b Work) HashNoNonce() common.Hash { return b.hashNoNonce }
// func (b Work) Nonce() uint64            { return b.nonce }
// func (b Work) MixDigest() common.Hash   { return b.mixDigest }
// func (b Work) NumberU64() uint64        { return b.number }

func (s *ProxyServer) fetchBlockTemplate() {
	rpc := s.rpc()
	t := s.currentBlockTemplate()
	var reply BlockTemplate
	err := rpc.GetBlockTemplate(&reply)
	if err != nil {
		log.Printf("Error while refreshing block template on %s: %s", rpc.Name, err)
		return
	}
	// No need to update, we have fresh job
	if t != nil && t.prevBlockHash == reply.prevBlockHash {
		return
	}

	// TODO calc merkle root etc
	// if (blockHeight.toString(16).length % 2 === 0) {
	//     var blockHeightSerial = blockHeight.toString(16);
	// } else {
	//     var blockHeightSerial = '0' + blockHeight.toString(16);
	// }
	// var height = Math.ceil((blockHeight << 1).toString(2).length / 8);
	// var lengthDiff = blockHeightSerial.length/2 - height;
	// for (var i = 0; i < lengthDiff; i++) {
	//     blockHeightSerial = blockHeightSerial + '00';
	// }
	// length = '0' + height;
	// var serializedBlockHeight = new Buffer.concat([
	//     new Buffer(length, 'hex'),
	//     util.reverseBuffer(new Buffer(blockHeightSerial, 'hex')),
	//     new Buffer('00', 'hex') // OP_0
	// ]);

	// tx.addInput(new Buffer('0000000000000000000000000000000000000000000000000000000000000000', 'hex'),
	//     4294967295,
	//     4294967295,
	//     new Buffer.concat([serializedBlockHeight,
	//         Buffer('5a2d4e4f4d50212068747470733a2f2f6769746875622e636f6d2f6a6f7368756179616275742f7a2d6e6f6d70', 'hex')]) //Z-NOMP! https://github.com/joshuayabut/z-nomp
	// );

	// generatedTxHash := CreateRawTransaction(inputs, outputs).TxHash()
	txHashes := make([][32]byte, len(reply.transactions)+1)
	// txHashes[0] = util.ReverseBuffer(generatedTxHash)
	copy(txHashes[0][:], util.HexToBytes(reply.coinbaseTxn.hash)[:32])
	for i, transaction := range reply.transactions {
		copy(txHashes[i+1][:], util.HexToBytes(transaction.hash)[:32])
	}

	mtBottomRow := txHashes
	mt := merkleTree.NewMerkleTree(mtBottomRow)
	mtr := mt.MerkleRoot()

	newWork := Work{
		version:            util.BytesToHex(util.PackUInt32LE(reply.version)),
		prevHashReversed:   util.BytesToHex(util.ReverseBuffer(util.HexToBytes(reply.prevBlockHash))),
		merkleRootReversed: util.BytesToHex(util.ReverseBuffer(mtr[:])),
		reservedField:      "0000000000000000000000000000000000000000000000000000000000000000",
		time:               util.BytesToHex(util.PackUInt32LE(reply.curTime)),
		bits:               util.BytesToHex(util.ReverseBuffer(util.HexToBytes(reply.bits))),
		cleanJobs:          true,
	}

	// // Copy job backlog and add current one
	// newBlock.headers[reply[0]] = heightDiffPair{
	// 	diff:   util.TargetHexToDiff(reply[2]),
	// 	height: height,
	// }
	// if t != nil {
	// 	for k, v := range t.headers {
	// 		if v.height > height-maxBacklog {
	// 			newBlock.headers[k] = v
	// 		}
	// 	}
	// }
	s.workTemplate.Store(&newWork)
	log.Printf("New block to mine on %s at height %d", rpc.Name, reply.height)

	// Stratum
	if s.config.Proxy.Stratum.Enabled {
		go s.broadcastNewJobs()
	}
}
