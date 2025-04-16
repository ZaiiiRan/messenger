const getTime: (date: Date) => string = (date) => {
    const hours = date.getHours()
    const minutes = date.getMinutes()

    const formattedTime = `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
    return formattedTime
}

export default getTime