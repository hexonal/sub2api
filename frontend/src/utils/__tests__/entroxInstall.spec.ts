import { describe, expect, it } from 'vitest'
import {
  detectEntroxClientPlatform,
  getDefaultEntroxInstallMethod,
  getDefaultEntroxInstallPlatformTab,
  getEntroxScriptInstallCommand
} from '../entroxInstall'

describe('entroxInstall', () => {
  it('detects common client platforms', () => {
    expect(detectEntroxClientPlatform({ platform: 'Win32' })).toBe('windows')
    expect(detectEntroxClientPlatform({ platform: 'MacIntel' })).toBe('macos')
    expect(detectEntroxClientPlatform({ platform: 'Linux x86_64' })).toBe('linux')
    expect(detectEntroxClientPlatform({ platform: 'X11' })).toBe('linux')
    expect(detectEntroxClientPlatform({ platform: 'Unknown' })).toBe('unknown')
  })

  it('defaults Windows users to PowerShell and Unix users to curl installer', () => {
    expect(getDefaultEntroxInstallMethod({ platform: 'Win32' })).toBe('powershell')
    expect(getDefaultEntroxInstallMethod({ platform: 'MacIntel' })).toBe('script')
    expect(getDefaultEntroxInstallMethod({ platform: 'Linux x86_64' })).toBe('script')
    expect(getDefaultEntroxInstallPlatformTab({ platform: 'Win32' })).toBe('windows')
    expect(getDefaultEntroxInstallPlatformTab({ platform: 'MacIntel' })).toBe('macos')
    expect(getDefaultEntroxInstallPlatformTab({ platform: 'Linux x86_64' })).toBe('linux')
    expect(getEntroxScriptInstallCommand(false)).toBe('curl -fsSL https://entrox.996icu.wiki/install | bash')
  })
})
