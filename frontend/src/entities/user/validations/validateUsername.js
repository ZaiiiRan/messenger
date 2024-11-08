const validateUsername = (username) => {
    if (username === '') {
        return { valid: false, message: 'Username is empty' }
    } else if (username.length < 5) {
        return { valid: false, message: 'Username must contain at least 5 characters' }
    }
    return { valid: true }
}

export default validateUsername