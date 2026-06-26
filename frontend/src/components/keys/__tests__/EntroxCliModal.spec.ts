import { beforeEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

const localeRef = vi.hoisted(() => ({ value: 'en' }))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: localeRef
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn().mockResolvedValue(true)
  })
}))

import EntroxCliModal from '../EntroxCliModal.vue'

function mountModal() {
  return mount(EntroxCliModal, {
    props: {
      show: true
    },
    global: {
      stubs: {
        BaseDialog: {
          template: '<div><slot /><slot name="footer" /></div>'
        },
        Icon: {
          template: '<span />'
        }
      }
    }
  })
}

describe('EntroxCliModal', () => {
  beforeEach(() => {
    localeRef.value = 'en'
  })

  it('renders Entrox CLI install methods and login command', async () => {
    const wrapper = mountModal()

    expect(wrapper.text()).toContain('keys.entroxCli.description')

    const macosTab = wrapper.findAll('button').find((button) => button.text() === 'macOS')
    expect(macosTab).toBeDefined()
    await macosTab!.trigger('click')
    await nextTick()

    let codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    expect(codeBlocks[0]).toBe('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[0]).not.toContain('HOMEBREW_NO_AUTO_UPDATE')
    expect(codeBlocks[0]).not.toContain('irm https://entrox.996icu.wiki/install.ps1 | iex')
    expect(codeBlocks[1]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew tap hexonal/entrox')
    expect(codeBlocks[1]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew trust hexonal/entrox')
    expect(codeBlocks[1]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox')
    expect(codeBlocks[1]).not.toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[1]).not.toContain('irm https://entrox.996icu.wiki/install.ps1 | iex')
    expect(codeBlocks[1]).not.toContain('ENTROX_HOMEBREW_TAP')
    expect(codeBlocks[1]).not.toContain('reset --hard origin/main')
    expect(codeBlocks[2]).toBe('entrox login')

    const linuxTab = wrapper.findAll('button').find((button) => button.text() === 'Linux')
    expect(linuxTab).toBeDefined()
    await linuxTab!.trigger('click')
    await nextTick()

    codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    expect(codeBlocks[0]).toBe('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[0]).not.toContain('HOMEBREW_NO_AUTO_UPDATE')
    expect(codeBlocks[0]).not.toContain('irm https://entrox.996icu.wiki/install.ps1 | iex')
    expect(codeBlocks[1]).toBe('entrox login')

    const windowsTab = wrapper.findAll('button').find((button) => button.text() === 'Windows')
    expect(windowsTab).toBeDefined()
    await windowsTab!.trigger('click')
    await nextTick()

    codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    expect(codeBlocks[0]).toBe('irm https://entrox.996icu.wiki/install.ps1 | iex')
    expect(codeBlocks[1]).toContain('Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser')
    expect(codeBlocks[1]).toContain('Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression')
    expect(codeBlocks[1]).toContain('scoop bucket add entrox https://github.com/hexonal/scoop-entrox')
    expect(codeBlocks[1]).toContain('scoop install entrox')
    expect(codeBlocks[1]).not.toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[2]).toBe('entrox login')
    expect(codeBlocks.join('\n')).not.toContain('opencode.json')
  })

  it('keeps official CN Entrox install commands concise', async () => {
    localeRef.value = 'zh'

    const wrapper = mountModal()

    const macosTab = wrapper.findAll('button').find((button) => button.text() === 'macOS')
    expect(macosTab).toBeDefined()
    await macosTab!.trigger('click')
    await nextTick()

    let codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    expect(codeBlocks[0]).toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[0]).not.toContain('Aliyun OSS/CDN')
    expect(codeBlocks[0]).not.toContain('export HTTPS_PROXY')
    expect(codeBlocks[1]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox')

    const windowsTab = wrapper.findAll('button').find((button) => button.text() === 'Windows')
    expect(windowsTab).toBeDefined()
    await windowsTab!.trigger('click')
    await nextTick()

    codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    expect(codeBlocks[0]).toContain('irm https://entrox.996icu.wiki/install.ps1 | iex')
    expect(codeBlocks[1]).not.toContain('scoop config proxy')
    expect(codeBlocks[1]).not.toContain('你的代理端口')
    expect(codeBlocks[1]).not.toContain('127.0.0.1:7890')
  })
})
