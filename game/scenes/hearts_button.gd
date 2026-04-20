extends Button

func _ready() -> void:
	visible = false

func _on_pressed() -> void:
	EventManager.send_request(EventManager.setmode_request("HEARTS")
	# on error
	,func (error: String) -> void: print_debug(error)
	)
