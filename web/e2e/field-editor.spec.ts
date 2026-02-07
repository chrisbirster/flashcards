import { test, expect } from '@playwright/test'

test.describe('Field Editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck
    await page.fill('input[placeholder="Deck name"]', `Field Editor Test ${Date.now()}`)
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

    // Click close button in the modal footer
    await page.locator('.bg-gray-50 button:has-text("Close")').click()

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
})

test.describe('Field Options', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck
    await page.fill('input[placeholder="Deck name"]', `Field Options Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')

    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })

    // Open field editor
    await page.click('[data-testid="edit-fields-button"]')
  })

  test('shows field options button for each field', async ({ page }) => {
    await expect(page.locator('[data-testid="field-options-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-options-Back"]')).toBeVisible()
  })

  test('opens field options panel when clicking options button', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    await expect(page.locator('[data-testid="field-options-panel-Front"]')).toBeVisible()
    await expect(page.locator('text=Field Options: Front')).toBeVisible()
  })

  test('shows font, size, and RTL options in panel', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    await expect(page.locator('[data-testid="field-font-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-size-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-rtl-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-html-Front"]')).toBeVisible()
  })

  test('can set font option', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    // Select Arial font
    await page.selectOption('[data-testid="field-font-Front"]', 'Arial')

    // Wait for the mutation to complete
    await page.waitForTimeout(1000)

    // Verify the font is set
    await expect(page.locator('[data-testid="field-font-Front"]')).toHaveValue('Arial')
  })

  test('can set font size option', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    // Select 20px font size
    await page.selectOption('[data-testid="field-size-Front"]', '20')

    // Wait for mutation
    await page.waitForTimeout(1000)

    // Verify the font size is set
    await expect(page.locator('[data-testid="field-size-Front"]')).toHaveValue('20')
  })

  test('can enable RTL option', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    // Enable RTL by clicking the checkbox
    await page.click('[data-testid="field-rtl-Front"]')

    // Wait for mutation
    await page.waitForTimeout(1000)

    // Checkbox should be checked
    await expect(page.locator('[data-testid="field-rtl-Front"]')).toBeChecked()

    // Should show RTL badge on field
    await expect(page.locator('span:has-text("RTL")')).toBeVisible()
  })

  test('can enable HTML editor default option', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')

    // Enable HTML editor mode by default for this field
    await page.click('[data-testid="field-html-Front"]')
    await page.waitForTimeout(1000)

    await expect(page.locator('[data-testid="field-html-Front"]')).toBeChecked()
  })

  test('closes options panel when clicking button again', async ({ page }) => {
    await page.click('[data-testid="field-options-Front"]')
    await expect(page.locator('[data-testid="field-options-panel-Front"]')).toBeVisible()

    // Click again to close
    await page.click('[data-testid="field-options-Front"]')
    await expect(page.locator('[data-testid="field-options-panel-Front"]')).not.toBeVisible()
  })
})
