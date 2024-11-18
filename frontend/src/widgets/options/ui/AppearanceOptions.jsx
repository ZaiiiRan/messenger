/* eslint-disable react/prop-types */
import { useTranslation } from 'react-i18next'
import { MainWidget } from '../../../shared/ui/MainWidget'
import { ThemeChanging, LanguageChanging } from '../../../features/options'

const AppearanceOptions = ({ goBack }) => {
    const { t } = useTranslation('optionsWidget')

    return (
        <MainWidget key={'Appearance'} title={ t('Appearance') } goBack={goBack}>
            <div className='flex flex-col gap-8 2k:gap-12 4k:gap-20'>
                <LanguageChanging goBack={goBack} />
                <ThemeChanging />
            </div>
        </MainWidget>
    )
}

export default AppearanceOptions
