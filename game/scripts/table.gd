extends Node2D
class_name Table

func _ready() -> void:
	Globals.table = self
	EventManager.GAMESTART_received.connect(_on_gamestart)
	EventManager.TRUMPSTART_received.connect(_on_trumpstart)

func _on_gamestart():
	rotate_table()

	var leaveButton = $"LeaveButton"
	leaveButton.hide()

	var readyButton = $"ReadyButton"
	readyButton.hide()

	Globals.my_player.hide()
	Globals.my_player.seat.hide()


func _on_trumpstart():
	pass	

func rotate_table() -> void:
	var _offset = 4 - Globals.my_player.seat.seat_num

	for i in range(4):
		var next_anchor_str = "anchor" + str((i + _offset) % 4)
		var next_anchor = get_node(next_anchor_str) as Node2D

		var current_seat_str = "Seat" + str(i)
		var current_seat = get_node(current_seat_str) as Seat

		current_seat.global_position = next_anchor.global_position
		
		if current_seat.sitter:
			current_seat.sitter.global_position = current_seat.global_position
			var hand_str = "Hand" + str((i + _offset) % 4)
			var hand = get_node(hand_str) as Hand
			current_seat.sitter.hand = hand
			hand.player = current_seat.sitter

func unRotate_table():
	for i in range(4):
		var current_seat_str = "Seat" + str(i)
		var current_seat = get_node(current_seat_str) as Seat

		var anchor_str = "anchor" + str(i)
		var anchor = get_node(anchor_str) as Node2D

		current_seat.global_position = anchor.global_position

		if current_seat.sitter:
			current_seat.sitter.global_position = current_seat.global_position
