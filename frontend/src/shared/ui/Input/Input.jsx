/* eslint-disable react/prop-types */
import { useState } from 'react'
import styles from './Input.module.css'

const Input = ({ className, placeholder, onChange, value, password=false }) => {
    const [showPassword, setShowPassword] = useState(password ? false : true)

    const toggleShowPassword = () => setShowPassword((prev) => !prev)

    return (
        <div className={`${styles.inputWrapper} ${className}`}>
            <input 
                type={showPassword ? 'text' : 'password'} 
                placeholder={placeholder} 
                className={styles.Input + ' ' + className}
                value={value}
                onChange={onChange}
            />
            { password && ( 
                <span onClick={toggleShowPassword} className={styles.eyeIcon}>
                    {showPassword ? '🙈' : '👁️'}
                </span>
            )}
        </div>
        
    )
}

export default Input
