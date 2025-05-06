const apiErrors = {
    "invalid request format": "An unexpected error occurred",
    "unauthorized": "Unauthorized",
    "internal server error": "Internal server error",

    "error occured while sending activation code": "Error occured while sending activation code",

    // registration
    "username is empty": "Username is empty",
    "user with the same username already exists": "User with the same username already exists",
    "username must be at least 5 characters": "Username must be at least 5 characters",

    "email is empty": "Email is empty",
    "user with the same email already exists": "User with the same email already exists",
    "invalid email format": "Invalid email format",

    "user with the same phone number already exists": "User with the same phone number already exists",
    "phone must be in format +7(9xx)-xxx-xx-xx or empty": "Phone must be in format +7(9xx)-xxx-xx-xx or empty",

    "lastname is empty": "Lastname is empty",
    "lastname must be at least 2 characters": "Lastname must be at least 2 characters",
    "lastname must start with a capital letter": "Lastname must start with a capital letter",

    "firstname is empty": "Firstname is empty",
    "firstname must be at least 2 characters": "Firstname must be at least 2 characters",
    "firstname must start with a capital letter": "Firstname must start with a capital letter",

    "password is empty": "Password is empty",
    "password must be at least 8 characters": "Password must be at least 8 characters",
    "password must contain at least one uppercase letter": "Password must contain at least one uppercase letter",
    "password must contain at least one lowercase letter": "Password must contain at least one lowercase letter",
    "password must contain at least one digit": "Password must contain at least one digit",
    "password must contain at least one special character": "Password must contain at least one special character",

    // login
    "invalid login or password": "Invalid login or password",
    "login is empty": "Login is empty",

    // activation
    "user already activated": "User already activated",
    "failed to retrieve activation code": "Internal server error",
    "activation code not found": "Activation code not found",
    "activation code has expired": "Activation code has expired (request a new one)",
    "invalid activation code": "Invalid activation code",
    "failed to activate user account": "Failed to activate user account",

    // users fetching
    "search parameter is empty": "Search parameter is empty",
    "search parameter is very short": "Search parameter is very short",
    "users not found": "Users not found",

    // social functions
    "invalid user id": "User not found",
    "user not found": "User not found",
    "user is banned": "User is banned from the service",
    "you are blocked by this user": "You are blocked by this user",
    "you blocked this user": "You blocked this user",
    "you are already friends": "You are already friends",
    "friend request has already been sent": "Friend request has already been sent",

    // chat list fetching
    "chats not found": "Chats not found",
    "message not found": "An unexpected error occurred",
    "messages not found": "Messages not found",

    // chat creating
    "chat name is empty": "Chat name is empty",
    "chat name must be at least 5 characters": "Chat name must be at least 5 characters",
    "need at least 2 members for group chat": "Need at least 2 members for group chat",
    "maximum number of chat members: 1000": "Maximum number of chat members: 1000",
    "need at least 1 member for private chat": "Need at least 1 member for private chat",
    "max 1 member for private chat": "Max 1 member for private chat",


    // chat fetching
    "chat is not a group chat": "Chat is not a group chat",
    "you don't have enough rights": "You don't have enough rights",
    "the names are the same": "The names are the same",
    "you cannot access this chat": "You cannot access this chat",
    "chat not found": "Chat not found",
    "private chat not found": "Private chat not found",
    "you have been removed from the chat": "You have been removed from the chat",
    "you are already in chat": "You are already in chat",
    "you can't add yourself": "You can't add yourself",
    "you can't remove yourself": "You can't remove yourself",
    "trying to delete a member with a higher role": "Trying to delete a member with a higher role",
    "you can't change your role": "You can't change your role",
    "you cannot assign a role to an excluded member": "You cannot assign a role to an excluded member",
    "unknown role": "Unknown role",
    "owner role cannot be assigned": "Owner role cannot be assigned",
} as const

type ApiErrorsKey = keyof typeof apiErrors

export default apiErrors
export type { ApiErrorsKey }