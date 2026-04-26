extends CanvasLayer

func _ready() -> void:
	hide()
	EventManager.MODE_received.connect(_on_modeSet)
	EventManager.FFAROUNDEND_received.connect(_on_roundend)
	EventManager.FFATOTALSCORE_received.connect(_on_totalscore)

func _on_modeSet(mode: String) -> void:
	if mode == "HEARTS":
		setup()

func setup() -> void:
	$"HeaderRow/Label".text = $"../Table/Seat0".sitter.username
	$"HeaderRow/Label2".text = $"../Table/Seat1".sitter.username
	$"HeaderRow2/Label".text = $"../Table/Seat2".sitter.username
	$"HeaderRow2/Label2".text = $"../Table/Seat3".sitter.username
	$"TotalScoreRow/Label".text = "0"
	$"TotalScoreRow/Label2".text = "0"
	$"TotalScoreRow2/Label".text = "0"
	$"TotalScoreRow2/Label2".text = "0"
	for row in get_node("ScoreContainer").get_children():
		row.queue_free()
	for row in get_node("ScoreContainer2").get_children():
		row.queue_free()
	show()

func _on_roundend(s0: String, s1: String, s2: String, s3: String) -> void:
	var score_row_scene := preload("res://scenes/score_row.tscn")
	
	var row1 := score_row_scene.instantiate()
	get_node("ScoreContainer").add_child(row1)
	row1.get_node("Label").text = s0
	row1.get_node("Label2").text = s1
	
	var row2 := score_row_scene.instantiate()
	get_node("ScoreContainer2").add_child(row2)
	row2.get_node("Label").text = s2
	row2.get_node("Label2").text = s3

func _on_totalscore(s0: String, s1: String, s2: String, s3: String) -> void:
	get_node("TotalScoreRow/Label").text = s0
	get_node("TotalScoreRow/Label2").text = s1
	get_node("TotalScoreRow2/Label").text = s2
	get_node("TotalScoreRow2/Label2").text = s3
