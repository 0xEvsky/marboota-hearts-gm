extends Node

var card_dict = {
    "S_14": "res://assets/graphics/cards/cardSpadesA.png",
    "S_13": "res://assets/graphics/cards/cardSpadesK.png",
    "S_12": "res://assets/graphics/cards/cardSpadesQ.png",
    "S_11": "res://assets/graphics/cards/cardSpadesJ.png",
    "S_10": "res://assets/graphics/cards/cardSpades10.png",
    "S_9": "res://assets/graphics/cards/cardSpades9.png",
    "S_8": "res://assets/graphics/cards/cardSpades8.png",
    "S_7": "res://assets/graphics/cards/cardSpades7.png",
    "S_6": "res://assets/graphics/cards/cardSpades6.png",
    "S_5": "res://assets/graphics/cards/cardSpades5.png",
    "S_4": "res://assets/graphics/cards/cardSpades4.png",
    "S_3": "res://assets/graphics/cards/cardSpades3.png",
    "S_2": "res://assets/graphics/cards/cardSpades2.png",

    "H_14": "res://assets/graphics/cards/cardHeartsA.png",
    "H_13": "res://assets/graphics/cards/cardHeartsK.png",
    "H_12": "res://assets/graphics/cards/cardHeartsQ.png",
    "H_11": "res://assets/graphics/cards/cardHeartsJ.png",
    "H_10": "res://assets/graphics/cards/cardHearts10.png",
    "H_9": "res://assets/graphics/cards/cardHearts9.png",
    "H_8": "res://assets/graphics/cards/cardHearts8.png",
    "H_7": "res://assets/graphics/cards/cardHearts7.png",
    "H_6": "res://assets/graphics/cards/cardHearts6.png",
    "H_5": "res://assets/graphics/cards/cardHearts5.png",
    "H_4": "res://assets/graphics/cards/cardHearts4.png",
    "H_3": "res://assets/graphics/cards/cardHearts3.png",
    "H_2": "res://assets/graphics/cards/cardHearts2.png",

    "C_14": "res://assets/graphics/cards/cardClubsA.png",
    "C_13": "res://assets/graphics/cards/cardClubsK.png",
    "C_12": "res://assets/graphics/cards/cardClubsQ.png",
    "C_11": "res://assets/graphics/cards/cardClubsJ.png",
    "C_10": "res://assets/graphics/cards/cardClubs10.png",
    "C_9": "res://assets/graphics/cards/cardClubs9.png",
    "C_8": "res://assets/graphics/cards/cardClubs8.png",
    "C_7": "res://assets/graphics/cards/cardClubs7.png",
    "C_6": "res://assets/graphics/cards/cardClubs6.png",
    "C_5": "res://assets/graphics/cards/cardClubs5.png",
    "C_4": "res://assets/graphics/cards/cardClubs4.png",
    "C_3": "res://assets/graphics/cards/cardClubs3.png",
    "C_2": "res://assets/graphics/cards/cardClubs2.png",

    "D_14": "res://assets/graphics/cards/cardDiamondsA.png",
    "D_13": "res://assets/graphics/cards/cardDiamondsK.png",
    "D_12": "res://assets/graphics/cards/cardDiamondsQ.png",
    "D_11": "res://assets/graphics/cards/cardDiamondsJ.png",
    "D_10": "res://assets/graphics/cards/cardDiamonds10.png",
    "D_9": "res://assets/graphics/cards/cardDiamonds9.png",
    "D_8": "res://assets/graphics/cards/cardDiamonds8.png",
    "D_7": "res://assets/graphics/cards/cardDiamonds7.png",
    "D_6": "res://assets/graphics/cards/cardDiamonds6.png",
    "D_5": "res://assets/graphics/cards/cardDiamonds5.png",
    "D_4": "res://assets/graphics/cards/cardDiamonds4.png",
    "D_3": "res://assets/graphics/cards/cardDiamonds3.png",
    "D_2": "res://assets/graphics/cards/cardDiamonds2.png"
}


func _ready():
    for k in card_dict:
        load(card_dict[k])