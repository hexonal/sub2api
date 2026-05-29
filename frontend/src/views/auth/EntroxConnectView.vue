<template>
  <AuthLayout>
    <div class="space-y-5">
      <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full" :class="statusClass">
        <Icon v-if="status === 'success'" name="checkCircle" size="lg" />
        <Icon v-else-if="status === 'error'" name="xCircle" size="lg" />
        <Icon v-else-if="status === 'ready'" name="key" size="lg" />
        <svg
          v-else
          class="h-6 w-6 animate-spin"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      </div>

      <div class="text-center">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ title }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ message }}
        </p>
      </div>

      <div v-if="status === 'ready' || status === 'approving'" class="space-y-3">
        <div v-if="apiKeys.length > 0" class="space-y-2">
          <label
            v-for="key in apiKeys"
            :key="key.id"
            class="flex cursor-pointer items-start gap-3 rounded-lg border p-3 text-left transition"
            :class="selectedMode === 'existing' && selectedAPIKeyID === key.id
              ? 'border-primary-400 bg-primary-50/70 dark:border-primary-500/70 dark:bg-primary-500/10'
              : 'border-gray-200 bg-white hover:border-primary-200 dark:border-dark-700 dark:bg-dark-900 dark:hover:border-primary-700'"
          >
            <input
              type="radio"
              name="entrox-api-key"
              class="mt-1 h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-900"
              :checked="selectedMode === 'existing' && selectedAPIKeyID === key.id"
              :disabled="status === 'approving'"
              @change="selectExistingAPIKey(key.id)"
            />
            <span class="min-w-0 flex-1">
              <span class="block truncate text-sm font-medium text-gray-900 dark:text-white">
                {{ key.name }}
              </span>
              <span class="mt-1 block truncate font-mono text-xs text-gray-500 dark:text-dark-400">
                {{ maskAPIKey(key.key) }}
              </span>
              <span class="mt-1 block text-xs text-gray-500 dark:text-dark-400">
                {{ keyMeta(key) }}
              </span>
            </span>
          </label>
        </div>

        <label
          class="flex cursor-pointer items-start gap-3 rounded-lg border p-3 text-left transition"
          :class="selectedMode === 'create'
            ? 'border-primary-400 bg-primary-50/70 dark:border-primary-500/70 dark:bg-primary-500/10'
            : 'border-gray-200 bg-white hover:border-primary-200 dark:border-dark-700 dark:bg-dark-900 dark:hover:border-primary-700'"
        >
          <input
            type="radio"
            name="entrox-api-key"
            class="mt-1 h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-900"
            :checked="selectedMode === 'create'"
            :disabled="status === 'approving'"
            @change="selectCreateNew"
          />
          <span>
            <span class="block text-sm font-medium text-gray-900 dark:text-white">
              创建新的 API Key
            </span>
            <span class="mt-1 block text-xs text-gray-500 dark:text-dark-400">
              将为本次 entrox CLI 登录生成一个新的 SK。
            </span>
          </span>
        </label>
      </div>

      <button
        v-if="status === 'ready' || status === 'approving'"
        type="button"
        data-testid="entrox-approve-button"
        class="btn btn-primary w-full"
        :disabled="status === 'approving' || !canApprove"
        @click="approve"
      >
        <span v-if="status === 'approving'">正在授权...</span>
        <span v-else>授权 entrox CLI</span>
      </button>

      <button
        v-else-if="status === 'error'"
        type="button"
        class="btn btn-primary w-full"
        @click="loadAPIKeys"
      >
        重试
      </button>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { AuthLayout } from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import { apiClient, keysAPI } from '@/api'
import { useAppStore } from '@/stores'
import type { ApiKey } from '@/types'
import { formatDateTime } from '@/utils/format'

type Status = 'loading' | 'ready' | 'approving' | 'success' | 'error'
type SelectionMode = 'existing' | 'create'

