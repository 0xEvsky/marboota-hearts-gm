extends Node

# Signals that don't need requests (make sure to subscribe to those)
signal JOIN_received
signal LEAVE_received
signal SIT_received
signal UNSIT_received

var request_queue = []
var response_queue = []

func _process(_delta: float) -> void:
    if !response_queue.is_empty():
        _process_response_queue()

func send_request(msg: Dictionary, on_success: Callable, on_error: Callable) -> void:
    var message = msg
    message["REQUESTID"] = "123" # TODO: Generate
    request_queue.append({"message": message, "on_success": on_success, "on_error": on_error})
    NetworkManager._write_json(msg)

func _handle_message(msg: Dictionary) -> void:
    response_queue.append(msg)
    _process_response_queue()

func _process_response_queue() -> void:
    for res in response_queue:
        if res["REQUESTID"] == request_queue.front()["REQUESTID"]:
            if res["ACTION"] == "OK":
                request_queue.front()["on_success"].call()
            elif res["ACTION"] == "ERROR":
                request_queue.front()["on_error"].call(res["MESSAGE"])
            request_queue.pop_front()
            response_queue.erase(res)
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
