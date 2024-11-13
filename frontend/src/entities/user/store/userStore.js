import { makeAutoObservable } from 'mobx'
import Auth from '../api/auth'
import Activation from '../api/activation'

class UserStore {
    user = {}
    isAuth = false
    isLoading = false
    isBegin = true
    isOpen = false

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

    setOpen(bool) {
        this.isOpen = bool
    }

    async login(username, password) {
        const response = await Auth.login(username, password)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        return response.data
    }

    async register(username, email, password, firstname, lastname, phone, birthdate) {
        const response = await Auth.register(username, email, password, firstname, lastname, phone, birthdate)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        return response.data
    }

    async logout() {
        const response = await Auth.logout()
        localStorage.removeItem('token', response.data.accessToken)
        this.setAuth(false)
        this.setUser({})
        return response.data
    }

    async checkAuth() {
        const response = await Auth.refresh()
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        return response.data
    }

    async activate(code) {
        if (this.user.isActivated) throw Error('Аккаунт уже активирован')
        const response = await Activation.activate(code)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        return response.data
    }

    async resend() {
        if (this.user.isActivated) throw Error('Аккаунт уже активирован')
        const response = await Activation.resend()
        return response.data
    }
}

export default (new UserStore())
