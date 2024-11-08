const validateEmail = (email) => {
    if (email === '') {
        return { valid: false, message: 'Email is empty' }
    }
    const re = /[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!re.test(email)) {
        return { valid: false, message: 'Incorrect Email format' }
    }
    return { valid: true }
}

export default validateEmail