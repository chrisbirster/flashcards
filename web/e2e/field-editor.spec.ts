import { test, expect, type Page } from '@playwright/test'

async function openFieldOptions(page: Page, fieldName: string) {
  const optionsButton = page.locator(`[data-testid="field-options-${fieldName}"]`)
  await optionsButton.evaluate((el) => {
    ;(el as HTMLButtonElement).click()
  })
}

async function clickFieldEditorAdd(page: Page) {
  const addButton = page.locator('[data-testid="field-editor-add-button"]')
  await addButton.evaluate((el) => {
    ;(el as HTMLButtonElement).click()
  })
}

test.describe('Field Editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')

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

  test('opens field editor route when clicking edit fields button', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')

    // Field editor route + UI should be visible
    await expect(page).toHaveURL(/\/notes\/add\/note-types\/[^/]+\/fields(?:\?.*)?$/)
    await expect(page.locator('text=Edit Fields:')).toBeVisible()
  })

  test('displays existing fields in field editor', async ({ page }) => {
    // Select Basic note type which has Front and Back fields
    await page.waitForFunction(() =>
      Array.from(document.querySelectorAll('select option')).some(
        (option) => option.textContent?.trim() === 'Basic',
      ),
    )
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-fields-button"]')

    // Should show Front and Back fields
    await expect(page.getByText('Front', { exact: true })).toBeVisible()
    await expect(page.getByText('Back', { exact: true })).toBeVisible()
  })

  test('can add a new field', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')
    const fieldName = `MyNewField${Date.now()}`

    // Add a new field
    await page.fill('input[placeholder="New field name..."]', fieldName)
    await clickFieldEditorAdd(page)

    // Should show the new field
    await expect(
      page.locator('.bg-white.rounded-lg.shadow-xl').getByText(fieldName, { exact: true }).first(),
    ).toBeVisible({ timeout: 3000 })
  })

  test('shows error for reserved field names', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')

    // Try to add a reserved field name
    await page.fill('input[placeholder="New field name..."]', 'Tags')
    await clickFieldEditorAdd(page)

    // Should show error about reserved name
    await expect(page.locator('text=reserved field name')).toBeVisible({ timeout: 3000 })
  })

  test('can close field editor and return to add note route', async ({ page }) => {
    await page.click('[data-testid="edit-fields-button"]')
    await expect(page).toHaveURL(/\/notes\/add\/note-types\/[^/]+\/fields(?:\?.*)?$/)
    await expect(page.locator('text=Edit Fields:')).toBeVisible()

    // Click close button in the modal footer
    await page.locator('[data-testid="close-field-editor-footer"]').evaluate((el) => {
      ;(el as HTMLButtonElement).click()
    })

    // Modal should be closed
    await expect(page).toHaveURL(/\/notes\/add(?:\?.*)?$/)
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
    await clickFieldEditorAdd(page)

    // Should show error
    await expect(page.locator('text=already exists')).toBeVisible({ timeout: 3000 })
  })
})

test.describe('Field Options', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')

    // Create a test deck
    await page.fill('input[placeholder="Deck name"]', `Field Options Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')

    // Select Basic note type
    await page.waitForFunction(() =>
      Array.from(document.querySelectorAll('select option')).some(
        (option) => option.textContent?.trim() === 'Basic',
      ),
    )
    await page.selectOption('select', { label: 'Basic' })

    // Open field editor
    await page.click('[data-testid="edit-fields-button"]')
  })

  test('shows field options button for each field', async ({ page }) => {
    await expect(page.locator('[data-testid="field-options-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-options-Back"]')).toBeVisible()
  })

  test('opens field options panel when clicking options button', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    await expect(page.locator('[data-testid="field-options-panel-Front"]')).toBeVisible()
    await expect(page.locator('text=Field Options: Front')).toBeVisible()
  })

  test('shows font, size, and RTL options in panel', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    await expect(page.locator('[data-testid="field-font-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-size-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-rtl-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="field-html-Front"]')).toBeVisible()
  })

  test('can set font option', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    // Select Arial font
    await page.selectOption('[data-testid="field-font-Front"]', 'Arial')

    // Wait for the mutation to complete
    await page.waitForTimeout(1000)

    // Verify the font is set
    await expect(page.locator('[data-testid="field-font-Front"]')).toHaveValue('Arial')
  })

  test('can set font size option', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    // Select 20px font size
    await page.selectOption('[data-testid="field-size-Front"]', '20')

    // Wait for mutation
    await page.waitForTimeout(1000)

    // Verify the font size is set
    await expect(page.locator('[data-testid="field-size-Front"]')).toHaveValue('20')
  })

  test('can enable RTL option', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    // Enable RTL by clicking the checkbox
    const rtlInput = page.locator('[data-testid="field-rtl-Front"]')
    if (!(await rtlInput.isChecked())) {
      await rtlInput.click()
    }

    // Wait for mutation
    await page.waitForTimeout(1000)

    // Should show RTL badge on field
    await expect(page.locator('[data-testid="rtl-chip-Front"]')).toBeVisible()
  })

  test('can enable HTML editor default option', async ({ page }) => {
    await openFieldOptions(page, 'Front')

    // Enable HTML editor mode by default for this field
    const htmlInput = page.locator('[data-testid="field-html-Front"]')
    if (!(await htmlInput.isChecked())) {
      await htmlInput.click()
    }
    await page.waitForTimeout(1000)

    await expect(htmlInput).toBeChecked()
  })

  test('closes options panel when clicking button again', async ({ page }) => {
    const optionsButton = page.locator('[data-testid="field-options-Front"]')
    await optionsButton.evaluate((el) => {
      ;(el as HTMLButtonElement).click()
    })
    await expect(page.locator('[data-testid="field-options-panel-Front"]')).toBeVisible()

    // Click again to close
    await optionsButton.evaluate((el) => {
      ;(el as HTMLButtonElement).click()
    })
    await expect(page.locator('[data-testid="field-options-panel-Front"]')).not.toBeVisible()
  })
})
