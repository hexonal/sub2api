<template>
  <AuthLayout>
    <div class="space-y-5 text-center">
      <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full" :class="statusClass">
        <Icon v-if="status === 'success'" name="checkCircle" size="lg" />
        <Icon v-else-if="status === 'error'" name="xCircle" size="lg" />
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

      <div>
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          {{ title }}
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ message }}
        </p>
      </div>

      <button
        v-if="status === 'error'"
        type="button"
        class="btn btn-primary w-full"
        @click="approve"
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
import { apiClient } from '@/api'
import { useAppStore } from '@/stores'

type Status = 'loading' | 'success' | 'error'

const route = useRoute()
const appStore = useAppStore()
const status = ref<Status>('loading')
const message = ref('正在连接 entrox CLI...')

const title = computed(() => {
  if (status.value === 'success') return 'entrox 已连接'
  if (status.value === 'error') return '连接失败'
  return '连接 entrox CLI'
})

const statusClass = computed(() => {
  if (status.value === 'success') return 'bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-300'
  if (status.value === 'error') return 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-300'
  return 'bg-primary-100 text-primary-600 dark:bg-primary-900/30 dark:text-primary-300'
})

async function approve(): Promise<void> {
  const sessionID = typeof route.query.session_id === 'string' ? route.query.session_id : ''
  if (!sessionID) {
    status.value = 'error'
    message.value = '缺少 entrox 登录会话。'
    return
  }

  status.value = 'loading'
  message.value = '正在授权 entrox CLI...'

  try {
    await apiClient.post('/auth/entrox/approve', { session_id: sessionID })
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

onMounted(() => {
  approve()
})
</script>
