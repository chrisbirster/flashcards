import { createSignal, onSettled } from "solid-js";
import { For, Show } from "@solidjs/web";
import * as stylex from "@stylexjs/stylex";
import {
  ApiError,
  type Candidate,
  type Deck,
  type StudyNext,
  listDecks,
  nextStudyCard,
  savedToken,
  saveToken,
  submitReview,
} from "./api";
import { Button, Stack, Surface, Text } from "./ui";

function formatInterval(candidate: Candidate): string {
  const days = candidate.interval_days;
  if (days < 1 / 24) return `${Math.max(1, Math.round(days * 24 * 60))}m`;
  if (days < 1) return `${Math.max(1, Math.round(days * 24))}h`;
  if (days < 60) return `${Math.max(1, Math.round(days))}d`;
  if (days < 365) return `${Math.max(1, Math.round(days / 30.44))}mo`;
  return `${(days / 365).toFixed(1)}y`;
}

export default function App() {
  const [token, setToken] = createSignal(savedToken());
  const [tokenDraft, setTokenDraft] = createSignal(savedToken());
  const [needsToken, setNeedsToken] = createSignal(false);
  const [decks, setDecks] = createSignal<Deck[]>([]);
  const [selectedDeck, setSelectedDeck] = createSignal<Deck | null>(null);
  const [study, setStudy] = createSignal<StudyNext | null>(null);
  const [revealed, setRevealed] = createSignal(false);
  const [busy, setBusy] = createSignal(false);
  const [error, setError] = createSignal("");

  async function handleApiError(reason: unknown) {
    if (reason instanceof ApiError && reason.status === 401) {
      setNeedsToken(true);
      setError("");
      return;
    }
    setError(reason instanceof Error ? reason.message : "Something went wrong.");
  }

  async function refreshDecks() {
    setBusy(true);
    setError("");
    try {
      const next = await listDecks(token());
      setDecks(next);
      setNeedsToken(false);
    } catch (reason) {
      await handleApiError(reason);
    } finally {
      setBusy(false);
    }
  }

  async function loadNext(deck: Deck) {
    setBusy(true);
    setError("");
    try {
      const next = await nextStudyCard(deck.id, token());
      setStudy(next);
      setRevealed(false);
    } catch (reason) {
      await handleApiError(reason);
    } finally {
      setBusy(false);
    }
  }

  async function startStudy(deck: Deck) {
    setSelectedDeck(deck);
    await loadNext(deck);
  }

  async function rate(rating: 1 | 2 | 3 | 4) {
    const current = study()?.card;
    const deck = selectedDeck();
    if (!current || !deck) return;

    setBusy(true);
    setError("");
    try {
      await submitReview(current.id, rating, current.review_count, token());
      const next = await nextStudyCard(deck.id, token());
      setStudy(next);
      setRevealed(false);
    } catch (reason) {
      await handleApiError(reason);
    } finally {
      setBusy(false);
    }
  }

  async function useToken() {
    const nextToken = tokenDraft().trim();
    saveToken(nextToken);
    setToken(nextToken);
    setNeedsToken(false);
    await refreshDecks();
  }

  function leaveStudy() {
    setSelectedDeck(null);
    setStudy(null);
    setRevealed(false);
    void refreshDecks();
  }

  onSettled(() => {
    void refreshDecks();
  });

  return (
    <main {...stylex.props(styles.app)}>
      <div {...stylex.props(styles.shell)}>
        <header {...stylex.props(styles.header)}>
          <Stack gap="xs">
            <Text tone="accent" size="sm" weight="bold">PLANDALF</Text>
            <Text size="xl" weight="bold">Study what matters.</Text>
          </Stack>
          <Show when={selectedDeck()}>
            <Button tone="secondary" onClick={leaveStudy}>Decks</Button>
          </Show>
        </header>

        <Show when={error()}>
          <div {...stylex.props(styles.error)}>
            <Text size="sm">{error()}</Text>
          </div>
        </Show>

        <Show when={needsToken()}>
          <Surface>
            <Stack gap="md">
              <Stack gap="xs">
                <Text size="lg" weight="bold">Connect to Plandalf</Text>
                <Text tone="muted" size="sm">Enter the API token configured on your hosted Plandalf server.</Text>
              </Stack>
              <input
                type="password"
                autocomplete="current-password"
                value={tokenDraft()}
                onInput={(event) => setTokenDraft(event.currentTarget.value)}
                placeholder="API token"
                {...stylex.props(styles.input)}
              />
              <Button wide disabled={busy()} onClick={() => void useToken()}>Connect</Button>
            </Stack>
          </Surface>
        </Show>

        <Show when={!needsToken() && !selectedDeck()}>
          <section {...stylex.props(styles.section)}>
            <Stack gap="md">
              <Stack gap="xs">
                <Text size="lg" weight="bold">Your decks</Text>
                <Text tone="muted" size="sm">Due reviews are shown first. New cards remain available after them.</Text>
              </Stack>

              <Show when={!busy() && decks().length === 0}>
                <Surface>
                  <Stack gap="sm">
                    <Text weight="bold">No decks yet</Text>
                    <Text tone="muted" size="sm">The API is connected. Add or sync a deck to this Turso database to begin studying.</Text>
                  </Stack>
                </Surface>
              </Show>

              <div {...stylex.props(styles.deckGrid)}>
                <For each={decks()}>
                  {(deck) => (
                    <Surface>
                      <Stack gap="md">
                        <Stack gap="xs">
                          <Text size="lg" weight="bold">{deck.name}</Text>
                          <Text tone="muted" size="sm">{deck.card_count} cards</Text>
                        </Stack>
                        <div {...stylex.props(styles.stats)}>
                          <div {...stylex.props(styles.stat)}>
                            <Text size="lg" weight="bold">{deck.due_count}</Text>
                            <Text tone="muted" size="sm">Due</Text>
                          </div>
                          <div {...stylex.props(styles.stat)}>
                            <Text size="lg" weight="bold">{deck.new_count}</Text>
                            <Text tone="muted" size="sm">New</Text>
                          </div>
                        </div>
                        <Button wide disabled={busy() || deck.due_count + deck.new_count === 0} onClick={() => void startStudy(deck)}>
                          {deck.due_count + deck.new_count === 0 ? "Caught up" : "Study"}
                        </Button>
                      </Stack>
                    </Surface>
                  )}
                </For>
              </div>
            </Stack>
          </section>
        </Show>

        <Show when={!needsToken() && selectedDeck()}>
          {(deck) => (
            <section {...stylex.props(styles.studySection)}>
              <Stack gap="md">
                <Stack gap="xs">
                  <Text tone="muted" size="sm">{deck().name}</Text>
                  <Show when={study()?.card} fallback={
                    <Surface>
                      <Stack gap="md" align="center">
                        <Text size="lg" weight="bold">You are caught up.</Text>
                        <Text tone="muted" size="sm">There are no due or new cards available in this deck right now.</Text>
                        <Button tone="secondary" onClick={leaveStudy}>Back to decks</Button>
                      </Stack>
                    </Surface>
                  }>
                    {(card) => (
                      <Surface>
                        <div {...stylex.props(styles.studyCard)}>
                          <Stack gap="lg">
                            <Stack gap="sm">
                              <Text tone="muted" size="sm">QUESTION</Text>
                              <Text size="xl" weight="medium">{card().question}</Text>
                            </Stack>

                            <Show when={revealed()} fallback={
                              <Button wide disabled={busy()} onClick={() => setRevealed(true)}>Show answer</Button>
                            }>
                              <Stack gap="lg">
                                <div {...stylex.props(styles.answer)}>
                                  <Stack gap="sm">
                                    <Text tone="muted" size="sm">ANSWER</Text>
                                    <Text size="lg">{card().answer}</Text>
                                  </Stack>
                                </div>

                                <Show when={study()?.schedule}>
                                  {(schedule) => (
                                    <div {...stylex.props(styles.ratingGrid)}>
                                      <Button tone="danger" disabled={busy()} onClick={() => void rate(1)}>
                                        <Stack gap="xs" align="center"><Text weight="bold">Again</Text><Text size="sm">{formatInterval(schedule().again)}</Text></Stack>
                                      </Button>
                                      <Button tone="secondary" disabled={busy()} onClick={() => void rate(2)}>
                                        <Stack gap="xs" align="center"><Text weight="bold">Hard</Text><Text size="sm">{formatInterval(schedule().hard)}</Text></Stack>
                                      </Button>
                                      <Button tone="secondary" disabled={busy()} onClick={() => void rate(3)}>
                                        <Stack gap="xs" align="center"><Text weight="bold">Good</Text><Text size="sm">{formatInterval(schedule().good)}</Text></Stack>
                                      </Button>
                                      <Button disabled={busy()} onClick={() => void rate(4)}>
                                        <Stack gap="xs" align="center"><Text weight="bold">Easy</Text><Text size="sm">{formatInterval(schedule().easy)}</Text></Stack>
                                      </Button>
                                    </div>
                                  )}
                                </Show>
                              </Stack>
                            </Show>
                          </Stack>
                        </div>
                      </Surface>
                    )}
                  </Show>
                </Stack>
              </Stack>
            </section>
          )}
        </Show>

        <Show when={busy() && !needsToken()}>
          <div {...stylex.props(styles.busy)}><Text tone="muted" size="sm">Working…</Text></div>
        </Show>
      </div>
    </main>
  );
}

