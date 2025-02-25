extends Node

# Signals that don't need requests (make sure to subscribe to those)
signal JOIN_received
signal LEAVE_received
signal SIT_received
signal UNSIT_received

var _request_queue: Array[Dictionary] = []
var _response_queue: Array[Dictionary] = []

func _process(_delta: float) -> void:
    if !_response_queue.is_empty():
        _process_response_queue()

func send_request(msg: Dictionary, on_success: Callable, on_error: Callable) -> void:
    var request_id = _generate_request_id()
    _request_queue.append({"message": msg, "request_id": request_id, "on_success": on_success, "on_error": on_error})
    NetworkManager._write_json(msg)

func _handle_message(msg: Dictionary) -> void:
    _response_queue.append(msg)
    _process_response_queue()

func _process_response_queue() -> void:
    for res in _response_queue:
        if res["REQUESTID"] == _request_queue.front()["request_id"]:
            if res["ACTION"] == "OK":
                _request_queue.front()["on_success"].call()
            elif res["ACTION"] == "ERROR":
                _request_queue.front()["on_error"].call(res["MESSAGE"])
            _request_queue.pop_front()
            _response_queue.erase(res)
            break

func _dispatch(action: String) -> void:
    match action:
        "JOIN":
            JOIN_received.emit()
        "LEAVE":
            LEAVE_received.emit()
        "SIT":
            SIT_received.emit()
        "UNSIT":
            UNSIT_received.emit()
        _:
            print("Invalid action")

func _generate_request_id() -> String:
    return "request" + "-" + NetworkManager.user_id + "-" + str(RandomNumberGenerator.new().randi())

func sit_request(seat: int) -> Dictionary:
    return {"ACTION": "SIT", "SEAT": str(seat)}
