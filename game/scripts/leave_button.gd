extends Button


func _ready() -> void:
	visible = false


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

		,func():
			pass
		,func(error):
			print_debug(error)
			seat.seat_player(player.name)
		)
