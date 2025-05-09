extends Node2D
class_name Table

func _ready() -> void:
	Globals.table = self
	EventManager.TRUMPSTART_recevied.connect(_on_trumpstart)


func rotate_table() -> void:
	var _offset = 4 - Globals.my_player.seat.seat_num

	for i in range(4):
		var next_anchor_str = "anchor" + str((i + _offset) % 4)
		var nextAnchor = get_node(next_anchor_str) as Node2D

		var current_seat_str = "Seat" + str(i)
		var currentSeat = get_node(current_seat_str) as Seat

		print(current_seat_str, " ",  next_anchor_str)
		currentSeat.global_position = nextAnchor.global_position
		
		if currentSeat.sitter:
			currentSeat.sitter.global_position = currentSeat.global_position



func _on_trumpstart():
	rotate_table()

	var leaveButton = $"LeaveButton"
	leaveButton.hide()

	var readyButton = $"ReadyButton"
	readyButton.hide()

func unRotate_table():

	for i in range(4):
		var current_seat_str = "Seat" + str(i)
		var currentSeat = get_node(current_seat_str) as Seat

		var anchor_str = "anchor" + str(i)
		var anchor = get_node(anchor_str) as Node2D

		currentSeat.global_position = anchor.global_position

		if currentSeat.sitter:
			currentSeat.sitter.global_position = currentSeat.global_position
