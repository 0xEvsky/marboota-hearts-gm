extends Node2D
class_name Table

enum TableState {TABLE_IDLE, TABLE_READY, TABLE_TRUMPING, TABLE_PLAYING}
var state: TableState = TableState.TABLE_IDLE

func _ready() -> void:
	Globals.table = self
	EventManager.GAMESTART_received.connect(_on_gamestart)
	EventManager.TRUMPSTART_received.connect(_on_trumpstart)
	EventManager.PLAYSTART_received.connect(_on_playstart)
	EventManager.YOURPLAY_received.connect(_on_yourplay)
	EventManager.PLAY_received.connect(_on_play)
	EventManager.PLAYEND_received.connect(_on_playend)

func _on_gamestart():
	# ! reset at GAMEEND
	state = TableState.TABLE_READY
	if Globals.my_player.state > Globals.player_manager.PLAYER_IDLE:
		rotate_table()

		var leaveButton = $"LeaveButton"
		leaveButton.hide()

		var readyButton = $"ReadyButton"
		readyButton.hide()

		Globals.my_player.hide()
		Globals.my_player.seat.hide()
		Globals.my_player.hand.scale = Vector2(1, 1)
		Globals.my_player.hand.position.y = 260
	else:
		for i in range(4):
			var hand = get_node("Hand" + str(i)) as Hand
			hand._on_deal("")


func _on_trumpstart():
	# ! reset at ROUNDEND
	state = TableState.TABLE_TRUMPING
	for i in range(4):
		var seat = get_node("Seat" + str(i)) as Seat
		if seat.sitter:
			seat.sitter.state = Globals.player_manager.PLAYER_TRUMPING

func _on_playstart():
	# ! reset at GAMEEND?
	state = TableState.TABLE_PLAYING
	for i in range(4):
		var seat = get_node("Seat" + str(i)) as Seat
		if seat.sitter:
			seat.sitter.state = Globals.player_manager.PLAYER_PLAYING

func _on_yourplay(playable: String):
	Globals.my_player.hand.playable = playable
	for cardStr in playable.split(","):
		Globals.my_player.hand.set_playable(cardStr)

func _on_play(user_id: String, card_str: String):
	var player = Globals.player_manager.get_node(user_id) as Player
	var hand = player.hand
	hand.on_play(card_str)

func _on_playend(winner_id: String):
	#var cards: Array[Card] = []
	var winning_hand = Globals.player_manager.get_player_by_id(winner_id).hand as Hand
	for i in range(4):
		var card = get_node("CardAnchor" + str(i)).get_child(0)
		var tween = card.create_tween().set_trans(Tween.TRANS_SINE).set_ease(Tween.EASE_OUT)
		tween.tween_property(card, "global_position", winning_hand.global_position, 0.25)
		tween.parallel().tween_property(card, "rotation", 0, 0.5)
		tween.tween_callback(func():
			# TODO: Use this for score counter display
			card.queue_free()
		)
	
	

func rotate_table() -> void:
	var _offset = 4 - Globals.my_player.seat.seat_num

	for i in range(4):
		var next_anchor_str = "Anchor" + str((i + _offset) % 4)
		var next_anchor = get_node(next_anchor_str) as Node2D

		var current_seat_str = "Seat" + str(i)
		var current_seat = get_node(current_seat_str) as Seat

		current_seat.global_position = next_anchor.global_position
		
		if current_seat.sitter:
			current_seat.sitter.global_position = current_seat.global_position
			var hand_str = "Hand" + str((i + _offset) % 4)
			var hand = get_node(hand_str) as Hand
			current_seat.sitter.hand = hand
			hand.player = current_seat.sitter
			if current_seat.sitter == Globals.my_player:
				hand.is_mine = true
				Globals.my_player.hand = hand
				hand.player = Globals.my_player

func unRotate_table():
	for i in range(4):
		var current_seat_str = "Seat" + str(i)
		var current_seat = get_node(current_seat_str) as Seat

		var anchor_str = "Anchor" + str(i)
		var anchor = get_node(anchor_str) as Node2D

		current_seat.global_position = anchor.global_position

		if current_seat.sitter:
			current_seat.sitter.global_position = current_seat.global_position
