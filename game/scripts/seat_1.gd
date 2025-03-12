extends Node2D
class_name Seat

@export var seat_num = 1

var sitter: Player
var is_taken = false

func _disable_button() -> void:
	is_taken = true
	$Button.disabled = true

func _enable_button() -> void:
	is_taken = false
	$Button.disabled = false

func seat_player(id: String) -> void:
	var player_manager = Globals.player_manager
	player_manager.move_player(id, global_position)
	var player = player_manager.get_node(id) as Player

	if player.seat != null:
		var old_seat = player.seat as Seat
		old_seat.unseat_player()

	player.state = player_manager.PLAYER_WAITING # TODO: Change depending on game state
	player.seat = self

	sitter = player
	_disable_button()

func unseat_player() -> void:
	# TODO: Move player back to player list
	if sitter != null:
		sitter.seat = null
	sitter = null
	_enable_button()

func _on_button_button_up() -> void:
	EventManager.send_request(
		EventManager.sit_request(seat_num)
	,func():
		seat_player("Me")
	,func(err: String):
		push_error(err)
	)
