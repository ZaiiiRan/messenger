import { makeAutoObservable } from 'mobx'
import Auth from '../api/auth'

class UserStore {
    user = {}
    isAuth = false
    isLoading = false
    isBegin = true

    constructor() {
        makeAutoObservable(this)
    }

    setAuth(bool) {
        this.isAuth = bool
    }

    setUser(user) {
        this.user = user
    }

    setLoading(bool) {
        this.isLoading = bool
    }

    setBegin(bool) {
        this.isBegin = bool
    }

    async login(username, password) {
        const response = await Auth.login(username, password)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
    }

    async register(username, email, password, firstname, lastname, phone, birthdate) {
        const response = await Auth.register(username, email, password, firstname, lastname, phone, birthdate)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
    }

    async logout() {
        const response = await Auth.logout()
        localStorage.removeItem('token', response.data.accessToken)
        this.setAuth(false)
        this.setUser({})
    }

    async checkAuth() {
        const response = await Auth.refresh()
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
    }
}

export default (new UserStore())
