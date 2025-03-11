extends Node

# Signals that don't need requests (make sure to subscribe to those)
signal JOIN_received
signal LEAVE_received
signal SIT_received
signal UNSIT_received
signal READY_received
signal UNREADY_received

var _request_queue: Array[Dictionary] = []
var _response_queue: Array[Dictionary] = []

func _process(_delta: float) -> void:
    #if !_response_queue.is_empty():
        #_process_response_queue()
    pass

func send_request(msg: Dictionary, on_success: Callable, on_error: Callable) -> void:
    var request_id = _generate_request_id()
    msg["REQUESTID"] = request_id
    _request_queue.append({"message": msg, "request_id": request_id, "on_success": on_success, "on_error": on_error})
    print_debug("request: ", msg)
    NetworkManager._write_json(msg)

func _generate_request_id() -> String:
    return "request" + "-" + NetworkManager.user_id + "-" + str(RandomNumberGenerator.new().randi())

func _handle_message(msg: Dictionary) -> void:
    _response_queue.append(msg)
    _process_response_queue()

func _process_response_queue() -> void:
    var res = _response_queue.pop_front()
    print_debug("request: ", res)
    for req in _request_queue:
        if res["REQUESTID"] == req["request_id"]:
            if res["ACTION"] == "OK":
                req["on_success"].call()
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
            UNSIT_received.emit()
        "READY":
            READY_received.emit()
        "UNREADY":
            UNREADY_received.emit()
        _:
            push_error("Invalid or unknown action received from server")

# Requests

func sit_request(seat: int) -> Dictionary:
    return {"ACTION": "SIT", "SEAT": str(seat)}

func unsit_request() -> Dictionary:
    return {"ACTION": "UNSIT"}

func ready_request() -> Dictionary:
    return {"ACTION": "READY"}

func unready_request() -> Dictionary:
    return {"ACTION": "UNREADY"}