const styles = stylex.create({
  app: {
    backgroundColor: "#11110f",
    color: "#f5f1e8",
    minHeight: "100vh",
    paddingBlock: 24,
    paddingInline: 16,
  },
  shell: {
    marginInline: "auto",
    maxWidth: 760,
    width: "100%",
  },
  header: {
    alignItems: "center",
    display: "flex",
    justifyContent: "space-between",
    marginBottom: 24,
  },
  section: { paddingBottom: 32 },
  studySection: { paddingBottom: 48 },
  error: {
    backgroundColor: "#3a211f",
    borderColor: "#69423d",
    borderRadius: 14,
    borderStyle: "solid",
    borderWidth: 1,
    marginBottom: 16,
    padding: 12,
  },
  input: {
    backgroundColor: "#141411",
    borderColor: {
      default: "#424038",
      ":focus": "#d9bf77",
    },
    borderRadius: 14,
    borderStyle: "solid",
    borderWidth: 1,
    color: "#f5f1e8",
    fontFamily: "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, sans-serif",
    fontSize: 16,
    minHeight: 48,
    outline: "none",
    paddingInline: 14,
  },
  deckGrid: {
    display: "grid",
    gap: 14,
    gridTemplateColumns: {
      default: "1fr",
      "@media (min-width: 680px)": "repeat(2, minmax(0, 1fr))",
    },
  },
  stats: {
    display: "grid",
    gap: 10,
    gridTemplateColumns: "1fr 1fr",
  },
  stat: {
    backgroundColor: "#151512",
    borderRadius: 14,
    display: "flex",
    flexDirection: "column",
    gap: 2,
    padding: 12,
  },
  studyCard: {
    minHeight: {
      default: 420,
      "@media (min-width: 680px)": 480,
    },
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
  },
  answer: {
    borderTopColor: "#35332d",
    borderTopStyle: "solid",
    borderTopWidth: 1,
    paddingTop: 20,
  },
  ratingGrid: {
    display: "grid",
    gap: 10,
    gridTemplateColumns: {
      default: "1fr 1fr",
      "@media (min-width: 680px)": "repeat(4, minmax(0, 1fr))",
    },
  },
  busy: {
    display: "flex",
    justifyContent: "center",
    paddingBlock: 16,
  },
});
