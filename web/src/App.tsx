import { BrowserRouter, Routes, Route, Navigate, Outlet, useLocation } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import { Layout } from '#/layouts/root-layout'
import { DecksPage } from '#/pages/DecksPage'
import { StudyPage } from '#/pages/StudyPage'
import { AddNotePage } from '#/pages/AddNotePage'
import { TemplatesPage } from '#/pages/TemplatesPage'
import { EmptyCardsPage } from '#/pages/EmptyCardsPage'
import { LoginPage } from '#/pages/LoginPage'
import { useAppRepository } from '#/lib/app-repository'
import {
  AddNoteFieldEditorRoutePage,
  AddNoteTemplateEditorRoutePage,
  TemplatesTemplateEditorRoutePage,
} from '#/pages/NoteTypeEditorRoutes'

function RequireAuthLayout() {
  const repository = useAppRepository()
  const location = useLocation()
  const sessionQuery = useQuery({
    queryKey: ['auth-session'],
    queryFn: () => repository.fetchSession(),
  })

  if (sessionQuery.isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 text-sm text-gray-500">
        Checking session...
      </div>
    )
  }

  if (!sessionQuery.data?.authenticated) {
    return <Navigate to="/login" replace state={{from: `${location.pathname}${location.search}`}} />
  }

  return (
    <Layout>
      <Outlet />
    </Layout>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<RequireAuthLayout />}>
          <Route index element={<Navigate to="/decks" replace />} />
          <Route path="decks" element={<DecksPage />} />
          <Route path="study/:deckId" element={<StudyPage />} />
          <Route path="notes/add" element={<AddNotePage />}>
            <Route path="note-types/:noteTypeName/fields" element={<AddNoteFieldEditorRoutePage />} />
            <Route path="note-types/:noteTypeName/templates" element={<AddNoteTemplateEditorRoutePage />} />
          </Route>
          <Route path="templates" element={<TemplatesPage />}>
            <Route path=":noteTypeName" element={<TemplatesTemplateEditorRoutePage />} />
          </Route>
          <Route path="tools/empty-cards" element={<EmptyCardsPage />} />
          {/* Future routes */}
          {/* <Route path="browse" element={<BrowsePage />} /> */}
          {/* <Route path="notes/:noteId" element={<EditNotePage />} /> */}
          {/* <Route path="card/:cardId" element={<CardPage />} /> */}
          {/* <Route path="templates" element={<TemplatesPage />} /> */}
          {/* <Route path="stats" element={<StatsPage />} /> */}
          {/* <Route path="settings" element={<SettingsPage />} /> */}
        </Route>
        <Route path="*" element={<Navigate to="/decks" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
