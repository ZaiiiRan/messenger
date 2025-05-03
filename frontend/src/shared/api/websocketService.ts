import { transformKeysToCamelCase } from "../../utils/transformKeysToCamelCase"

export const WS_URL = import.meta.env.VITE_WS_URL

class WebSocketService {
    socket: WebSocket | null = null
    retries: number = 0
    maxRetries: number = 2
    messageHandler: ((data: any) => void) | null = null

    async connect() {
        if (this.socket) return

        const token = localStorage.getItem('token')

        try {
            this.socket = new WebSocket(WS_URL)

            this.socket.onopen = () => {
                console.log("WebSocket connected")
                this.sendMessage({ token })
            }

            this.socket.onmessage = (event) => {
                const data = JSON.parse(event.data)
                if (data.type === "error" && data.content === "unauthorized") {
                    console.warn("Authorization error, try again...")

                    if (this.retries < this.maxRetries) {
                        this.retries++
                        this.disconnect()
                        setTimeout(() => this.connect(), 1000)
                    } else {
                        throw new Error("Access to WebSocket Denied: Invalid Token")
                    }
                } else {
                    if (this.messageHandler) {
                        this.messageHandler(transformKeysToCamelCase(data))
                    }
                }
            }

            this.socket.onclose = () => console.log("WebSocket closed")
            this.socket.onerror = (error) => console.error("WebSocket Error:", error)
        } catch (error) {
            console.error("WebSocket authentication error:", error)
        }
    }

    disconnect() {
        if (this.socket) {
            this.socket.close()
            this.socket = null
            console.log("WebSocket is disabled")
        }
    }

    sendMessage(message: any) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify(message))
        } else {
            console.error("WebSocket not connected")
        }
    }

    setHandler(handler: (message: any) => void) {
        this.messageHandler = handler
    }
}

export const webSocketService = new WebSocketService()