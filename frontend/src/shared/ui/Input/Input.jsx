/* eslint-disable react/prop-types */
import { useState } from 'react'
import InputMask from 'react-input-mask'
import styles from './Input.module.css'

const Input = ({ className, placeholder, onChange, value, password=false, phone=false, date=false, disabled=false, error }) => {
    const [showPassword, setShowPassword] = useState(password ? false : true)

    const toggleShowPassword = () => setShowPassword((prev) => !prev)

    let mask = null
    if (phone) {
        mask = '+7(999)-999-99-99'
    } else if (date) {
        mask = '99.99.9999'
    }

    return (
        <div className={`${styles.inputWrapper} ${className} ${error ? styles.error : ''}`}>
            { 
                mask ? (
                    <InputMask
                        className={styles.Input + ' ' + className}
                        placeholder={placeholder}
                        mask={mask}
                        value={value}
                        onChange={onChange}
                        maskChar="_"
                    />
                ) : (
                    <input 
                        disabled={disabled}
                        type={showPassword ? 'text' : 'password'} 
                        placeholder={placeholder} 
                        className={styles.Input + ' ' + className}
                        value={value}
                        onChange={onChange}
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
}

export default Input
