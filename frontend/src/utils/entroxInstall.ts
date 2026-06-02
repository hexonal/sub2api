export type EntroxInstallMethodId = 'script' | 'powershell' | 'homebrew' | 'scoop'

export interface EntroxInstallMethod {
  id: EntroxInstallMethodId
  labelKey: string
  command: string
}

export function isCnInstallLocale(locale: unknown): boolean {
  return String(locale || '').toLowerCase().startsWith('zh')
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

export function getEntroxMacLinuxInstallGuide(isCn: boolean): string {
  return `${getEntroxScriptInstallCommand(isCn)}

${getEntroxHomebrewInstallCommand(isCn)}`
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
  return `${getEntroxMacLinuxInstallGuide(isCn)}

${getEntroxWindowsInstallGuide(isCn)}`
}
