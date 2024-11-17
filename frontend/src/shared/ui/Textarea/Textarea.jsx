/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable react-hooks/rules-of-hooks */
/* eslint-disable react/prop-types */
import { forwardRef, useRef, useEffect } from 'react'
import styles from './Textarea.module.css'

const Textarea = forwardRef(({ className, placeholder, onChange, value, disabled = false, maxRows = 4, ...props }, ref) => {
    const textareaRef = ref || useRef()

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
                {...props}
            />
        </div>
    )
})

Textarea.displayName = 'Textarea'

export default Textarea
