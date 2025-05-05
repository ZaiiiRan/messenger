import { IValidationError } from "../../../shared/validationError"

function validateChatName(chatName: string): IValidationError {
    const trimmedName = chatName.trim()
    if (trimmedName === "") {
        return {
            message: "Chat name is empty",
            error: true
        }
    } else if (trimmedName.length < 5) {
        return {
            message: "Chat name must be at least 5 characters",
            error: true
        }
    }
    return {
        error: false
    }
}

export default validateChatName