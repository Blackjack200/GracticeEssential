package handler

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"net"
	"time"
)

type AttackEntityHandler interface {
	HandleAttackEntity(param1 *event.Context, param2 world.Entity, param3 *float64, param4 *float64, param5 *bool)
}
type BlockBreakHandler interface {
	HandleBlockBreak(param1 *event.Context, param2 cube.Pos, param3 *[]item.Stack)
}
type BlockPickHandler interface {
	HandleBlockPick(param1 *event.Context, param2 cube.Pos, param3 world.Block)
}
type BlockPlaceHandler interface {
	HandleBlockPlace(param1 *event.Context, param2 cube.Pos, param3 world.Block)
}
type ChangeWorldHandler interface {
	HandleChangeWorld(param1 *world.World, param2 *world.World)
}
type ChatHandler interface {
	HandleChat(param1 *event.Context, param2 *string)
}
type CommandExecutionHandler interface {
	HandleCommandExecution(param1 *event.Context, param2 cmd.Command, param3 []string)
}
type DeathHandler interface {
	HandleDeath(param1 damage.Source)
}
type ExperienceGainHandler interface {
	HandleExperienceGain(param1 *event.Context, param2 *int)
}
type FoodLossHandler interface {
	HandleFoodLoss(param1 *event.Context, param2 int, param3 int)
}
type HealHandler interface {
	HandleHeal(param1 *event.Context, param2 *float64, param3 healing.Source)
}
type HurtHandler interface {
	HandleHurt(param1 *event.Context, param2 *float64, param3 *time.Duration, param4 damage.Source)
}
type ItemConsumeHandler interface {
	HandleItemConsume(param1 *event.Context, param2 item.Stack)
}
type ItemDamageHandler interface {
	HandleItemDamage(param1 *event.Context, param2 item.Stack, param3 int)
}
type ItemDropHandler interface {
	HandleItemDrop(param1 *event.Context, param2 *entity.Item)
}
type ItemPickupHandler interface {
	HandleItemPickup(param1 *event.Context, param2 item.Stack)
}
type ItemUseHandler interface {
	HandleItemUse(param1 *event.Context)
}
type ItemUseOnBlockHandler interface {
	HandleItemUseOnBlock(param1 *event.Context, param2 cube.Pos, param3 cube.Face, param4 mgl64.Vec3)
}
type ItemUseOnEntityHandler interface {
	HandleItemUseOnEntity(param1 *event.Context, param2 world.Entity)
}
type JumpHandler interface {
	HandleJump()
}
type MoveHandler interface {
	HandleMove(param1 *event.Context, param2 mgl64.Vec3, param3 float64, param4 float64)
}
type PunchAirHandler interface {
	HandlePunchAir(param1 *event.Context)
}
type QuitHandler interface {
	HandleQuit()
}
type RespawnHandler interface {
	HandleRespawn(param1 *mgl64.Vec3, param2 **world.World)
}
type SignEditHandler interface {
	HandleSignEdit(param1 *event.Context, param2 string, param3 string)
}
type SkinChangeHandler interface {
	HandleSkinChange(param1 *event.Context, param2 *skin.Skin)
}
type StartBreakHandler interface {
	HandleStartBreak(param1 *event.Context, param2 cube.Pos)
}
type TeleportHandler interface {
	HandleTeleport(param1 *event.Context, param2 mgl64.Vec3)
}
type ToggleSneakHandler interface {
	HandleToggleSneak(param1 *event.Context, param2 bool)
}
type ToggleSprintHandler interface {
	HandleToggleSprint(param1 *event.Context, param2 bool)
}
type TransferHandler interface {
	HandleTransfer(param1 *event.Context, param2 *net.UDPAddr)
}

