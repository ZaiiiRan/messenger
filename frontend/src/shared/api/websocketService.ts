export const WS_URL = import.meta.env.VITE_WS_URL

class WebSocketService {
    socket: WebSocket | null = null
    retries: number = 0
    maxRetries: number = 2

    async connect() {
        if (this.socket) return

        const token = localStorage.getItem('token')

        try {
            this.socket = new WebSocket(WS_URL)

            this.socket.onopen = () => {
                console.log("✅ WebSocket подключен!")
                this.sendMessage({ token })
            }

            this.socket.onmessage = (event) => {
                const data = JSON.parse(event.data)
                if (data.type === "error" && data.content === "unauthorized") {
                    console.warn("⚠️ Ошибка авторизации, пробуем ещё раз...")

                    if (this.retries < this.maxRetries) {
                        this.retries++
                        this.disconnect()
                        setTimeout(() => this.connect(), 1000)
                    } else {
                        throw new Error("🚫 Доступ запрещён: неверный токен")
                    }
                } else {
                    console.log("📩 Сообщение:", event.data)
                }
            }

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

    sendMessage(message: any) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message));
        } else {
            console.error("❌ WebSocket не подключен!");
        }
    }
}

export const webSocketService = new WebSocketService()