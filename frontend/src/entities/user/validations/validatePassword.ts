import ValidateResponse from './validateResponse'

const validatePassword: (password: string) => ValidateResponse = (password) => {
    if (password === '') {
        return { valid: false, message: 'Password is empty' }
    }
    if (password.length < 8) {
        return { valid: false, message: 'Password must contain at least 8 characters' }
    }
    const hasUpperCase = /[A-ZА-ЯЁ]/
    const hasLowerCase = /[a-zа-яё]/
    const hasNumber = /[0-9]/
    const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/

    if (!hasUpperCase.test(password)) {
        return { valid: false, message: 'The password must contain at least one capital letter' }
    }
    if (!hasLowerCase.test(password)) {
        return { valid: false, message: 'The password must contain at least one lowercase letter' }
    }
    if (!hasNumber.test(password)) {
        return { valid: false, message: 'The password must contain at least one number' }
    }
    if (!hasSpecialChar.test(password)) {
        return { valid: false, message: 'The password must contain at least one special character' }
    }
    return { valid: true }
}

export default validatePassword