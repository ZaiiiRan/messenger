import ValidateResponse from './validateResponse'

const validateBirthdate: (date: string) => ValidateResponse = (date) => {
    if (date === '') {
        return { valid: true }
    }

    const re = /^(0[1-9]|[12][0-9]|3[01])\.(0[1-9]|1[0-2])\.(\d{4})$/
    if (!re.test(date)) {
        return { valid: false, message: 'Invalid date format. Expected format is dd.mm.yyyy' }
    }

    const [day, month, year] = date.split('.').map(Number)
    const parsedDate = new Date(year, month - 1, day)
    const today = new Date()

    if (parsedDate > today) {
        return { valid: false, message: 'The date must not be in the future' }
    }

    if (parsedDate.getDate() !== day || parsedDate.getMonth() + 1 !== month || parsedDate.getFullYear() !== year) {
        return { valid: false, message: 'Invalid date' }
    }

    return { valid: true }
}

export default validateBirthdate
