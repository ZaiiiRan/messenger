import { useTranslation } from 'react-i18next'
import { MainWidget } from '../../../shared/ui/MainWidget'
import { ThemeChanging, LanguageChanging } from '../../../features/options'

interface AppearanceOptionsProps {
    goBack: () => void,
}

const AppearanceOptions: React.FC<AppearanceOptionsProps> = ({ goBack }) => {
    const { t } = useTranslation('optionsWidget')

    return (
        <MainWidget key={'Appearance'} title={ t('Appearance') } goBack={goBack}>
            <div className='flex flex-col gap-8 2k:gap-12 4k:gap-20'>
                <LanguageChanging />
                <ThemeChanging />
            </div>
        </MainWidget>
    )
}

export default AppearanceOptions
