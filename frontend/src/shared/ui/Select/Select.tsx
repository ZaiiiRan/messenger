import { motion, AnimatePresence } from 'framer-motion'
import { useEffect, useState } from 'react'
import styles from './Select.module.css'

interface SelectProps {
    className?: string,
    options: Array<any>,
    defaultValue?: any,
    onChange?: (option: any) => void
}

const Select: React.FC<SelectProps> = ({ options, defaultValue=options[0], onChange, className }) => {
    const [isOpen, setIsOpen] = useState<boolean>(false)
    const [selectedOption, setSelectedOption] = useState<any>(defaultValue)

    const toggleOpen = () => setIsOpen(prev => !prev)

    const handleSelect = (option: any) => {
        setSelectedOption(option)
        setIsOpen(false)

        if (onChange) {
            onChange(option)
        }
    }

    useEffect(() => {
        setSelectedOption(defaultValue)
    }, [defaultValue])

    return (
        <div className={`${styles.SelectContainer} ${className}`}>
            <div onClick={toggleOpen} className={`${styles.selectHeader} rounded-3xl`}>
                <span>{selectedOption.label}</span>
            </div>

        <AnimatePresence>
            {isOpen && (
                <motion.ul
                    className={`${styles.optionsList} rounded-3xl`}
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    transition={{ duration: 0.3 }}
                >
                    {options.map((option, index) => (
                        <motion.li
                            key={index}
                            className={`${styles.optionItem} rounded-3xl`}
                            onClick={() => handleSelect(option)}
                        >
                            {option.label}
                        </motion.li>
                    ))}
                </motion.ul>
            )}
        </AnimatePresence>
        </div>
    )
}

export default Select
