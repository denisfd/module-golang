package main

type Block struct {
	PrevHash string
	CurHash  string
	Data     []string
}

type BlockChain struct {
	Blocks []Block
}

func (b *Block) Reset() {
	b.Data = b.Data[:0]
}

func (b *Block) Append(data string) {
	b.Data = append(b.Data, data)
}

func (chain *BlockChain) Drop() {
	chain.Blocks = chain.Blocks[:0]
}

func (chain *BlockChain) Add(block Block) bool {
	if block.PrevHash == chain.Blocks[len(chain.Blocks)-1].CurHash {
		chain.Blocks = append(chain.Blocks, block)
		return true
	}
	return false
}
