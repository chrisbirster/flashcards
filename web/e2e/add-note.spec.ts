import { test, expect } from '@playwright/test'

test.describe('Add Note Screen', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck first
    await page.fill('input[placeholder="Enter deck name..."]', `Add Note Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)
  })

  test('opens add note screen when clicking Add Cards button', async ({ page }) => {
    // Click Add Cards button on the first deck
    await page.click('button:has-text("Add Cards")')

    // Should show Add Note screen
    await expect(page.locator('h1:has-text("Add Note")')).toBeVisible()
  })

  test('displays note type selector', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Should have note type dropdown
    const noteTypeSelect = page.locator('select').first()
    await expect(noteTypeSelect).toBeVisible()

    // Should have a note type selected (any built-in type is valid)
    const value = await noteTypeSelect.inputValue()
    expect(value.length).toBeGreaterThan(0)
    expect(value).toContain('Basic') // All built-in types contain "Basic" or "Cloze"
  })

  test('displays deck selector', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Should have deck dropdown
    const deckSelect = page.locator('select').nth(1)
    await expect(deckSelect).toBeVisible()
  })

  test('displays field inputs for selected note type', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Should show Front and Back fields for Basic note type
    await expect(page.locator('label:has-text("Front")')).toBeVisible()
    await expect(page.locator('label:has-text("Back")')).toBeVisible()
  })

  test('displays tags input', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Should have tags input
    await expect(page.locator('label:has-text("Tags")')).toBeVisible()
    await expect(page.locator('input[placeholder*="tags"]')).toBeVisible()
  })

  test('Add Note button is disabled when no content', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Add Note button should be disabled with empty fields
    const addButton = page.locator('button:has-text("Add Note")')
    await expect(addButton).toBeDisabled()
  })

  test('Add Note button is enabled when fields have content', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Fill in the Front field
    await page.locator('textarea').first().fill('Test question')

    // Add Note button should now be enabled
    const addButton = page.locator('button:has-text("Add Note")')
    await expect(addButton).toBeEnabled()
  })

  test('successfully creates a note', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Fill in the fields
    await page.locator('textarea').first().fill('What is 2+2?')
    await page.locator('textarea').nth(1).fill('4')

    // Click Add Note
    await page.click('button:has-text("Add Note")')

    // Should show success message
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })
  })

  test('clears fields after successful note creation', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Fill in the fields
    const frontField = page.locator('textarea').first()
    const backField = page.locator('textarea').nth(1)
    await frontField.fill('What is 2+2?')
    await backField.fill('4')

    // Click Add Note
    await page.click('button:has-text("Add Note")')

    // Wait for success
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Fields should be cleared
    await expect(frontField).toHaveValue('')
    await expect(backField).toHaveValue('')
  })

  test('Close button returns to deck list', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Click Close button
    await page.click('button:has-text("Close")')

    // Should return to deck list
    await expect(page.locator('h1:has-text("Microdote")')).toBeVisible()
  })

  test('Cancel button returns to deck list', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Click Cancel button
    await page.click('button:has-text("Cancel")')

    // Should return to deck list
    await expect(page.locator('h1:has-text("Microdote")')).toBeVisible()
  })

  test('shows preview when content is entered', async ({ page }) => {
    await page.click('button:has-text("Add Cards")')

    // Fill in the fields
    await page.locator('textarea').first().fill('Test question')

    // Should show preview section
    await expect(page.locator('h2:has-text("Preview")')).toBeVisible()
  })

  test('updates deck stats after adding note', async ({ page }) => {
    // Find the deck we just created
    const deckItem = page.locator('li').first()
    await expect(deckItem).toBeVisible()

    // Add a note
    await page.click('button:has-text("Add Cards")')
    await page.locator('textarea').first().fill('Stats Test Question')
    await page.locator('textarea').nth(1).fill('Stats Test Answer')
    await page.click('button:has-text("Add Note")')

    // Wait for success
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Close and verify stats show at least 1 card
    await page.click('button:has-text("Close")')

    // Stats should be visible showing some total (at least the card we just added)
    await expect(deckItem.locator('text=total')).toBeVisible({ timeout: 5000 })
  })

  test('shows duplicate warning when entering existing content', async ({ page }) => {
    // First create a note with known content
    await page.click('button:has-text("Add Cards")')
    const uniqueContent = `Duplicate Test ${Date.now()}`
    await page.locator('textarea').first().fill(uniqueContent)
    await page.locator('textarea').nth(1).fill('Answer')
    await page.click('button:has-text("Add Note")')
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Now try to create another note with the same Front content
    await page.locator('textarea').first().fill(uniqueContent)

    // Wait for duplicate check (debounced 500ms)
    await page.waitForTimeout(700)

    // Should show duplicate warning
    await expect(page.locator('[data-testid="duplicate-warning"]')).toBeVisible({ timeout: 3000 })
    await expect(page.locator('text=Possible duplicate found')).toBeVisible()
  })

  test('can still add note despite duplicate warning', async ({ page }) => {
    // First create a note with known content
    await page.click('button:has-text("Add Cards")')
    const uniqueContent = `Duplicate Override ${Date.now()}`
    await page.locator('textarea').first().fill(uniqueContent)
    await page.locator('textarea').nth(1).fill('Answer 1')
    await page.click('button:has-text("Add Note")')
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Enter the same content again
    await page.locator('textarea').first().fill(uniqueContent)
    await page.locator('textarea').nth(1).fill('Answer 2')

    // Wait for duplicate check
    await page.waitForTimeout(700)

    // Should show duplicate warning
    await expect(page.locator('[data-testid="duplicate-warning"]')).toBeVisible({ timeout: 3000 })

    // Should still be able to add the note
    await page.click('button:has-text("Add Note")')
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })
  })

  test('duplicate warning clears after successful add', async ({ page }) => {
    // First create a note with known content
    await page.click('button:has-text("Add Cards")')
    const uniqueContent = `Duplicate Clear ${Date.now()}`
    await page.locator('textarea').first().fill(uniqueContent)
    await page.locator('textarea').nth(1).fill('Answer')
    await page.click('button:has-text("Add Note")')
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Enter the same content again to trigger warning
    await page.locator('textarea').first().fill(uniqueContent)
    await page.waitForTimeout(700)
    await expect(page.locator('[data-testid="duplicate-warning"]')).toBeVisible({ timeout: 3000 })

    // Add the note anyway
    await page.locator('textarea').nth(1).fill('Answer 2')
    await page.click('button:has-text("Add Note")')
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })

    // Warning should be cleared after successful add (fields are cleared)
    await expect(page.locator('[data-testid="duplicate-warning"]')).not.toBeVisible()
  })
})

test.describe('Cloze Note Editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck first
    await page.fill('input[placeholder="Enter deck name..."]', `Cloze Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')
  })

  test('shows cloze button when Cloze note type is selected', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Cloze button should be visible
    await expect(page.locator('[data-testid="add-cloze-button"]')).toBeVisible()
  })

  test('does not show cloze button for Basic note type', async ({ page }) => {
    // Select Basic note type (should be default)
    const noteTypeSelect = page.locator('select').first()
    const value = await noteTypeSelect.inputValue()
    expect(value).toContain('Basic')

    // Cloze button should not be visible
    await expect(page.locator('[data-testid="add-cloze-button"]')).not.toBeVisible()
  })

  test('Cloze note type shows Text and Extra fields', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Should have Text and Extra fields
    await expect(page.locator('label:has-text("Text")')).toBeVisible()
    await expect(page.locator('label:has-text("Extra")')).toBeVisible()
  })

  test('can insert cloze deletion using button', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Focus on Text field and type some text
    const textField = page.locator('textarea').first()
    await textField.focus()
    await textField.fill('The capital of France is Paris')

    // Select "Paris"
    await textField.evaluate((el) => {
      const textarea = el as HTMLTextAreaElement
      const text = textarea.value
      const start = text.indexOf('Paris')
      textarea.setSelectionRange(start, start + 'Paris'.length)
    })

    // Click cloze button
    await page.click('[data-testid="add-cloze-button"]')

    // Text should now contain cloze syntax
    const value = await textField.inputValue()
    expect(value).toContain('{{c1::Paris}}')
  })

  test('auto-increments cloze number', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Enter text with existing c1
    const textField = page.locator('textarea').first()
    await textField.fill('{{c1::France}} is in {{c1::Europe}}. The capital is Paris')
    await textField.focus()

    // Select "Paris"
    await textField.evaluate((el) => {
      const textarea = el as HTMLTextAreaElement
      const text = textarea.value
      const start = text.indexOf('Paris')
      textarea.setSelectionRange(start, start + 'Paris'.length)
    })

    // Click cloze button - should use c2 since c1 already exists
    await page.click('[data-testid="add-cloze-button"]')

    // Text should now contain c2
    const value = await textField.inputValue()
    expect(value).toContain('{{c2::Paris}}')
  })

  test('shows preview for cloze cards', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Enter cloze text
    const textField = page.locator('textarea').first()
    await textField.fill('The {{c1::capital}} of {{c2::France}} is Paris')

    // Preview should show 2 cards (c1 and c2)
    await expect(page.locator('text=Card 1: Cloze 1')).toBeVisible()
    await expect(page.locator('text=Card 2: Cloze 2')).toBeVisible()
  })

  test('can create cloze note successfully', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Enter cloze text
    const textField = page.locator('textarea').first()
    await textField.fill('The capital of {{c1::France}} is {{c2::Paris}}')

    // Click Add Note
    await page.click('button:has-text("Add Note")')

    // Should show success message
    await expect(page.locator('text=Note added successfully')).toBeVisible({ timeout: 5000 })
  })

  test('shows hint when no cloze deletions found', async ({ page }) => {
    // Select Cloze note type
    await page.selectOption('select', { label: 'Cloze' })

    // Enter text without cloze syntax
    const textField = page.locator('textarea').first()
    await textField.fill('This text has no cloze deletions')

    // Should show hint about cloze syntax
    await expect(page.locator('text=No cloze deletions found')).toBeVisible()
  })
})
