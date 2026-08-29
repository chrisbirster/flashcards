export type Deck = {
  id: string;
  name: string;
  card_count: number;
  due_count: number;
  new_count: number;
};

export type Candidate = {
  rating: 1 | 2 | 3 | 4;
  due_at_ms: number;
  interval_days: number;
};

export type Schedule = {
  again: Candidate;
  hard: Candidate;
  good: Candidate;
  easy: Candidate;
};

export type StudyCard = {
  id: string;
  deck_id: string;
  question: string;
  answer: string;
  due_at_ms: number | null;
  review_count: number;
};

export type StudyNext = {
  card: StudyCard | null;
  schedule?: Schedule;
};

type ApiErrorPayload = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

const tokenKey = "plandalf.api-token";

export function savedToken(): string {
  return localStorage.getItem(tokenKey) ?? "";
}

export function saveToken(token: string): void {
  const trimmed = token.trim();
  if (trimmed) localStorage.setItem(tokenKey, trimmed);
  else localStorage.removeItem(tokenKey);
}

async function request<T>(path: string, token: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (init.body) headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  const response = await fetch(path, { ...init, headers });
  if (!response.ok) {
    let payload: ApiErrorPayload = {};
    try {
      payload = (await response.json()) as ApiErrorPayload;
    } catch {
      // Use the HTTP status below when the server did not return JSON.
    }
    throw new ApiError(
      response.status,
      payload.error?.code ?? "request_failed",
      payload.error?.message ?? `Request failed with HTTP ${response.status}`,
    );
  }
  return (await response.json()) as T;
}

export async function listDecks(token: string): Promise<Deck[]> {
  const response = await request<{ decks: Deck[] }>("/api/v1/decks", token);
  return response.decks;
}

export async function nextStudyCard(deckId: string, token: string): Promise<StudyNext> {
  return request<StudyNext>(`/api/v1/decks/${encodeURIComponent(deckId)}/study/next`, token);
}

export async function submitReview(
  cardId: string,
  rating: 1 | 2 | 3 | 4,
  expectedReviewCount: number,
  token: string,
): Promise<void> {
  await request(`/api/v1/cards/${encodeURIComponent(cardId)}/reviews`, token, {
    method: "POST",
    body: JSON.stringify({
      rating,
      expected_review_count: expectedReviewCount,
    }),
  });
}
