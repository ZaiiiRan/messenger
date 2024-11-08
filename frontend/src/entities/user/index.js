import validateUsername from "./validations/validateUsername"
import validateEmail from "./validations/validateEmail"
import validatePhone from "./validations/validatePhone"
import validatePassword from "./validations/validatePassword"
import { validateFirstName, validateLastName } from "./validations/validateName"
import validateBirthdate from './validations/validateBirthdate'
import useAuth from './hook/useAuth'

export { validateEmail, validateFirstName, validateLastName, validateUsername, 
    validatePhone, validatePassword, validateBirthdate, useAuth }