const validateEmail = (email) => {
    if (email === '') {
        return { valid: false, message: 'Email пуст' }
    }
    const re = /[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!re.test(email)) {
        return { valid: false, message: 'Некорректный формат Email' }
    }
    return { valid: true }
}

export default validateEmail