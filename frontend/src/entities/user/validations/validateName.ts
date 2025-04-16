import ValidateResponse from './validateResponse'

const validateFirstName: (name: string) => ValidateResponse = (name) => {
    if (name === '') {
        return { valid: false, message: 'Firstname is empty' }
    }
    if (name.length < 2) {
        return { valid: false, message: 'Firstname must contain at least 2 letters' }
    }
    const re = /^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$/
    if (!re.test(name)) {
        return { valid: false, message: 'Firstname must begin with a capital letter and not contain numbers or special characters' }
    }
    return { valid: true }
}

const validateLastName: (name: string) => ValidateResponse = (name) => {
    if (name === '') {
        return { valid: false, message: 'Lastname is empty' }
    }
    if (name.length < 2) {
        return { valid: false, message: 'Lastname must contain at least 2 letters' }
    }
    const re = /^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$/
    if (!re.test(name)) {
        return { valid: false, message: 'Lastname must begin with a capital letter and not contain numbers or special characters' }
    }
    return { valid: true }
}

export { validateFirstName, validateLastName }