import { IWebSocketMessage, webSocketService } from '../../../shared/api'

function sendMessage(chatId: string | number, message: string) {
    const wsMessage: IWebSocketMessage = {
        type: "send_message",
        content: {
            chatId: chatId,
            messageContent: message
        }
    }

    webSocketService.sendMessage(wsMessage)
}

export default sendMessage