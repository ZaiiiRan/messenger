interface ChatWidgetHeaderProps {
    goBack: () => void,
    chatName: string,
    isGroupChat: boolean,
    openProperties: () => void
}

const ChatWidgetHeader: React.FC<ChatWidgetHeaderProps> = ({ goBack, chatName, isGroupChat, openProperties }) => {
    return (
        <>
            <div className='flex items-center gap-5 2k:gap-7 4k:gap-9'>
                <div 
                    className='backBtn 2xl:w-10 2xl:h-10 xl:w-9 xl:h-9 lg:w-9 lg:h-8 2k:w-12 2k:h-12 4k:w-14 4k:h-14 
                        mobile:w-8 mobile:h-8 md:w-9 md:h-9 
                        rounded-3xl flex items-center justify-center'
                    onClick={goBack}
                >
                    <div className='Icon flex items-center justify-center h-1/4 aspect-square'>
                        <svg viewBox="0 0 7.424 12.02" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <defs/>
                            <path id="Vector" d="M0 6.01L6 12.02L7.42 10.6L2.82 6L7.42 1.4L6 0L0 6.01Z" fillOpacity="1.000000" fillRule="nonzero"/>
                        </svg>
                    </div>
                </div>
                <div 
                    className='ChatHeader flex items-center h-full gap-5 py-3 px-4 2k:py-5 2k:px-6 2k:h-24 4k:py-7 4k:px-8 4k:h-32 rounded-3xl max-w-[90%]'
                    onClick={openProperties}
                >
                    {/* Avatar */}
                    <div className='h-full rounded-3xl aspect-square'>
                        <div className='flex items-center justify-center w-full h-full Avatar-standart xl:rounded-3xl lg:rounded-2xl mobile:rounded-2xl md:rounded-3xl'>
                            <div className='flex items-center justify-center w-1/2 aspect-square'>
                                {
                                    isGroupChat ? (
                                        <svg viewBox="0 0 26.6666 24" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="coolicon" d="M2.66 6.66C2.66 2.98 5.65 0 9.33 0C13.01 0 16 2.98 16 6.66C16 10.34 13.01 13.33 9.33 13.33C5.65 13.33 2.66 10.34 2.66 6.66ZM9.33 2.66C7.12 2.66 5.33 4.45 5.33 6.66C5.33 8.87 7.12 10.66 9.33 10.66C11.54 10.66 13.33 8.87 13.33 6.66C13.33 4.45 11.54 2.66 9.33 2.66ZM18.66 6.66C19.08 6.66 19.5 6.76 19.87 6.95C20.25 7.14 20.57 7.42 20.82 7.76C21.07 8.1 21.23 8.5 21.3 8.91C21.36 9.33 21.33 9.75 21.2 10.16C21.07 10.56 20.84 10.92 20.54 11.22C20.25 11.51 19.88 11.74 19.48 11.87C19.08 12 18.66 12.03 18.24 11.96C17.82 11.89 17.43 11.73 17.09 11.48L15.52 13.64L15.52 13.64C16.2 14.13 16.99 14.46 17.82 14.59C17.91 14.61 18 14.62 18.09 14.63C18.84 14.71 19.59 14.63 20.3 14.4C21.1 14.14 21.83 13.7 22.43 13.1C23.02 12.51 23.47 11.78 23.73 10.98C23.99 10.18 24.06 9.33 23.93 8.5C23.8 7.67 23.47 6.88 22.98 6.2C22.54 5.59 21.98 5.08 21.33 4.71C21.25 4.66 21.17 4.62 21.08 4.58C20.33 4.19 19.5 4 18.66 4L18.66 6.66ZM16 24L18.66 24C18.66 18.84 14.48 14.66 9.33 14.66C4.17 14.66 0 18.84 0 24L2.66 24C2.66 20.31 5.65 17.33 9.33 17.33C13.01 17.33 16 20.31 16 24ZM23.99 24C23.99 23.3 23.85 22.6 23.59 21.96C23.32 21.31 22.93 20.72 22.43 20.23C21.94 19.73 21.35 19.34 20.7 19.07C20.05 18.8 19.36 18.66 18.66 18.66L18.66 16C19.57 16 20.47 16.15 21.33 16.45C21.46 16.5 21.59 16.55 21.72 16.6C22.69 17.01 23.58 17.6 24.32 18.34C25.06 19.08 25.65 19.96 26.05 20.93C26.11 21.06 26.16 21.2 26.2 21.33C26.51 22.18 26.66 23.09 26.66 24L23.99 24Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    ) : (
                                        <svg viewBox="0 0 16 19" xmlns="http://www.w3.org/2000/svg">
                                            <defs/>
                                            <path id="Vector" d="M8 0C5.23 0 3 2.23 3 5C3 7.76 5.23 10 8 10C10.76 10 13 7.76 13 5C13 2.23 10.76 0 8 0ZM11 5C11 6.65 9.65 8 8 8C6.34 8 5 6.65 5 5C5 3.34 6.34 2 8 2C9.65 2 11 3.34 11 5ZM0 19C0 16.87 0.84 14.84 2.34 13.34C3.84 11.84 5.87 11 8 11C10.12 11 12.15 11.84 13.65 13.34C15.15 14.84 16 16.87 16 19L14 19C14 17.4 13.36 15.88 12.24 14.75C11.11 13.63 9.59 13 8 13C6.4 13 4.88 13.63 3.75 14.75C2.63 15.88 2 17.4 2 19L0 19Z" fill="#0F1828" fillOpacity="1.000000" fillRule="evenodd"/>
                                        </svg>
                                    )
                                }
                            </div>
                        </div>
                    </div>
                    {/* Title */}
                    <h1 className='font-extrabold 2xl:text-3xl xl:text-2xl lg:text-xl 2k:text-4xl 4k:text-5xl
                        md:text-3xl sm:text-2xl mobile:text-xl whitespace-nowrap text-ellipsis overflow-hidden max-w-[80%]'
                    >
                        { chatName }
                    </h1>
                </div>
            </div>
        </>
    )
    {/* Back btn */}
}

export default ChatWidgetHeader