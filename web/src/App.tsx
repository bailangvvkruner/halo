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
  type Tag,
  type Theme,
  type User
} from './api'

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
  const [stats, setStats] = useState<DashboardStats | null>(null)

  useEffect(() => {
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
    api.get<DashboardStats>('/dashboard/stats').then((response) => setStats(response.data))
  }, [])

  return (
    <main className="app-shell">
      <section className="hero">
        <div>
          <span className="badge">Halo Go</span>
          <h1>现代化 Go 重构进行中</h1>
          <p>当前版本已完成旧项目归档、Go 服务启动、基础 API、Docker 构建与新前端骨架。</p>
          <p>健康状态：{health?.status ?? '加载中'}</p>
        </div>
      </section>

      <section className="grid-panels dashboard-grid">
        <section className="panel compact-panel">
          <h2>概览</h2>
          <ul className="mini-list">
            <li>文章：{stats?.posts ?? 0}</li>
            <li>页面：{stats?.pages ?? 0}</li>
            <li>分类：{stats?.categories ?? 0}</li>
            <li>标签：{stats?.tags ?? 0}</li>
            <li>用户：{stats?.users ?? 0}</li>
            <li>主题：{stats?.themes ?? 0}</li>
            <li>插件：{stats?.plugins ?? 0}</li>
          </ul>
        </section>
      </section>

      <section className="panel">
        <h2>文章列表</h2>
        {posts.length === 0 ? (
          <p>当前还没有文章。</p>
        ) : (
          <ul className="list">
            {posts.map((post) => (
              <li key={post.id} className="card">
                <h3>{post.title}</h3>
                <p>{post.excerpt || '暂无摘要'}</p>
                <span>{post.slug} · {post.category} · {post.tags}</span>
              </li>
            ))}
          </ul>
        )}
      </section>

      <section className="grid-panels">
        <section className="panel compact-panel">
          <h2>主题</h2>
          <ul className="mini-list">
            {themes.map((theme) => (
              <li key={theme.id}>{theme.displayName} {theme.activated ? '已启用' : '未启用'}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel">
          <h2>插件</h2>
          <ul className="mini-list">
            {plugins.map((plugin) => (
              <li key={plugin.id}>{plugin.displayName} {plugin.enabled ? '已启用' : '未启用'}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel">
          <h2>用户</h2>
          <ul className="mini-list">
            {users.map((user) => (
              <li key={user.id}>{user.username} · {user.role}</li>
            ))}
          </ul>
        </section>
      </section>

      <section className="grid-panels secondary-grid">
        <section className="panel compact-panel">
          <h2>页面</h2>
          <ul className="mini-list">
            {pages.map((page) => (
              <li key={page.id}>{page.title} · {page.slug}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel">
          <h2>分类</h2>
          <ul className="mini-list">
            {categories.map((category) => (
              <li key={category.id}>{category.displayName} · {category.slug}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel">
          <h2>标签</h2>
          <ul className="mini-list">
            {tags.map((tag) => (
              <li key={tag.id}>{tag.name}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel full-span">
          <h2>菜单</h2>
          <ul className="mini-list">
            {menus.map((menu) => (
              <li key={menu.id}>{menu.name} · {menu.items}</li>
            ))}
          </ul>
        </section>

        <section className="panel compact-panel full-span">
          <h2>评论</h2>
          <ul className="mini-list">
            {comments.map((comment) => (
              <li key={comment.id}>{comment.author} · {comment.content}</li>
            ))}
          </ul>
        </section>
      </section>
    </main>
  )
}
