import { useEffect, useState } from 'react'
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
    password: 'admin123',
    confirmPassword: 'admin123'
  })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const externalUrl = `${window.location.protocol}//${window.location.host}`
    setForm((current) => ({
      ...current,
      baseURL: externalUrl
    }))
  }, [])

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
      setError('初始化失败，请检查输入内容与密码确认')
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="setup-shell">
      <section className="setup-card">
        <div className="setup-header">
          <h1>系统初始化</h1>
        </div>

        <div className="setup-warning">
          <strong>警告：正在使用 SQLite 数据库</strong>
          <p>SQLite 数据库仅适用于开发环境和测试环境，不推荐在生产环境中使用。如果必须使用，请按时进行数据备份。</p>
        </div>

        <form className="setup-form" onSubmit={submit}>
          <div className="setup-field">
            <label htmlFor="language">语言</label>
            <input id="language" value={form.language} onChange={(event) => updateField('language', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="baseURL">外部访问地址</label>
            <input id="baseURL" value={form.baseURL} onChange={(event) => updateField('baseURL', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="siteTitle">站点标题</label>
            <input id="siteTitle" value={form.siteTitle} onChange={(event) => updateField('siteTitle', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="username">用户名</label>
            <input id="username" value={form.username} onChange={(event) => updateField('username', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="email">电子邮箱</label>
            <input id="email" value={form.email} onChange={(event) => updateField('email', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="password">密码</label>
            <input id="password" type="password" value={form.password} onChange={(event) => updateField('password', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="confirmPassword">确认密码</label>
            <input id="confirmPassword" type="password" value={form.confirmPassword} onChange={(event) => updateField('confirmPassword', event.target.value)} />
          </div>
          {error ? <p className="error-text">{error}</p> : null}
          <div className="setup-actions">
            <button className="setup-submit" type="submit" disabled={loading}>{loading ? '初始化中' : '初始化'}</button>
          </div>
        </form>
      </section>
    </main>
  )
}
