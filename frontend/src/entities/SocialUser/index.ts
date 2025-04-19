import { i18n } from '../../shared/config/i18n'
import en from './locales/en.json'
import ru from './locales/ru.json'
i18n.addResourceBundle('en', 'socialUser', en)
i18n.addResourceBundle('ru', 'socialUser', ru)

import SocialUser from './ui/SocialUser'
import ISocialUser from './models/ISocialUser'
import ISocialUserData from './models/ISocialUserData'
export { SocialUser }
export type { ISocialUser, ISocialUserData }