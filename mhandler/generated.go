package mhandler

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/exp/slices"
	"net"
	"time"
)

type MoveHandler interface {
	// HandleMove handles the movement of a player. ctx.Cancel() may be called to cancel the movement event.
	// The new position, yaw and pitch are passed.
	HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64)
}
type JumpHandler interface {
	// HandleJump handles the player jumping.
	HandleJump()
}
type TeleportHandler interface {
	// HandleTeleport handles the teleportation of a player. ctx.Cancel() may be called to cancel it.
	HandleTeleport(ctx *event.Context, pos mgl64.Vec3)
}
type ChangeWorldHandler interface {
	// HandleChangeWorld handles when the player is added to a new world. before may be nil.
	HandleChangeWorld(before, after *world.World)
}
type ToggleSprintHandler interface {
	// HandleToggleSprint handles when the player starts or stops sprinting.
	// After is true if the player is sprinting after toggling (changing their sprinting state).
	HandleToggleSprint(ctx *event.Context, after bool)
}
type ToggleSneakHandler interface {
	// HandleToggleSneak handles when the player starts or stops sneaking.
	// After is true if the player is sneaking after toggling (changing their sneaking state).
	HandleToggleSneak(ctx *event.Context, after bool)
}
type ChatHandler interface {
	// HandleChat handles a message sent in the chat by a player. ctx.Cancel() may be called to cancel the
	// message being sent in chat.
	// The message may be changed by assigning to *message.
	HandleChat(ctx *event.Context, message *string)
}
type FoodLossHandler interface {
	// HandleFoodLoss handles the food bar of a player depleting naturally, for example because the player was
	// sprinting and jumping. ctx.Cancel() may be called to cancel the food points being lost.
	HandleFoodLoss(ctx *event.Context, from int, to *int)
}
type HealHandler interface {
	// HandleHeal handles the player being healed by a healing source. ctx.Cancel() may be called to cancel
	// the healing.
	// The health added may be changed by assigning to *health.
	HandleHeal(ctx *event.Context, health *float64, src world.HealingSource)
}
type HurtHandler interface {
	// HandleHurt handles the player being hurt by any damage source. ctx.Cancel() may be called to cancel the
	// damage being dealt to the player.
	// The damage dealt to the player may be changed by assigning to *damage.
	HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource)
}
type DeathHandler interface {
	// HandleDeath handles the player dying to a particular damage cause.
	HandleDeath(src world.DamageSource, keepInv *bool)
}
type RespawnHandler interface {
	// HandleRespawn handles the respawning of the player in the world. The spawn position passed may be
	// changed by assigning to *pos. The world.World in which the Player is respawned may be modifying by assigning to
	// *w. This world may be the world the Player died in, but it might also point to a different world (the overworld)
	// if the Player died in the nether or end.
	HandleRespawn(pos *mgl64.Vec3, w **world.World)
}
type SkinChangeHandler interface {
	// HandleSkinChange handles the player changing their skin. ctx.Cancel() may be called to cancel the skin
	// change.
	HandleSkinChange(ctx *event.Context, skin *skin.Skin)
}
type StartBreakHandler interface {
	// HandleStartBreak handles the player starting to break a block at the position passed. ctx.Cancel() may
	// be called to stop the player from breaking the block completely.
	HandleStartBreak(ctx *event.Context, pos cube.Pos)
}
type BlockBreakHandler interface {
	// HandleBlockBreak handles a block that is being broken by a player. ctx.Cancel() may be called to cancel
	// the block being broken. A pointer to a slice of the block's drops is passed, and may be altered
	// to change what items will actually be dropped.
	HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int)
}
type BlockPlaceHandler interface {
	// HandleBlockPlace handles the player placing a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being placed.
	HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block)
}
type BlockPickHandler interface {
	// HandleBlockPick handles the player picking a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being picked.
	HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block)
}
type ItemUseHandler interface {
	// HandleItemUse handles the player using an item in the air. It is called for each item, although most
	// will not actually do anything. Items such as snowballs may be thrown if HandleItemUse does not cancel
	// the context using ctx.Cancel(). It is not called if the player is holding no item.
	HandleItemUse(ctx *event.Context)
}
type ItemUseOnBlockHandler interface {
	// HandleItemUseOnBlock handles the player using the item held in its main hand on a block at the block
	// position passed. The face of the block clicked is also passed, along with the relative click position.
	// The click position has X, Y and Z values which are all in the range 0.0-1.0. It is also called if the
	// player is holding no item.
	HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
}
type ItemUseOnEntityHandler interface {
	// HandleItemUseOnEntity handles the player using the item held in its main hand on an entity passed to
	// the method.
	// HandleItemUseOnEntity is always called when a player uses an item on an entity, regardless of whether
	// the item actually does anything when used on an entity. It is also called if the player is holding no
	// item.
	HandleItemUseOnEntity(ctx *event.Context, e world.Entity)
}
type ItemConsumeHandler interface {
	// HandleItemConsume handles the player consuming an item. This is called whenever a consumable such as
	// food is consumed.
	HandleItemConsume(ctx *event.Context, item item.Stack)
}
type AttackEntityHandler interface {
	// HandleAttackEntity handles the player attacking an entity using the item held in its hand. ctx.Cancel()
	// may be called to cancel the attack, which will cancel damage dealt to the target and will stop the
	// entity from being knocked back.
	// The entity attacked may not be alive (implements entity.Living), in which case no damage will be dealt
	// and the target won't be knocked back.
	// The entity attacked may also be immune when this method is called, in which case no damage and knock-
	// back will be dealt.
	// The knock back force and height is also provided which can be modified.
	// The attack can be a critical attack, which would increase damage by a factor of 1.5 and
	// spawn critical hit particles around the target entity. These particles will not be displayed
	// if no damage is dealt.
	HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool)
}
type ExperienceGainHandler interface {
	// HandleExperienceGain handles the player gaining experience. ctx.Cancel() may be called to cancel
	// the gain.
	// The amount is also provided which can be modified.
	HandleExperienceGain(ctx *event.Context, amount *int)
}
type PunchAirHandler interface {
	// HandlePunchAir handles the player punching air.
	HandlePunchAir(ctx *event.Context)
}
type SignEditHandler interface {
	// HandleSignEdit handles the player editing a sign. It is called for every keystroke while editing a sign and
	// has both the old text passed and the text after the edit. This typically only has a change of one character.
	HandleSignEdit(ctx *event.Context, oldText, newText string)
}
type ItemDamageHandler interface {
	// HandleItemDamage handles the event wherein the item either held by the player or as armour takes
	// damage through usage.
	// The type of the item may be checked to determine whether it was armour or a tool used. The damage to
	// the item is passed.
	HandleItemDamage(ctx *event.Context, i item.Stack, damage int)
}
type ItemPickupHandler interface {
	// HandleItemPickup handles the player picking up an item from the ground. The item stack laying on the
	// ground is passed. ctx.Cancel() may be called to prevent the player from picking up the item.
	HandleItemPickup(ctx *event.Context, i item.Stack)
}
type ItemDropHandler interface {
	// HandleItemDrop handles the player dropping an item on the ground. The dropped item entity is passed.
	// ctx.Cancel() may be called to prevent the player from dropping the entity.Item passed on the ground.
	// e.Item() may be called to obtain the item stack dropped.
	HandleItemDrop(ctx *event.Context, e *entity.Item)
}
type TransferHandler interface {
	// HandleTransfer handles a player being transferred to another server. ctx.Cancel() may be called to
	// cancel the transfer.
	HandleTransfer(ctx *event.Context, addr *net.UDPAddr)
}
type CommandExecutionHandler interface {
	// HandleCommandExecution handles the command execution of a player, who wrote a command in the chat.
	// ctx.Cancel() may be called to cancel the command execution.
	HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string)
}
type QuitHandler interface {
	// HandleQuit handles the closing of a player. It is always called when the player is disconnected,
	// regardless of the reason.
	HandleQuit()
}

