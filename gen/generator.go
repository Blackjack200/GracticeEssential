package gen

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

type Void struct {
	Spawn cube.Pos
}

func New(spawn cube.Pos) *Void {
	return &Void{Spawn: spawn}
}

func (w *Void) GenerateChunk(pos world.ChunkPos, c *chunk.Chunk) {
	id := world.BlockRuntimeID(block.Stone{})
	s := w.Spawn
	if int32(s.X())>>4 == pos.X() && int32(s.Z())>>4 == pos.Z() {
		c.SetBlock(uint8(w.Spawn.X()), 1, uint8(w.Spawn.Z()), 0, id)
	}
}
