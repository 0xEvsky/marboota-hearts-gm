extends Node2D
class_name Card

enum Suit {SPADES, HEARTS, CLUBS, DIAMONDS}

var is_face_up = false
var suit: Suit
var value: int
var is_playable = false
var hover_index: int
var is_played = false

func set_card(cardStr: String):
	if cardStr == "":
		return
	name = cardStr
	var params = cardStr.split("_")
	match params[0]:
		"S":
			suit = Suit.SPADES
		"H":
			suit = Suit.HEARTS
		"C":
			suit = Suit.CLUBS
		"D":
			suit = Suit.DIAMONDS

	value = int(params[1])

	# var image = Image.load_from_file(CardMap.card_dict[cardStr])
	# var texture = ImageTexture.create_from_image(image)
	var texture = load(CardMap.card_dict[cardStr])
	$"Sprite2D".texture = texture
	is_face_up = true
	$"Sprite2DBack".hide()

func set_shroud(shroud: bool):
	$"Panel".visible = shroud

func set_playable(playable: bool):
	set_shroud(!playable)
	is_playable = playable

func _on_area_2d_mouse_entered() -> void:
	if !is_played && $"..".is_mine:
		$"..".card_hovered(self)

func _on_area_2d_mouse_exited() -> void:
	if !is_played && $"..".is_mine:
		$"..".card_unhovered(self)

func hover(b: bool):
	if !is_played:
		if b:
			if is_playable:
				position.y -= 25
			else:
				position.y -= 5
		else:
			position.y = 0

func _unhandled_input(event: InputEvent) -> void:
	if event.is_action_released("click"):
		if $"Sprite2D".is_pixel_opaque(get_local_mouse_position()):
			if is_playable && !is_played:
				$"..".play(self)
			get_viewport().set_input_as_handled()
