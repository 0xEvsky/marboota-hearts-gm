extends CanvasLayer


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	EventManager.GAMESTART_received_received.connect(_on_gamestart)
	EventManager.ROUNDEND_received.connect(_on_roundend)
	EventManager.TOTALSCORE_received.connect(_on_totalscore)

func _on_gamestart():
	for row in get_node("ScoreContainer").get_children():
		row.queue_free()
	show()

func _on_roundend(team_a_score: String, team_b_score: String):
	var score_row_scene = preload("res://scenes/score_row.tscn")
	var score_row = score_row_scene.instantiate()
	get_node("ScoreContainer").add_child(score_row)
	score_row.get_node("Label").text = team_a_score
	score_row.get_node("Label2").text = team_b_score
	

func _on_totalscore(team_a_score: String, team_b_score: String):
	get_node("TotalScoreRow/Label").text = team_a_score
	get_node("TotalScoreRow/Label2").text = team_b_score
