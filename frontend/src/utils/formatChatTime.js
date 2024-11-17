import { format } from 'date-fns'
import { differenceInCalendarDays } from 'date-fns'

const formatChatTime = (date) => {
    const now = new Date()
    const messageDate = new Date(date)

    const dayDifference = differenceInCalendarDays(now, messageDate)

    if (dayDifference === 0) {
        return format(messageDate, 'HH:mm')
    }

    if (dayDifference === 1) {
        return 'Yesterday'
    }

    if (dayDifference <= 365) {
        return format(messageDate, 'dd.MM')
    }

    return format(messageDate, 'dd.MM.yyyy')
}

export default formatChatTime