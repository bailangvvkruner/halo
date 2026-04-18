import { useEffect, useState } from 'react'
import {
  api,
  type Category,
  type Comment,
  type DashboardStats,
  type Health,
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

export function App() {
  const [health, setHealth] = useState<Health | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [themes, setThemes] = useState<Theme[]>([])
  const [plugins, setPlugins] = useState<Plugin[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [pages, setPages] = useState<Page[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [menus, setMenus] = useState<Menu[]>([])
  const [comments, setComments] = useState<Comment[]>([])
  const [settings, setSettings] = useState<Setting[]>([])
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [authenticated, setAuthenticated] = useState(isLoggedIn())

  useEffect(() => {
    if (!authenticated) {
      return
    }

    api.get<Health>('/health').then((response) => setHealth(response.data))
    api.get<Post[]>('/posts').then((response) => setPosts(response.data))
    api.get<Theme[]>('/themes').then((response) => setThemes(response.data))
    api.get<Plugin[]>('/plugins').then((response) => setPlugins(response.data))
    api.get<User[]>('/users').then((response) => setUsers(response.data))
    api.get<Page[]>('/pages').then((response) => setPages(response.data))
    api.get<Category[]>('/categories').then((response) => setCategories(response.data))
    api.get<Tag[]>('/tags').then((response) => setTags(response.data))
    api.get<Menu[]>('/menus').then((response) => setMenus(response.data))
    api.get<Comment[]>('/comments').then((response) => setComments(response.data))
    api.get<Setting[]>('/settings').then((response) => setSettings(response.data))
    api.get<DashboardStats>('/dashboard/stats').then((response) => setStats(response.data))
  }, [authenticated])

  if (!authenticated) {
    return <LoginForm onSuccess={() => setAuthenticated(true)} />
  }

  return (
    <Dashboard
      posts={posts}
      themes={themes}
      plugins={plugins}
      users={users}
      pages={pages}
      categories={categories}
      tags={tags}
      menus={menus}
      comments={comments}
      settings={settings}
      stats={stats}
      onLogout={() => {
        clearToken()
        setAuthenticated(false)
      }}
    />
  )
}
