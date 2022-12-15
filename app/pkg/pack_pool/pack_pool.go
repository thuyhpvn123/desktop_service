package pack_pool

import (
	"sort"
	"sync"

	"gitlab.com/meta-node/meta-node/pkg/pack"
)

type PackPool struct {
	mu    sync.Mutex
	packs []pack.IPack
}

func NewPackPool() *PackPool {
	return &PackPool{
		packs: make([]pack.IPack, 0),
	}
}

func (p *PackPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.packs)
}

func (p *PackPool) AddPack(pack pack.IPack) {
	p.mu.Lock()
	p.packs = append(p.packs, pack)
	p.mu.Unlock()
}

func (p *PackPool) AddPacks(packs []pack.IPack) {
	p.mu.Lock()
	p.packs = append(p.packs, packs...)
	p.mu.Unlock()
}

func (p *PackPool) TakePack(numberOfPack uint64) []pack.IPack {
	p.mu.Lock()
	sort.Slice(p.packs, func(i int, u int) bool {
		return p.packs[i].GetTimestamp() < p.packs[u].GetTimestamp()
	})
	if int(numberOfPack) > len(p.packs) {
		numberOfPack = uint64(len(p.packs))
	}
	rs := p.packs[:numberOfPack]
	p.packs = p.packs[numberOfPack:]
	p.mu.Unlock()
	return rs
}

func (p *PackPool) Copy() *PackPool {
	rs := NewPackPool()
	rs.packs = append(rs.packs, p.packs...)
	return rs
}