func (h *MultipleHandler) HandleAttackEntity(param1 *event.Context, param2 world.Entity, param3 *float64, param4 *float64, param5 *bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(AttackEntityHandler); ok {
			hdr.HandleAttackEntity(param1, param2, param3, param4, param5)
		}
	}
}
func (h *MultipleHandler) HandleBlockBreak(param1 *event.Context, param2 cube.Pos, param3 *[]item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockBreakHandler); ok {
			hdr.HandleBlockBreak(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleBlockPick(param1 *event.Context, param2 cube.Pos, param3 world.Block) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockPickHandler); ok {
			hdr.HandleBlockPick(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleBlockPlace(param1 *event.Context, param2 cube.Pos, param3 world.Block) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockPlaceHandler); ok {
			hdr.HandleBlockPlace(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleChangeWorld(param1 *world.World, param2 *world.World) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ChangeWorldHandler); ok {
			hdr.HandleChangeWorld(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleChat(param1 *event.Context, param2 *string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ChatHandler); ok {
			hdr.HandleChat(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleCommandExecution(param1 *event.Context, param2 cmd.Command, param3 []string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(CommandExecutionHandler); ok {
			hdr.HandleCommandExecution(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleDeath(param1 damage.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(DeathHandler); ok {
			hdr.HandleDeath(param1)
		}
	}
}
func (h *MultipleHandler) HandleExperienceGain(param1 *event.Context, param2 *int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ExperienceGainHandler); ok {
			hdr.HandleExperienceGain(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleFoodLoss(param1 *event.Context, param2 int, param3 int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(FoodLossHandler); ok {
			hdr.HandleFoodLoss(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleHeal(param1 *event.Context, param2 *float64, param3 healing.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(HealHandler); ok {
			hdr.HandleHeal(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleHurt(param1 *event.Context, param2 *float64, param3 *time.Duration, param4 damage.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(HurtHandler); ok {
			hdr.HandleHurt(param1, param2, param3, param4)
		}
	}
}
func (h *MultipleHandler) HandleItemConsume(param1 *event.Context, param2 item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemConsumeHandler); ok {
			hdr.HandleItemConsume(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleItemDamage(param1 *event.Context, param2 item.Stack, param3 int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemDamageHandler); ok {
			hdr.HandleItemDamage(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleItemDrop(param1 *event.Context, param2 *entity.Item) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemDropHandler); ok {
			hdr.HandleItemDrop(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleItemPickup(param1 *event.Context, param2 item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemPickupHandler); ok {
			hdr.HandleItemPickup(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleItemUse(param1 *event.Context) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseHandler); ok {
			hdr.HandleItemUse(param1)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnBlock(param1 *event.Context, param2 cube.Pos, param3 cube.Face, param4 mgl64.Vec3) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
			hdr.HandleItemUseOnBlock(param1, param2, param3, param4)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnEntity(param1 *event.Context, param2 world.Entity) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
			hdr.HandleItemUseOnEntity(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleJump() {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(JumpHandler); ok {
			hdr.HandleJump()
		}
	}
}
func (h *MultipleHandler) HandleMove(param1 *event.Context, param2 mgl64.Vec3, param3 float64, param4 float64) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(MoveHandler); ok {
			hdr.HandleMove(param1, param2, param3, param4)
		}
	}
}
func (h *MultipleHandler) HandlePunchAir(param1 *event.Context) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(PunchAirHandler); ok {
			hdr.HandlePunchAir(param1)
		}
	}
}
func (h *MultipleHandler) HandleQuit() {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(QuitHandler); ok {
			hdr.HandleQuit()
		}
	}
}
func (h *MultipleHandler) HandleRespawn(param1 *mgl64.Vec3, param2 **world.World) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(RespawnHandler); ok {
			hdr.HandleRespawn(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleSignEdit(param1 *event.Context, param2 string, param3 string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(SignEditHandler); ok {
			hdr.HandleSignEdit(param1, param2, param3)
		}
	}
}
func (h *MultipleHandler) HandleSkinChange(param1 *event.Context, param2 *skin.Skin) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(SkinChangeHandler); ok {
			hdr.HandleSkinChange(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleStartBreak(param1 *event.Context, param2 cube.Pos) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(StartBreakHandler); ok {
			hdr.HandleStartBreak(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleTeleport(param1 *event.Context, param2 mgl64.Vec3) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(TeleportHandler); ok {
			hdr.HandleTeleport(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleToggleSneak(param1 *event.Context, param2 bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ToggleSneakHandler); ok {
			hdr.HandleToggleSneak(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleToggleSprint(param1 *event.Context, param2 bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ToggleSprintHandler); ok {
			hdr.HandleToggleSprint(param1, param2)
		}
	}
}
func (h *MultipleHandler) HandleTransfer(param1 *event.Context, param2 *net.UDPAddr) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(TransferHandler); ok {
			hdr.HandleTransfer(param1, param2)
		}
	}
}
