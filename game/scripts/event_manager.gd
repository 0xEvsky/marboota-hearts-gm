extends Node

# Signals that don't need requests (make sure to subscribe to those)
signal JOIN_received
signal LEAVE_received
signal SIT_received
signal UNSIT_received
signal READY_received
signal UNREADY_received
signal GAMESTART_received
signal DEAL_received
signal TRUMPSTART_received
signal TRUMPCALL_received
signal YOURTRUMPCALL_received
signal PLAYSTART_received
signal YOURPLAY_received
signal PLAY_received
signal PLAYEND_received



var _request_queue: Array[Dictionary] = []
var _response_queue: Array[Dictionary] = []

func _process(_delta: float) -> void:
	#if !_response_queue.is_empty():
		#_process_response_queue()
	pass

func send_request(msg: Dictionary, on_error: Callable) -> void:
	var request_id = _generate_request_id()
	msg["REQUESTID"] = request_id
	_request_queue.append({"message": msg, "request_id": request_id, "on_error": on_error})
	print_debug("request: ", msg)
	NetworkManager._write_json(msg)

func _generate_request_id() -> String:
	return "request" + "-" + NetworkManager.user_id + "-" + str(RandomNumberGenerator.new().randi())

func _handle_message(msg: Dictionary) -> void:
	_response_queue.append(msg)
	_process_response_queue()

func _process_response_queue() -> void:
	var res = _response_queue.pop_front() as Dictionary
	print_debug("response: ", res)
	for req in _request_queue:
		if !res.has("REQUESTID"):
			_dispatch(res["ACTION"], res)
			return
		if res["REQUESTID"] == req["request_id"]:
			if res["ACTION"] == "OK":
				print_debug("OK")
			elif res["ACTION"] == "ERROR":
				req["on_error"].call(res["MESSAGE"])
			_request_queue.pop_front()
			return
	_dispatch(res["ACTION"], res)

func _dispatch(action: String, msg: Dictionary) -> void:
	match action:
		"JOIN":
			JOIN_received.emit(msg["USERID"], msg["USERNAME"], msg["ICONURL"])
		"LEAVE":
			LEAVE_received.emit(msg["USERID"])
		"SIT":
			SIT_received.emit(msg["USERID"], msg["SEAT"])
		"UNSIT":
			UNSIT_received.emit(msg["USERID"])
		"READY":
			READY_received.emit()
		"UNREADY":
			UNREADY_received.emit()
		"GAMESTART":
			GAMESTART_received.emit()
		"DEAL": 
			DEAL_received.emit(msg["CARDS"])
		"TRUMPSTART":
			TRUMPSTART_received.emit()
		"TRUMPCALL":
			TRUMPCALL_received.emit(msg["USERID"], msg["SCORE"])
		"YOURTRUMPCALL":
			YOURTRUMPCALL_received.emit(msg["MINSCORE"], msg["MAXSCORE"])
		"PLAYSTART":
			PLAYSTART_received.emit()
		"YOURPLAY":
			YOURPLAY_received.emit(msg["PLAYABLE"])
		"PLAY":
			PLAY_received.emit(msg["USERID"], msg["CARD"])
		"PLAYEND":
			PLAYEND_received.emit(msg["WINNERID"])

		_:
			push_error("Invalid or unknown action received from server:" + str(action))

# Requests

func sit_request(seat: int) -> Dictionary:
	return {"ACTION": "SIT", "SEAT": str(seat)}

func unsit_request() -> Dictionary:
	return {"ACTION": "UNSIT"}

func ready_request() -> Dictionary:
	return {"ACTION": "READY"}

func unready_request() -> Dictionary:
	return {"ACTION": "UNREADY"}

func trumpcall_request(score: String) -> Dictionary:
	return {"ACTION": "TRUMPCALL", "SCORE": score}

func play_request(card: String) -> Dictionary:
	return {"ACTION": "PLAY", "CARD": card}
