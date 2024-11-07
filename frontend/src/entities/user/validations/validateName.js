const validateFirstName = (name) => {
    if (name === '') {
        return { valid: false, message: 'Имя пусто' }
    }
    if (name.length < 2) {
        return { valid: false, message: 'Имя должно содержать хотя бы 2 буквы' }
    }
    const re = /^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$/
    if (!re.test(name)) {
        return { valid: false, message: 'Имя должно начинаться с заглавной буквы и не содержать цифр и спец. символов' }
    }
    return { valid: true }
}

const validateLastName = (name) => {
    if (name === '') {
        return { valid: false, message: 'Фамилия пустая' }
    }
    if (name.length < 2) {
        return { valid: false, message: 'Фамилия должна содержать хотя бы 2 буквы' }
    }
    const re = /^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$/
    if (!re.test(name)) {
        return { valid: false, message: 'Фамииля должна начинаться с заглавной буквы и не содержать цифр и спец. символов' }
    }
    return { valid: true }
}

export { validateFirstName, validateLastName }