export type EntroxInstallMethodId = 'script' | 'powershell' | 'homebrew' | 'scoop'
export type EntroxClientPlatform = 'windows' | 'macos' | 'linux' | 'unknown'
export type EntroxInstallPlatformTab = 'macos' | 'linux' | 'windows'

export interface EntroxInstallMethod {
  id: EntroxInstallMethodId
  labelKey: string
  command: string
}

interface EntroxClientPlatformInput {
  platform?: string
  userAgent?: string
}

export function isCnInstallLocale(locale: unknown): boolean {
  return String(locale || '').toLowerCase().startsWith('zh')
}

export function detectEntroxClientPlatform(input?: EntroxClientPlatformInput): EntroxClientPlatform {
  const platform = input
    ? input.platform ?? ''
    : typeof navigator === 'undefined'
      ? ''
      : navigator.platform
  const userAgent = input
    ? input.userAgent ?? ''
    : typeof navigator === 'undefined'
      ? ''
      : navigator.userAgent
  const value = `${platform} ${userAgent}`.toLowerCase()

  if (value.includes('win')) return 'windows'
  if (value.includes('mac')) return 'macos'
  if (value.includes('linux') || value.includes('x11')) return 'linux'
  return 'unknown'
}

export function getDefaultEntroxInstallMethod(input?: EntroxClientPlatformInput): EntroxInstallMethodId {
  return detectEntroxClientPlatform(input) === 'windows' ? 'powershell' : 'script'
}

export function getDefaultEntroxInstallPlatformTab(input?: EntroxClientPlatformInput): EntroxInstallPlatformTab {
  const platform = detectEntroxClientPlatform(input)
  if (platform === 'windows') return 'windows'
  if (platform === 'linux') return 'linux'
  return 'macos'
}

export function getEntroxHomebrewInstallCommand(_isCn: boolean): string {
  return `HOMEBREW_NO_AUTO_UPDATE=1 brew tap hexonal/entrox
(HOMEBREW_NO_AUTO_UPDATE=1 brew trust hexonal/entrox || true)
HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox`
}

export function getEntroxScriptInstallCommand(_isCn: boolean): string {
  return 'curl -fsSL https://entrox.996icu.wiki/install | bash'
}

export function getEntroxPowerShellInstallCommand(_isCn: boolean): string {
  return 'irm https://entrox.996icu.wiki/install.ps1 | iex'
}

export function getEntroxScoopInstallCommand(_isCn: boolean): string {
  return `Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

scoop bucket add entrox https://github.com/hexonal/scoop-entrox
scoop install entrox`
}

export function getEntroxMacosInstallGuide(isCn: boolean): string {
  return `${getEntroxScriptInstallCommand(isCn)}

${getEntroxHomebrewInstallCommand(isCn)}`
}

export function getEntroxLinuxInstallGuide(isCn: boolean): string {
  return getEntroxScriptInstallCommand(isCn)
}

export function getEntroxWindowsInstallGuide(isCn: boolean): string {
  return `${getEntroxPowerShellInstallCommand(isCn)}

${getEntroxScoopInstallCommand(isCn)}`
}

export function getEntroxInstallMethods(isCn: boolean): EntroxInstallMethod[] {
  return [
    {
      id: 'script',
      labelKey: 'home.install.script',
      command: getEntroxScriptInstallCommand(isCn)
    },
    {
      id: 'powershell',
      labelKey: 'home.install.powershell',
      command: getEntroxPowerShellInstallCommand(isCn)
    },
    {
      id: 'homebrew',
      labelKey: 'home.install.homebrew',
      command: getEntroxHomebrewInstallCommand(isCn)
    },
    {
      id: 'scoop',
      labelKey: 'home.install.scoop',
      command: getEntroxScoopInstallCommand(isCn)
    }
  ]
}

export function getEntroxInstallGuide(isCn: boolean): string {
  return `${getEntroxMacosInstallGuide(isCn)}

${getEntroxLinuxInstallGuide(isCn)}

${getEntroxWindowsInstallGuide(isCn)}`
}
