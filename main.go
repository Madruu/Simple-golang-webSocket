package main;

import(
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

//Upgrades HTTP connections to websocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "<http://localhost:8080>";
	},
}


//Connecting clients, create broadcast channel and protect clients map
var clients = make(map[*websocket.Conn]bool);
var broadcast = make(chan []byte);
var mutex = &sync.Mutex{};

func wsHandler(w  http.ResponseWriter, r *http.Request) {
	//Opening connection and upgrading connections
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading:", err);
		return;
	}

	defer conn.Close(); //Close right at end of function

	mutex.Lock();
	clients[conn] = true;
	mutex.Unlock();

	for {
		_, message, err := conn.ReadMessage();
		if err != nil {
			mutex.Lock(); //Locks the mutex for secutity
			delete(clients, conn); //Deletes information in case of emergency
			mutex.Unlock(); //Unlocks mutex
			break;
		}
		broadcast <- message;
	}
	
}

func handleMessages() {
	for {
		//Grabs next message from broadcast channel
		message := <-broadcast;
		//Sends message to all clients connected
		mutex.Lock();
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message);
			if err != nil{
				client.Close();
				delete(clients, client);
			}
		}
		mutex.Unlock();
	}
}

// func handleConnection(conn *websocket.Conn){
// 	for {
// 		//Read message from client
// 		_, message, err := conn.ReadMessage();//Message comes in the form of bytes
// 		if err != nil {
// 			fmt.Println("Error reading message:", err);
// 			break;
// 		}

// 		//Converting from bytes to string
// 		fmt.Println("Received: %s\\n", string(message));
		

// 		//Echoes message back to client
// 		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
// 			fmt.Println("Error writing message:", err);
// 			break;
// 		}
// 	}
// }

func main(){

	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler);
	fmt.Println("WebSocket server started on: 8080");

	err := http.ListenAndServe(":8080", nil);

	if err != nil {
		fmt.Println("Error starting server:", err);
	}
}
