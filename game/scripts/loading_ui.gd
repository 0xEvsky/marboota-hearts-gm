extends CanvasLayer


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	show()
	NetworkManager.AUTH_accepted.connect(func() -> void:
		hide()
	)


