import { describe, expect, it, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import HomeView from '../HomeView.vue'

const { checkAuth, fetchPublicSettings } = vi.hoisted(() => ({
  checkAuth: vi.fn(),
  fetchPublicSettings: vi.fn(),
}))
const localeRef = vi.hoisted(() => ({ value: 'en' }))

const messages: Record<string, string> = {
  'home.getStarted': 'Get Started',
  'home.login': 'Login',
  'home.goToDashboard': 'Dashboard',
  'home.viewDocs': 'Docs',
  'home.switchToLight': 'Light',
  'home.switchToDark': 'Dark',
  'home.install.title': 'Install Entrox CLI',
  'home.install.script': 'macOS / Linux',
  'home.install.powershell': 'Windows',
  'home.install.homebrew': 'Homebrew',
  'home.install.scoop': 'Scoop',
  'home.desktop.badge': 'Desktop Client',
  'home.desktop.title': 'Entrox Desktop',
  'home.desktop.description':
    'A desktop AI coding workspace. Sign in with Entrox to sync models and certificates.',
  'home.desktop.download': 'Download',
  'home.desktop.downloadPending': 'Download URL coming soon',
  'home.desktop.platforms.mac.title': 'macOS Client',
  'home.desktop.platforms.mac.description': 'For Apple Silicon and Intel Mac',
  'home.desktop.platforms.windows.title': 'Windows Client',
  'home.desktop.platforms.windows.description': 'For Windows 10/11 desktop environments',
  'home.desktop.platforms.linux.title': 'Linux Client',
  'home.desktop.platforms.linux.description': 'For mainstream AppImage distributions',
  'home.pricing.badge': 'Subscription plans',
  'home.pricing.title': 'Choose your Entrox plan',
  'home.pricing.description': 'Plans include Entrox Desktop, Entrox CLI, gateway access, and monthly model usage.',
  'home.pricing.recommended': 'Recommended',
  'home.pricing.month': '/ month',
  'home.pricing.cta': 'Subscribe',
  'home.pricing.note': 'Monthly usage resets each billing cycle for self-service plans. Enterprise terms are confirmed separately.',
  'home.pricing.plans.pro.name': 'PRO',
  'home.pricing.plans.pro.subtitle': 'For light personal coding',
  'home.pricing.plans.pro.usage': '1x monthly usage included',
  'home.pricing.plans.pro.concurrency': '1 concurrent agent task',
  'home.pricing.plans.pro.queue': 'Standard queue',
  'home.pricing.plans.plus.name': 'PLUS',
  'home.pricing.plans.plus.subtitle': 'For daily AI coding work',
  'home.pricing.plans.plus.usage': '3.5x monthly usage included',
  'home.pricing.plans.plus.concurrency': '3 concurrent agent tasks',
  'home.pricing.plans.plus.queue': 'Priority queue',
  'home.pricing.plans.ultra.name': 'Ultra',
  'home.pricing.plans.ultra.subtitle': 'For heavy agent workflows',
  'home.pricing.plans.ultra.usage': '10x monthly usage included',
  'home.pricing.plans.ultra.concurrency': '6-8 concurrent agent tasks',
  'home.pricing.plans.ultra.queue': 'Highest priority queue',
  'home.pricing.plans.enterprise.name': 'Enterprise Custom',
  'home.pricing.plans.enterprise.subtitle': 'For teams and dedicated deployments',
  'home.pricing.plans.enterprise.price': 'Custom pricing',
  'home.pricing.plans.enterprise.usage': 'Custom monthly usage pool',
  'home.pricing.plans.enterprise.concurrency': 'Dedicated concurrency and priority capacity',
  'home.pricing.plans.enterprise.support': 'Deployment, onboarding, and compliance support',
  'home.pricing.plans.enterprise.cta': 'Contact us',
  'home.features.unifiedGateway': 'Unified Gateway',
  'home.features.unifiedGatewayDesc': 'Unified gateway description',
  'home.features.multiAccount': 'Account Pool',
  'home.features.multiAccountDesc': 'Account pool description',
  'home.features.balanceQuota': 'Billing',
  'home.features.balanceQuotaDesc': 'Billing description',
  'home.providers.title': 'Providers',
  'home.providers.description': 'Supported providers',
  'home.providers.claude': 'Claude',
  'home.providers.gemini': 'Gemini',
  'home.providers.supported': 'Supported',
  'home.providers.more': 'More',
  'home.providers.soon': 'Soon',
  'home.footer.allRightsReserved': 'All rights reserved.',
  'home.docs': 'Docs',
}

const defaultHeroSubtitle =
  'Entrox 把模型网关、CLI 与 Desktop 工作台连成一套 AI 编程系统：一键登录同步 Claude/GPT/Gemini 等模型资源，在桌面同时管理多个项目与 Agent 会话，实时查看流式输出、排队提示词、切换模型/后端并接入 MCP 工具；在终端用 entrox 快速启动同一套能力，让日常改动、自动化任务和长程编码都保持连续、清晰、可控。'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => messages[key] ?? key,
      locale: localeRef,
    }),
  }
})

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    isAuthenticated: false,
    isAdmin: false,
    user: null,
    checkAuth,
  }),
  useAppStore: () => ({
    cachedPublicSettings: null,
    siteName: 'Entrox',
    siteLogo: '',
    siteSubtitle: 'Entrox CLI',
    docUrl: '',
    homeContent: '',
    publicSettingsLoaded: true,
    fetchPublicSettings,
  }),
}))

