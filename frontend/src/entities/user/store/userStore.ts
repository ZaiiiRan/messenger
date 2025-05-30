import { makeAutoObservable } from 'mobx'
import Auth from '../api/auth'
import Activation from '../api/activation'
import { webSocketService } from '../../../shared/api'
import IUser from '../models/IUser'

class UserStore {
    user: IUser | null = null
    isAuth: boolean = false
    isLoading: boolean = false
    isBegin: boolean = true
    isOpen: boolean = false

    constructor() {
        makeAutoObservable(this)
    }

    setAuth(bool: boolean) {
        this.isAuth = bool
    }

    setUser(user: any) {
        this.user = user
    }

    setLoading(bool: boolean) {
        this.isLoading = bool
    }

    setBegin(bool: boolean) {
        this.isBegin = bool
    }

    setOpen(bool: boolean) {
        this.isOpen = bool
    }

    async login(username: string, password: string) {
        const response = await Auth.login(username, password)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        await this.connectToWS()
        return response.data
    }

    async register(username: string, email: string, password: string, firstname: string, 
    lastname: string, phone: string | undefined | null, birthdate: string | undefined | null) {
        const response = await Auth.register(username, email, password, firstname, lastname, phone, birthdate)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        return response.data
    }

    async logout() {
        const response = await Auth.logout()
        localStorage.removeItem('token')
        this.setAuth(false)
        this.setUser(null)
        webSocketService.disconnect()
        this.setOpen(false)
        return response.data
    }

    async checkAuth() {
        try {
            const response = await Auth.refresh()
            localStorage.setItem('token', response.data.accessToken)
            this.setAuth(true)
            this.setUser(response.data.user)
            await this.connectToWS()
            return response.data
        } catch (e: any) {
            if (e.response.status === 401) {
                this.setAuth(false)
                this.setUser(null)
                webSocketService.disconnect()
            }
        }
    }

    async activate(code: string) {
        if (this.user && this.user.isActivated) throw Error('Аккаунт уже активирован')
        const response = await Activation.activate(code)
        localStorage.setItem('token', response.data.accessToken)
        this.setAuth(true)
        this.setUser(response.data.user)
        await this.connectToWS()
        return response.data
    }

    async resend() {
        if (this.user && this.user.isActivated) throw Error('Аккаунт уже активирован')
        const response = await Activation.resend()
        return response.data
    }

    private async connectToWS() {
        this.checkConditionsForDisconnectWS()
        if (this.user && this.user.isActivated && !this.user.isDeleted 
                && !this.user.isBanned && !webSocketService.isConnected()) {
            await webSocketService.connect()
        }
    }

    private checkConditionsForDisconnectWS() {
        if (!this.isAuth || !this.user || !this.user.isActivated || this.user.isDeleted || this.user.isBanned) {
            webSocketService.disconnect()
        }
    }
}

export default (new UserStore())
