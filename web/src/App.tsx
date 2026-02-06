import { BrowserRouter, Routes, Route, Navigate } from 'react-router'
import { Layout } from '#/layouts/root-layout'
import { DecksPage } from '#/pages/DecksPage'
import { StudyPage } from '#/pages/StudyPage'
import { AddNotePage } from '#/pages/AddNotePage'
import { TemplatesPage } from '#/pages/TemplatesPage'
import { EmptyCardsPage } from '#/pages/EmptyCardsPage'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/decks" replace />} />
          <Route path="decks" element={<DecksPage />} />
          <Route path="study/:deckId" element={<StudyPage />} />
          <Route path="notes/add" element={<AddNotePage />} />
          <Route path="templates" element={<TemplatesPage />} />
          <Route path="tools/empty-cards" element={<EmptyCardsPage />} />
          {/* Future routes */}
          {/* <Route path="browse" element={<BrowsePage />} /> */}
          {/* <Route path="notes/:noteId" element={<EditNotePage />} /> */}
          {/* <Route path="card/:cardId" element={<CardPage />} /> */}
          {/* <Route path="templates" element={<TemplatesPage />} /> */}
          {/* <Route path="stats" element={<StatsPage />} /> */}
          {/* <Route path="settings" element={<SettingsPage />} /> */}
          <Route path="*" element={<Navigate to="/decks" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
