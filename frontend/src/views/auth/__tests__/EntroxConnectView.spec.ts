import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import EntroxConnectView from '../EntroxConnectView.vue'
import type { ApiKey, Group } from '@/types'

const {
  apiClientPostMock,
  keysListMock,
  keysCreateMock,
  groupsGetAvailableMock,
  showSuccessMock,
  showErrorMock,
  routeState,
} = vi.hoisted(() => ({
  apiClientPostMock: vi.fn(),
  keysListMock: vi.fn(),
  keysCreateMock: vi.fn(),
  groupsGetAvailableMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  routeState: {
    query: {
      session_id: 'session-1',
    } as Record<string, string>,
  },
}))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
}))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    showSuccess: (...args: unknown[]) => showSuccessMock(...args),
    showError: (...args: unknown[]) => showErrorMock(...args),
  }),
}))

vi.mock('@/api', () => ({
  apiClient: {
    post: (...args: unknown[]) => apiClientPostMock(...args),
  },
  keysAPI: {
    list: (...args: unknown[]) => keysListMock(...args),
    create: (...args: unknown[]) => keysCreateMock(...args),
  },
  userGroupsAPI: {
    getAvailable: (...args: unknown[]) => groupsGetAvailableMock(...args),
  },
}))

function makeAPIKey(overrides: Partial<ApiKey> = {}): ApiKey {
  return {
    id: 11,
    user_id: 7,
    key: 'sk-existing-1234567890',
    name: 'Existing key',
    group_id: 3,
    status: 'active',
    ip_whitelist: [],
    ip_blacklist: [],
    last_used_at: null,
    quota: 0,
    quota_used: 0,
    expires_at: null,
    created_at: '2026-05-29T01:00:00Z',
    updated_at: '2026-05-29T01:00:00Z',
    rate_limit_5h: 0,
    rate_limit_1d: 0,
    rate_limit_7d: 0,
    usage_5h: 0,
    usage_1d: 0,
    usage_7d: 0,
    window_5h_start: null,
    window_1d_start: null,
    window_7d_start: null,
    reset_5h_at: null,
    reset_1d_at: null,
    reset_7d_at: null,
    ...overrides,
  }
}

function makeGroup(overrides: Partial<Group> = {}): Group {
  return {
    id: 3,
    name: 'Entrox Pro',
    description: null,
    platform: 'entrox',
    rate_multiplier: 1,
    is_exclusive: false,
    status: 'active',
    subscription_type: 'subscription',
    daily_limit_usd: null,
    weekly_limit_usd: null,
    monthly_limit_usd: null,
    allow_image_generation: false,
    image_rate_independent: false,
    image_rate_multiplier: 1,
    image_price_1k: null,
    image_price_2k: null,
    image_price_4k: null,
    claude_code_only: false,
    fallback_group_id: null,
    fallback_group_id_on_invalid_request: null,
    require_oauth_only: false,
    require_privacy_set: false,
    created_at: '2026-05-29T01:00:00Z',
    updated_at: '2026-05-29T01:00:00Z',
    ...overrides,
  }
}

function mountView() {
  return mount(EntroxConnectView, {
    global: {
      stubs: {
        AuthLayout: { template: '<div><slot /></div>' },
        Icon: true,
      },
    },
  })
}

describe('EntroxConnectView', () => {
  beforeEach(() => {
    apiClientPostMock.mockReset()
    keysListMock.mockReset()
    keysCreateMock.mockReset()
    groupsGetAvailableMock.mockReset()
    showSuccessMock.mockReset()
    showErrorMock.mockReset()
    routeState.query = { session_id: 'session-1' }
    groupsGetAvailableMock.mockResolvedValue([makeGroup()])
    keysCreateMock.mockResolvedValue(makeAPIKey({ id: 22, name: 'Created key', key: 'sk-created-1234567890' }))
    keysListMock.mockResolvedValue({
      items: [makeAPIKey()],
      total: 1,
      page: 1,
      page_size: 100,
      pages: 1,
    })
  })

  it('loads active api keys and waits for explicit confirmation', async () => {
    mountView()

    await flushPromises()

    expect(keysListMock).toHaveBeenCalledWith(1, 100, {
      status: 'active',
      sort_by: 'created_at',
      sort_order: 'desc',
    })
    expect(groupsGetAvailableMock).toHaveBeenCalled()
    expect(apiClientPostMock).not.toHaveBeenCalled()
  })

  it('approves the selected existing api key', async () => {
    const wrapper = mountView()

    await flushPromises()
    await wrapper.find('[data-testid="entrox-approve-button"]').trigger('click')
    await flushPromises()

    expect(apiClientPostMock).toHaveBeenCalledWith('/auth/entrox/approve', {
      session_id: 'session-1',
      api_key_id: 11,
    })
    expect(showSuccessMock).toHaveBeenCalledWith('entrox CLI 已连接')
  })

  it('approves explicit new api key creation when no key exists', async () => {
    keysListMock.mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 100,
      pages: 1,
    })
    const wrapper = mountView()

    await flushPromises()
    await wrapper.find('[data-testid="entrox-approve-button"]').trigger('click')
    await flushPromises()

    expect(keysCreateMock).toHaveBeenCalledWith(expect.stringMatching(/^entrox CLI /), 3)
    expect(apiClientPostMock).toHaveBeenCalledWith('/auth/entrox/approve', {
      session_id: 'session-1',
      api_key_id: 22,
    })
  })

  it('allows choosing new api key creation even when existing keys are available', async () => {
    const wrapper = mountView()

    await flushPromises()
    await wrapper.findAll('input[name="entrox-api-key"]')[1].setValue(true)
    await wrapper.find('[data-testid="entrox-approve-button"]').trigger('click')
    await flushPromises()

    expect(keysCreateMock).toHaveBeenCalledWith(expect.stringMatching(/^entrox CLI /), 3)
    expect(apiClientPostMock).toHaveBeenCalledWith('/auth/entrox/approve', {
      session_id: 'session-1',
      api_key_id: 22,
    })
  })
})
