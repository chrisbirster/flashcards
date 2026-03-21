import { useQuery } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'

export function StudyGroupsPage() {
  const repository = useAppRepository()
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })

  const canCreate = entitlementsQuery.data?.features.studyGroups ?? false

  return (
    <div className="mx-auto max-w-4xl">
      <div className="rounded-[2rem] border border-slate-200 bg-white p-6 shadow-sm md:p-8">
        <p className="text-xs uppercase tracking-[0.28em] text-slate-400">Coming soon</p>
        <h2 className="mt-4 text-3xl font-semibold tracking-tight text-slate-950">Study Groups</h2>
        <p className="mt-4 max-w-2xl text-sm leading-7 text-slate-500">
          Study Groups will let you attach a collaborative space to a primary deck, invite members, and manage deck-specific participation without mixing that workflow into basic deck CRUD.
        </p>

        <div className="mt-8 grid gap-4 md:grid-cols-2">
          <div className="rounded-2xl bg-stone-50 p-5">
            <h3 className="text-lg font-semibold text-slate-950">What lands in the next tranche</h3>
            <ul className="mt-4 space-y-2 text-sm text-slate-600">
              <li>Deck-linked study groups with owner, admin, and member roles.</li>
              <li>Invite and remove members with explicit membership status.</li>
              <li>Team and Enterprise plan gating for group creation.</li>
            </ul>
          </div>

          <div className="rounded-2xl border border-slate-200 p-5">
            <p className="text-xs uppercase tracking-[0.18em] text-slate-400">Plan access</p>
            <p className="mt-3 text-2xl font-semibold text-slate-950">
              {canCreate ? 'Eligible to create later' : 'Upgrade required'}
            </p>
            <p className="mt-3 text-sm text-slate-500">
              {canCreate
                ? 'This workspace can create study groups once the feature is shipped.'
                : 'Study group creation is reserved for Team and Enterprise workspaces.'}
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
