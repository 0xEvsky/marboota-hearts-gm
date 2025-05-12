extends Node2D
class_name Card

enum Suit {SPADES, HEARTS, CLUBS, DIAMONDS}

var is_face_up = false
var suit: Suit
var value: int

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
