const getDateLabel = (messageDate) => {
    const messageDateObj = new Date(messageDate)
    const today = new Date()
    const yesterday = new Date(today)
    yesterday.setDate(today.getDate() - 1)

    if (messageDateObj.toDateString() === today.toDateString()) {
        return "Today"
    }

    if (messageDateObj.toDateString() === yesterday.toDateString()) {
        return "Yesterday"
    }

    return messageDateObj.toLocaleDateString()
}

export default getDateLabel