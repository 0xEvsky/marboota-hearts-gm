extends CanvasLayer


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	hide()
	EventManager.GAMESTART_received.connect(_on_gamestart)
	EventManager.FFAROUNDEND_received.connect(_on_roundend)
	EventManager.FFATOTALSCORE_received.connect(_on_totalscore)
	
func _on_gamestart() -> void:
	pass
	
func _on_roundend() -> void:
	pass
	
func _on_totalscore() -> void:
	pass
