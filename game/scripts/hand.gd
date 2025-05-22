extends Node2D
class_name Hand

const GAP = 50
var player: Player
var cards: Array[Card] = []
var is_mine := false
var hovered_cards: Array[Card] = []
@export var anchor: Node2D
@export var score: Panel
var playable: String = ""


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	EventManager.DEAL_received.connect(_on_deal)

func _on_deal(cardStr: String) -> void:
	if player == Globals.my_player:
		var cardStrArr := cardStr.split(",")
		var i := 0
		for card in cardStrArr:
			add_card(card, i)
			i += 1
	else:
		for i in range(13):
			add_card("", -1)
	rearrange()

func add_card(cardStr: String, index: int) -> void:
	var cardScene := preload("res://scenes/card.tscn")
	var card := cardScene.instantiate() as Card
	cards.append(card)
	add_child(card)
	card.set_card(cardStr)
	card.hover_index = index
	if index == -1:
		card.set_shroud(true)
	card.global_position = global_position

func sort_cards() -> void:
	cards.sort_custom(func(a: Card, b: Card) -> bool:
		if a.hover_index > b.hover_index:
			return true
		return false
	)

func rearrange() -> void:
	var half := (len(cards) - 1) / 2.0
	for i in range(len(cards)):
		cards[i].position.x = (i - half) * GAP

func set_playable(cardStr: String) -> void:
	var card := get_node(cardStr) as Card
	card.set_playable(true)

func card_hovered(card: Card) -> void:
	hovered_cards.append(card)
	rehover()

func card_unhovered(card: Card) -> void:
	card.hover(false)
	hovered_cards.erase(card)
	rehover()

func rehover() -> void:
	var id := -1
	var hover_id := -1
	for i in range(len(hovered_cards)):
		hovered_cards[i].hover(false)
		if hovered_cards[i].hover_index > hover_id:
			id = i
			hover_id = hovered_cards[i].hover_index
	if id > -1:
		hovered_cards[id].hover(true)
	
func play(card: Card) -> void:
	card.reparent(anchor)
	cards.erase(card)
	hovered_cards.erase(card)
	card.global_position = anchor.global_position
	card.position = Vector2.ZERO
	card.scale = Vector2(0.75, 0.75)
	card.is_played = true
	rearrange()
	for c in cards:
		c.set_playable(false)
	EventManager.send_request(
		EventManager.play_request(card.name),
		func(error: String) -> void:
			print_debug(error)
			var c := anchor.get_child(0) as Card
			add_card(c.name, c.hover_index)
			c.queue_free()
			sort_cards()
			rearrange()
			get_parent()._on_yourplay(playable)
)
	if !get_parent().play_started:
		get_parent().play_started = true
		for i in range(4):
			var score_node := get_parent().get_node("Score" + str(i))
			score_node.get_node("Label").text = "0"
			score_node.show()

func on_play(card_str: String) -> void:
	var card := get_child(-1) as Card
	card.set_card(card_str)
	card.reparent(anchor)
	cards.erase(card)
	card.global_position = anchor.global_position
	card.position = Vector2.ZERO
	card.scale = Vector2(0.75, 0.75)
	card.set_shroud(false)
	card.is_played = true
	rearrange()
	if !get_parent().play_started:
		get_parent().play_started = true
		for i in range(4):
			var score_node := get_parent().get_node("Score" + str(i))
			score_node.get_node("Label").text = "0"
			score_node.show()