import type { ReactNode } from 'react'

type NavKey = 'overview' | 'posts' | 'pages' | 'taxonomies' | 'comments' | 'users' | 'plugins' | 'settings' | 'register' | 'search' | 'attachments' | 'post-editor' | 'page-editor'

type Props = {
  active: NavKey
  onNavigate: (key: NavKey) => void
  onLogout: () => void
  children: ReactNode
}

const items: Array<{ key: NavKey; label: string }> = [
  { key: 'overview', label: '概览' },
  { key: 'posts', label: '文章' },
  { key: 'pages', label: '页面' },
  { key: 'taxonomies', label: '分类标签' },
  { key: 'comments', label: '评论' },
  { key: 'users', label: '用户' },
  { key: 'plugins', label: '插件主题' },
  { key: 'settings', label: '设置' },
  { key: 'register', label: '注册验证' },
  { key: 'search', label: '搜索' },
  { key: 'attachments', label: '附件' },
  { key: 'post-editor', label: '文章编辑' },
  { key: 'page-editor', label: '页面编辑' }
]

export function AdminShell({ active, onNavigate, onLogout, children }: Props) {
  return (
    <div className="admin-layout">
      <aside className="admin-sidebar">
        <div className="admin-brand">
          <span className="badge">Halo Go</span>
          <h1>控制台</h1>
        </div>
        <nav className="admin-nav">
          {items.map((item) => (
            <button
              key={item.key}
              className={item.key === active ? 'nav-item active' : 'nav-item'}
              onClick={() => onNavigate(item.key)}
            >
              {item.label}
            </button>
          ))}
        </nav>
        <button className="logout-button sidebar-logout" onClick={onLogout}>退出登录</button>
      </aside>
      <main className="admin-main">
        <header className="admin-topbar">
          <div>
            <p className="topbar-label">Console</p>
            <h2>{items.find((item) => item.key === active)?.label ?? '控制台'}</h2>
          </div>
          <div className="topbar-actions">
            <button className="topbar-button">查看站点</button>
            <button className="topbar-button">快捷操作</button>
          </div>
        </header>
        {children}
      </main>
    </div>
  )
}
