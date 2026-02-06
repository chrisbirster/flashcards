import { test, expect } from '@playwright/test'

test.describe('Template Editor', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck
    await page.fill('input[placeholder="Deck name"]', `Template Editor Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')
  })

  test('shows edit templates button next to note type selector', async ({ page }) => {
    await expect(page.locator('[data-testid="edit-templates-button"]')).toBeVisible()
  })

  test('opens template editor modal when clicking edit templates button', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    // Template editor modal should be visible
    await expect(page.locator('text=Edit Templates:')).toBeVisible()
  })

  test('displays three tabs: Front, Back, and Styling', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    await expect(page.locator('[data-testid="tab-front"]')).toBeVisible()
    await expect(page.locator('[data-testid="tab-back"]')).toBeVisible()
    await expect(page.locator('[data-testid="tab-styling"]')).toBeVisible()
  })

  test('displays front template content by default', async ({ page }) => {
    // Select Basic note type for predictable template content
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Front tab should be active (has blue border)
    const frontTab = page.locator('[data-testid="tab-front"]')
    await expect(frontTab).toHaveClass(/border-blue-500/)

    // Editor should show front template content
    const editor = page.locator('[data-testid="editor-front"]')
    await expect(editor).toBeVisible()
    // Basic template has {{Front}}
    await expect(editor).toHaveValue(/\{\{Front\}\}/)
  })

  test('can switch between tabs', async ({ page }) => {
    // Select Basic note type for predictable template content
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Click Back tab
    await page.click('[data-testid="tab-back"]')
    const backEditor = page.locator('[data-testid="editor-back"]')
    await expect(backEditor).toBeVisible()
    await expect(backEditor).toHaveValue(/\{\{Back\}\}/)

    // Click Styling tab
    await page.click('[data-testid="tab-styling"]')
    const stylingEditor = page.locator('[data-testid="editor-styling"]')
    await expect(stylingEditor).toBeVisible()
  })

  test('shows live preview pane', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    await expect(page.locator('[data-testid="preview-content"]')).toBeVisible()
  })

  test('updates preview when editing template', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    // Type custom content
    const editor = page.locator('[data-testid="editor-front"]')
    await editor.fill('Custom: {{Front}}!')

    // Preview should update
    await expect(page.locator('[data-testid="preview-content"]')).toContainText('Custom:')
  })

  test('shows sample field inputs for preview', async ({ page }) => {
    // Select Basic note type which has Front and Back fields
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Should have sample inputs for Front and Back fields
    await expect(page.locator('[data-testid="sample-Front"]')).toBeVisible()
    await expect(page.locator('[data-testid="sample-Back"]')).toBeVisible()
  })

  test('updates preview when changing sample values', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal to load
    await expect(page.locator('[data-testid="sample-Front"]')).toBeVisible()

    // Change sample value
    await page.locator('[data-testid="sample-Front"]').fill('My Custom Question')

    // Preview should show the custom value
    await expect(page.locator('[data-testid="preview-content"]')).toContainText('My Custom Question')
  })

  test('shows save button disabled when no changes', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal to load
    await expect(page.locator('[data-testid="editor-front"]')).toBeVisible()

    const saveButton = page.locator('[data-testid="save-template"]')
    await expect(saveButton).toBeDisabled()
  })

  test('enables save button when template is modified', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal to load
    await expect(page.locator('[data-testid="editor-front"]')).toBeVisible()

    // Modify the template
    const editor = page.locator('[data-testid="editor-front"]')
    await editor.fill('Modified: {{Front}}')

    // Save button should be enabled
    const saveButton = page.locator('[data-testid="save-template"]')
    await expect(saveButton).toBeEnabled()
  })

  test('shows unsaved changes indicator when template is modified', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal to load
    await expect(page.locator('[data-testid="editor-front"]')).toBeVisible()

    // Modify the template
    const editor = page.locator('[data-testid="editor-front"]')
    await editor.fill('Modified: {{Front}}')

    // Should show unsaved changes indicator
    await expect(page.locator('text=Unsaved changes')).toBeVisible()
  })

  test('can save template changes', async ({ page }) => {
    // Select Basic note type which has a simple template
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal and editor to fully load with initial value
    const editor = page.locator('[data-testid="editor-front"]')
    await expect(editor).toBeVisible()
    await expect(editor).toHaveValue(/\{\{Front\}\}/, { timeout: 5000 })

    // Use a unique value to trigger change detection
    const uniqueValue = `Updated ${Date.now()}: {{Front}}`
    await editor.fill(uniqueValue)

    // Verify unsaved changes is shown (proves the change was detected)
    await expect(page.locator('text=Unsaved changes')).toBeVisible({ timeout: 5000 })

    // Save changes
    await page.click('[data-testid="save-template"]')

    // Wait for save to complete - the mutation clears hasChanges on success
    await expect(page.locator('text=Unsaved changes')).not.toBeVisible({ timeout: 10000 })

    // Save button should be disabled again
    await expect(page.locator('[data-testid="save-template"]')).toBeDisabled()
  })

  test('can close template editor modal', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')
    await expect(page.locator('text=Edit Templates:')).toBeVisible()

    // Click close button in the modal footer (the one next to Save Changes)
    await page.locator('[data-testid="save-template"]').locator('..').locator('button:has-text("Close")').click()

    // Modal should be closed
    await expect(page.locator('text=Edit Templates:')).not.toBeVisible()
  })

  test('can close template editor using X button', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')
    await expect(page.locator('text=Edit Templates:')).toBeVisible()

    // Click X button
    await page.click('[data-testid="close-template-editor"]')

    // Modal should be closed
    await expect(page.locator('text=Edit Templates:')).not.toBeVisible()
  })

  test('can edit CSS styling', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal and front editor to be loaded
    const frontEditor = page.locator('[data-testid="editor-front"]')
    await expect(frontEditor).toBeVisible()
    await expect(frontEditor).toHaveValue(/\{\{Front\}\}/, { timeout: 5000 })

    // Switch to styling tab
    await page.click('[data-testid="tab-styling"]')

    // Wait for styling editor to be visible
    const editor = page.locator('[data-testid="editor-styling"]')
    await expect(editor).toBeVisible()

    // Add CSS with unique content
    const uniqueCSS = `.card { color: red; font-size: ${Date.now()}px; }`
    await editor.fill(uniqueCSS)

    // Should show unsaved changes
    await expect(page.locator('text=Unsaved changes')).toBeVisible({ timeout: 5000 })

    // Save
    await page.click('[data-testid="save-template"]')
    await expect(page.locator('[data-testid="save-template"]')).toBeDisabled({ timeout: 10000 })
  })

  test('persists template changes after closing and reopening', async ({ page }) => {
    // Select Basic note type
    await page.selectOption('select', { label: 'Basic' })
    await page.click('[data-testid="edit-templates-button"]')

    // Wait for modal and editor to be fully loaded
    const editor = page.locator('[data-testid="editor-front"]')
    await expect(editor).toBeVisible()
    await expect(editor).toHaveValue(/\{\{Front\}\}/) // Wait for initial value to load

    // Use a unique value to ensure there are changes
    const uniqueValue = `Persisted ${Date.now()}: {{Front}}`
    await editor.fill(uniqueValue)

    // Verify the change was detected
    await expect(page.locator('text=Unsaved changes')).toBeVisible()

    // Save and close
    await page.click('[data-testid="save-template"]')
    await expect(page.locator('[data-testid="save-template"]')).toBeDisabled({ timeout: 10000 })

    // Close via X button which is in the modal header
    await page.click('[data-testid="close-template-editor"]')

    // Reopen template editor
    await page.click('[data-testid="edit-templates-button"]')

    // Should still have the modified content
    await expect(page.locator('[data-testid="editor-front"]')).toHaveValue(uniqueValue)
  })
})

test.describe('Template Editor with Multiple Templates', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')

    // Create a test deck
    await page.fill('input[placeholder="Deck name"]', `Multi Template Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    // Open add note screen
    await page.click('button:has-text("Add Cards")')

    // Select Basic (and reversed card) which has 2 templates
    await page.selectOption('select', { label: 'Basic (and reversed card)' })
  })

  test('shows template selector for note types with multiple templates', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    await expect(page.locator('[data-testid="template-selector"]')).toBeVisible()
  })

  test('can switch between templates', async ({ page }) => {
    await page.click('[data-testid="edit-templates-button"]')

    // Select second template
    await page.selectOption('[data-testid="template-selector"]', { label: 'Card 2' })

    // Editor should now show Card 2's content (which uses {{Back}} for question)
    const editor = page.locator('[data-testid="editor-front"]')
    await expect(editor).toHaveValue(/\{\{Back\}\}/)
  })
})
