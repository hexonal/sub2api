import { describe, expect, it, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import HomeView from '../HomeView.vue'

const { checkAuth, fetchPublicSettings } = vi.hoisted(() => ({
  checkAuth: vi.fn(),
  fetchPublicSettings: vi.fn(),
}))

const messages: Record<string, string> = {
  'home.getStarted': 'Get Started',
  'home.login': 'Login',
  'home.goToDashboard': 'Dashboard',
  'home.viewDocs': 'Docs',
  'home.switchToLight': 'Light',
  'home.switchToDark': 'Dark',
  'home.install.title': 'Install Entrox CLI',
  'home.install.script': 'Script',
  'home.install.homebrew': 'Homebrew',
  'home.install.scoop': 'Scoop',
  'home.tags.subscriptionToApi': 'Subscription to API',
  'home.tags.stickySession': 'Sticky session',
  'home.tags.realtimeBilling': 'Realtime billing',
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

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => messages[key] ?? key,
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

    expect(commandLine).toEqual(['$', 'entrox', 'login'])
    expect(terminal.text()).toContain('Opening browser authorization')
    expect(terminal.text()).toContain('SIGNED IN')
    expect(terminal.text()).toContain('Entrox CLI ready')
    expect(terminal.text()).not.toContain('curl')
    expect(terminal.text()).not.toContain('/v1/messages')
  })

  it('shows Entrox CLI install commands on the home hero', async () => {
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

    const homebrewTab = wrapper.findAll('button').find((button) => button.text() === 'Homebrew')
    expect(homebrewTab).toBeDefined()
    await homebrewTab!.trigger('click')
    expect(wrapper.text()).toContain('brew tap hexonal/entrox')
    expect(wrapper.text()).toContain('brew install entrox')

    const scoopTab = wrapper.findAll('button').find((button) => button.text() === 'Scoop')
    expect(scoopTab).toBeDefined()
    await scoopTab!.trigger('click')
    expect(wrapper.text()).toContain('scoop bucket add entrox https://github.com/hexonal/scoop-entrox')
    expect(wrapper.text()).toContain('scoop install entrox')
  })
})
