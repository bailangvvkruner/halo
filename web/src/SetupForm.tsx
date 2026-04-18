import { useState } from 'react'
import { api } from './api'

type Props = {
  onSuccess: () => void
}

export function SetupForm({ onSuccess }: Props) {
  const [form, setForm] = useState({
    language: 'zh-CN',
    baseURL: 'http://localhost:8091',
    siteTitle: 'Halo Go',
    username: 'admin',
    email: 'admin@example.com',
    password: 'admin123'
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const updateField = (key: keyof typeof form, value: string) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)
    setError('')

    try {
      await api.post('/setup', form)
      onSuccess()
    } catch {
      setError('初始化失败，请检查输入内容')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="panel login-panel">
      <h2>系统初始化</h2>
      <p className="warning-text">警告：当前正在使用 SQLite，本地开发可用，生产环境请谨慎选择数据库与备份策略。</p>
      <form className="login-form" onSubmit={submit}>
        <label>
          <span>语言</span>
          <input value={form.language} onChange={(event) => updateField('language', event.target.value)} />
        </label>
        <label>
          <span>外部访问地址</span>
          <input value={form.baseURL} onChange={(event) => updateField('baseURL', event.target.value)} />
        </label>
        <label>
          <span>站点标题</span>
          <input value={form.siteTitle} onChange={(event) => updateField('siteTitle', event.target.value)} />
        </label>
        <label>
          <span>用户名</span>
          <input value={form.username} onChange={(event) => updateField('username', event.target.value)} />
        </label>
        <label>
          <span>电子邮箱</span>
          <input value={form.email} onChange={(event) => updateField('email', event.target.value)} />
        </label>
        <label>
          <span>密码</span>
          <input type="password" value={form.password} onChange={(event) => updateField('password', event.target.value)} />
        </label>
        {error ? <p className="error-text">{error}</p> : null}
        <button type="submit" disabled={loading}>{loading ? '初始化中' : '初始化'}</button>
      </form>
    </section>
  )
}
