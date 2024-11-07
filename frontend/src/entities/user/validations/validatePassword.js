const validatePassword = (password) => {
    if (password === '') {
        return { valid: false, message: 'Пароль пуст' }
    }
    if (password.length < 8) {
        return { valid: false, message: 'Пароль должен содержать хотя бы 8 символов' }
    }
    const hasUpperCase = /[A-ZА-ЯЁ]/
    const hasLowerCase = /[a-zа-яё]/
    const hasNumber = /[0-9]/
    const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/

    if (!hasUpperCase.test(password)) {
        return { valid: false, message: 'Пароль должен содержать хотя бы одну заглавную букву' }
    }
    if (!hasLowerCase.test(password)) {
        return { valid: false, message: 'Пароль должен содержать хотя бы одну строчную букву' }
    }
    if (!hasNumber.test(password)) {
        return { valid: false, message: 'Пароль должен содержать хотя бы одну цифру' }
    }
    if (!hasSpecialChar.test(password)) {
        return { valid: false, message: 'Пароль должен содержать хотя бы один специальный символ' }
    }
    return { valid: true }
}

export default validatePassword