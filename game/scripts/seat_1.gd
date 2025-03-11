extends Node2D
class_name Seat

@export var seat_num = 1

var sitter: Player
var is_taken = false

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.

func disable_button() -> void:
	is_taken = true
	$Button.disabled = true

func enable_button() -> void:
	is_taken = false
	$Button.disabled = false

func seat_player(id: String) -> void:
	var player_manager = Globals.player_manager
	player_manager.move_player(id, global_position)
	var player = player_manager.get_node(id) as Player

	if player.seat != null:
		var old_seat = get_node("../Seat" + str(player.seat)) as Seat
		old_seat.unseat_player()

	player.state = player_manager.PLAYER_READY # TODO: Change depending on game state
	player.seat = self

	sitter = player
	disable_button()

func unseat_player() -> void:
	# TODO: Move player back to player list
	sitter = null
	enable_button()

func _on_button_button_up() -> void:
	EventManager.send_request(
		EventManager.sit_request(seat_num)
	,func():
		seat_player("Me")
	,func(err: String):
		push_error(err)
	)
