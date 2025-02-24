extends Node

var instance_id = "1234"
var username = "Player"
var user_id = "11223344"
var icon_url = "https://google.com"

var _backend_url = "ws://localhost:3000/ws"

var _socket = WebSocketPeer.new()
var authenticated = false
var _auth_requested


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	# TODO: Get data from JSbridge
	var err = _socket.connect_to_url(_backend_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)



# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(_delta: float) -> void:
	_socket.poll()

	if _socket.get_ready_state() != WebSocketPeer.STATE_OPEN:
		return

	if !authenticated:
		_handle_auth()
		return

	_read_loop()

func _handle_auth() -> void:
	if !_auth_requested:
		var authMsg = {"ACTION": "AUTH", "INSTANCEID": instance_id,"USERID": user_id, "USERNAME": username, "ICONURL": icon_url}
		_write_json(authMsg)
		_auth_requested = true
		return
	
	var msg = _read_json()
	if msg == {}:
		return
	
	if msg["ACTION"] == "OK":
		authenticated = true
		print("Authenticated!")
	else:
		# TODO: handle error
		print("Authentication failed: " + msg["MESSAGE"] + ", exiting...")
		set_process(false)

func _read_loop() -> void:
	var msg = _read_json()
	if msg == {}:
		return
		
	EventManager._handle_message(msg)

func _write_json(msg: Dictionary) -> void:
	var msgJson = JSON.stringify(msg)
	_socket.send_text(msgJson)

func _read_json() -> Dictionary:
	if _socket.get_available_packet_count():
		return JSON.parse_string(_socket.get_packet().get_string_from_utf8())
	return {}
