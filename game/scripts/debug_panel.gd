extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(_delta: float) -> void:
	var label: RichTextLabel = $"Panel/RichTextLabel"
	label.text = str(
		"[color=red]instance id:[/color] ", NetworkManager.instance_id, "\n",
		"[color=red]userrname:[/color] ", NetworkManager.username, "\n",
		"[color=red]user id:[/color] ", NetworkManager.user_id, "\n",
		"[color=red]icon url:[/color] ", NetworkManager.icon_url, "\n",
		"[color=red]authenticated:[/color] ", NetworkManager.authenticated
	)


func _on_button_button_up() -> void:
	EventManager.send_request(
		EventManager.unsit_request(),
		func():
			print_debug("sit 0 success"),
		func():
			print_debug("sit 0 fail"),
	)


func _on_button_2_button_up() -> void:
	EventManager.send_request(
		EventManager.sit_request(1),
		func():
			print_debug("sit 1 success"),
		func():
			print_debug("sit 1 fail"),
	)


func _on_button_3_button_up() -> void:
	EventManager.send_request(
		EventManager.sit_request(2),
		func():
			print_debug("sit 2 success"),
		func():
			print_debug("sit 2 fail"),
	)
