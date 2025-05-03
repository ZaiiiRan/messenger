import IWebSocketMessage from './models/IWebSocketMessage'
import handleNewMessageNotification from './handlers/handleNewMessageNotification'

const handleWSMessage = (data: any) => {
    const wsMessage = data as IWebSocketMessage
    handle(wsMessage)
}

function handle(wsMessage: IWebSocketMessage) {
    switch(wsMessage.type) {
        case "new_message_notification":
            handleNewMessageNotification(wsMessage)
            break
    }
}

export default handleWSMessage
