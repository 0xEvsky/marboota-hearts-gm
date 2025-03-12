extends Node2D
class_name Seat


@export var seat_num = 1
var sitter: Player
@onready var seat_ready_button : Button = $"../ReadyButton"
var is_taken = false

func _disable_button() -> void:
	is_taken = true
	$Button.disabled = true

func _enable_button() -> void:
	is_taken = false
	$Button.disabled = false

func seat_player(id: String) -> void:
	var player_manager = Globals.player_manager
	var player = player_manager.get_node(id) as Player

	if player.seat != null:
		var old_seat = player.seat as Seat
		old_seat.unseat_player()

	player_manager.unpin_player(player)
	player_manager.move_player(id, global_position)

	player.state = player_manager.PLAYER_WAITING # TODO: Change depending on game state
	player.seat = self

	sitter = player
	_disable_button()

func unseat_player() -> void:
	if sitter != null:
		sitter.unseat()
	
	sitter = null
	_enable_button()

func _on_button_button_up() -> void:
	seat_player("Me")
	seat_ready_button.show()

	EventManager.send_request(
		EventManager.sit_request(seat_num)
	,func():
		pass
	,func(err: String):
		var me = Globals.my_player
		seat_ready_button.hide()
		if sitter == me:
			unseat_player()
		else:
			me.unseat()
		push_error(err)
	)
