extends Button

func _ready() -> void:
	visible = false

func _on_pressed() -> void:
	Globals.my_player.state = Globals.player_manager.PLAYER_SELECTING
	EventManager.send_request(EventManager.setmode_request("HEARTS")
	# on error
	,func (error: String) -> void: print_debug(error)
	)
