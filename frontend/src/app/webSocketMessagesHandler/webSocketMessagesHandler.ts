import { IWebSocketMessage } from '../../shared/api'
import handleNewMessageNotification from './handlers/handleNewMessageNotification'
import handleError from './handlers/handleError'

function handleWSMessage(wsMessage: IWebSocketMessage) {
    switch(wsMessage.type) {
        case "new_message_notification":
            handleNewMessageNotification(wsMessage)
            break
        case "error":
            handleError(wsMessage)
            break
    }
}

export default handleWSMessage
