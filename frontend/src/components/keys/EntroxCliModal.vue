<template>
  <BaseDialog
    :show="show"
    :title="t('keys.entroxCli.title')"
    width="wide"
    @close="emit('close')"
  >
    <div class="space-y-4">
      <p class="text-sm text-gray-600 dark:text-gray-400">
        {{ t('keys.entroxCli.description') }}
      </p>

      <div class="border-b border-gray-200 dark:border-dark-700">
        <nav class="-mb-px flex space-x-4" aria-label="Entrox CLI install platform">
          <button
            v-for="tab in platformTabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'whitespace-nowrap py-2.5 px-1 border-b-2 font-medium text-sm transition-colors',
              activeTab === tab.id
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
            ]"
          >
            <span class="flex items-center gap-2">
              <component :is="tab.icon" class="w-4 h-4" />
              {{ tab.label }}
            </span>
          </button>
        </nav>
      </div>

      <div class="space-y-4">
        <div
          v-for="(file, index) in currentFiles"
          :key="index"
          class="relative"
        >
          <p v-if="file.hint" class="text-xs text-amber-600 dark:text-amber-400 mb-1.5 flex items-center gap-1">
            <Icon name="exclamationCircle" size="sm" class="flex-shrink-0" />
            {{ file.hint }}
          </p>
          <div class="bg-gray-900 dark:bg-dark-900 rounded-xl overflow-hidden">
            <div class="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-dark-800 border-b border-gray-700 dark:border-dark-700">
              <span class="text-xs text-gray-400 font-mono">{{ file.path }}</span>
              <button
                @click="copyContent(file.content, index)"
                class="flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium rounded-lg transition-colors"
                :class="copiedIndex === index
                  ? 'bg-green-500/20 text-green-400'
                  : 'bg-gray-700 hover:bg-gray-600 text-gray-300 hover:text-white'"
              >
                <Icon
                  :name="copiedIndex === index ? 'check' : 'clipboard'"
                  size="xs"
                  :stroke-width="copiedIndex === index ? 2 : 1.5"
                />
                {{ copiedIndex === index ? t('keys.useKeyModal.copied') : t('keys.useKeyModal.copy') }}
              </button>
            </div>
            <pre class="p-4 text-sm font-mono text-gray-100 overflow-x-auto"><code v-text="file.content"></code></pre>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end">
        <button
          @click="emit('close')"
          class="btn btn-secondary"
        >
          {{ t('common.close') }}
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, h, ref, type Component } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import {
  getDefaultEntroxInstallPlatformTab,
  getEntroxHomebrewInstallCommand,
  getEntroxPowerShellInstallCommand,
  getEntroxScoopInstallCommand,
  getEntroxScriptInstallCommand,
  isCnInstallLocale,
  type EntroxInstallPlatformTab
} from '@/utils/entroxInstall'

interface Props {
  show: boolean
}

interface Emits {
  (e: 'close'): void
}

interface TabConfig {
  id: EntroxInstallPlatformTab
  label: string
  icon: Component
}

interface FileConfig {
  path: string
  content: string
  hint?: string
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const { t, locale } = useI18n()
const { copyToClipboard: clipboardCopy } = useClipboard()

const copiedIndex = ref<number | null>(null)
const activeTab = ref<EntroxInstallPlatformTab>(getDefaultEntroxInstallPlatformTab())

const AppleIcon = {
  render() {
    return h('svg', {
      fill: 'currentColor',
      viewBox: '0 0 24 24',
      class: 'w-4 h-4'
    }, [
      h('path', { d: 'M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z' })
    ])
  }
}

const TerminalIcon = {
  render() {
    return h('svg', {
      fill: 'none',
      stroke: 'currentColor',
      viewBox: '0 0 24 24',
      'stroke-width': '1.5',
      class: 'w-4 h-4'
    }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        d: 'm6.75 7.5 3 2.25-3 2.25m4.5 0h3m-9 8.25h13.5A2.25 2.25 0 0 0 21 17.25V6.75A2.25 2.25 0 0 0 18.75 4.5H5.25A2.25 2.25 0 0 0 3 6.75v10.5A2.25 2.25 0 0 0 5.25 20.25Z'
      })
    ])
  }
}

const WindowsIcon = {
  render() {
    return h('svg', {
      fill: 'currentColor',
      viewBox: '0 0 24 24',
      class: 'w-4 h-4'
    }, [
      h('path', { d: 'M3 12V6.75l6-1.32v6.48L3 12zm17-9v8.75l-10 .15V5.21L20 3zM3 13l6 .09v6.81l-6-1.15V13zm7 .25l10 .15V21l-10-1.91v-5.84z' })
    ])
  }
}

const platformTabs: TabConfig[] = [
  { id: 'macos', label: 'macOS', icon: AppleIcon },
  { id: 'linux', label: 'Linux', icon: TerminalIcon },
  { id: 'windows', label: 'Windows', icon: WindowsIcon }
]

const currentFiles = computed((): FileConfig[] => {
  return buildFiles(activeTab.value, isCnInstallLocale(locale.value))
})

function buildFiles(platform: EntroxInstallPlatformTab, isCn: boolean): FileConfig[] {
  const installTitle = t('keys.entroxCli.installTitle')
  const loginFile = {
    path: t('keys.entroxCli.loginTitle'),
    content: 'entrox login',
    hint: t('keys.entroxCli.hint')
  }

  if (platform === 'windows') {
    return [
      {
        path: `${installTitle} - Windows - PowerShell`,
        content: getEntroxPowerShellInstallCommand(isCn),
        hint: t('keys.entroxCli.installHint')
      },
      {
        path: `${installTitle} - Windows - Scoop`,
        content: getEntroxScoopInstallCommand(isCn)
      },
      loginFile
    ]
  }

  if (platform === 'linux') {
    return [
      {
        path: `${installTitle} - Linux`,
        content: getEntroxScriptInstallCommand(isCn),
        hint: t('keys.entroxCli.installHint')
      },
      loginFile
    ]
  }

  return [
    {
      path: `${installTitle} - macOS - curl`,
      content: getEntroxScriptInstallCommand(isCn),
      hint: t('keys.entroxCli.installHint')
    },
    {
      path: `${installTitle} - macOS - Homebrew`,
      content: getEntroxHomebrewInstallCommand(isCn)
    },
    loginFile
  ]
}

const copyContent = async (content: string, index: number) => {
  const success = await clipboardCopy(content, t('keys.copied'))
  if (success) {
    copiedIndex.value = index
    setTimeout(() => {
      copiedIndex.value = null
    }, 2000)
  }
}
</script>
