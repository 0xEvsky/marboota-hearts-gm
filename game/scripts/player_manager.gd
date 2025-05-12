extends Node2D
class_name PlayerManager

enum {PLAYER_UNAVAILABLE, PLAYER_IDLE, PLAYER_WAITING, PLAYER_READY, PLAYER_TRUMPING, PLAYER_PLAYING}
var pinned_players: Array[Player] = []

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	Globals.player_manager = self
	NetworkManager.AUTH_accepted.connect(_init_my_player)
	EventManager.JOIN_received.connect(_on_player_join)
	EventManager.LEAVE_received.connect(_on_player_leave)
	EventManager.SIT_received.connect(_on_player_sit)
	EventManager.UNSIT_received.connect(_on_player_unsit)


func _init_my_player() -> void:
	_on_player_join("Me", NetworkManager.username, NetworkManager.icon_url)


func _on_player_join(id: String, username: String, url: String) -> void:
	var player_scene = preload("res://scenes/player.tscn")
	var new_player = player_scene.instantiate() as Player

	new_player.name = id
	new_player.username = username

	# Get icon from URL
	var http_request = HTTPRequest.new()
	add_child(http_request)
	http_request.request_completed.connect(func(result, _response_code, _headers, body):
		if result != HTTPRequest.RESULT_SUCCESS:
			push_error("Image couldn't be downloaded. Try a different image.")

		var image = Image.new()
		var err = image.load_png_from_buffer(body)
		if err != OK:
			push_error("Couldn't load the image.")

		var texture = ImageTexture.create_from_image(image)
		new_player.icon.texture = texture
	)

	var error = http_request.request(url)
	if error != OK:
		push_error("An error occurred in the HTTP request.")
	
	add_child(new_player)
	pin_player(new_player)
	if id == "Me":
		Globals.my_player = new_player


func _on_player_leave(id: String) -> void:
	var leaving_player = get_node(id) as Player

	if leaving_player.seat != null:
		leaving_player.seat.unseat_player()
	
	unpin_player(leaving_player)

	leaving_player.queue_free()


func _on_player_sit(id: String, seat_num: String) -> void:
	var seat_str = "../Table/Seat" + seat_num
	var seat = get_node(seat_str) as Seat
	seat.seat_player(id)


func _on_player_unsit(id: String) -> void:
	var player = get_node(id) as Player
	player.seat.unseat_player()


func move_player(id: String, pos: Vector2) -> void:
	var player = get_node(id) as Player
	player.global_position = pos


func _update_pinned_players() -> void:
	for i in pinned_players.size():
		var player = pinned_players[i]
		move_player(player.name, Vector2(global_position.x, global_position.y + (i * 120)))
	pass


func pin_player(player: Player) -> void:
	if pinned_players.has(player):
		return
	pinned_players.append(player)
	_update_pinned_players()


func unpin_player(player: Player) -> void:
	if !pinned_players.has(player):
		return
	pinned_players.erase(player)
	_update_pinned_players()
