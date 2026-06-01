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
  const repairTapCommand = `ENTROX_HOMEBREW_TAP="$(brew --repository)/Library/Taps/hexonal/homebrew-entrox"
if [ -d "$ENTROX_HOMEBREW_TAP/.git" ]; then
  git -C "$ENTROX_HOMEBREW_TAP" fetch origin
  git -C "$ENTROX_HOMEBREW_TAP" reset --hard origin/main
  git -C "$ENTROX_HOMEBREW_TAP" clean -fd
elif [ -d "$ENTROX_HOMEBREW_TAP" ]; then
  rm -rf "$ENTROX_HOMEBREW_TAP"
fi`
  const installCommand = `${repairTapCommand}
HOMEBREW_NO_AUTO_UPDATE=1 brew tap hexonal/entrox
(HOMEBREW_NO_AUTO_UPDATE=1 brew trust hexonal/entrox || true)
HOMEBREW_NO_AUTO_UPDATE=1 brew upgrade hexonal/entrox/entrox || HOMEBREW_NO_AUTO_UPDATE=1 brew install hexonal/entrox/entrox`

  if (!isCn) {
    return installCommand
  }

  return `# 中国大陆网络：如 GitHub 连接失败，先按你的本机代理地址修改并执行下面两行
# export HTTPS_PROXY="http://127.0.0.1:<你的代理端口>"
# export HTTP_PROXY="$HTTPS_PROXY"

${installCommand}`
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
