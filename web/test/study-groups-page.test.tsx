import { afterEach, describe, expect, it, vi } from "vitest";
import { cleanup, render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter, Route, Routes } from "react-router";
import { StudyGroupDetailPage } from "#/pages/StudyGroupsPage";
import {
  AppRepositoryProvider,
  type AppRepository,
} from "#/lib/app-repository";

afterEach(() => {
  cleanup();
});

function renderStudyGroupDetail(repository: AppRepository) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  render(
    <MemoryRouter initialEntries={["/study-groups/sg_1"]}>
      <QueryClientProvider client={queryClient}>
        <AppRepositoryProvider repository={repository}>
          <Routes>
            <Route path="/study-groups/:groupId" element={<StudyGroupDetailPage />} />
          </Routes>
        </AppRepositoryProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

describe("StudyGroupDetailPage", () => {
  it("renders richer dashboard analytics for the group", async () => {
    const repository = {
      fetchStudyGroup: vi.fn().mockResolvedValue({
        group: {
          id: "sg_1",
          workspaceId: "ws_1",
          primaryDeckId: 1,
          name: "USMLE Cohort",
          description: "Source deck and member installs.",
          visibility: "private",
          joinPolicy: "invite",
          createdByUserId: "usr_owner",
          createdAt: "2026-03-25T12:00:00.000Z",
          updatedAt: "2026-03-25T12:00:00.000Z",
        },
        role: "owner",
        membershipStatus: "active",
        sourceDeckName: "USMLE Source",
        latestVersion: {
          id: "sgv_1",
          studyGroupId: "sg_1",
          versionNumber: 3,
          sourceDeckId: 1,
          publishedByUserId: "usr_owner",
          changeSummary: "Added physiology refresh cards.",
          noteCount: 12,
          cardCount: 30,
          createdAt: "2026-03-25T12:00:00.000Z",
        },
        versions: [],
        members: [
          {
            id: "sgm_1",
            studyGroupId: "sg_1",
            userId: "usr_owner",
            email: "owner@example.com",
            role: "owner",
            status: "active",
            createdAt: "2026-03-25T12:00:00.000Z",
          },
        ],
        currentUserInstall: {
          id: "sgi_1",
          studyGroupId: "sg_1",
          studyGroupMemberId: "sgm_1",
          destinationWorkspaceId: "ws_1",
          installedDeckId: 22,
          installedDeckName: "USMLE Personal Copy",
          sourceVersionNumber: 3,
          status: "active",
          syncState: "clean",
          createdAt: "2026-03-25T12:00:00.000Z",
          updatedAt: "2026-03-25T12:00:00.000Z",
        },
        updateAvailable: false,
        canEdit: true,
        canManageMembers: true,
        canInvite: true,
        canPublishVersion: true,
        dashboard: {
          memberCount: 6,
          activeMembers7d: 4,
          activeInstalls: 5,
          reviewsToday: 18,
          reviews7d: 121,
          sessions7d: 9,
          minutesStudied7d: 142,
          latestVersionNumber: 3,
          latestVersionAdoption: 3,
          latestVersionAdoptionPercent: 60,
          dailyActivity: [
            { date: "2026-03-19", sessions: 1, cardsReviewed: 11, minutesStudied: 15 },
            { date: "2026-03-20", sessions: 0, cardsReviewed: 0, minutesStudied: 0 },
            { date: "2026-03-21", sessions: 2, cardsReviewed: 34, minutesStudied: 40 },
            { date: "2026-03-22", sessions: 1, cardsReviewed: 8, minutesStudied: 12 },
            { date: "2026-03-23", sessions: 2, cardsReviewed: 29, minutesStudied: 35 },
            { date: "2026-03-24", sessions: 1, cardsReviewed: 17, minutesStudied: 20 },
            { date: "2026-03-25", sessions: 2, cardsReviewed: 22, minutesStudied: 20 },
          ],
          leaderboard: [
            {
              email: "owner@example.com",
              displayName: "Owner",
              reviews7d: 48,
              sessions7d: 3,
              minutes7d: 54,
            },
          ],
        },
        recentEvents: [],
        availableWorkspaces: [
          {
            id: "ws_1",
            name: "Owner Workspace",
            slug: "owner-workspace",
            collectionId: "col_1",
          },
        ],
      }),
    } as unknown as AppRepository;

    renderStudyGroupDetail(repository);

    expect(await screen.findByText("Weekly activity")).toBeInTheDocument();
    expect(await screen.findByText("Reviews today")).toBeInTheDocument();
    expect(await screen.findByText("Active installs")).toBeInTheDocument();
    expect(await screen.findByText("Minutes 7d")).toBeInTheDocument();
    expect(await screen.findByText("(60%)")).toBeInTheDocument();
    expect(await screen.findByText(/3 sessions/i)).toBeInTheDocument();
    expect(await screen.findByText(/54 min/i)).toBeInTheDocument();
  });
});
