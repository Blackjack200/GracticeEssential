package mhandler

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"net"
	"time"
)

type MoveHandler interface {
	// HandleMove handles the movement of a player. ctx.Cancel() may be called to cancel the movement event.
	// The new position, yaw and pitch are passed.
	HandleMove(ctx *event.Context[*player.Player], newPos mgl64.Vec3, newRot cube.Rotation)
}
type JumpHandler interface {
	// HandleJump handles the player jumping.
	HandleJump(p *player.Player)
}
type TeleportHandler interface {
	// HandleTeleport handles the teleportation of a player. ctx.Cancel() may be called to cancel it.
	HandleTeleport(ctx *event.Context[*player.Player], pos mgl64.Vec3)
}
type ChangeWorldHandler interface {
	// HandleChangeWorld handles when the player is added to a new world. before may be nil.
	HandleChangeWorld(p *player.Player, before, after *world.World)
}
type ToggleSprintHandler interface {
	// HandleToggleSprint handles when the player starts or stops sprinting.
	// After is true if the player is sprinting after toggling (changing their sprinting state).
	HandleToggleSprint(ctx *event.Context[*player.Player], after bool)
}
type ToggleSneakHandler interface {
	// HandleToggleSneak handles when the player starts or stops sneaking.
	// After is true if the player is sneaking after toggling (changing their sneaking state).
	HandleToggleSneak(ctx *event.Context[*player.Player], after bool)
}
type ChatHandler interface {
	// HandleChat handles a message sent in the chat by a player. ctx.Cancel() may be called to cancel the
	// message being sent in chat.
	// The message may be changed by assigning to *message.
	HandleChat(ctx *event.Context[*player.Player], message *string)
}
type FoodLossHandler interface {
	// HandleFoodLoss handles the food bar of a player depleting naturally, for example because the player was
	// sprinting and jumping. ctx.Cancel() may be called to cancel the food points being lost.
	HandleFoodLoss(ctx *event.Context[*player.Player], from int, to *int)
}
type HealHandler interface {
	// HandleHeal handles the player being healed by a healing source. ctx.Cancel() may be called to cancel
	// the healing.
	// The health added may be changed by assigning to *health.
	HandleHeal(ctx *event.Context[*player.Player], health *float64, src world.HealingSource)
}
type HurtHandler interface {
	// HandleHurt handles the player being hurt by any damage source. ctx.Cancel() may be called to cancel the
	// damage being dealt to the player.
	// The damage dealt to the player may be changed by assigning to *damage.
	// *damage is the final damage dealt to the player. Immune is set to true
	// if the player was hurt during an immunity frame with higher damage than
	// the original cause of the immunity frame. In this case, the damage is
	// reduced but the player is still knocked back.
	HandleHurt(ctx *event.Context[*player.Player], damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource)
}
type DeathHandler interface {
	// HandleDeath handles the player dying to a particular damage cause.
	HandleDeath(p *player.Player, src world.DamageSource, keepInv *bool)
}
type RespawnHandler interface {
	// HandleRespawn handles the respawning of the player in the world. The spawn position passed may be
	// changed by assigning to *pos. The world.World in which the Player is respawned may be modifying by assigning to
	// *w. This world may be the world the Player died in, but it might also point to a different world (the overworld)
	// if the Player died in the nether or end.
	HandleRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World)
}
type SkinChangeHandler interface {
	// HandleSkinChange handles the player changing their skin. ctx.Cancel() may be called to cancel the skin
	// change.
	HandleSkinChange(ctx *event.Context[*player.Player], skin *skin.Skin)
}
type FireExtinguishHandler interface {
	// HandleFireExtinguish handles the player extinguishing a fire at a specific position. ctx.Cancel() may
	// be called to cancel the fire being extinguished.
	// cube.Pos can be used to see where was the fire extinguished, may be used to cancel this on specific positions.
	HandleFireExtinguish(ctx *event.Context[*player.Player], pos cube.Pos)
}
type StartBreakHandler interface {
	// HandleStartBreak handles the player starting to break a block at the position passed. ctx.Cancel() may
	// be called to stop the player from breaking the block completely.
	HandleStartBreak(ctx *event.Context[*player.Player], pos cube.Pos)
}
type BlockBreakHandler interface {
	// HandleBlockBreak handles a block that is being broken by a player. ctx.Cancel() may be called to cancel
	// the block being broken. A pointer to a slice of the block's drops is passed, and may be altered
	// to change what items will actually be dropped.
	HandleBlockBreak(ctx *event.Context[*player.Player], pos cube.Pos, drops *[]item.Stack, xp *int)
}
type BlockPlaceHandler interface {
	// HandleBlockPlace handles the player placing a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being placed.
	HandleBlockPlace(ctx *event.Context[*player.Player], pos cube.Pos, b world.Block)
}
type BlockPickHandler interface {
	// HandleBlockPick handles the player picking a specific block at a position in its world. ctx.Cancel()
	// may be called to cancel the block being picked.
	HandleBlockPick(ctx *event.Context[*player.Player], pos cube.Pos, b world.Block)
}
type ItemUseHandler interface {
	// HandleItemUse handles the player using an item in the air. It is called for each item, although most
	// will not actually do anything. Items such as snowballs may be thrown if HandleItemUse does not cancel
	// the context using ctx.Cancel(). It is not called if the player is holding no item.
	HandleItemUse(ctx *event.Context[*player.Player])
}
type ItemUseOnBlockHandler interface {
	// HandleItemUseOnBlock handles the player using the item held in its main hand on a block at the block
	// position passed. The face of the block clicked is also passed, along with the relative click position.
	// The click position has X, Y and Z values which are all in the range 0.0-1.0. It is also called if the
	// player is holding no item.
	HandleItemUseOnBlock(ctx *event.Context[*player.Player], pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
}
type ItemUseOnEntityHandler interface {
	// HandleItemUseOnEntity handles the player using the item held in its main hand on an entity passed to
	// the method.
	// HandleItemUseOnEntity is always called when a player uses an item on an entity, regardless of whether
	// the item actually does anything when used on an entity. It is also called if the player is holding no
	// item.
	HandleItemUseOnEntity(ctx *event.Context[*player.Player], e world.Entity)
}
type ItemReleaseHandler interface {
	// HandleItemRelease handles the player releasing an item after using it for
	// a particular duration. These include items such as bows.
	HandleItemRelease(ctx *event.Context[*player.Player], item item.Stack, dur time.Duration)
}
type ItemConsumeHandler interface {
	// HandleItemConsume handles the player consuming an item. This is called whenever a consumable such as
	// food is consumed.
	HandleItemConsume(ctx *event.Context[*player.Player], item item.Stack)
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
	HandleAttackEntity(ctx *event.Context[*player.Player], e world.Entity, force, height *float64, critical *bool)
}
type ExperienceGainHandler interface {
	// HandleExperienceGain handles the player gaining experience. ctx.Cancel() may be called to cancel
	// the gain.
	// The amount is also provided which can be modified.
	HandleExperienceGain(ctx *event.Context[*player.Player], amount *int)
}
type PunchAirHandler interface {
	// HandlePunchAir handles the player punching air.
	HandlePunchAir(ctx *event.Context[*player.Player])
}
type SignEditHandler interface {
	// HandleSignEdit handles the player editing a sign. It is called for every keystroke while editing a sign and
	// has both the old text passed and the text after the edit. This typically only has a change of one character.
	HandleSignEdit(ctx *event.Context[*player.Player], pos cube.Pos, frontSide bool, oldText, newText string)
}
type LecternPageTurnHandler interface {
	// HandleLecternPageTurn handles the player turning a page in a lectern. ctx.Cancel() may be called to cancel the
	// page turn. The page number may be changed by assigning to *page.
	HandleLecternPageTurn(ctx *event.Context[*player.Player], pos cube.Pos, oldPage int, newPage *int)
}
type ItemDamageHandler interface {
	// HandleItemDamage handles the event wherein the item either held by the player or as armour takes
	// damage through usage.
	// The type of the item may be checked to determine whether it was armour or a tool used. The damage to
	// the item is passed.
	HandleItemDamage(ctx *event.Context[*player.Player], i item.Stack, damage int)
}
type ItemPickupHandler interface {
	// HandleItemPickup handles the player picking up an item from the ground. The item stack laying on the
	// ground is passed. ctx.Cancel() may be called to prevent the player from picking up the item.
	HandleItemPickup(ctx *event.Context[*player.Player], i *item.Stack)
}
type HeldSlotChangeHandler interface {
	// HandleHeldSlotChange handles the player changing the slot they are currently holding.
	HandleHeldSlotChange(ctx *event.Context[*player.Player], from, to int)
}
type ItemDropHandler interface {
	// HandleItemDrop handles the player dropping an item on the ground. The dropped item entity is passed.
	// ctx.Cancel() may be called to prevent the player from dropping the entity.Item passed on the ground.
	// e.Item() may be called to obtain the item stack dropped.
	HandleItemDrop(ctx *event.Context[*player.Player], s item.Stack)
}
type TransferHandler interface {
	// HandleTransfer handles a player being transferred to another server. ctx.Cancel() may be called to
	// cancel the transfer.
	HandleTransfer(ctx *event.Context[*player.Player], addr *net.UDPAddr)
}
type CommandExecutionHandler interface {
	// HandleCommandExecution handles the command execution of a player, who wrote a command in the chat.
	// ctx.Cancel() may be called to cancel the command execution.
	HandleCommandExecution(ctx *event.Context[*player.Player], command cmd.Command, args []string)
}
type QuitHandler interface {
	// HandleQuit handles the closing of a player. It is always called when the player is disconnected,
	// regardless of the reason.
	HandleQuit(p *player.Player)
}
type DiagnosticsHandler interface {
	// HandleDiagnostics handles the latest diagnostics data that the player has sent to the server. This is
	// not sent by every client however, only those with the "Creator > Enable Client Diagnostics" setting
	// enabled.
	HandleDiagnostics(p *player.Player, d session.Diagnostics)
}

func (h *MultipleHandler) HandleMove(ctx *event.Context[*player.Player], newPos mgl64.Vec3, newRot cube.Rotation) {
	for _, hdr := range h._MoveHandler {
		if hdr, ok := hdr.(MoveHandler); ok {
			hdr.HandleMove(ctx, newPos, newRot)
		}
	}
}
func (h *MultipleHandler) HandleJump(p *player.Player) {
	for _, hdr := range h._JumpHandler {
		if hdr, ok := hdr.(JumpHandler); ok {
			hdr.HandleJump(p)
		}
	}
}
func (h *MultipleHandler) HandleTeleport(ctx *event.Context[*player.Player], pos mgl64.Vec3) {
	for _, hdr := range h._TeleportHandler {
		if hdr, ok := hdr.(TeleportHandler); ok {
			hdr.HandleTeleport(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleChangeWorld(p *player.Player, before, after *world.World) {
	for _, hdr := range h._ChangeWorldHandler {
		if hdr, ok := hdr.(ChangeWorldHandler); ok {
			hdr.HandleChangeWorld(p, before, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSprint(ctx *event.Context[*player.Player], after bool) {
	for _, hdr := range h._ToggleSprintHandler {
		if hdr, ok := hdr.(ToggleSprintHandler); ok {
			hdr.HandleToggleSprint(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleToggleSneak(ctx *event.Context[*player.Player], after bool) {
	for _, hdr := range h._ToggleSneakHandler {
		if hdr, ok := hdr.(ToggleSneakHandler); ok {
			hdr.HandleToggleSneak(ctx, after)
		}
	}
}
func (h *MultipleHandler) HandleChat(ctx *event.Context[*player.Player], message *string) {
	for _, hdr := range h._ChatHandler {
		if hdr, ok := hdr.(ChatHandler); ok {
			hdr.HandleChat(ctx, message)
		}
	}
}
func (h *MultipleHandler) HandleFoodLoss(ctx *event.Context[*player.Player], from int, to *int) {
	for _, hdr := range h._FoodLossHandler {
		if hdr, ok := hdr.(FoodLossHandler); ok {
			hdr.HandleFoodLoss(ctx, from, to)
		}
	}
}
func (h *MultipleHandler) HandleHeal(ctx *event.Context[*player.Player], health *float64, src world.HealingSource) {
	for _, hdr := range h._HealHandler {
		if hdr, ok := hdr.(HealHandler); ok {
			hdr.HandleHeal(ctx, health, src)
		}
	}
}
func (h *MultipleHandler) HandleHurt(ctx *event.Context[*player.Player], damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	for _, hdr := range h._HurtHandler {
		if hdr, ok := hdr.(HurtHandler); ok {
			hdr.HandleHurt(ctx, damage, immune, attackImmunity, src)
		}
	}
}
func (h *MultipleHandler) HandleDeath(p *player.Player, src world.DamageSource, keepInv *bool) {
	for _, hdr := range h._DeathHandler {
		if hdr, ok := hdr.(DeathHandler); ok {
			hdr.HandleDeath(p, src, keepInv)
		}
	}
}
func (h *MultipleHandler) HandleRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World) {
	for _, hdr := range h._RespawnHandler {
		if hdr, ok := hdr.(RespawnHandler); ok {
			hdr.HandleRespawn(p, pos, w)
		}
	}
}
func (h *MultipleHandler) HandleSkinChange(ctx *event.Context[*player.Player], skin *skin.Skin) {
	for _, hdr := range h._SkinChangeHandler {
		if hdr, ok := hdr.(SkinChangeHandler); ok {
			hdr.HandleSkinChange(ctx, skin)
		}
	}
}
func (h *MultipleHandler) HandleFireExtinguish(ctx *event.Context[*player.Player], pos cube.Pos) {
	for _, hdr := range h._FireExtinguishHandler {
		if hdr, ok := hdr.(FireExtinguishHandler); ok {
			hdr.HandleFireExtinguish(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleStartBreak(ctx *event.Context[*player.Player], pos cube.Pos) {
	for _, hdr := range h._StartBreakHandler {
		if hdr, ok := hdr.(StartBreakHandler); ok {
			hdr.HandleStartBreak(ctx, pos)
		}
	}
}
func (h *MultipleHandler) HandleBlockBreak(ctx *event.Context[*player.Player], pos cube.Pos, drops *[]item.Stack, xp *int) {
	for _, hdr := range h._BlockBreakHandler {
		if hdr, ok := hdr.(BlockBreakHandler); ok {
			hdr.HandleBlockBreak(ctx, pos, drops, xp)
		}
	}
}
func (h *MultipleHandler) HandleBlockPlace(ctx *event.Context[*player.Player], pos cube.Pos, b world.Block) {
	for _, hdr := range h._BlockPlaceHandler {
		if hdr, ok := hdr.(BlockPlaceHandler); ok {
			hdr.HandleBlockPlace(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleBlockPick(ctx *event.Context[*player.Player], pos cube.Pos, b world.Block) {
	for _, hdr := range h._BlockPickHandler {
		if hdr, ok := hdr.(BlockPickHandler); ok {
			hdr.HandleBlockPick(ctx, pos, b)
		}
	}
}
func (h *MultipleHandler) HandleItemUse(ctx *event.Context[*player.Player]) {
	for _, hdr := range h._ItemUseHandler {
		if hdr, ok := hdr.(ItemUseHandler); ok {
			hdr.HandleItemUse(ctx)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnBlock(ctx *event.Context[*player.Player], pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	for _, hdr := range h._ItemUseOnBlockHandler {
		if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
			hdr.HandleItemUseOnBlock(ctx, pos, face, clickPos)
		}
	}
}
func (h *MultipleHandler) HandleItemUseOnEntity(ctx *event.Context[*player.Player], e world.Entity) {
	for _, hdr := range h._ItemUseOnEntityHandler {
		if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
			hdr.HandleItemUseOnEntity(ctx, e)
		}
	}
}
func (h *MultipleHandler) HandleItemRelease(ctx *event.Context[*player.Player], item item.Stack, dur time.Duration) {
	for _, hdr := range h._ItemReleaseHandler {
		if hdr, ok := hdr.(ItemReleaseHandler); ok {
			hdr.HandleItemRelease(ctx, item, dur)
		}
	}
}
func (h *MultipleHandler) HandleItemConsume(ctx *event.Context[*player.Player], item item.Stack) {
	for _, hdr := range h._ItemConsumeHandler {
		if hdr, ok := hdr.(ItemConsumeHandler); ok {
			hdr.HandleItemConsume(ctx, item)
		}
	}
}
func (h *MultipleHandler) HandleAttackEntity(ctx *event.Context[*player.Player], e world.Entity, force, height *float64, critical *bool) {
	for _, hdr := range h._AttackEntityHandler {
		if hdr, ok := hdr.(AttackEntityHandler); ok {
			hdr.HandleAttackEntity(ctx, e, force, height, critical)
		}
	}
}
func (h *MultipleHandler) HandleExperienceGain(ctx *event.Context[*player.Player], amount *int) {
	for _, hdr := range h._ExperienceGainHandler {
		if hdr, ok := hdr.(ExperienceGainHandler); ok {
			hdr.HandleExperienceGain(ctx, amount)
		}
	}
}
func (h *MultipleHandler) HandlePunchAir(ctx *event.Context[*player.Player]) {
	for _, hdr := range h._PunchAirHandler {
		if hdr, ok := hdr.(PunchAirHandler); ok {
			hdr.HandlePunchAir(ctx)
		}
	}
}
func (h *MultipleHandler) HandleSignEdit(ctx *event.Context[*player.Player], pos cube.Pos, frontSide bool, oldText, newText string) {
	for _, hdr := range h._SignEditHandler {
		if hdr, ok := hdr.(SignEditHandler); ok {
			hdr.HandleSignEdit(ctx, pos, frontSide, oldText, newText)
		}
	}
}
func (h *MultipleHandler) HandleLecternPageTurn(ctx *event.Context[*player.Player], pos cube.Pos, oldPage int, newPage *int) {
	for _, hdr := range h._LecternPageTurnHandler {
		if hdr, ok := hdr.(LecternPageTurnHandler); ok {
			hdr.HandleLecternPageTurn(ctx, pos, oldPage, newPage)
		}
	}
}
func (h *MultipleHandler) HandleItemDamage(ctx *event.Context[*player.Player], i item.Stack, damage int) {
	for _, hdr := range h._ItemDamageHandler {
		if hdr, ok := hdr.(ItemDamageHandler); ok {
			hdr.HandleItemDamage(ctx, i, damage)
		}
	}
}
func (h *MultipleHandler) HandleItemPickup(ctx *event.Context[*player.Player], i *item.Stack) {
	for _, hdr := range h._ItemPickupHandler {
		if hdr, ok := hdr.(ItemPickupHandler); ok {
			hdr.HandleItemPickup(ctx, i)
		}
	}
}
func (h *MultipleHandler) HandleHeldSlotChange(ctx *event.Context[*player.Player], from, to int) {
	for _, hdr := range h._HeldSlotChangeHandler {
		if hdr, ok := hdr.(HeldSlotChangeHandler); ok {
			hdr.HandleHeldSlotChange(ctx, from, to)
		}
	}
}
func (h *MultipleHandler) HandleItemDrop(ctx *event.Context[*player.Player], s item.Stack) {
	for _, hdr := range h._ItemDropHandler {
		if hdr, ok := hdr.(ItemDropHandler); ok {
			hdr.HandleItemDrop(ctx, s)
		}
	}
}
func (h *MultipleHandler) HandleTransfer(ctx *event.Context[*player.Player], addr *net.UDPAddr) {
	for _, hdr := range h._TransferHandler {
		if hdr, ok := hdr.(TransferHandler); ok {
			hdr.HandleTransfer(ctx, addr)
		}
	}
}
func (h *MultipleHandler) HandleCommandExecution(ctx *event.Context[*player.Player], command cmd.Command, args []string) {
	for _, hdr := range h._CommandExecutionHandler {
		if hdr, ok := hdr.(CommandExecutionHandler); ok {
			hdr.HandleCommandExecution(ctx, command, args)
		}
	}
}
func (h *MultipleHandler) HandleQuit(p *player.Player) {
	for _, hdr := range h._QuitHandler {
		if hdr, ok := hdr.(QuitHandler); ok {
			hdr.HandleQuit(p)
		}
	}
}
func (h *MultipleHandler) HandleDiagnostics(p *player.Player, d session.Diagnostics) {
	for _, hdr := range h._DiagnosticsHandler {
		if hdr, ok := hdr.(DiagnosticsHandler); ok {
			hdr.HandleDiagnostics(p, d)
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
	_FireExtinguishHandler   []FireExtinguishHandler
	_StartBreakHandler       []StartBreakHandler
	_BlockBreakHandler       []BlockBreakHandler
	_BlockPlaceHandler       []BlockPlaceHandler
	_BlockPickHandler        []BlockPickHandler
	_ItemUseHandler          []ItemUseHandler
	_ItemUseOnBlockHandler   []ItemUseOnBlockHandler
	_ItemUseOnEntityHandler  []ItemUseOnEntityHandler
	_ItemReleaseHandler      []ItemReleaseHandler
	_ItemConsumeHandler      []ItemConsumeHandler
	_AttackEntityHandler     []AttackEntityHandler
	_ExperienceGainHandler   []ExperienceGainHandler
	_PunchAirHandler         []PunchAirHandler
	_SignEditHandler         []SignEditHandler
	_LecternPageTurnHandler  []LecternPageTurnHandler
	_ItemDamageHandler       []ItemDamageHandler
	_ItemPickupHandler       []ItemPickupHandler
	_HeldSlotChangeHandler   []HeldSlotChangeHandler
	_ItemDropHandler         []ItemDropHandler
	_TransferHandler         []TransferHandler
	_CommandExecutionHandler []CommandExecutionHandler
	_QuitHandler             []QuitHandler
	_DiagnosticsHandler      []DiagnosticsHandler
}

func (h *MultipleHandler) Register(hdr any) func() {
	reg := false
	var funcs []func()
	if hdr, ok := hdr.(MoveHandler); ok {
		h._MoveHandler = append(h._MoveHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._MoveHandler = deleteVal(h._MoveHandler, hdr)
		})
	}
	if hdr, ok := hdr.(JumpHandler); ok {
		h._JumpHandler = append(h._JumpHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._JumpHandler = deleteVal(h._JumpHandler, hdr)
		})
	}
	if hdr, ok := hdr.(TeleportHandler); ok {
		h._TeleportHandler = append(h._TeleportHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._TeleportHandler = deleteVal(h._TeleportHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ChangeWorldHandler); ok {
		h._ChangeWorldHandler = append(h._ChangeWorldHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ChangeWorldHandler = deleteVal(h._ChangeWorldHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ToggleSprintHandler); ok {
		h._ToggleSprintHandler = append(h._ToggleSprintHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ToggleSprintHandler = deleteVal(h._ToggleSprintHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ToggleSneakHandler); ok {
		h._ToggleSneakHandler = append(h._ToggleSneakHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ToggleSneakHandler = deleteVal(h._ToggleSneakHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ChatHandler); ok {
		h._ChatHandler = append(h._ChatHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ChatHandler = deleteVal(h._ChatHandler, hdr)
		})
	}
	if hdr, ok := hdr.(FoodLossHandler); ok {
		h._FoodLossHandler = append(h._FoodLossHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._FoodLossHandler = deleteVal(h._FoodLossHandler, hdr)
		})
	}
	if hdr, ok := hdr.(HealHandler); ok {
		h._HealHandler = append(h._HealHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._HealHandler = deleteVal(h._HealHandler, hdr)
		})
	}
	if hdr, ok := hdr.(HurtHandler); ok {
		h._HurtHandler = append(h._HurtHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._HurtHandler = deleteVal(h._HurtHandler, hdr)
		})
	}
	if hdr, ok := hdr.(DeathHandler); ok {
		h._DeathHandler = append(h._DeathHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._DeathHandler = deleteVal(h._DeathHandler, hdr)
		})
	}
	if hdr, ok := hdr.(RespawnHandler); ok {
		h._RespawnHandler = append(h._RespawnHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._RespawnHandler = deleteVal(h._RespawnHandler, hdr)
		})
	}
	if hdr, ok := hdr.(SkinChangeHandler); ok {
		h._SkinChangeHandler = append(h._SkinChangeHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._SkinChangeHandler = deleteVal(h._SkinChangeHandler, hdr)
		})
	}
	if hdr, ok := hdr.(FireExtinguishHandler); ok {
		h._FireExtinguishHandler = append(h._FireExtinguishHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._FireExtinguishHandler = deleteVal(h._FireExtinguishHandler, hdr)
		})
	}
	if hdr, ok := hdr.(StartBreakHandler); ok {
		h._StartBreakHandler = append(h._StartBreakHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._StartBreakHandler = deleteVal(h._StartBreakHandler, hdr)
		})
	}
	if hdr, ok := hdr.(BlockBreakHandler); ok {
		h._BlockBreakHandler = append(h._BlockBreakHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._BlockBreakHandler = deleteVal(h._BlockBreakHandler, hdr)
		})
	}
	if hdr, ok := hdr.(BlockPlaceHandler); ok {
		h._BlockPlaceHandler = append(h._BlockPlaceHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._BlockPlaceHandler = deleteVal(h._BlockPlaceHandler, hdr)
		})
	}
	if hdr, ok := hdr.(BlockPickHandler); ok {
		h._BlockPickHandler = append(h._BlockPickHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._BlockPickHandler = deleteVal(h._BlockPickHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemUseHandler); ok {
		h._ItemUseHandler = append(h._ItemUseHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemUseHandler = deleteVal(h._ItemUseHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemUseOnBlockHandler); ok {
		h._ItemUseOnBlockHandler = append(h._ItemUseOnBlockHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemUseOnBlockHandler = deleteVal(h._ItemUseOnBlockHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemUseOnEntityHandler); ok {
		h._ItemUseOnEntityHandler = append(h._ItemUseOnEntityHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemUseOnEntityHandler = deleteVal(h._ItemUseOnEntityHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemReleaseHandler); ok {
		h._ItemReleaseHandler = append(h._ItemReleaseHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemReleaseHandler = deleteVal(h._ItemReleaseHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemConsumeHandler); ok {
		h._ItemConsumeHandler = append(h._ItemConsumeHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemConsumeHandler = deleteVal(h._ItemConsumeHandler, hdr)
		})
	}
	if hdr, ok := hdr.(AttackEntityHandler); ok {
		h._AttackEntityHandler = append(h._AttackEntityHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._AttackEntityHandler = deleteVal(h._AttackEntityHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ExperienceGainHandler); ok {
		h._ExperienceGainHandler = append(h._ExperienceGainHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ExperienceGainHandler = deleteVal(h._ExperienceGainHandler, hdr)
		})
	}
	if hdr, ok := hdr.(PunchAirHandler); ok {
		h._PunchAirHandler = append(h._PunchAirHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._PunchAirHandler = deleteVal(h._PunchAirHandler, hdr)
		})
	}
	if hdr, ok := hdr.(SignEditHandler); ok {
		h._SignEditHandler = append(h._SignEditHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._SignEditHandler = deleteVal(h._SignEditHandler, hdr)
		})
	}
	if hdr, ok := hdr.(LecternPageTurnHandler); ok {
		h._LecternPageTurnHandler = append(h._LecternPageTurnHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._LecternPageTurnHandler = deleteVal(h._LecternPageTurnHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemDamageHandler); ok {
		h._ItemDamageHandler = append(h._ItemDamageHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemDamageHandler = deleteVal(h._ItemDamageHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemPickupHandler); ok {
		h._ItemPickupHandler = append(h._ItemPickupHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemPickupHandler = deleteVal(h._ItemPickupHandler, hdr)
		})
	}
	if hdr, ok := hdr.(HeldSlotChangeHandler); ok {
		h._HeldSlotChangeHandler = append(h._HeldSlotChangeHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._HeldSlotChangeHandler = deleteVal(h._HeldSlotChangeHandler, hdr)
		})
	}
	if hdr, ok := hdr.(ItemDropHandler); ok {
		h._ItemDropHandler = append(h._ItemDropHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._ItemDropHandler = deleteVal(h._ItemDropHandler, hdr)
		})
	}
	if hdr, ok := hdr.(TransferHandler); ok {
		h._TransferHandler = append(h._TransferHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._TransferHandler = deleteVal(h._TransferHandler, hdr)
		})
	}
	if hdr, ok := hdr.(CommandExecutionHandler); ok {
		h._CommandExecutionHandler = append(h._CommandExecutionHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._CommandExecutionHandler = deleteVal(h._CommandExecutionHandler, hdr)
		})
	}
	if hdr, ok := hdr.(QuitHandler); ok {
		h._QuitHandler = append(h._QuitHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._QuitHandler = deleteVal(h._QuitHandler, hdr)
		})
	}
	if hdr, ok := hdr.(DiagnosticsHandler); ok {
		h._DiagnosticsHandler = append(h._DiagnosticsHandler, hdr)
		reg = true
		funcs = append(funcs, func() {
			h._DiagnosticsHandler = deleteVal(h._DiagnosticsHandler, hdr)
		})
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
	h._FireExtinguishHandler = nil
	h._StartBreakHandler = nil
	h._BlockBreakHandler = nil
	h._BlockPlaceHandler = nil
	h._BlockPickHandler = nil
	h._ItemUseHandler = nil
	h._ItemUseOnBlockHandler = nil
	h._ItemUseOnEntityHandler = nil
	h._ItemReleaseHandler = nil
	h._ItemConsumeHandler = nil
	h._AttackEntityHandler = nil
	h._ExperienceGainHandler = nil
	h._PunchAirHandler = nil
	h._SignEditHandler = nil
	h._LecternPageTurnHandler = nil
	h._ItemDamageHandler = nil
	h._ItemPickupHandler = nil
	h._HeldSlotChangeHandler = nil
	h._ItemDropHandler = nil
	h._TransferHandler = nil
	h._CommandExecutionHandler = nil
	h._QuitHandler = nil
	h._DiagnosticsHandler = nil
}
