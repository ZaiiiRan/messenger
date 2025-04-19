export function snakeToCamel(str: string): string {
    return str.replace(/(_\w)/g, (match) => match[1].toUpperCase())
}

export function transformKeysToCamelCase(obj: any): any {
    if (Array.isArray(obj)) {
        return obj.map((item) => transformKeysToCamelCase(item))
    }
    if (obj !== null && typeof obj === 'object') {
        return Object.keys(obj).reduce((acc, key) => {
            const camelKey = snakeToCamel(key)
            acc[camelKey] = transformKeysToCamelCase(obj[key])
            return acc
        }, {} as any)
    }
    return obj
}