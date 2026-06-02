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

export function getEntroxScoopInstallCommand(isCn: boolean): string {
  if (!isCn) {
    return '# Install Scoop first if it is not already installed\nSet-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser\nInvoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression\n\nscoop bucket add entrox https://github.com/hexonal/scoop-entrox\nscoop install entrox'
  }

  return `# 如果尚未安装 Scoop，先用 PowerShell 执行下面两行
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

# 中国大陆网络：如 GitHub 连接失败，先按你的本机代理地址修改并执行下面一行
# scoop config proxy 127.0.0.1:<你的代理端口>

scoop bucket add entrox https://github.com/hexonal/scoop-entrox
scoop install entrox`
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
  return `# macOS / Linux (official installer)
${getEntroxScriptInstallCommand(isCn)}

# Windows PowerShell (official installer)
${getEntroxPowerShellInstallCommand(isCn)}

# macOS fallback (Homebrew)
${getEntroxHomebrewInstallCommand(isCn)}

# Windows fallback (Scoop)
${getEntroxScoopInstallCommand(isCn)}`
}
