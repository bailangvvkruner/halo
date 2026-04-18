import { useEffect, useState } from 'react'
import {
  api,
  type Category,
  type Comment,
  type DashboardStats,
  type Health,
  type Attachment,
  type Menu,
  type Page,
  type Plugin,
  type Post,
  type Setting,
  type Tag,
  type Theme,
  type User
} from './api'
import { clearToken, isLoggedIn } from './auth'
import { Dashboard } from './Dashboard'
import { LoginForm } from './LoginForm'
import { SetupForm } from './SetupForm'

export function App() {
  const [health, setHealth] = useState<Health | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [themes, setThemes] = useState<Theme[]>([])
  const [plugins, setPlugins] = useState<Plugin[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [attachments, setAttachments] = useState<Attachment[]>([])
  const [pages, setPages] = useState<Page[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [menus, setMenus] = useState<Menu[]>([])
  const [comments, setComments] = useState<Comment[]>([])
  const [settings, setSettings] = useState<Setting[]>([])
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [authenticated, setAuthenticated] = useState(isLoggedIn())
  const [initialized, setInitialized] = useState<boolean | null>(null)
  const [requestError, setRequestError] = useState('')

  useEffect(() => {
    api.get<{ initialized: boolean }>('/setup/status').then((response) => {
      setInitialized(response.data.initialized)
    })
  }, [])

  useEffect(() => {
	if (initialized === false && window.location.pathname !== '/system/setup') {
		window.history.replaceState({}, '', '/system/setup')
	}
	if (initialized && window.location.pathname === '/system/setup') {
		window.history.replaceState({}, '', '/')
	}
  }, [initialized])

  useEffect(() => {
    if (!authenticated || !initialized) {
      return
    }

    setRequestError('')
    api.get<Health>('/health').then((response) => setHealth(response.data)).catch(handleRequestError)
    api.get<Post[]>('/posts').then((response) => setPosts(response.data)).catch(handleRequestError)
    api.get<Theme[]>('/themes').then((response) => setThemes(response.data)).catch(handleRequestError)
    api.get<Plugin[]>('/plugins').then((response) => setPlugins(response.data)).catch(handleRequestError)
    api.get<User[]>('/users').then((response) => setUsers(response.data)).catch(handleRequestError)
    api.get<Attachment[]>('/attachments').then((response) => setAttachments(response.data)).catch(handleRequestError)
    api.get<Page[]>('/pages').then((response) => setPages(response.data)).catch(handleRequestError)
    api.get<Category[]>('/categories').then((response) => setCategories(response.data)).catch(handleRequestError)
    api.get<Tag[]>('/tags').then((response) => setTags(response.data)).catch(handleRequestError)
    api.get<Menu[]>('/menus').then((response) => setMenus(response.data)).catch(handleRequestError)
    api.get<Comment[]>('/comments').then((response) => setComments(response.data)).catch(handleRequestError)
    api.get<Setting[]>('/settings').then((response) => setSettings(response.data)).catch(handleRequestError)
    api.get<DashboardStats>('/dashboard/stats').then((response) => setStats(response.data)).catch(handleRequestError)
  }, [authenticated, initialized])

  const handleRequestError = (error: unknown) => {
    const status = (error as { response?: { status?: number } })?.response?.status
    if (status === 401) {
      setRequestError('请求被拒绝：请重新登录。')
      clearToken()
      setAuthenticated(false)
      return
    }
    if (status === 403) {
      setRequestError('请求被拒绝：当前账号没有对应权限。')
      return
    }
    setRequestError('请求失败，请检查后端服务。')
  }

  if (initialized === null) {
    return <section className="panel login-panel"><p>加载初始化状态中...</p></section>
  }

  if (!initialized) {
    return <SetupForm onSuccess={() => setInitialized(true)} />
  }

  if (!authenticated) {
    return <LoginForm onSuccess={() => setAuthenticated(true)} />
  }

  return (
    <Dashboard
      posts={posts}
      themes={themes}
      plugins={plugins}
      users={users}
      attachments={attachments}
      pages={pages}
      categories={categories}
      tags={tags}
      menus={menus}
      comments={comments}
      settings={settings}
      stats={stats}
      requestError={requestError}
      onLogout={() => {
        clearToken()
        setAuthenticated(false)
      }}
    />
  )
}
