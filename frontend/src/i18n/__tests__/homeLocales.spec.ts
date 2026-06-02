import { describe, expect, it } from 'vitest'

import zh from '../locales/zh'

describe('home locale copy', () => {
  it('describes unified gateway access as product certificate in zh', () => {
    expect(zh.home.features.unifiedGatewayDesc).toBe(
      '获取一个产品证书，即可调用所有已接入的 AI 模型，无需分别申请。'
    )
    expect(zh.home.features.unifiedGatewayDesc).not.toContain('API 密钥')
  })
})
