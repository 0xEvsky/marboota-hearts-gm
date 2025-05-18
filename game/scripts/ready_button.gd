extends Button

func _ready() -> void:
	visible = false
	EventManager.GAMEEND_received.connect(_on_gameend)


func _on_toggled(toggled_on: bool) -> void:
	if toggled_on:
		Globals.my_player.state = Globals.player_manager.PLAYER_READY
		EventManager.send_request(EventManager.ready_request()
		# on error
		,func(error):
			print_debug(error)
			Globals.my_player.state = Globals.player_manager.PLAYER_WAITING
		)
	else:
		if Globals.my_player.state == Globals.player_manager.PLAYER_READY:
			Globals.my_player.state = Globals.player_manager.PLAYER_WAITING
			EventManager.send_request(EventManager.unready_request()
			,func(error):
				print_debug(error)
				Globals.my_player.state = Globals.player_manager.PLAYER_READY
			)

func _on_gameend(_1, _2):
	button_pressed = false