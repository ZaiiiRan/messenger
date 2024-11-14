import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'userCard', en)
i18n.addResourceBundle('ru', 'userCard', ru)

import validateUsername from "./validations/validateUsername"
import validateEmail from "./validations/validateEmail"
import validatePhone from "./validations/validatePhone"
import validatePassword from "./validations/validatePassword"
import { validateFirstName, validateLastName } from "./validations/validateName"
import validateBirthdate from './validations/validateBirthdate'
import useAuth from './hook/useAuth'
import User from "./ui/User"

export { validateEmail, validateFirstName, validateLastName, validateUsername, 
    validatePhone, validatePassword, validateBirthdate, useAuth, User }