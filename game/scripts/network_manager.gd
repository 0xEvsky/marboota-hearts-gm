extends Node

var instance_id = "1234"
var username = "Player"
var user_id = str(randi_range(1,1000000))
var icon_url = "https://placehold.co/128x128.png?text=" + username + user_id

var _backend_url_suffix = "/.proxy/backend/ws"

var _socket = WebSocketPeer.new()
var authenticated = false
var _auth_requested

signal AUTH_accepted


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	# TODO: Get data from JSbridge
	var _backend_url = JavaScriptBridge.eval("window.location.hostname", true) # true = safe
	var full_url = "wss://" + _backend_url + _backend_url_suffix
	#$"../Game/LoadingUI/Label".text = full_url
	#var discord = JavaScriptBridge.get_interface("discord")
	#icon_url = "https://cdn.discordapp.com/avatars/" + discord.session.id + "/" + discord.session.avatar + ".png?size=128"
	var err = _socket.connect_to_url(full_url)
	if err != OK:
		push_error("Unable to connect")
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
		AUTH_accepted.emit()
		print_debug("Authenticated!")
	elif msg["ACTION"] == "ERROR":
		# TODO: handle error
		push_error("Authentication failed: " + msg["MESSAGE"] + ", exiting...")
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
