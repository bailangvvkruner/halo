import { useEffect, useState } from 'react'
import { api } from './api'

const messages = {
  'zh-CN': {
    title: '系统初始化',
    language: '语言',
    externalUrl: '外部访问地址',
    siteTitle: '站点标题',
    username: '用户名',
    email: '电子邮箱',
    password: '密码',
    confirmPassword: '确认密码',
    submit: '初始化',
    warningTitle: '警告：正在使用 SQLite 数据库',
    warningContent: 'SQLite 数据库仅适用于开发环境和测试环境，不推荐在生产环境中使用。如果必须使用，请按时进行数据备份。',
    setupError: '初始化失败，请检查输入内容与密码确认'
  },
  'zh-TW': {
    title: '系統初始化',
    language: '語言',
    externalUrl: '外部訪問地址',
    siteTitle: '站點標題',
    username: '使用者名稱',
    email: '電子郵件',
    password: '密碼',
    confirmPassword: '確認密碼',
    submit: '初始化',
    warningTitle: '警告：正在使用 SQLite 資料庫',
    warningContent: 'SQLite 資料庫僅適用於開發環境和測試環境，不建議在生產環境中使用。如果必須使用，請按時進行資料備份。',
    setupError: '初始化失敗，請檢查輸入內容與密碼確認'
  },
  en: {
    title: 'Setup',
    language: 'Language',
    externalUrl: 'External URL',
    siteTitle: 'Site title',
    username: 'Username',
    email: 'Email',
    password: 'Password',
    confirmPassword: 'Confirm Password',
    submit: 'Setup',
    warningTitle: 'Warning: Using SQLite Database',
    warningContent: 'SQLite is suitable for development and testing environments. It is not recommended for production environments. If you must use it, please back up your data regularly.',
    setupError: 'Setup failed. Please check your inputs and password confirmation.'
  },
  es: {
    title: 'Configuración',
    language: 'Idioma',
    externalUrl: 'URL Externa',
    siteTitle: 'Título del Sitio',
    username: 'Nombre de Usuario',
    email: 'Correo Electrónico',
    password: 'Contraseña',
    confirmPassword: 'Confirmar Contraseña',
    submit: 'Configurar',
    warningTitle: 'Advertencia: usando SQLite',
    warningContent: 'SQLite solo es adecuada para entornos de desarrollo y prueba. No se recomienda en producción. Si debe usarla, haga copias de seguridad regularmente.',
    setupError: 'La configuración falló. Revise los datos y la confirmación de contraseña.'
  }
} as const

const languageOptions = [
  { value: 'en', label: 'English' },
  { value: 'es', label: 'Español' },
  { value: 'zh-CN', label: '简体中文' },
  { value: 'zh-TW', label: '繁体中文' }
]

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

  const currentMessages = messages[form.language as keyof typeof messages] ?? messages['zh-CN']

  useEffect(() => {
    const currentUrl = new URL(window.location.href)
    const language = currentUrl.searchParams.get('language')
    const externalUrl = `${window.location.protocol}//${window.location.host}`
    setForm((current) => ({
      ...current,
      language: language || current.language,
      baseURL: externalUrl
    }))
  }, [])

  const updateField = (key: keyof typeof form, value: string) => {
    setForm((current) => ({ ...current, [key]: value }))
  }

  const updateLanguage = (value: string) => {
    updateField('language', value)
    const currentUrl = new URL(window.location.href)
    currentUrl.searchParams.set('language', value)
    window.history.replaceState(null, '', currentUrl.toString())
    window.location.reload()
  }

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)
    setError('')

    try {
      await api.post('/setup', form)
      onSuccess()
    } catch {
      setError(currentMessages.setupError)
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="setup-shell">
      <section className="setup-card">
        <div className="setup-header">
          <h1>{currentMessages.title}</h1>
        </div>

        <div className="setup-warning">
          <strong>{currentMessages.warningTitle}</strong>
          <p>{currentMessages.warningContent}</p>
        </div>

        <form className="setup-form" onSubmit={submit}>
          <div className="setup-field">
            <label htmlFor="language">{currentMessages.language}</label>
            <select id="language" value={form.language} onChange={(event) => updateLanguage(event.target.value)}>
              {languageOptions.map((option) => (
                <option key={option.value} value={option.value}>{option.label}</option>
              ))}
            </select>
          </div>
          <div className="setup-field">
            <label htmlFor="baseURL">{currentMessages.externalUrl}</label>
            <input id="baseURL" value={form.baseURL} onChange={(event) => updateField('baseURL', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="siteTitle">{currentMessages.siteTitle}</label>
            <input id="siteTitle" value={form.siteTitle} onChange={(event) => updateField('siteTitle', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="username">{currentMessages.username}</label>
            <input id="username" value={form.username} onChange={(event) => updateField('username', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="email">{currentMessages.email}</label>
            <input id="email" value={form.email} onChange={(event) => updateField('email', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="password">{currentMessages.password}</label>
            <input id="password" type="password" value={form.password} onChange={(event) => updateField('password', event.target.value)} />
          </div>
          <div className="setup-field">
            <label htmlFor="confirmPassword">{currentMessages.confirmPassword}</label>
            <input id="confirmPassword" type="password" value={form.confirmPassword} onChange={(event) => updateField('confirmPassword', event.target.value)} />
          </div>
          {error ? <p className="error-text">{error}</p> : null}
          <div className="setup-actions">
            <button className="setup-submit" type="submit" disabled={loading}>{loading ? `${currentMessages.submit}...` : currentMessages.submit}</button>
          </div>
        </form>
      </section>
    </main>
  )
}
