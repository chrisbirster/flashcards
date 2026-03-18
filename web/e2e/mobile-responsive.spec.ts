import { test, expect, type Page } from '@playwright/test'

async function expectNoHorizontalOverflow(page: Page) {
  const hasOverflow = await page.evaluate(() => {
    const viewportWidth = document.documentElement.clientWidth
    return document.documentElement.scrollWidth > viewportWidth + 1
  })

  expect(hasOverflow).toBeFalsy()
}

test.describe('Mobile Responsive Layouts', () => {
  test.use({ viewport: { width: 390, height: 844 } })

  test('renders deck and add-note screens without overflow', async ({ page }) => {
    await page.goto('/')

    await page.fill('input[placeholder="Deck name"]', `Mobile Test ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    await expect(page.locator('h1:has-text("Vutadex")')).toBeVisible()
    await expectNoHorizontalOverflow(page)

    await page.click('button:has-text("Add Cards")')

    await expect(page.locator('h1:has-text("Add Note")')).toBeVisible()
    await page.locator('textarea').first().fill('Mobile front content')
    await expect(page.locator('h2:has-text("Preview")')).toBeVisible()
    await expectNoHorizontalOverflow(page)
  })

  test('renders template editor in mobile viewport', async ({ page }) => {
    await page.goto('/')

    await page.fill('input[placeholder="Deck name"]', `Mobile Template ${Date.now()}`)
    await page.click('button:has-text("Create")')
    await page.waitForTimeout(500)

    await page.click('button:has-text("Add Cards")')
    await page.click('[data-testid="edit-templates-button"]')

    await expect(page).toHaveURL(/\/notes\/add\/note-types\/[^/]+\/templates(?:\?.*)?$/)
    await expect(page.locator('[data-testid="editor-front"]')).toBeVisible()
    await expect(page.locator('[data-testid="save-template"]')).toBeVisible()
    await expectNoHorizontalOverflow(page)
  })
})
