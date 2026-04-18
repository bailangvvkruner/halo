import { useState } from 'react'
import { api, type Category, type Comment, type DashboardStats, type Menu, type Page, type Plugin, type Post, type RegistrationPayload, type Reply, type Setting, type Tag, type Theme, type User } from './api'

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
      <section className="panel compact-panel">
        <h2>评论管理</h2>
        <ul className="mini-list">
          {comments.map((item) => (
            <li key={item.id}>
              <button className="link-button" onClick={() => loadReplies(item.id)}>{item.author} · {item.content}</button>
            </li>
          ))}
        </ul>
      </section>
      <section className="panel compact-panel">
        <h2>回复</h2>
        <ul className="mini-list">
          {replies.map((reply) => (
            <li key={reply.id} className="reply-item">
              <span>{reply.author} · {reply.content} · {reply.status}</span>
              <div className="reply-actions">
                <button onClick={() => approveReply(reply.id)}>通过</button>
                <button className="danger-button" onClick={() => rejectReply(reply.id)}>拒绝</button>
              </div>
            </li>
          ))}
        </ul>
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
