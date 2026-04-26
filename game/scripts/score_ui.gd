extends CanvasLayer


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	hide()
	EventManager.GAMESTART_received.connect(_on_gamestart)
	EventManager.TEAMROUNDEND_received.connect(_on_roundend)
	EventManager.TEAMTOTALSCORE_received.connect(_on_totalscore)

func _on_gamestart() -> void:
	if Globals.table.state <= Globals.table.TableState.TABLE_READY:
		$"HeaderRow/Label".text = $"../Table/Seat0".sitter.username + "\n" + $"../Table/Seat2".sitter.username
		$"HeaderRow/Label2".text = $"../Table/Seat1".sitter.username + "\n" + $"../Table/Seat3".sitter.username
		$"TotalScoreRow/Label".text = "0"
		$"TotalScoreRow/Label2".text = "0"
		for row in get_node("ScoreContainer").get_children():
			row.queue_free()
		show()

func _on_roundend(team_a_score: String, team_b_score: String) -> void:
	var score_row_scene := preload("res://scenes/score_row.tscn")
	var score_row := score_row_scene.instantiate()
	get_node("ScoreContainer").add_child(score_row)
	score_row.get_node("Label").text = team_a_score
	score_row.get_node("Label2").text = team_b_score
	

func _on_totalscore(team_a_score: String, team_b_score: String) -> void:
	get_node("TotalScoreRow/Label").text = team_a_score
	get_node("TotalScoreRow/Label2").text = team_b_score
