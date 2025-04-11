import { useState, forwardRef } from 'react'
import InputMask from 'react-input-mask'
import styles from './Input.module.css'

interface InputProps {
    className?: string,
    placeholder?: string,
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void,
    value?: string,
    password?: boolean,
    phone?: boolean,
    date?: boolean,
    disabled?: boolean,
    oneDigit?: boolean,
    error?: boolean,
    name?: string,
    onKeyDown?: (event: React.KeyboardEvent<HTMLInputElement>) => void,
}

const Input = forwardRef<HTMLInputElement, InputProps>(({ className, placeholder, onChange, value, password = false, phone = false, 
date = false, disabled = false, oneDigit = false, error, name, onKeyDown }, ref) => {
    const [showPassword, setShowPassword] = useState<boolean>(password ? false : true)

    const toggleShowPassword = () => setShowPassword((prev) => !prev)

    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (oneDigit && !/[0-9]/.test(e.key)) {
            e.preventDefault()
        }
    }

    let mask = null
    if (phone) {
        mask = '+7(999)-999-99-99'
    } else if (date) {
        mask = '99.99.9999'
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault()
        if (onChange) {
            onChange(e)
        }
    }

    return (
        <div className={`${styles.inputWrapper} ${className} ${error ? styles.error : ''}`}>
            { 
                mask ? (
                    <InputMask
                        name={name}
                        disabled={disabled}
                        mask={mask}
                        value={value}
                        onChange={handleInputChange}
                        maskChar="_"
                        onKeyDown={onKeyDown}
                        onInput={handleInputChange}
                    >
                        {() => <input
                            type="tel"
                            placeholder={placeholder}
                            className={`${styles.Input} ${className}`}
                            ref={ref}
                        />}
                    </InputMask>
                ) : (
                    <input 
                        ref={ref}
                        name={name}
                        disabled={disabled}
                        type={showPassword ? (oneDigit ? 'tel' : 'text') : 'password'} 
                        placeholder={placeholder} 
                        className={`${styles.Input} ${className}`}
                        value={value}
                        onChange={onChange}
                        onKeyDown={onKeyDown}
                        onKeyPress={handleKeyPress}
                        maxLength={oneDigit ? 1 : undefined}
                    />
                )
            }

            { password && ( 
                <span onClick={toggleShowPassword} className={styles.eyeIcon}>
                    {showPassword ? '🙈' : '👁️'}
                </span>
            )}
        </div>
    )
})

Input.displayName = 'Input'

export default Input

