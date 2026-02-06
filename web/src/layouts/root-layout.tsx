import { Headersection } from '../components/header/header-section'
import { ContentSection } from '../components/content-section'

export function Layout() {

    return (
        <div className="min-h-screen bg-gray-50">
            <Headersection />
            <ContentSection />
        </div>
    )
}
