import { makeAutoObservable } from 'mobx'

class ThemeStore {
    theme: string = localStorage.getItem('theme') || 'system'

    constructor() {
        makeAutoObservable(this)
    }

    setTheme(theme: string) {
        this.theme = theme
        localStorage.setItem('theme', theme)
        this.applyTheme(theme)
    }

    applyTheme(theme: string) {
        const root = document.documentElement

        this.removeSystemThemeListener()

        if (theme === 'system') {
            const prefersLightScheme = window.matchMedia("(prefers-color-scheme: light)").matches
            root.setAttribute("data-color-scheme", prefersLightScheme ? "light" : "dark")
            this.addSystemThemeListener()
        } else {
            root.setAttribute("data-color-scheme", theme)
        }
    }

    handleSystemThemeChange = (e: MediaQueryListEvent) => {
        this.applyTheme(e.matches ? "light" : "dark")
    }

    removeSystemThemeListener() {
        window.matchMedia("(prefers-color-scheme: light)").removeEventListener("change", this.handleSystemThemeChange)
    }

    addSystemThemeListener() {
        window.matchMedia("(prefers-color-scheme: light)").addEventListener("change", this.handleSystemThemeChange)
    }
}

export const themeStore = new ThemeStore()