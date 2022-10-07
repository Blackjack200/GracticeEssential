package mhandler

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
	HandleFoodLoss(ctx *event.Context, from, to int)
}
type HealHandler interface {
	// HandleHeal handles the player being healed by a healing source. ctx.Cancel() may be called to cancel
	// the healing.
	// The health added may be changed by assigning to *health.
	HandleHeal(ctx *event.Context, health *float64, src healing.Source)
}
type HurtHandler interface {
	// HandleHurt handles the player being hurt by any damage source. ctx.Cancel() may be called to cancel the
	// damage being dealt to the player.
	// The damage dealt to the player may be changed by assigning to *damage.
	HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src damage.Source)
}
type DeathHandler interface {
	// HandleDeath handles the player dying to a particular damage cause.
	HandleDeath(src damage.Source)
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
	HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack)
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
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(MoveHandler); ok {
			hdr.HandleMove(ctx, newPos, newYaw, newPitch)
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
func (h *MultipleHandler) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(TeleportHandler); ok {
			hdr.HandleTeleport(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleChangeWorld(before, after *world.World) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ChangeWorldHandler); ok {
			hdr.HandleChangeWorld(before, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSprint(ctx *event.Context, after bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ToggleSprintHandler); ok {
			hdr.HandleToggleSprint(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSneak(ctx *event.Context, after bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ToggleSneakHandler); ok {
			hdr.HandleToggleSneak(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleChat(ctx *event.Context, message *string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ChatHandler); ok {
			hdr.HandleChat(ctx, message)
		}
	}
}
func (h *MultipleHandler) HandleFoodLoss(ctx *event.Context, from, to int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(FoodLossHandler); ok {
			hdr.HandleFoodLoss(ctx, from, to)
		}
	}
}
func (h *MultipleHandler) HandleHeal(ctx *event.Context, health *float64, src healing.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(HealHandler); ok {
			hdr.HandleHeal(ctx, health, src)
		}
	}
}
func (h *MultipleHandler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src damage.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(HurtHandler); ok {
			hdr.HandleHurt(ctx, damage, attackImmunity, src)
		}
	}
}
func (h *MultipleHandler) HandleDeath(src damage.Source) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(DeathHandler); ok {
			hdr.HandleDeath(src)
		}
	}
}
func (h *MultipleHandler) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(RespawnHandler); ok {
			hdr.HandleRespawn(pos, w)
		}
	}
}
func (h *MultipleHandler) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(SkinChangeHandler); ok {
			hdr.HandleSkinChange(ctx, skin)
		}
	}
}
func (h *MultipleHandler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(StartBreakHandler); ok {
			hdr.HandleStartBreak(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockBreakHandler); ok {
			hdr.HandleBlockBreak(ctx, pos, drops)
		}
	}
}
func (h *MultipleHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockPlaceHandler); ok {
			hdr.HandleBlockPlace(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(BlockPickHandler); ok {
			hdr.HandleBlockPick(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleItemUse(ctx *event.Context) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseHandler); ok {
			hdr.HandleItemUse(ctx)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
			hdr.HandleItemUseOnBlock(ctx, pos, face, clickPos)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnEntity(ctx *event.Context, e world.Entity) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
			hdr.HandleItemUseOnEntity(ctx, e)
		}
	}
}
func (h *MultipleHandler) HandleItemConsume(ctx *event.Context, item item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemConsumeHandler); ok {
			hdr.HandleItemConsume(ctx, item)
		}
	}
}
func (h *MultipleHandler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(AttackEntityHandler); ok {
			hdr.HandleAttackEntity(ctx, e, force, height, critical)
		}
	}
}
func (h *MultipleHandler) HandleExperienceGain(ctx *event.Context, amount *int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ExperienceGainHandler); ok {
			hdr.HandleExperienceGain(ctx, amount)
		}
	}
}
func (h *MultipleHandler) HandlePunchAir(ctx *event.Context) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(PunchAirHandler); ok {
			hdr.HandlePunchAir(ctx)
		}
	}
}
func (h *MultipleHandler) HandleSignEdit(ctx *event.Context, oldText, newText string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(SignEditHandler); ok {
			hdr.HandleSignEdit(ctx, oldText, newText)
		}
	}
}
func (h *MultipleHandler) HandleItemDamage(ctx *event.Context, i item.Stack, damage int) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemDamageHandler); ok {
			hdr.HandleItemDamage(ctx, i, damage)
		}
	}
}
func (h *MultipleHandler) HandleItemPickup(ctx *event.Context, i item.Stack) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemPickupHandler); ok {
			hdr.HandleItemPickup(ctx, i)
		}
	}
}
func (h *MultipleHandler) HandleItemDrop(ctx *event.Context, e *entity.Item) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(ItemDropHandler); ok {
			hdr.HandleItemDrop(ctx, e)
		}
	}
}
func (h *MultipleHandler) HandleTransfer(ctx *event.Context, addr *net.UDPAddr) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(TransferHandler); ok {
			hdr.HandleTransfer(ctx, addr)
		}
	}
}
func (h *MultipleHandler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	for hdr, _ := range h.handlers {
		if hdr, ok := hdr.(CommandExecutionHandler); ok {
			hdr.HandleCommandExecution(ctx, command, args)
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
