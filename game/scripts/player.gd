extends Node2D
class_name Player

var username = "Player"

@onready var manager: PlayerManager = get_parent()
@onready var icon: Sprite2D = $Icon

var state = manager.PLAYER_IDLE
var seat: Seat = null

func unseat() -> void:
    state = manager.PLAYER_IDLE
    seat = null
    # Move player back to player list
    manager.pin_player(self)