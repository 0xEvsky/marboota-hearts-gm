extends Button

func _ready() -> void:
	visible = false

func _on_button_up() -> void:
	Globals.my_player.state = Globals.player_manager.PLAYER_SELECTING
	EventManager.send_request(EventManager.setmode_request("WHIST")
	# on error
	,func (error: String) -> void: print_debug(error)
	)
