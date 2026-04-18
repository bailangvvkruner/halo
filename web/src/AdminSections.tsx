import { useState } from 'react'
import { api, type Attachment, type Category, type Comment, type DashboardStats, type Menu, type Page, type Plugin, type PluginScanItem, type Post, type RegistrationPayload, type Reply, type SearchResult, type Setting, type Tag, type Theme, type ThemeScanItem, type ThemeSyncItem, type User } from './api'

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
  const [keyword, setKeyword] = useState('')
  const filtered = posts.filter((post) => {
    if (!keyword.trim()) {
      return true
    }
    const lower = keyword.toLowerCase()
    return post.title.toLowerCase().includes(lower) || post.slug.toLowerCase().includes(lower)
  })

  return (
    <section className="panel posts-panel">
      <div className="section-toolbar">
        <div>
          <h2>文章管理</h2>
          <p className="section-subtitle">更接近原版后台的文章列表视图。</p>
        </div>
        <button>新建文章</button>
      </div>
      <div className="table-toolbar">
        <input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="搜索标题或别名" />
        <span className="table-count">共 {filtered.length} 篇</span>
      </div>
      <div className="table-shell">
        <table className="admin-table">
          <thead>
            <tr>
              <th>标题</th>
              <th>别名</th>
              <th>分类</th>
              <th>标签</th>
              <th>状态</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((post) => (
              <tr key={post.id}>
                <td>
                  <strong>{post.title}</strong>
                  <p className="table-excerpt">{post.excerpt || '暂无摘要'}</p>
                </td>
                <td>{post.slug}</td>
                <td>{post.category || '未分类'}</td>
                <td>{post.tags || '无标签'}</td>
                <td>{post.published ? '已发布' : '草稿'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}

export function PagesSection({ pages }: { pages: Page[] }) {
  return <ContentTableSection title="页面管理" subtitle="页面列表与别名管理。" primaryAction="新建页面" rows={pages.map((page) => ({
    title: page.title,
    description: page.slug,
    metaA: page.slug,
    metaB: page.published ? '已发布' : '草稿',
    metaC: '页面'
  }))} searchPlaceholder="搜索页面标题或别名" columnA="标题" columnB="别名" columnC="状态" columnD="类型" />
}

export function TaxonomiesSection(props: { categories: Category[]; tags: Tag[]; menus: Menu[] }) {
  const { categories, tags, menus } = props
  return (
    <section className="grid-panels dashboard-grid">
      <ContentTableSection title="分类管理" subtitle="分类名称、别名与展示名。" primaryAction="新建分类" rows={categories.map((item) => ({
        title: item.displayName,
        description: item.name,
        metaA: item.slug,
        metaB: item.name,
        metaC: '分类'
      }))} searchPlaceholder="搜索分类名称或别名" columnA="标题" columnB="别名" columnC="名称" columnD="类型" />
      <ContentTableSection title="标签管理" subtitle="标签管理与命名空间。" primaryAction="新建标签" rows={tags.map((item) => ({
        title: item.name,
        description: item.slug,
        metaA: item.slug,
        metaB: item.name,
        metaC: '标签'
      }))} searchPlaceholder="搜索标签名称或别名" columnA="标题" columnB="别名" columnC="名称" columnD="类型" />
      <section className="panel compact-panel full-span">
        <h2>菜单管理</h2>
        <ul className="mini-list">
          {menus.map((item) => (
            <li key={item.id}>{item.name} · {item.items}</li>
          ))}
        </ul>
      </section>
    </section>
  )
}

export function CommentsSection({ comments }: { comments: Comment[] }) {
  const [selectedCommentId, setSelectedCommentId] = useState<number | null>(comments[0]?.id ?? null)
  const [replies, setReplies] = useState<Reply[]>([])
  const [replyForm, setReplyForm] = useState({ author: '管理员', email: 'admin@example.com', content: '' })
  const [actionMessage, setActionMessage] = useState('')

  const loadReplies = async (commentId: number) => {
    setSelectedCommentId(commentId)
    const response = await api.get<Reply[]>(`/comments/${commentId}/replies`)
    setReplies(response.data)
  }

  const createReply = async () => {
    if (!selectedCommentId || !replyForm.content.trim()) {
      return
    }
    await api.post(`/comments/${selectedCommentId}/replies`, replyForm)
    setReplyForm((current) => ({ ...current, content: '' }))
    await loadReplies(selectedCommentId)
    setActionMessage('回复已创建')
  }

  const approveReply = async (replyId: number) => {
    await api.post(`/replies/${replyId}/approve`)
    if (selectedCommentId) {
      await loadReplies(selectedCommentId)
    }
    setActionMessage('回复已通过')
  }

  const rejectReply = async (replyId: number) => {
    await api.post(`/replies/${replyId}/reject`)
    if (selectedCommentId) {
      await loadReplies(selectedCommentId)
    }
    setActionMessage('回复已拒绝')
  }

  return (
    <section className="grid-panels dashboard-grid">
      <section className="panel posts-panel compact-panel">
        <div className="section-toolbar">
          <div>
            <h2>评论管理</h2>
            <p className="section-subtitle">评论列表、内容预览与回复审核。</p>
          </div>
          <button>导出评论</button>
        </div>
        <div className="table-shell">
          <table className="admin-table">
            <thead>
              <tr>
                <th>作者</th>
                <th>邮箱</th>
                <th>内容</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {comments.map((item) => (
                <tr key={item.id}>
                  <td>{item.author}</td>
                  <td>{item.email}</td>
                  <td>{item.content}</td>
                  <td><button className="link-button" onClick={() => loadReplies(item.id)}>查看回复</button></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
      <section className="panel posts-panel compact-panel">
        <div className="section-toolbar">
          <div>
            <h2>回复审核</h2>
            <p className="section-subtitle">按评论查看并审核回复内容。</p>
          </div>
          <button>创建回复</button>
        </div>
        <div className="table-shell">
          <table className="admin-table">
            <thead>
              <tr>
                <th>作者</th>
                <th>内容</th>
                <th>状态</th>
                <th>审核</th>
              </tr>
            </thead>
            <tbody>
              {replies.map((reply) => (
                <tr key={reply.id}>
                  <td>{reply.author}</td>
                  <td>{reply.content}</td>
                  <td>{reply.status}</td>
                  <td>
                    <div className="reply-actions">
                      <button onClick={() => approveReply(reply.id)}>通过</button>
                      <button className="danger-button" onClick={() => rejectReply(reply.id)}>拒绝</button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="inline-form">
          <input value={replyForm.author} onChange={(event) => setReplyForm((current) => ({ ...current, author: event.target.value }))} placeholder="作者" />
          <input value={replyForm.email} onChange={(event) => setReplyForm((current) => ({ ...current, email: event.target.value }))} placeholder="邮箱" />
          <textarea value={replyForm.content} onChange={(event) => setReplyForm((current) => ({ ...current, content: event.target.value }))} placeholder="回复内容" />
          <button onClick={createReply}>创建回复</button>
          {actionMessage ? <p>{actionMessage}</p> : null}
        </div>
      </section>
    </section>
  )
}

export function UsersSection({ users }: { users: User[] }) {
  return <ContentTableSection title="用户管理" subtitle="用户账号与角色管理。" primaryAction="新建用户" rows={users.map((item) => ({
    title: item.username,
    description: item.role,
    metaA: item.username,
    metaB: item.role,
    metaC: '用户'
  }))} searchPlaceholder="搜索用户名" columnA="用户名" columnB="账号" columnC="角色" columnD="类型" />
}

export function PluginsSection({ plugins, themes }: { plugins: Plugin[]; themes: Theme[] }) {
  const [scannedThemes, setScannedThemes] = useState<ThemeScanItem[]>([])
  const [scannedPlugins, setScannedPlugins] = useState<PluginScanItem[]>([])
  const [syncedThemes, setSyncedThemes] = useState<ThemeSyncItem[]>([])

  const scanThemes = async () => {
    const response = await api.get<ThemeScanItem[]>('/themes/scan')
    setScannedThemes(response.data)
  }

  const syncThemes = async () => {
    const response = await api.post<ThemeSyncItem[]>('/themes/sync')
    setSyncedThemes(response.data)
  }

  const scanPlugins = async () => {
    const response = await api.get<PluginScanItem[]>('/plugins/scan')
    setScannedPlugins(response.data)
  }

  return (
    <section className="grid-panels dashboard-grid">
      <section className="panel posts-panel compact-panel">
        <div className="section-toolbar">
          <div>
            <h2>插件管理</h2>
            <p className="section-subtitle">插件列表与目录扫描。</p>
          </div>
          <button onClick={scanPlugins}>扫描插件目录</button>
        </div>
        <div className="table-shell">
          <table className="admin-table">
            <thead>
              <tr>
                <th>插件</th>
                <th>名称</th>
                <th>状态</th>
                <th>扫描结果</th>
              </tr>
            </thead>
            <tbody>
              {plugins.map((item) => (
                <tr key={item.id}>
                  <td>{item.displayName}</td>
                  <td>{item.name}</td>
                  <td>{item.enabled ? '已启用' : '未启用'}</td>
                  <td>{scannedPlugins.find((plugin) => plugin.name === item.name)?.path || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
      <section className="panel posts-panel compact-panel">
        <div className="section-toolbar">
          <div>
            <h2>主题管理</h2>
            <p className="section-subtitle">主题列表、目录扫描与同步。</p>
          </div>
          <div className="reply-actions">
            <button onClick={scanThemes}>扫描主题目录</button>
            <button onClick={syncThemes}>同步到主题列表</button>
          </div>
        </div>
        <div className="table-shell">
          <table className="admin-table">
            <thead>
              <tr>
                <th>主题</th>
                <th>名称</th>
                <th>状态</th>
                <th>扫描结果</th>
              </tr>
            </thead>
            <tbody>
              {themes.map((item) => (
                <tr key={item.id}>
                  <td>{item.displayName}</td>
                  <td>{item.name}</td>
                  <td>{item.activated ? '已启用' : '未启用'}</td>
                  <td>{scannedThemes.find((theme) => theme.name === item.name)?.displayName || syncedThemes.find((theme) => theme.name === item.name)?.displayName || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </section>
  )
}

export function SettingsSection({ settings }: { settings: Setting[] }) {
  return <SimpleSection title="系统设置" items={settings.map((item) => `${item.key} · ${item.value}`)} />
}

export function AttachmentsSection({ attachments }: { attachments: Attachment[] }) {
  return <ContentTableSection title="附件管理" subtitle="附件列表与公开 URL。" primaryAction="上传附件" rows={attachments.map((item) => ({
    title: item.filename,
    description: item.url,
    metaA: item.path,
    metaB: `${item.size} bytes`,
    metaC: '附件'
  }))} searchPlaceholder="搜索附件名称" columnA="文件名" columnB="路径" columnC="大小" columnD="类型" />
}

export function RegistrationSection() {
  const [form, setForm] = useState<RegistrationPayload>({ username: '', email: '', password: '' })
  const [result, setResult] = useState('')
  const [verifyToken, setVerifyToken] = useState('')
  const [verifyResult, setVerifyResult] = useState('')

  const submit = async () => {
    const response = await api.post<{ token: string; message: string }>('/register', form)
    setResult(`注册成功，验证 token：${response.data.token}`)
    setVerifyToken(response.data.token)
  }

  const verify = async () => {
    if (!verifyToken.trim()) {
      return
    }
    const response = await api.get<{ message: string }>(`/register/verify?token=${encodeURIComponent(verifyToken)}`)
    setVerifyResult(response.data.message)
  }

  return (
    <section className="panel compact-panel">
      <h2>注册验证</h2>
      <div className="inline-form">
        <input value={form.username} onChange={(event) => setForm((current) => ({ ...current, username: event.target.value }))} placeholder="用户名" />
        <input value={form.email} onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))} placeholder="邮箱" />
        <input type="password" value={form.password} onChange={(event) => setForm((current) => ({ ...current, password: event.target.value }))} placeholder="密码" />
        <button onClick={submit}>提交注册</button>
        {result ? <p>{result}</p> : null}
        <input value={verifyToken} onChange={(event) => setVerifyToken(event.target.value)} placeholder="验证 token" />
        <button onClick={verify}>执行验证</button>
        {verifyResult ? <p>{verifyResult}</p> : null}
      </div>
    </section>
  )
}

export function SearchSection() {
  const [keyword, setKeyword] = useState('')
  const [result, setResult] = useState<SearchResult>({ posts: [], pages: [] })

  const search = async () => {
    const response = await api.get<SearchResult>(`/search?keyword=${encodeURIComponent(keyword)}`)
    setResult(response.data)
  }

  return (
    <section className="panel compact-panel">
      <h2>搜索</h2>
      <div className="inline-form">
        <input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder="输入关键词" />
        <button onClick={search}>执行搜索</button>
      </div>
      <div className="grid-panels dashboard-grid">
        <SimpleSection title="文章结果" items={result.posts.map((item) => `${item.title} · ${item.slug}`)} />
        <SimpleSection title="页面结果" items={result.pages.map((item) => `${item.title} · ${item.slug}`)} />
      </div>
    </section>
  )
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

function ContentTableSection(props: {
  title: string
  subtitle: string
  primaryAction: string
  rows: Array<{ title: string; description: string; metaA: string; metaB: string; metaC: string }>
  searchPlaceholder: string
  columnA: string
  columnB: string
  columnC: string
  columnD: string
}) {
  const { title, subtitle, primaryAction, rows, searchPlaceholder, columnA, columnB, columnC, columnD } = props
  const [keyword, setKeyword] = useState('')
  const filtered = rows.filter((row) => {
    if (!keyword.trim()) {
      return true
    }
    const lower = keyword.toLowerCase()
    return row.title.toLowerCase().includes(lower) || row.metaA.toLowerCase().includes(lower)
  })

  return (
    <section className="panel posts-panel compact-panel">
      <div className="section-toolbar">
        <div>
          <h2>{title}</h2>
          <p className="section-subtitle">{subtitle}</p>
        </div>
        <button>{primaryAction}</button>
      </div>
      <div className="table-toolbar">
        <input value={keyword} onChange={(event) => setKeyword(event.target.value)} placeholder={searchPlaceholder} />
        <span className="table-count">共 {filtered.length} 条</span>
      </div>
      <div className="table-shell">
        <table className="admin-table">
          <thead>
            <tr>
              <th>{columnA}</th>
              <th>{columnB}</th>
              <th>{columnC}</th>
              <th>{columnD}</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((row) => (
              <tr key={`${row.title}-${row.metaA}`}>
                <td>
                  <strong>{row.title}</strong>
                  <p className="table-excerpt">{row.description}</p>
                </td>
                <td>{row.metaA}</td>
                <td>{row.metaB}</td>
                <td>{row.metaC}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </section>
  )
}
