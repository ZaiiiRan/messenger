/* eslint-disable react/prop-types */
import styles from './Input.module.css'

const Input = ({ className, placeholder, onChange, value, type='text' }) => {
    return (
        <input 
            type={type} 
            placeholder={placeholder} 
            className={styles.Input + ' ' + className}
            value={value}
            onChange={onChange}
        />
    )
}

export default Input
