import { useSearchParams, useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { AddNoteScreen } from '#/components/AddNoteScreen'

export function AddNotePage() {
    const [searchParams] = useSearchParams()
    const navigate = useNavigate()
    const queryClient = useQueryClient()

    const deckId = searchParams.get('deckId')
    const onSuccess = () => {
        if (deckId) {
            queryClient.invalidateQueries({ queryKey: ['deck-stats', Number(deckId)] })
        }
        navigate(-1)
    }

    return (
        <AddNoteScreen
            deckId={deckId ? Number(deckId) : undefined}
            onClose={() => navigate(-1)}
            onSuccess={onSuccess}
        />
    )
}