describe('HomeView', () => {
  beforeEach(() => {
    checkAuth.mockReset()
    fetchPublicSettings.mockReset()
    localeRef.value = 'en'
    localStorage.clear()

    Object.defineProperty(window, 'matchMedia', {
      configurable: true,
      value: vi.fn().mockReturnValue({ matches: false }),
    })
  })

  it('renders the hero terminal as an Entrox CLI login flow', () => {
    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    const terminal = wrapper.find('.terminal-body')
    const commandLine = wrapper.findAll('.line-1 span').map((span) => span.text())

    expect(wrapper.text()).toContain(defaultHeroSubtitle)
    expect(commandLine).toEqual(['$', 'entrox', 'login'])
    expect(terminal.text()).toContain('Opening browser authorization')
    expect(terminal.text()).toContain('SIGNED IN')
    expect(terminal.text()).toContain('Entrox CLI ready')
    expect(terminal.text()).not.toContain('curl')
    expect(terminal.text()).not.toContain('/v1/messages')
  })

  it('shows official Entrox CLI installer commands on the home hero', async () => {
    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).toContain('Install Entrox CLI')
    expect(wrapper.text()).toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')
    expect(wrapper.text()).not.toContain('Script')

    const scriptTab = wrapper.findAll('button').find((button) => button.text() === 'macOS / Linux')
    expect(scriptTab).toBeDefined()

    const powerShellTab = wrapper.findAll('button').find((button) => button.text() === 'Windows')
    expect(powerShellTab).toBeDefined()
    await powerShellTab!.trigger('click')
    expect(wrapper.text()).toContain('irm https://entrox.996icu.wiki/install.ps1 | iex')

    const homebrewTab = wrapper.findAll('button').find((button) => button.text() === 'Homebrew')
    expect(homebrewTab).toBeDefined()
    await homebrewTab!.trigger('click')
    expect(wrapper.text()).toContain('ENTROX_HOMEBREW_TAP="$(brew --repository)/Library/Taps/hexonal/homebrew-entrox"')
    expect(wrapper.text()).toContain('git -C "$ENTROX_HOMEBREW_TAP" reset --hard origin/main')
    expect(wrapper.text()).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew tap hexonal/entrox')
    expect(wrapper.text()).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew trust hexonal/entrox')
    expect(wrapper.text()).toContain('HOMEBREW_NO_AUTO_UPDATE=1 brew upgrade hexonal/entrox/entrox || HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox')

    const scoopTab = wrapper.findAll('button').find((button) => button.text() === 'Scoop')
    expect(scoopTab).toBeDefined()
    await scoopTab!.trigger('click')
    expect(wrapper.text()).toContain('Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser')
    expect(wrapper.text()).toContain('Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression')
    expect(wrapper.text()).toContain('scoop bucket add entrox https://github.com/hexonal/scoop-entrox')
    expect(wrapper.text()).toContain('scoop install entrox')
  })

  it('adds optional CN network proxy guidance for Chinese locale', async () => {
    localeRef.value = 'zh'

    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).toContain('Aliyun OSS/CDN')
    expect(wrapper.text()).toContain('curl -fsSL https://entrox.996icu.wiki/install | bash')

    const homebrewTab = wrapper.findAll('button').find((button) => button.text() === 'Homebrew')
    expect(homebrewTab).toBeDefined()
    await homebrewTab!.trigger('click')
    expect(wrapper.text()).toContain('如 GitHub 连接失败')
    expect(wrapper.text()).toContain('export HTTPS_PROXY')

    const scoopTab = wrapper.findAll('button').find((button) => button.text() === 'Scoop')
    expect(scoopTab).toBeDefined()
    await scoopTab!.trigger('click')
    expect(wrapper.text()).toContain('scoop config proxy')
    expect(wrapper.text()).toContain('你的代理端口')
    expect(wrapper.text()).not.toContain('127.0.0.1:7890')
  })

  it('shows Entrox Desktop download cards with pending download buttons', () => {
    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).toContain('Entrox Desktop')
    expect(wrapper.text()).toContain('A desktop AI coding workspace. Sign in with Entrox to sync models and certificates.')
    expect(wrapper.text()).toContain('macOS Client')
    expect(wrapper.text()).toContain('Windows Client')
    expect(wrapper.text()).toContain('Linux Client')

    const downloadButtons = wrapper.findAll('button[disabled]').filter((button) => button.text() === 'Download')
    expect(downloadButtons).toHaveLength(3)
    expect(downloadButtons.every((button) => button.attributes('title') === 'Download URL coming soon')).toBe(true)
  })

  it('shows approved subscription pricing tiers with enterprise custom option', () => {
    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).toContain('Subscription plans')
    expect(wrapper.text()).toContain('Choose your Entrox plan')
    expect(wrapper.text()).toContain('PRO')
    expect(wrapper.text()).toContain('¥59')
    expect(wrapper.text()).toContain('PLUS')
    expect(wrapper.text()).toContain('¥159')
    expect(wrapper.text()).toContain('Ultra')
    expect(wrapper.text()).toContain('¥399')
    expect(wrapper.text()).toContain('Enterprise Custom')
    expect(wrapper.text()).toContain('Custom pricing')
    expect(wrapper.text()).toContain('Custom monthly usage pool')
    expect(wrapper.text()).toContain('Dedicated concurrency and priority capacity')
    expect(wrapper.text()).toContain('Deployment, onboarding, and compliance support')
    expect(wrapper.text()).toContain('Contact us')
    const enterpriseContactLink = wrapper.find('a[href="mailto:ailun@matchtrio.com"]')
    expect(enterpriseContactLink.exists()).toBe(true)
    expect(enterpriseContactLink.text()).toBe('Contact us')
    expect(wrapper.text()).toContain('Recommended')
    expect(wrapper.text()).toContain(
      'Monthly usage resets each billing cycle for self-service plans. Enterprise terms are confirmed separately.'
    )
  })

  it('does not show the legacy feature tag row above pricing', () => {
    const wrapper = mount(HomeView, {
      global: {
        stubs: {
          LocaleSwitcher: true,
          Icon: true,
          RouterLink: {
            template: '<a><slot /></a>',
          },
        },
      },
    })

    expect(wrapper.text()).not.toContain('Subscription to API')
    expect(wrapper.text()).not.toContain('Sticky session')
    expect(wrapper.text()).not.toContain('Monthly quota')
  })
})
