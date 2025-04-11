import ValidateResponse from './validateResponse'

const validatePhone: (phone: string) => ValidateResponse = (phone) => {
    if (phone === '') {
        return { valid: true }
    }
    const re = /^\+7\(9\d{2}\)-\d{3}-\d{2}-\d{2}$/
    if (!re.test(phone)) {
        return { valid: false, message: 'Invalid phone number format' }
    }
    return { valid: true }
}

export default validatePhone