func (h *MultipleHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	for _, hdr := range h._MoveHandler {
		if hdr, ok := hdr.(MoveHandler); ok {
			hdr.HandleMove(ctx, newPos, newYaw, newPitch)
		}
	}
}
func (h *MultipleHandler) HandleJump() {
	for _, hdr := range h._JumpHandler {
		if hdr, ok := hdr.(JumpHandler); ok {
			hdr.HandleJump()
		}
	}
}
func (h *MultipleHandler) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	for _, hdr := range h._TeleportHandler {
		if hdr, ok := hdr.(TeleportHandler); ok {
			hdr.HandleTeleport(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleChangeWorld(before, after *world.World) {
	for _, hdr := range h._ChangeWorldHandler {
		if hdr, ok := hdr.(ChangeWorldHandler); ok {
			hdr.HandleChangeWorld(before, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSprint(ctx *event.Context, after bool) {
	for _, hdr := range h._ToggleSprintHandler {
		if hdr, ok := hdr.(ToggleSprintHandler); ok {
			hdr.HandleToggleSprint(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSneak(ctx *event.Context, after bool) {
	for _, hdr := range h._ToggleSneakHandler {
		if hdr, ok := hdr.(ToggleSneakHandler); ok {
			hdr.HandleToggleSneak(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleChat(ctx *event.Context, message *string) {
	for _, hdr := range h._ChatHandler {
		if hdr, ok := hdr.(ChatHandler); ok {
			hdr.HandleChat(ctx, message)
		}
	}
}
func (h *MultipleHandler) HandleFoodLoss(ctx *event.Context, from int, to *int) {
	for _, hdr := range h._FoodLossHandler {
		if hdr, ok := hdr.(FoodLossHandler); ok {
			hdr.HandleFoodLoss(ctx, from, to)
		}
	}
}
func (h *MultipleHandler) HandleHeal(ctx *event.Context, health *float64, src world.HealingSource) {
	for _, hdr := range h._HealHandler {
		if hdr, ok := hdr.(HealHandler); ok {
			hdr.HandleHeal(ctx, health, src)
		}
	}
}
func (h *MultipleHandler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	for _, hdr := range h._HurtHandler {
		if hdr, ok := hdr.(HurtHandler); ok {
			hdr.HandleHurt(ctx, damage, attackImmunity, src)
		}
	}
}
func (h *MultipleHandler) HandleDeath(src world.DamageSource, keepInv *bool) {
	for _, hdr := range h._DeathHandler {
		if hdr, ok := hdr.(DeathHandler); ok {
			hdr.HandleDeath(src, keepInv)
		}
	}
}
func (h *MultipleHandler) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	for _, hdr := range h._RespawnHandler {
		if hdr, ok := hdr.(RespawnHandler); ok {
			hdr.HandleRespawn(pos, w)
		}
	}
}
func (h *MultipleHandler) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	for _, hdr := range h._SkinChangeHandler {
		if hdr, ok := hdr.(SkinChangeHandler); ok {
			hdr.HandleSkinChange(ctx, skin)
		}
	}
}
func (h *MultipleHandler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	for _, hdr := range h._StartBreakHandler {
		if hdr, ok := hdr.(StartBreakHandler); ok {
			hdr.HandleStartBreak(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	for _, hdr := range h._BlockBreakHandler {
		if hdr, ok := hdr.(BlockBreakHandler); ok {
			hdr.HandleBlockBreak(ctx, pos, drops, xp)
		}
	}
}
func (h *MultipleHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	for _, hdr := range h._BlockPlaceHandler {
		if hdr, ok := hdr.(BlockPlaceHandler); ok {
			hdr.HandleBlockPlace(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block) {
	for _, hdr := range h._BlockPickHandler {
		if hdr, ok := hdr.(BlockPickHandler); ok {
			hdr.HandleBlockPick(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleItemUse(ctx *event.Context) {
	for _, hdr := range h._ItemUseHandler {
		if hdr, ok := hdr.(ItemUseHandler); ok {
			hdr.HandleItemUse(ctx)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	for _, hdr := range h._ItemUseOnBlockHandler {
		if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
			hdr.HandleItemUseOnBlock(ctx, pos, face, clickPos)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnEntity(ctx *event.Context, e world.Entity) {
	for _, hdr := range h._ItemUseOnEntityHandler {
		if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
			hdr.HandleItemUseOnEntity(ctx, e)
		}
	}
}
func (h *MultipleHandler) HandleItemConsume(ctx *event.Context, item item.Stack) {
	for _, hdr := range h._ItemConsumeHandler {
		if hdr, ok := hdr.(ItemConsumeHandler); ok {
			hdr.HandleItemConsume(ctx, item)
		}
	}
}
func (h *MultipleHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	for _, hdr := range h._AttackEntityHandler {
		if hdr, ok := hdr.(AttackEntityHandler); ok {
			hdr.HandleAttackEntity(ctx, e, force, height, critical)
		}
	}
}
func (h *MultipleHandler) HandleExperienceGain(ctx *event.Context, amount *int) {
	for _, hdr := range h._ExperienceGainHandler {
		if hdr, ok := hdr.(ExperienceGainHandler); ok {
			hdr.HandleExperienceGain(ctx, amount)
		}
	}
}
func (h *MultipleHandler) HandlePunchAir(ctx *event.Context) {
	for _, hdr := range h._PunchAirHandler {
		if hdr, ok := hdr.(PunchAirHandler); ok {
			hdr.HandlePunchAir(ctx)
		}
	}
}
func (h *MultipleHandler) HandleSignEdit(ctx *event.Context, oldText, newText string) {
	for _, hdr := range h._SignEditHandler {
		if hdr, ok := hdr.(SignEditHandler); ok {
			hdr.HandleSignEdit(ctx, oldText, newText)
		}
	}
}
func (h *MultipleHandler) HandleItemDamage(ctx *event.Context, i item.Stack, damage int) {
	for _, hdr := range h._ItemDamageHandler {
		if hdr, ok := hdr.(ItemDamageHandler); ok {
			hdr.HandleItemDamage(ctx, i, damage)
		}
	}
}
func (h *MultipleHandler) HandleItemPickup(ctx *event.Context, i item.Stack) {
	for _, hdr := range h._ItemPickupHandler {
		if hdr, ok := hdr.(ItemPickupHandler); ok {
			hdr.HandleItemPickup(ctx, i)
		}
	}
}
func (h *MultipleHandler) HandleItemDrop(ctx *event.Context, e *entity.Item) {
	for _, hdr := range h._ItemDropHandler {
		if hdr, ok := hdr.(ItemDropHandler); ok {
			hdr.HandleItemDrop(ctx, e)
		}
	}
}
func (h *MultipleHandler) HandleTransfer(ctx *event.Context, addr *net.UDPAddr) {
	for _, hdr := range h._TransferHandler {
		if hdr, ok := hdr.(TransferHandler); ok {
			hdr.HandleTransfer(ctx, addr)
		}
	}
}
func (h *MultipleHandler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	for _, hdr := range h._CommandExecutionHandler {
		if hdr, ok := hdr.(CommandExecutionHandler); ok {
			hdr.HandleCommandExecution(ctx, command, args)
		}
	}
}
func (h *MultipleHandler) HandleQuit() {
	for _, hdr := range h._QuitHandler {
		if hdr, ok := hdr.(QuitHandler); ok {
			hdr.HandleQuit()
		}
	}
}

type MultipleHandler struct {
	_MoveHandler             []MoveHandler
	_JumpHandler             []JumpHandler
	_TeleportHandler         []TeleportHandler
	_ChangeWorldHandler      []ChangeWorldHandler
	_ToggleSprintHandler     []ToggleSprintHandler
	_ToggleSneakHandler      []ToggleSneakHandler
	_ChatHandler             []ChatHandler
	_FoodLossHandler         []FoodLossHandler
	_HealHandler             []HealHandler
	_HurtHandler             []HurtHandler
	_DeathHandler            []DeathHandler
	_RespawnHandler          []RespawnHandler
	_SkinChangeHandler       []SkinChangeHandler
	_StartBreakHandler       []StartBreakHandler
	_BlockBreakHandler       []BlockBreakHandler
	_BlockPlaceHandler       []BlockPlaceHandler
	_BlockPickHandler        []BlockPickHandler
	_ItemUseHandler          []ItemUseHandler
	_ItemUseOnBlockHandler   []ItemUseOnBlockHandler
	_ItemUseOnEntityHandler  []ItemUseOnEntityHandler
	_ItemConsumeHandler      []ItemConsumeHandler
	_AttackEntityHandler     []AttackEntityHandler
	_ExperienceGainHandler   []ExperienceGainHandler
	_PunchAirHandler         []PunchAirHandler
	_SignEditHandler         []SignEditHandler
	_ItemDamageHandler       []ItemDamageHandler
	_ItemPickupHandler       []ItemPickupHandler
	_ItemDropHandler         []ItemDropHandler
	_TransferHandler         []TransferHandler
	_CommandExecutionHandler []CommandExecutionHandler
	_QuitHandler             []QuitHandler
}

func (h *MultipleHandler) Register(hdr any) func() {
	reg := false
	var funcs []func()
	if hdr, ok := hdr.(MoveHandler); ok {
		if k := slices.Contains(h._MoveHandler, hdr); !k {
			h._MoveHandler = append(h._MoveHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._MoveHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._MoveHandler = slices.Delete(h._MoveHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(JumpHandler); ok {
		if k := slices.Contains(h._JumpHandler, hdr); !k {
			h._JumpHandler = append(h._JumpHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._JumpHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._JumpHandler = slices.Delete(h._JumpHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(TeleportHandler); ok {
		if k := slices.Contains(h._TeleportHandler, hdr); !k {
			h._TeleportHandler = append(h._TeleportHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._TeleportHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._TeleportHandler = slices.Delete(h._TeleportHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ChangeWorldHandler); ok {
		if k := slices.Contains(h._ChangeWorldHandler, hdr); !k {
			h._ChangeWorldHandler = append(h._ChangeWorldHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ChangeWorldHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ChangeWorldHandler = slices.Delete(h._ChangeWorldHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ToggleSprintHandler); ok {
		if k := slices.Contains(h._ToggleSprintHandler, hdr); !k {
			h._ToggleSprintHandler = append(h._ToggleSprintHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ToggleSprintHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ToggleSprintHandler = slices.Delete(h._ToggleSprintHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ToggleSneakHandler); ok {
		if k := slices.Contains(h._ToggleSneakHandler, hdr); !k {
			h._ToggleSneakHandler = append(h._ToggleSneakHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ToggleSneakHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ToggleSneakHandler = slices.Delete(h._ToggleSneakHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ChatHandler); ok {
		if k := slices.Contains(h._ChatHandler, hdr); !k {
			h._ChatHandler = append(h._ChatHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ChatHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ChatHandler = slices.Delete(h._ChatHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(FoodLossHandler); ok {
		if k := slices.Contains(h._FoodLossHandler, hdr); !k {
			h._FoodLossHandler = append(h._FoodLossHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._FoodLossHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._FoodLossHandler = slices.Delete(h._FoodLossHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(HealHandler); ok {
		if k := slices.Contains(h._HealHandler, hdr); !k {
			h._HealHandler = append(h._HealHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._HealHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._HealHandler = slices.Delete(h._HealHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(HurtHandler); ok {
		if k := slices.Contains(h._HurtHandler, hdr); !k {
			h._HurtHandler = append(h._HurtHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._HurtHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._HurtHandler = slices.Delete(h._HurtHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(DeathHandler); ok {
		if k := slices.Contains(h._DeathHandler, hdr); !k {
			h._DeathHandler = append(h._DeathHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._DeathHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._DeathHandler = slices.Delete(h._DeathHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(RespawnHandler); ok {
		if k := slices.Contains(h._RespawnHandler, hdr); !k {
			h._RespawnHandler = append(h._RespawnHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._RespawnHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._RespawnHandler = slices.Delete(h._RespawnHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(SkinChangeHandler); ok {
		if k := slices.Contains(h._SkinChangeHandler, hdr); !k {
			h._SkinChangeHandler = append(h._SkinChangeHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._SkinChangeHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._SkinChangeHandler = slices.Delete(h._SkinChangeHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(StartBreakHandler); ok {
		if k := slices.Contains(h._StartBreakHandler, hdr); !k {
			h._StartBreakHandler = append(h._StartBreakHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._StartBreakHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._StartBreakHandler = slices.Delete(h._StartBreakHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(BlockBreakHandler); ok {
		if k := slices.Contains(h._BlockBreakHandler, hdr); !k {
			h._BlockBreakHandler = append(h._BlockBreakHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._BlockBreakHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._BlockBreakHandler = slices.Delete(h._BlockBreakHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(BlockPlaceHandler); ok {
		if k := slices.Contains(h._BlockPlaceHandler, hdr); !k {
			h._BlockPlaceHandler = append(h._BlockPlaceHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._BlockPlaceHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._BlockPlaceHandler = slices.Delete(h._BlockPlaceHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(BlockPickHandler); ok {
		if k := slices.Contains(h._BlockPickHandler, hdr); !k {
			h._BlockPickHandler = append(h._BlockPickHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._BlockPickHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._BlockPickHandler = slices.Delete(h._BlockPickHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemUseHandler); ok {
		if k := slices.Contains(h._ItemUseHandler, hdr); !k {
			h._ItemUseHandler = append(h._ItemUseHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemUseHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemUseHandler = slices.Delete(h._ItemUseHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
		if k := slices.Contains(h._ItemUseOnBlockHandler, hdr); !k {
			h._ItemUseOnBlockHandler = append(h._ItemUseOnBlockHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemUseOnBlockHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemUseOnBlockHandler = slices.Delete(h._ItemUseOnBlockHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
		if k := slices.Contains(h._ItemUseOnEntityHandler, hdr); !k {
			h._ItemUseOnEntityHandler = append(h._ItemUseOnEntityHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemUseOnEntityHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemUseOnEntityHandler = slices.Delete(h._ItemUseOnEntityHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemConsumeHandler); ok {
		if k := slices.Contains(h._ItemConsumeHandler, hdr); !k {
			h._ItemConsumeHandler = append(h._ItemConsumeHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemConsumeHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemConsumeHandler = slices.Delete(h._ItemConsumeHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(AttackEntityHandler); ok {
		if k := slices.Contains(h._AttackEntityHandler, hdr); !k {
			h._AttackEntityHandler = append(h._AttackEntityHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._AttackEntityHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._AttackEntityHandler = slices.Delete(h._AttackEntityHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ExperienceGainHandler); ok {
		if k := slices.Contains(h._ExperienceGainHandler, hdr); !k {
			h._ExperienceGainHandler = append(h._ExperienceGainHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ExperienceGainHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ExperienceGainHandler = slices.Delete(h._ExperienceGainHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(PunchAirHandler); ok {
		if k := slices.Contains(h._PunchAirHandler, hdr); !k {
			h._PunchAirHandler = append(h._PunchAirHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._PunchAirHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._PunchAirHandler = slices.Delete(h._PunchAirHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(SignEditHandler); ok {
		if k := slices.Contains(h._SignEditHandler, hdr); !k {
			h._SignEditHandler = append(h._SignEditHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._SignEditHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._SignEditHandler = slices.Delete(h._SignEditHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemDamageHandler); ok {
		if k := slices.Contains(h._ItemDamageHandler, hdr); !k {
			h._ItemDamageHandler = append(h._ItemDamageHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemDamageHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemDamageHandler = slices.Delete(h._ItemDamageHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemPickupHandler); ok {
		if k := slices.Contains(h._ItemPickupHandler, hdr); !k {
			h._ItemPickupHandler = append(h._ItemPickupHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemPickupHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemPickupHandler = slices.Delete(h._ItemPickupHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(ItemDropHandler); ok {
		if k := slices.Contains(h._ItemDropHandler, hdr); !k {
			h._ItemDropHandler = append(h._ItemDropHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._ItemDropHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._ItemDropHandler = slices.Delete(h._ItemDropHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(TransferHandler); ok {
		if k := slices.Contains(h._TransferHandler, hdr); !k {
			h._TransferHandler = append(h._TransferHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._TransferHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._TransferHandler = slices.Delete(h._TransferHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(CommandExecutionHandler); ok {
		if k := slices.Contains(h._CommandExecutionHandler, hdr); !k {
			h._CommandExecutionHandler = append(h._CommandExecutionHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._CommandExecutionHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._CommandExecutionHandler = slices.Delete(h._CommandExecutionHandler, idx, idx+1)
			})
		}
	}
	if hdr, ok := hdr.(QuitHandler); ok {
		if k := slices.Contains(h._QuitHandler, hdr); !k {
			h._QuitHandler = append(h._QuitHandler, hdr)
			reg = true
			funcs = append(funcs, func() {
				idx := slices.Index(h._QuitHandler, hdr)
				if idx == -1 {
					panic("this should not happened")
				}
				h._QuitHandler = slices.Delete(h._QuitHandler, idx, idx+1)
			})
		}
	}
	if !reg {
		panic("not a valid handler")
	}
	return func() {
		for _, f := range funcs {
			f()
		}
	}
}
func (h *MultipleHandler) Clear() {
	h._MoveHandler = nil
	h._JumpHandler = nil
	h._TeleportHandler = nil
	h._ChangeWorldHandler = nil
	h._ToggleSprintHandler = nil
	h._ToggleSneakHandler = nil
	h._ChatHandler = nil
	h._FoodLossHandler = nil
	h._HealHandler = nil
	h._HurtHandler = nil
	h._DeathHandler = nil
	h._RespawnHandler = nil
	h._SkinChangeHandler = nil
	h._StartBreakHandler = nil
	h._BlockBreakHandler = nil
	h._BlockPlaceHandler = nil
	h._BlockPickHandler = nil
	h._ItemUseHandler = nil
	h._ItemUseOnBlockHandler = nil
	h._ItemUseOnEntityHandler = nil
	h._ItemConsumeHandler = nil
	h._AttackEntityHandler = nil
	h._ExperienceGainHandler = nil
	h._PunchAirHandler = nil
	h._SignEditHandler = nil
	h._ItemDamageHandler = nil
	h._ItemPickupHandler = nil
	h._ItemDropHandler = nil
	h._TransferHandler = nil
	h._CommandExecutionHandler = nil
	h._QuitHandler = nil
}
