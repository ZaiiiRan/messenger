import { i18n }  from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'registerFeature', en)
i18n.addResourceBundle('ru', 'registerFeature', ru)

import StepAdditionalInfoRegister from "./ui/StepAdditionalInfoRegister"
import StepEmailUsername from "./ui/StepEmailUsername"
import StepNames from "./ui/StepNames"
import StepPassword from "./ui/StepPassword"

export { StepAdditionalInfoRegister, StepEmailUsername, StepNames, StepPassword }