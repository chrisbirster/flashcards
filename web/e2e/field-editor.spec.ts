import { test, expect } from '@playwright/test'

test.describe('Field Editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck
    await page.fill('input[placeholder="Enter deck name..."]', `Field Editor Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')
  })

  test('shows edit fields button next to note type selector', async ({ page }) => {
    await expect(page.locator('[data-testid="edit-fields-button"]')).toBeVisible()
  })

  test('opens field editor modal when clicking edit fields button', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')

    // Field editor modal should be visible
    await expect(page.locator('text=Edit Fields:')).toBeVisible()
  })

  test('displays existing fields in field editor', async ({ page }) => {
    // Select Basic note type which has Front and Back fields
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-fields-button"]')

    // Should show Front and Back fields
    await expect(page.getByText('Front', { exact: true })).toBeVisible()
    await expect(page.getByText('Back', { exact: true })).toBeVisible()
  })

  test('can add a new field', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')

    // Add a new field
    await page.fill('input[placeholder="New field name..."]', 'MyNewField')
    await page.click('button:has-text("Add")')

    // Should show the new field
    await expect(page.locator('text=MyNewField')).toBeVisible({ timeout: 3000 })
  })
})

test('shows error for reserved field names', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')

  // Try to add a reserved field name
  await page.fill('input[placeholder="New field name..."]', 'Tags')
  await page.click('button:has-text("Add")')

  // Should show error about reserved name
  await expect(page.locator('text=reserved field name')).toBeVisible({ timeout: 3000 })
})

test('can close field editor modal', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')
  await expect(page.locator('text=Edit Fields:')).toBeVisible()

  // Click close button
  await page.click('button:has-text("Close")')

  // Modal should be closed
  await expect(page.locator('text=Edit Fields:')).not.toBeVisible()
})

test('shows move up/down buttons for reordering', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')

  // Should have move buttons (chevron up/down)
  const upButtons = page.locator('button[title="Move up"]')
  const downButtons = page.locator('button[title="Move down"]')

  await expect(upButtons.first()).toBeVisible()
  await expect(downButtons.first()).toBeVisible()
})

test('first field move up button is disabled', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')

  // First field's move up button should be disabled
  const firstUpButton = page.locator('button[title="Move up"]').first()
  await expect(firstUpButton).toBeDisabled()
})

test('shows delete button for fields', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')

  // Should have delete buttons
  const deleteButtons = page.locator('button[title="Remove field"]')
  await expect(deleteButtons.first()).toBeVisible()
})

test('prevents duplicate field names', async ({ page }) => {
  await page.click('[data-testid="edit-fields-button"]')

  // Try to add a duplicate field name (Front already exists)
  await page.fill('input[placeholder="New field name..."]', 'Front')
  await page.click('button:has-text("Add")')

  // Should show error
  await expect(page.locator('text=already exists')).toBeVisible({ timeout: 3000 })
})
