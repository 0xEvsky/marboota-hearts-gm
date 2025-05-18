extends Button


func _ready() -> void:
	visible = false
	EventManager.GAMEEND_received.connect(_on_gameend)


func _on_button_up() -> void:
	var player = Globals.my_player
	var manager = Globals.player_manager
	var seat: Seat = player.seat
	var seat_ready_button: Button = $"../ReadyButton"

	if player.state == manager.PLAYER_READY or player.state == manager.PLAYER_WAITING:
		seat.unseat_player()
		self.hide()
		seat_ready_button.hide()
		EventManager.send_request(EventManager.unsit_request()
		,func(error):
			print_debug(error)
			seat.seat_player(player.name)
		)

func _on_gameend(_1, _2):
	button_pressed = false