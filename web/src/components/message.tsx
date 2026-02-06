import { forwardRef, type JSX } from "react"

type MessageType = 'success' | 'error'
type MessageProps = JSX.IntrinsicElements['div'] & {
    messageType: MessageType
    message: string,
}

const Message = forwardRef<HTMLDivElement, MessageProps>(function SuccessMessage(
    { messageType, message, ...props }, ref) {
    const messageColor = messageType === 'success' ? 'green' : 'red';
    const cn = `p-4 bg-${messageColor}-50 border border-${messageColor}-200 rounded-lg`;
    return (
        <div
            className={cn}
            ref={ref}
            {...props}
        >
            <p className={`text-${messageColor}-700`}>
                {message}
            </p>
        </div>
    )
})

function SuccessMessage() {
    return (
        <Message messageType="success" message="Note added successfully! Add another or close." />
    )
}

function ErrorMessage({ message }: { message?: string | undefined }) {
    const defaultMessage = 'Failed to create note. Please try again.';
    return (
        <Message messageType="error" message={message ? message : defaultMessage} />
    )
}

export {
    SuccessMessage,
    ErrorMessage,
}
