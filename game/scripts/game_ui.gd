extends Panel

var min_score_str_global: String
var max_score_str_global: String

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	EventManager.YOURTRUMPCALL_received.connect(_on_yourtrumpcall)

func _on_yourtrumpcall(min_score_str: String, max_score_str: String):
	min_score_str_global = min_score_str
	max_score_str_global = max_score_str
	var min_score = int(min_score_str)
	var max_score = int(max_score_str)
	var buttons = get_node("VBoxContainer/HBoxContainer").get_children() as Array[Button]
	for button in buttons:
		button.disabled = true
	for i in range(min_score - 7, max_score + 1 - 7):
		buttons[i].disabled = false
	show()


func _on_trump_button_up(score: String) -> void:
	EventManager.send_request(
		EventManager.trumpcall_request(score),
		func(error):
			print_debug(error)
			_on_yourtrumpcall(min_score_str_global, max_score_str_global)
	)
	hide()
