export const WS_URL = import.meta.env.VITE_WS_URL

class WebSocketService {
    socket = null

    async connect() {
        if (this.socket) return

        try {
            this.socket = new WebSocket(`ws://localhost:8080/ws?token=${localStorage.getItem('token')}`)

            this.socket.onopen = () => console.log("✅ WebSocket подключен!")
            this.socket.onmessage = (event) => console.log("📩 Сообщение:", event.data)
            this.socket.onclose = () => console.log("🔄 WebSocket закрыт!")
            this.socket.onerror = (error) => console.error("❌ Ошибка WebSocket:", error)
        } catch (error) {
            console.error("❌ Ошибка WebSocket-аутентификации:", error)
        }
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
            console.log("🔴 WebSocket отключен");
        }
    }

    sendMessage(message) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        } else {
            console.error("❌ WebSocket не подключен!");
        }
    }
}

export const webSocketService = new WebSocketService()