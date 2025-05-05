import { IValidationError } from "../../../shared/validationError"

function validateMembers(members: (string | number)[]): IValidationError {
    if (members.length < 2) {
        return {
            error: true,
            message: 'Need at least 2 members for group chat'
        }
    } else if (members.length > 1000) {
        return {
            error: true,
            message: 'Maximum number of chat members: 1000'
        }
    }
    return {
        error: false
    }
}

export default validateMembers