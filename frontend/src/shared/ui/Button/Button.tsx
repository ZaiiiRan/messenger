import styles from './Button.module.css'

interface ButtonProps {
    children?: React.ReactNode,
    className?: string,
    onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void,
    disabled?: boolean,
    id?: string
}

const Button: React.FC<ButtonProps> = ({ children, className, onClick, disabled=false, id }) => {
    return (
        <button 
            className={styles.Button + ' ' + className}
            disabled={disabled}
            onClick={onClick}
            id={id}
        >
            {children}
        </button>
    )
}

export default Button
