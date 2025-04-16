import ValidateResponse from './validateResponse'

const validateUsername: (username: string) => ValidateResponse = (username) => {
    if (username === '') {
        return { valid: false, message: 'Username is empty' }
    } else if (username.length < 5) {
        return { valid: false, message: 'Username must contain at least 5 characters' }
    }
    return { valid: true }
}

export default validateUsername