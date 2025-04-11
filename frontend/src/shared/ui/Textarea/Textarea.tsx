import { forwardRef, useRef, useEffect } from 'react'
import styles from './Textarea.module.css'

interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
    className?: string,
    placeholder?: string,
    onChange?: (event: React.ChangeEvent<HTMLTextAreaElement>) => void,
    value?: string,
    disabled?: boolean,
    maxRows?: number,
}

const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(({ className, placeholder, onChange, value, disabled = false, maxRows = 4, ...restProps }, ref) => {
    const textareaRef: any = ref || useRef()

    useEffect(() => {
        if (textareaRef.current) {
            autoResize()
        }
    }, [value])

    const autoResize = () => {
        const textarea = textareaRef.current
        textarea.style.height = 'auto'
        const computedHeight = textarea.scrollHeight

        const lineHeight = parseInt(getComputedStyle(textarea).lineHeight, 10) || 20
        const maxHeight = lineHeight * maxRows
        textarea.style.height = `${Math.min(computedHeight, maxHeight)}px`
    };

    return (
        <div className={`${styles.textareaWrapper} ${className}`}>
            <textarea
                ref={textareaRef}
                disabled={disabled}
                placeholder={placeholder}
                className={`${styles.Textarea}`}
                value={value}
                onChange={(e) => {
                    onChange?.(e);
                    autoResize();
                }}
                {...restProps}
            />
        </div>
    )
})

Textarea.displayName = 'Textarea'

export default Textarea
