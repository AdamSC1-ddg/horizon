package simplepath

import (
	"bytes"
	"fmt"
	"github.com/stellar/go-stellar-base/xdr"
	"github.com/stellar/horizon/db"
	"github.com/stellar/horizon/paths"
)

// pathNode implements the paths.Path interface and represents a path
// as a linked list pointing from destination to source.
type pathNode struct {
	Asset xdr.Asset
	Tail  *pathNode
	DB    db.SqlQuery
}

// check interface compatibility
var _ paths.Path = &pathNode{}

func (p *pathNode) String() string {
	if p == nil {
		return ""
	}

	var out bytes.Buffer
	fmt.Fprintf(&out, "%v", p.Asset)

	cur := p.Tail

	for cur != nil {
		fmt.Fprintf(&out, " -> %v", cur.Asset)
		cur = cur.Tail
	}

	return out.String()
}

func (p *pathNode) Destination() xdr.Asset {
	cur := p
	for cur.Tail != nil {
		cur = cur.Tail
	}
	return cur.Asset
}

func (p *pathNode) Source() xdr.Asset {
	// the destination for path is the head of the linked list
	return p.Asset
}

func (p *pathNode) Path() []xdr.Asset {
	path := p.Flatten()

	if len(path) < 2 {
		return nil
	}

	// return the flattened slice without the first and last elements
	// which are the source and the destination assets
	return path[1 : len(path)-1]
}

func (p *pathNode) Cost(amount xdr.Int64) (result xdr.Int64, err error) {
	result = amount

	if p.Tail == nil {
		return
	}

	cur := p

	for cur.Tail != nil {
		ob := cur.OrderBook()
		result, err = ob.Cost(cur.Asset, result)
		if err != nil {
			return
		}
		cur = cur.Tail
	}

	return
}

func (p *pathNode) Depth() int {
	depth := 0
	cur := p
	for {
		if cur == nil {
			return depth
		}
		cur = cur.Tail
		depth++
	}
}

func (p *pathNode) Flatten() (result []xdr.Asset) {
	cur := p

	for {
		if cur == nil {
			return
		}
		result = append(result, cur.Asset)
		cur = cur.Tail
	}

	return
}

func (p *pathNode) OrderBook() *orderBook {
	if p.Tail == nil {
		return nil
	}

	return &orderBook{
		Selling: p.Tail.Asset,
		Buying:  p.Asset,
		DB:      p.DB,
	}
}
