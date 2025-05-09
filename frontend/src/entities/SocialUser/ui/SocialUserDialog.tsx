import { useEffect, useState } from 'react'
import { Dialog } from '../../../shared/ui/Dialog'
import SocialUser from './SocialUser'
import IShortUser from '../models/IShortUser'
import { shortUserStore } from '..'
import { observer } from 'mobx-react'

interface SocialUserDialogProps {
    show: boolean,
    setShow: (show: boolean) => void,
    id: string | number,
    onMessageClick?: () => void
}

const SocialUserDialog: React.FC<SocialUserDialogProps> = ({ show, setShow, id, onMessageClick }) => {
    const [user, setUser] = useState<IShortUser | null>()

    useEffect(() => {
        let isMounted = true

        const loadUser = async () => {
            if (isMounted) {
                const user = await shortUserStore.get(id)
                setUser(user)
            }
        }

        loadUser()

        return () => {
            isMounted = false
        }
    }, [])

    return ( 
        <>
            {
                user && (
                    <Dialog
                        show={show}
                        setShow={setShow}
                        title={user.username}
                    >
                        <SocialUser 
                            id={id}
                            onError={() => setShow(false)}
                            onMessageClick={onMessageClick}
                        />
                    </Dialog>
                )
            }
        </>
    )
}

export default observer(SocialUserDialog)