const route = useRoute()
const appStore = useAppStore()
const status = ref<Status>('loading')
const message = ref('正在加载可用的 API Key...')
const apiKeys = ref<ApiKey[]>([])
const selectedMode = ref<SelectionMode>('create')
const selectedAPIKeyID = ref<number | null>(null)

const title = computed(() => {
  if (status.value === 'success') return 'entrox 已连接'
  if (status.value === 'error') return '连接失败'
  if (status.value === 'ready') return '选择 API Key'
  return '连接 entrox CLI'
})

const statusClass = computed(() => {
  if (status.value === 'success') return 'bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-300'
  if (status.value === 'error') return 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-300'
  return 'bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300'
})

const sessionID = computed(() => (typeof route.query.session_id === 'string' ? route.query.session_id : ''))

const canApprove = computed(() => {
  if (selectedMode.value === 'create') return true
  return selectedAPIKeyID.value !== null
})

async function loadAPIKeys(): Promise<void> {
  if (!sessionID.value) {
    status.value = 'error'
    message.value = '缺少 entrox 登录会话。'
    return
  }

  status.value = 'loading'
  message.value = '正在加载可用的 API Key...'

  try {
    const result = await keysAPI.list(1, 100, {
      status: 'active',
      sort_by: 'created_at',
      sort_order: 'desc',
    })
    apiKeys.value = result.items
    if (apiKeys.value.length > 0) {
      selectedMode.value = 'existing'
      selectedAPIKeyID.value = apiKeys.value[0].id
      message.value = '选择一个已有 SK，或创建新的 API Key。'
    } else {
      selectedMode.value = 'create'
      selectedAPIKeyID.value = null
      message.value = '当前账号没有可用 SK，可以创建新的 API Key。'
    }
    status.value = 'ready'
  } catch (error) {
    status.value = 'error'
    const err = error as { message?: string }
    message.value = err.message || '加载 API Key 失败，请重新发起 entrox 登录。'
    appStore.showError(message.value)
  }
}

async function approve(): Promise<void> {
  if (!sessionID.value) {
    status.value = 'error'
    message.value = '缺少 entrox 登录会话。'
    return
  }
  if (!canApprove.value) {
    status.value = 'error'
    message.value = '请选择一个 API Key。'
    return
  }

  status.value = 'approving'
  message.value = '正在授权 entrox CLI...'

  try {
    await apiClient.post('/auth/entrox/approve', buildApprovePayload())
    status.value = 'success'
    message.value = '可以回到终端继续使用 entrox。'
    appStore.showSuccess('entrox CLI 已连接')
    window.setTimeout(() => {
      window.close()
    }, 1200)
  } catch (error) {
    status.value = 'error'
    const err = error as { message?: string }
    message.value = err.message || '授权失败，请重新发起 entrox 登录。'
    appStore.showError(message.value)
  }
}

function buildApprovePayload(): { session_id: string; api_key_id?: number; create_new?: boolean } {
  if (selectedMode.value === 'existing' && selectedAPIKeyID.value !== null) {
    return {
      session_id: sessionID.value,
      api_key_id: selectedAPIKeyID.value,
    }
  }
  return {
    session_id: sessionID.value,
    create_new: true,
  }
}

function selectExistingAPIKey(id: number): void {
  selectedMode.value = 'existing'
  selectedAPIKeyID.value = id
}

function selectCreateNew(): void {
  selectedMode.value = 'create'
  selectedAPIKeyID.value = null
}

function maskAPIKey(key: string): string {
  if (key.length <= 12) return key
  return `${key.slice(0, 8)}...${key.slice(-4)}`
}

function keyMeta(key: ApiKey): string {
  if (key.last_used_at) return `上次使用 ${formatDateTime(key.last_used_at)}`
  return `创建于 ${formatDateTime(key.created_at)}`
}

onMounted(() => {
  loadAPIKeys()
})
</script>
