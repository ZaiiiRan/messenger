import { modalStore } from '../../../features/modal'
import { apiErrors, ApiErrorsKey, IWebSocketMessage } from '../../../shared/api'
import { i18n } from '../../../shared/config/i18n'

function handleError(wsMessage: IWebSocketMessage) {
    const errorKey = wsMessage.content as ApiErrorsKey

    const title = i18n.t('Error', { ns: 'app' })
    const message = i18n.t(apiErrors[errorKey], { ns: 'app' }) || i18n.t('Internal server error', { ns: 'app' })
    modalStore.openModal(title, message)
}

export default handleError