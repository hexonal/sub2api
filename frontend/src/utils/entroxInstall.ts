export type EntroxInstallMethodId = 'homebrew' | 'scoop'

export interface EntroxInstallMethod {
  id: EntroxInstallMethodId
  labelKey: string
  command: string
}

export function isCnInstallLocale(locale: unknown): boolean {
  return String(locale || '').toLowerCase().startsWith('zh')
}

export function getEntroxHomebrewInstallCommand(isCn: boolean): string {
  if (!isCn) {
    return 'brew tap hexonal/entrox\nbrew install entrox'
  }

  return `# 中国大陆网络：如 GitHub 连接失败，先按你的本机代理地址修改并执行下面两行
# export HTTPS_PROXY="http://127.0.0.1:<你的代理端口>"
# export HTTP_PROXY="$HTTPS_PROXY"

brew tap hexonal/entrox
brew install entrox`
}

export function getEntroxScoopInstallCommand(isCn: boolean): string {
  if (!isCn) {
    return 'scoop bucket add entrox https://github.com/hexonal/scoop-entrox\nscoop install entrox'
  }

  return `# 中国大陆网络：如 GitHub 连接失败，先按你的本机代理地址修改并执行下面一行
# scoop config proxy 127.0.0.1:<你的代理端口>

scoop bucket add entrox https://github.com/hexonal/scoop-entrox
scoop install entrox`
}

export function getEntroxInstallMethods(isCn: boolean): EntroxInstallMethod[] {
  return [
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
  return `# macOS / Linux (Homebrew)
${getEntroxHomebrewInstallCommand(isCn)}

# Windows (Scoop)
${getEntroxScoopInstallCommand(isCn)}`
}
