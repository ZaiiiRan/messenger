import styles from './LinkButton.module.css'

interface LinkButtonProps {
    children?: React.ReactNode,
    className?: string,
    onClick?: (event: React.MouseEvent<HTMLAnchorElement>) => void,
}

const LinkButton: React.FC<LinkButtonProps> = ({ onClick, children, className }) => {
    return (
        <a 
            className={styles.LinkButton + ' ' + className}
            onClick={onClick}
        >
            {children}
        </a>
    )
}

export default LinkButton
