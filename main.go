package main;

import(
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

//Upgrades HTTP connections to websocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true;
	},
}

func wsHandler(w  http.ResponseWriter, r *http.Request) {
	//Opening connection and upgrading connections
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading:", err);
		return;
	}

	defer conn.Close();//Close right at end of function

	go handleConnection(conn);
}

func handleConnection(conn *websocket.Conn){
	for {
		//Read message from client
		_, message, err := conn.ReadMessage();//Message comes in the form of bytes
		if err != nil {
			fmt.Println("Error reading message:", err);
			break;
		}

		//Converting from bytes to string
		fmt.Println("Received: %s\\n", string(message));
		

		//Echoes message back to client
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err);
			break;
		}
	}
}

func main(){

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler);
	fmt.Println("WebSocket server started on: 8080");

	err := http.ListenAndServe(":8080", nil);

	if err != nil {
		fmt.Println("Error starting server:", err);
	}
}
