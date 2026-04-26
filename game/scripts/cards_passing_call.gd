extends Panel
var passButton: Button
var cards_buttons: Array[Button] = []
var cards_names: Array[String] = []
var chosed_cards_count := 0
var chosen_cards: Array[String] = []

func _ready() -> void:
	EventManager.PASSCARDS_recevied.connect(_on_passcards)
	EventManager.FFAROUNDEND_received.connect(_on_ffaRoundEnd)
	passButton = $"VBoxContainer/PassButton"
	passButton.disabled = true
	for button in $VBoxContainer/HBoxContainer.get_children():
		cards_buttons.append(button as Button)
		button.toggle_mode = true
		button.connect("toggled", _on_card_button_toggled.bind(button))
		
func _on_ffaRoundEnd() -> void:
	pass
	
func _on_passcards() -> void:
	for i in range(Globals.my_player.hand.cards.size()):
		var card: Card = Globals.my_player.hand.cards[i]
		var card_texture: Texture2D = card.card_texture
		var card_button_texture_rect: TextureRect = cards_buttons[i].get_node("CardTexture")
		card_button_texture_rect.texture = card_texture
		cards_names.append(card.name)
	show()

func _on_card_button_toggled(pressed: bool, button: Button) -> void:
	if pressed:
		button.modulate = Color(0.5, 0.5, 0.5, 1)
		chosed_cards_count += 1
	else:
		button.modulate = Color(1, 1, 1, 1)
		chosed_cards_count -= 1
	
	print(chosed_cards_count)
	if chosed_cards_count >= 4:
		passButton.disabled = false
		for card_button in cards_buttons:
			if not card_button.button_pressed:
				card_button.modulate = Color(0.3, 0.3, 0.3, 1)
				card_button.disabled = true
	else:
		passButton.disabled = true
		for card_button in cards_buttons:
			if not card_button.button_pressed:
				card_button.modulate = Color(1, 1, 1, 1)
				card_button.disabled = false
				
func _on_pass_button_button_up() -> void:
	chosen_cards = []
	for i in range(cards_buttons.size()):
		if not cards_buttons[i].disabled:
			chosen_cards.append(cards_names[i])
	var msg := EventManager.passcards_request(", ".join(chosen_cards))
	EventManager.send_request(msg, func on_error(error: String) -> void: print_debug(error))
	cards_names = []
	chosen_cards = []
	hide()
	for card_button in cards_buttons:
		card_button.button_pressed = false
		card_button.modulate = Color(1, 1, 1, 1)
		card_button.disabled = false
	
