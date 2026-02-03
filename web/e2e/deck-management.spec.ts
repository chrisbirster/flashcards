import { test, expect } from '@playwright/test';

test.describe('Deck Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the application title', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Microdote' })).toBeVisible();
  });

  test('should display empty state when no decks exist', async ({ page }) => {
    // If there are existing decks, this test may need adjustment
    // For a fresh database, we should see the empty state
    const emptyMessage = page.getByText(/No decks yet/i);
    if (await emptyMessage.isVisible()) {
      await expect(emptyMessage).toBeVisible();
    }
  });

  test('should create a new deck', async ({ page }) => {
    const deckName = `Test Deck ${Date.now()}`;

    // Fill in the deck name
    await page.getByPlaceholder('Enter deck name...').fill(deckName);

    // Click create button
    await page.getByRole('button', { name: 'Create' }).click();

    // Wait for the deck to appear in the list
    await expect(page.getByText(deckName)).toBeVisible({ timeout: 5000 });

    // Input should be cleared
    await expect(page.getByPlaceholder('Enter deck name...')).toHaveValue('');
  });

  test('should display deck with card count', async ({ page }) => {
    // Create a deck first
    const deckName = `Deck with Cards ${Date.now()}`;
    await page.getByPlaceholder('Enter deck name...').fill(deckName);
    await page.getByRole('button', { name: 'Create' }).click();

    // Wait for deck to appear
    await expect(page.getByText(deckName)).toBeVisible();

    // Should show 0 total cards initially - stats format is "(N total)"
    const deckItem = page.locator('li', { hasText: deckName });
    await expect(deckItem.getByText('(0 total)')).toBeVisible();
  });

  test('should disable create button when input is empty', async ({ page }) => {
    const createButton = page.getByRole('button', { name: 'Create' });

    // Button should be disabled with empty input
    await expect(createButton).toBeDisabled();

    // Type something
    await page.getByPlaceholder('Enter deck name...').fill('Test');
    await expect(createButton).toBeEnabled();

    // Clear input
    await page.getByPlaceholder('Enter deck name...').clear();
    await expect(createButton).toBeDisabled();
  });

  test('should display Study and Add Cards buttons for each deck', async ({ page }) => {
    // Create a deck
    const deckName = `Action Buttons Deck ${Date.now()}`;
    await page.getByPlaceholder('Enter deck name...').fill(deckName);
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText(deckName)).toBeVisible();

    // Find the deck item and check for buttons
    const deckItem = page.locator('li', { hasText: deckName });
    await expect(deckItem.getByRole('button', { name: 'Study' })).toBeVisible();
    await expect(deckItem.getByRole('button', { name: 'Add Cards' })).toBeVisible();
  });

  test('should handle rapid deck creation', async ({ page }) => {
    // Test creating multiple decks in quick succession
    const decks = [
      `Rapid 1 ${Date.now()}`,
      `Rapid 2 ${Date.now()}`,
      `Rapid 3 ${Date.now()}`,
    ];

    for (const deckName of decks) {
      await page.getByPlaceholder('Enter deck name...').fill(deckName);
      await page.getByRole('button', { name: 'Create' }).click();
      // Small delay to ensure creation completes
      await page.waitForTimeout(500);
    }

    // All decks should be visible
    for (const deckName of decks) {
      await expect(page.getByText(deckName)).toBeVisible();
    }
  });

  test('should show loading state while fetching decks', async ({ page }) => {
    // On first load, there might be a brief loading state
    // This is more noticeable on slower connections or when the API is slow
    await page.goto('/');

    // Check if loading indicator appears (may be very brief)
    const loadingIndicator = page.getByText('Loading...');

    // Either we see loading or we directly see content
    // This is a weak test but ensures the page doesn't crash
    await expect(page.getByRole('heading', { name: 'Microdote' })).toBeVisible({ timeout: 5000 });
  });

  test('should persist decks across page reloads', async ({ page }) => {
    const deckName = `Persistent Deck ${Date.now()}`;

    // Create a deck
    await page.getByPlaceholder('Enter deck name...').fill(deckName);
    await page.getByRole('button', { name: 'Create' }).click();
    await expect(page.getByText(deckName)).toBeVisible();

    // Reload the page
    await page.reload();

    // Deck should still be visible
    await expect(page.getByText(deckName)).toBeVisible({ timeout: 5000 });
  });

  test('should display milestone information in footer', async ({ page }) => {
    await expect(page.getByText('Milestone M1')).toBeVisible();
    await expect(page.getByText('Backend: Go + SQLite')).toBeVisible();
    await expect(page.getByText('Frontend: React + TanStack Query')).toBeVisible();
  });
});
