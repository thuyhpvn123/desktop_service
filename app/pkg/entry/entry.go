package entry

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/pack"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IEntry interface {
	GetBlockCount() *uint256.Int
	GetHash() common.Hash
	GetProto() protoreflect.ProtoMessage
	GetProtoWithBlockCount() protoreflect.ProtoMessage
	Marshal() []byte
}

type Entry struct {
	blockCount *uint256.Int
	numHashes  uint64 // number of hash since previous entry
	hash       common.Hash
	packs      []pack.IPack
	proto      *pb.PohEntry
}

func NewEntry(
	blockCount *uint256.Int,
	lastHash common.Hash,
	numHashes uint64,
	packs []pack.IPack) IEntry {
	var hash common.Hash
	for i := uint64(0); i < numHashes-1; i++ {
		hash = createHash(lastHash, nil)
		lastHash = hash
	}
	hash = createHash(lastHash, packs)
	return &Entry{
		blockCount: blockCount,
		numHashes:  numHashes,
		hash:       hash,
		packs:      packs,
	}
}

func UnmarshalWithBlockCount(b []byte) (IEntry, error) {
	pb := &pb.EntryWithBlockCount{}
	err := proto.Unmarshal(b, pb)
	if err != nil {
		return nil, err
	}
	blockCount := uint256.NewInt(0).SetBytes(pb.BlockCount)
	packs := []pack.IPack{}
	for _, v := range pb.Entry.Packs {
		packs = append(packs, pack.NewPack(v))
	}
	return &Entry{
		blockCount: blockCount,
		hash:       common.BytesToHash(pb.Entry.Hash),
		numHashes:  pb.Entry.NumHashes,
		packs:      packs,
	}, nil
}

func (e *Entry) GetBlockCount() *uint256.Int {
	return e.blockCount
}

func (e *Entry) GetHash() common.Hash {
	return e.hash
}

func createHash(lastHash common.Hash, packs []pack.IPack) common.Hash {
	packHashes := [][]byte{}

	for _, v := range packs {
		packHashes = append(packHashes, v.GetHash().Bytes())
	}

	hashData := &pb.PohHashData{
		PreHash:    lastHash.Bytes(),
		PackHashes: packHashes,
	}

	b, _ := proto.Marshal(hashData)
	return crypto.Keccak256Hash(b)
}

func (e *Entry) GetProto() protoreflect.ProtoMessage {
	e.proto = &pb.PohEntry{
		Hash:      e.hash[:],
		NumHashes: e.numHashes,
	}
	for _, v := range e.packs {
		e.proto.Packs = append(e.proto.Packs, v.GetProto().Interface().(*pb.Pack))
	}

	return e.proto
}

func (e *Entry) GetProtoWithBlockCount() protoreflect.ProtoMessage {
	pbWithBlockCount := &pb.EntryWithBlockCount{
		Entry:      e.GetProto().(*pb.PohEntry),
		BlockCount: e.blockCount.Bytes(),
	}
	return pbWithBlockCount
}

func (e *Entry) Marshal() []byte {
	b, err := proto.Marshal(e.proto)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return b
}
