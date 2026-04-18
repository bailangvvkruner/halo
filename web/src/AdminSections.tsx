import type { Category, Comment, DashboardStats, Menu, Page, Plugin, Post, Setting, Tag, Theme, User } from './api'

export function OverviewSection({ stats }: { stats: DashboardStats | null }) {
  return (
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
  )
}

export function PostsSection({ posts }: { posts: Post[] }) {
  return (
    <section className="panel">
      <h2>文章管理</h2>
      <ul className="list">
        {posts.map((post) => (
          <li key={post.id} className="card">
            <h3>{post.title}</h3>
            <p>{post.excerpt || '暂无摘要'}</p>
            <span>{post.slug} · {post.category} · {post.tags}</span>
          </li>
        ))}
      </ul>
    </section>
  )
}

export function PagesSection({ pages }: { pages: Page[] }) {
  return <SimpleSection title="页面管理" items={pages.map((page) => `${page.title} · ${page.slug}`)} />
}

export function TaxonomiesSection(props: { categories: Category[]; tags: Tag[]; menus: Menu[] }) {
  const { categories, tags, menus } = props
  return (
    <section className="grid-panels secondary-grid">
      <SimpleSection title="分类" items={categories.map((item) => `${item.displayName} · ${item.slug}`)} />
      <SimpleSection title="标签" items={tags.map((item) => item.name)} />
      <SimpleSection title="菜单" items={menus.map((item) => `${item.name} · ${item.items}`)} fullSpan />
    </section>
  )
}

export function CommentsSection({ comments }: { comments: Comment[] }) {
  return <SimpleSection title="评论管理" items={comments.map((item) => `${item.author} · ${item.content}`)} />
}

export function UsersSection({ users }: { users: User[] }) {
  return <SimpleSection title="用户管理" items={users.map((item) => `${item.username} · ${item.role}`)} />
}

export function PluginsSection({ plugins, themes }: { plugins: Plugin[]; themes: Theme[] }) {
  return (
    <section className="grid-panels">
      <SimpleSection title="插件" items={plugins.map((item) => `${item.displayName} · ${item.enabled ? '已启用' : '未启用'}`)} />
      <SimpleSection title="主题" items={themes.map((item) => `${item.displayName} · ${item.activated ? '已启用' : '未启用'}`)} />
    </section>
  )
}

export function SettingsSection({ settings }: { settings: Setting[] }) {
  return <SimpleSection title="系统设置" items={settings.map((item) => `${item.key} · ${item.value}`)} />
}

function SimpleSection({ title, items, fullSpan = false }: { title: string; items: string[]; fullSpan?: boolean }) {
  return (
    <section className={fullSpan ? 'panel compact-panel full-span' : 'panel compact-panel'}>
      <h2>{title}</h2>
      <ul className="mini-list">
        {items.map((item) => (
          <li key={item}>{item}</li>
        ))}
      </ul>
    </section>
  )
}
