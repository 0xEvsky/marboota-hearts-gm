extends Node2D
class_name Hand

const GAP = 50
var player: Player
var cards: Array[Card] = []
var is_mine = false


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	EventManager.DEAL_received.connect(_on_deal)

func _on_deal(cardStr: String):
	if player == Globals.my_player:
		var cardStrArr = cardStr.split(",")
		for card in cardStrArr:
			add_card(card)
	else:
		for i in range(13):
			add_card("")
	rearrange()

func add_card(cardStr: String):
	var cardScene = preload("res://scenes/card.tscn")
	var card = cardScene.instantiate()
	cards.append(card)
	add_child(card)
	card.set_card(cardStr)
	card.global_position = global_position

func rearrange():
	var half = (len(cards) - 1) / 2.0
	for i in range(len(cards)):
		cards[i].position.x = (i - half) * GAP
