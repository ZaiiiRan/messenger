import { IChat } from "../../../entities/Chat"
import { api } from "../../../shared/api"
import { saveChat } from "../../../features/chatsFetching"

async function createChat(name: string, members: (number | string)[]): Promise<IChat> {
    const response = await api.post('/chats', { name, members, isGroup: true })

    const chat = saveChat(response.data)

    return chat
}

export default createChat