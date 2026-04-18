import { useState } from 'react'
import type { Attachment, Category, Comment, DashboardStats, Menu, Page, Plugin, Post, Setting, Tag, Theme, User } from './api'
import { AdminShell } from './AdminShell'
import {
  CommentsSection,
  OverviewSection,
  PagesSection,
  PluginsSection,
  PostsSection,
  RegistrationSection,
  SearchSection,
  SettingsSection,
  TaxonomiesSection,
  AttachmentsSection,
  UsersSection
} from './AdminSections'

type Props = {
  posts: Post[]
  themes: Theme[]
  plugins: Plugin[]
  users: User[]
  attachments: Attachment[]
  pages: Page[]
  categories: Category[]
  tags: Tag[]
  menus: Menu[]
  comments: Comment[]
  settings: Setting[]
  stats: DashboardStats | null
  requestError: string
  onLogout: () => void
}

export function Dashboard(props: Props) {
  const [active, setActive] = useState<'overview' | 'posts' | 'pages' | 'taxonomies' | 'comments' | 'users' | 'plugins' | 'settings' | 'register' | 'search' | 'attachments'>('overview')
  const {
    posts,
    themes,
    plugins,
    users,
    attachments,
    pages,
    categories,
    tags,
    menus,
    comments,
    settings,
    stats,
    requestError,
    onLogout
  } = props

  return (
    <AdminShell active={active} onNavigate={setActive} onLogout={onLogout}>
      {requestError ? <section className="panel compact-panel error-panel"><p>{requestError}</p></section> : null}
      {active === 'overview' ? <OverviewSection stats={stats} /> : null}
      {active === 'posts' ? <PostsSection posts={posts} /> : null}
      {active === 'pages' ? <PagesSection pages={pages} /> : null}
      {active === 'taxonomies' ? <TaxonomiesSection categories={categories} tags={tags} menus={menus} /> : null}
      {active === 'comments' ? <CommentsSection comments={comments} /> : null}
      {active === 'users' ? <UsersSection users={users} /> : null}
      {active === 'plugins' ? <PluginsSection plugins={plugins} themes={themes} /> : null}
      {active === 'settings' ? <SettingsSection settings={settings} /> : null}
      {active === 'register' ? <RegistrationSection /> : null}
      {active === 'search' ? <SearchSection /> : null}
      {active === 'attachments' ? <AttachmentsSection attachments={attachments} /> : null}
    </AdminShell>
  )
}
