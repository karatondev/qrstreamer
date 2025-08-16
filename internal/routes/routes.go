package routes

import (
	"fmt"
	"net/http"
	"qrstreamer/internal/handler"
	"qrstreamer/internal/service"
)

func RegisterRoutes(hub *handler.Hub, svc service.QRStreamer) {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("id")
		if deviceID == "" {
			// Jika tidak ada ID, ambil dari header
			deviceID = r.Header.Get("Device-ID")
		}
		if deviceID == "" {
			http.Error(w, "Device ID is required. Use ?id=your_device_id or Device-ID header", http.StatusBadRequest)
			return
		}

		handler.ServeWS(hub, w, r)

		if err := svc.StreamWhatsappQR(r.Context(), deviceID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	// Default root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `WebSocket server running at ws://localhost:8080/ws`)
	})
}
