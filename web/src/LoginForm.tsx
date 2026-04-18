import { useState } from 'react'
import { api, type LoginResponse } from './api'
import { setToken } from './auth'

type Props = {
  onSuccess: () => void
}

export function LoginForm({ onSuccess }: Props) {
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const submit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)
    setError('')

    try {
      const response = await api.post<LoginResponse>('/login', { username, password })
      setToken(response.data.token)
      onSuccess()
    } catch {
      setError('登录失败，请检查用户名和密码')
    } finally {
      setLoading(false)
    }
  }

  return (
    <section className="panel login-panel">
      <h2>登录后台</h2>
      <form className="login-form" onSubmit={submit}>
        <label>
          <span>用户名</span>
          <input value={username} onChange={(event) => setUsername(event.target.value)} />
        </label>
        <label>
          <span>密码</span>
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} />
        </label>
        {error ? <p className="error-text">{error}</p> : null}
        <button type="submit" disabled={loading}>{loading ? '登录中' : '登录'}</button>
      </form>
    </section>
  )
}
