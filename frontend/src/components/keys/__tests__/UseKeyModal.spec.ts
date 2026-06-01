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

import UseKeyModal from '../UseKeyModal.vue'

describe('UseKeyModal', () => {
  beforeEach(() => {
    localeRef.value = 'en'
  })

  it('renders GPT-5.5 and goals feature in OpenAI Codex config', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'openai'
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

    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    const configToml = codeBlocks.find((content) => content.includes('model_provider = "OpenAI"'))

    expect(configToml).toBeDefined()
    expect(configToml).toContain('model = "gpt-5.5"')
    expect(configToml).toContain('review_model = "gpt-5.5"')
    expect(configToml).not.toContain('model = "gpt-5.4"')
    expect(configToml).not.toContain('model_context_window')
    expect(configToml).not.toContain('model_auto_compact_token_limit')
    expect(configToml).toContain('[features]\ngoals = true')
  })

  it('renders GPT-5.5 and goals feature in OpenAI Codex WebSocket config', async () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'openai'
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

    const wsTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.codexCliWs')
    )

    expect(wsTab).toBeDefined()
    await wsTab!.trigger('click')
    await nextTick()

    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text())
    const configToml = codeBlocks.find((content) => content.includes('supports_websockets = true'))

    expect(configToml).toBeDefined()
    expect(configToml).toContain('model = "gpt-5.5"')
    expect(configToml).toContain('review_model = "gpt-5.5"')
    expect(configToml).not.toContain('model = "gpt-5.4"')
    expect(configToml).not.toContain('model_context_window')
    expect(configToml).not.toContain('model_auto_compact_token_limit')
    expect(configToml).toContain('[features]\nresponses_websockets_v2 = true\ngoals = true')
  })

  it('renders Entrox CLI install methods and login command', async () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'entrox'
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

    const opencodeTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.opencode')
    )

    expect(opencodeTab).toBeDefined()
    await opencodeTab!.trigger('click')
    await nextTick()

    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text())

    expect(wrapper.text()).toContain('keys.useKeyModal.opencode.loginTitle')
    expect(wrapper.text()).toContain('keys.useKeyModal.opencode.installTitle')
    expect(codeBlocks[0]).toContain('ENTROX_HOMEBREW_TAP="$(brew --repository)/Library/Taps/hexonal/homebrew-entrox"')
    expect(codeBlocks[0]).toContain('git -C "$ENTROX_HOMEBREW_TAP" reset --hard origin/main')
    expect(codeBlocks[0]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew tap hexonal/entrox')
    expect(codeBlocks[0]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew trust hexonal/entrox')
    expect(codeBlocks[0]).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox')
    expect(codeBlocks[0]).not.toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(codeBlocks[0]).toContain('Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser')
    expect(codeBlocks[0]).toContain('Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression')
    expect(codeBlocks[0]).toContain('scoop bucket add entrox https://github.com/hexonal/scoop-entrox')
    expect(codeBlocks[0]).toContain('scoop install entrox')
    expect(codeBlocks[1]).toBe('entrox login')
    expect(codeBlocks.join('\n')).not.toContain('opencode.json')
  })

  it('renders optional CN proxy guidance for Entrox install methods', async () => {
    localeRef.value = 'zh'

    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'entrox'
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

    const opencodeTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.opencode')
    )

    expect(opencodeTab).toBeDefined()
    await opencodeTab!.trigger('click')
    await nextTick()

    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text())

    expect(codeBlocks[0]).toContain('如 GitHub 连接失败')
    expect(codeBlocks[0]).toContain('export HTTPS_PROXY')
    expect(codeBlocks[0]).toContain('scoop config proxy')
    expect(codeBlocks[0]).toContain('你的代理端口')
    expect(codeBlocks[0]).not.toContain('127.0.0.1:7890')
  })

  it('does not show Entrox CLI tab for non-Entrox platforms', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'openai'
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

    const tabLabels = wrapper.findAll('button').map((button) => button.text())

    expect(tabLabels.some((label) => label.includes('keys.useKeyModal.cliTabs.opencode'))).toBe(false)
  })
})
