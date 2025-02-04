import { api } from '../../../shared/api'

export default class Auth {
    static async login(login, password) {
        return api.post('/auth/login', { login, password })
    }

    static async register(username, email, password, firstname, lastname, phone, birthdate) {
        return api.post('/auth/register', { username, email, password, firstname, lastname, phone, birthdate })
    }

    static async logout() {
        return api.delete('/auth/logout')
    }

    static async refresh() {
        return api.get('/auth/refresh')
    }
}
