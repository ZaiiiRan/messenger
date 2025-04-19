export function camelToSnake(str: string): string {
    return str.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`)
}

export function transformKeysToSnakeCase(obj: any): any {
    if (Array.isArray(obj)) {
        return obj.map((item) => transformKeysToSnakeCase(item))
    }
    if (obj !== null && typeof obj === 'object') {
        return Object.keys(obj).reduce((acc, key) => {
            const snakeKey = camelToSnake(key)
            acc[snakeKey] = transformKeysToSnakeCase(obj[key])
            return acc
        }, {} as any)
    }
    return obj
}
