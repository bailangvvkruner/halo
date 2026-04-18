import type {
  Category,
  Comment,
  DashboardStats,
  Menu,
  Page,
  Plugin,
  Post,
  Setting,
  Tag,
  Theme,
  User
} from './api'

type Props = {
  posts: Post[]
  themes: Theme[]
  plugins: Plugin[]
  users: User[]
  pages: Page[]
  categories: Category[]
  tags: Tag[]
  menus: Menu[]
  comments: Comment[]
  settings: Setting[]
  stats: DashboardStats | null
  onLogout: () => void
}

export function Dashboard(props: Props) {
  const {
    posts,
    themes,
    plugins,
    users,
    pages,
    categories,
    tags,
    menus,
    comments,
    settings,
    stats,
    onLogout
  } = props

  return (
    <main className="app-shell">
      <section className="hero">
        <div>
          <span className="badge">Halo Go</span>
          <h1>现代化 Go 控制台</h1>
          <p>当前阶段已具备基础后台能力与核心资源管理接口。</p>
        </div>
        <button className="logout-button" onClick={onLogout}>退出登录</button>
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

        <section className="panel compact-panel full-span">
          <h2>设置</h2>
          <ul className="mini-list">
            {settings.map((setting) => (
              <li key={setting.id}>{setting.key} · {setting.value}</li>
            ))}
          </ul>
        </section>
      </section>
    </main>
  )
}
