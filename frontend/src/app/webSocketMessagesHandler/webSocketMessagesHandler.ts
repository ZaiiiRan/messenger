import { IWebSocketMessage } from '../../shared/api'
import handleNewMessageNotification from './handlers/handleNewMessageNotification'

function handleWSMessage(wsMessage: IWebSocketMessage) {
    switch(wsMessage.type) {
        case "new_message_notification":
            handleNewMessageNotification(wsMessage)
            break
    }
}

export default handleWSMessage
