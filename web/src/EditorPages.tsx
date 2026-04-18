import { useState } from 'react'
import { api, type Page, type Post } from './api'

export function PostEditorSection({ posts }: { posts: Post[] }) {
  const current = posts[0]
  const [form, setForm] = useState({
    title: current?.title ?? '',
    slug: current?.slug ?? '',
    excerpt: current?.excerpt ?? '',
    content: current?.content ?? '',
    category: current?.category ?? '',
    tags: current?.tags ?? '',
    published: current?.published ?? false
  })
  const [message, setMessage] = useState('')

  const saveDraft = async () => {
    if (!current) {
      return
    }
    await api.put(`/posts/${current.id}`, {
      title: form.title,
      slug: form.slug,
      excerpt: form.excerpt,
      content: form.content,
      published: false
    })
    setMessage('草稿已保存')
  }

  const publish = async () => {
    if (!current) {
      return
    }
    await api.put(`/posts/${current.id}`, {
      title: form.title,
      slug: form.slug,
      excerpt: form.excerpt,
      content: form.content,
      published: true
    })
    setMessage('文章已发布')
  }

  return (
    <section className="editor-layout">
      <section className="panel editor-main">
        <div className="section-toolbar">
          <div>
            <h2>文章编辑器</h2>
            <p className="section-subtitle">用于模拟原版后台的内容编辑流。</p>
          </div>
          <div className="reply-actions">
            <button onClick={saveDraft}>保存草稿</button>
            <button onClick={publish}>发布</button>
          </div>
        </div>
        <div className="editor-form">
          <input value={form.title} onChange={(event) => setForm((currentState) => ({ ...currentState, title: event.target.value }))} placeholder="文章标题" />
          <input value={form.slug} onChange={(event) => setForm((currentState) => ({ ...currentState, slug: event.target.value }))} placeholder="文章别名" />
          <textarea value={form.excerpt} onChange={(event) => setForm((currentState) => ({ ...currentState, excerpt: event.target.value }))} placeholder="摘要" />
          <textarea className="editor-content" value={form.content} onChange={(event) => setForm((currentState) => ({ ...currentState, content: event.target.value }))} placeholder="正文内容" />
          {message ? <p>{message}</p> : null}
        </div>
      </section>
      <section className="panel editor-side compact-panel">
        <h2>发布设置</h2>
        <div className="inline-form">
          <input value={form.category} onChange={(event) => setForm((currentState) => ({ ...currentState, category: event.target.value }))} placeholder="分类" />
          <input value={form.tags} onChange={(event) => setForm((currentState) => ({ ...currentState, tags: event.target.value }))} placeholder="标签" />
          <button>更新属性</button>
        </div>
      </section>
    </section>
  )
}

export function PageEditorSection({ pages }: { pages: Page[] }) {
  const current = pages[0]
  const [form, setForm] = useState({
    title: current?.title ?? '',
    slug: current?.slug ?? '',
    content: current?.content ?? '',
    published: current?.published ?? false
  })
  const [message, setMessage] = useState('')

  const savePage = async () => {
    if (!current) {
      return
    }
    await api.put(`/pages/${current.id}`, {
      title: form.title,
      slug: form.slug,
      content: form.content,
      published: form.published
    })
    setMessage('页面已保存')
  }

  return (
    <section className="editor-layout">
      <section className="panel editor-main">
        <div className="section-toolbar">
          <div>
            <h2>页面编辑器</h2>
            <p className="section-subtitle">用于模拟原版后台的页面编辑流。</p>
          </div>
          <div className="reply-actions">
            <button onClick={savePage}>保存页面</button>
          </div>
        </div>
        <div className="editor-form">
          <input value={form.title} onChange={(event) => setForm((currentState) => ({ ...currentState, title: event.target.value }))} placeholder="页面标题" />
          <input value={form.slug} onChange={(event) => setForm((currentState) => ({ ...currentState, slug: event.target.value }))} placeholder="页面别名" />
          <textarea className="editor-content" value={form.content} onChange={(event) => setForm((currentState) => ({ ...currentState, content: event.target.value }))} placeholder="页面内容" />
          {message ? <p>{message}</p> : null}
        </div>
      </section>
      <section className="panel editor-side compact-panel">
        <h2>页面信息</h2>
        <div className="inline-form">
          <label className="toggle-row">
            <input type="checkbox" checked={form.published} onChange={(event) => setForm((currentState) => ({ ...currentState, published: event.target.checked }))} />
            <span>发布状态</span>
          </label>
          <p className="section-subtitle">这里可继续扩展 SEO、模板和发布状态。</p>
        </div>
      </section>
    </section>
  )
}
