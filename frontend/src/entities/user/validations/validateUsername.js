const validateUsername = (username) => {
    if (username === '') {
        return { valid: false, message: 'Имя пользователя пусто' }
    } else if (username.length < 5) {
        return { valid: false, message: 'Имя пользователя должно содержать хотя бы 5 символов' }
    }
    return { valid: true }
}

export default validateUsername