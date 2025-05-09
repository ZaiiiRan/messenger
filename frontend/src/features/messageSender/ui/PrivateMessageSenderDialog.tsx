import { useTranslation } from "react-i18next"
import { Dialog } from "../../../shared/ui/Dialog"
import { IShortUser, ISocialUser, ShortUser } from "../../../entities/SocialUser"
import { useEffect, useState } from "react"
import { Textarea } from "../../../shared/ui/Textarea"
import { Button } from "../../../shared/ui/Button"
import { useModal } from "../../modal"
import sendPrivateMessage from "../api/sendPrivateMessage"
import { apiErrors, ApiErrorsKey } from "../../../shared/api"
import { Loader } from "../../../shared/ui/Loader"
import { AxiosError } from "axios"

interface PrivateMessageSenderDialogProps {
    show: boolean,
    setShow: (show: boolean) => void,
    recipient: IShortUser,
    zIndex?: number
}

const PrivateMessageSenderDialog: React.FC<PrivateMessageSenderDialogProps> = ({ show, setShow, recipient, zIndex }) => {
    const { t } = useTranslation('messageSender')
    const [message, setMessage] = useState<string>('')
    const [isSending, setIsSending] = useState<boolean>(false)

    const { openModal } = useModal()

    const showModal = (title: string, message: string) => {
        openModal(title, message)
    }

    const send = async () => {
        const trimmedMessage = message.trim()
        if (message.trim().length === 0) {
            showModal(t('Error'), t('Message can\'t be empty'))
            return
        }

        setIsSending(true)
        try {
            await sendPrivateMessage(recipient.userId, trimmedMessage)
            showModal(t('Success'), t('The message has been sent'))
            setShow(false)
        } catch (e: any) {
            if (e instanceof AxiosError && e.response?.status === 400) {
                const errorKey: ApiErrorsKey = e.response?.data?.error
                showModal(t('Error'), `${t('You cannot write to the user')} ${recipient.username}`)
            } else {
                showModal(t('Error'), t('Internal server error'))
            }

        } finally {
            setIsSending(false)
        }
    }

    useEffect(() => {
        setMessage('')
    }, [show, recipient])
    
    return (
        <Dialog
            show={show}
            setShow={setShow}
            title={t('New message')}
            zIndex={zIndex}
            id='send-message'
        >
            <ShortUser 
                isClickable={false}
                user={recipient as IShortUser}
            />

            <Textarea 
                className='w-full px-5 py-3 2k:px-6 2k:py-4 4k:px-7 4k:py-5 
                    rounded-lg md:text-lg mobile:text-sm lg:text-sm 2xl:text-lg 2k:text-2xl 4k:text-4xl 2k:h-24 4k:h-40' 
                placeholder={t('Enter message')}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
            />

            <div className='self-end'>
                <Button
                    className='h-12 flex items-center justify-center 2k:h-16 4k:h-28 w-72 xl:w-60 lg:w-56 md:w-60 sm:w-56 mobile:w-56 2k:w-96
                        rounded-3xl font-semibold md:text-base mobile:text-sm 2k:text-xl 4k:text-2xl'
                    onClick={send}
                    disabled={isSending}
                >
                    {
                        isSending ? (
                            <Loader className='h-3 w-16 2k:h-4 2k:w-24 4k:h-6 4k:w-36'/>
                        ) : t('Send')
                    }
                </Button>
            </div>
        </Dialog>
    )
}

export default PrivateMessageSenderDialog