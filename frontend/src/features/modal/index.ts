import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'modal', en)
i18n.addResourceBundle('ru', 'modal', ru)

import Modal from "./ui/Modal"
import useModal from "./hooks/useModal"
import modalStore from "./store/modalStore"

export { Modal, useModal, modalStore }