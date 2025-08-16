package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // untuk demo, production sebaiknya dibatasi
	},
}

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[string]*Client // map[deviceID]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// EmitToAll mengirim pesan ke semua client yang terhubung
func (h *Hub) EmitToAll(message []byte) {
	h.broadcast <- message
}

// EmitToClient mengirim pesan ke client tertentu berdasarkan ID
func (h *Hub) EmitToClient(deviceID string, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, exists := h.clients[deviceID]; exists {
		select {
		case client.send <- message:
		default:
			// Client buffer penuh, disconnect client
			delete(h.clients, deviceID)
			close(client.send)
		}
	}
}

// GetClients mengembalikan daftar semua client yang terhubung
func (h *Hub) GetClients() []*Client {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := make([]*Client, 0, len(h.clients))
	for _, client := range h.clients {
		clients = append(clients, client)
	}
	return clients
}

// GetClientByID mengembalikan client berdasarkan ID
func (h *Hub) GetClientByID(deviceID string) (*Client, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clients[deviceID]
	return client, exists
}

// GetdeviceIDs mengembalikan daftar semua client ID yang terhubung
func (h *Hub) GetdeviceIDs() []string {
	h.mu.Lock()
	defer h.mu.Unlock()

	ids := make([]string, 0, len(h.clients))
	for id := range h.clients {
		ids = append(ids, id)
	}
	return ids
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.id] = client
			h.mu.Unlock()
			log.Printf("Client connected with ID: %s, Address: %s", client.id, client.conn.RemoteAddr())

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
				client.conn.Close()
				log.Printf("Client disconnected with ID: %s, Address: %s", client.id, client.conn.RemoteAddr())
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Received: %s", message)
		h.broadcast <- message
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func ServeWS(h *Hub, w http.ResponseWriter, r *http.Request) {
	// Ambil client ID dari query parameter
	deviceID := r.URL.Query().Get("id")

	// Cek apakah client ID sudah ada
	h.mu.Lock()
	if _, exists := h.clients[deviceID]; exists {
		h.mu.Unlock()
		http.Error(w, "Client ID already connected", http.StatusConflict)
		return
	}
	h.mu.Unlock()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		id:   deviceID,
		conn: conn,
		send: make(chan []byte, 256),
	}
	h.register <- client

	go client.readPump(h)
	go client.writePump()
}
