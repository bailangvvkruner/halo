import axios from 'axios'
import { getToken } from './auth'

export const api = axios.create({
  baseURL: '/api'
})

api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export type Health = {
  status: string
}

export type DashboardStats = {
  posts: number
  pages: number
  categories: number
  tags: number
  users: number
  themes: number
  plugins: number
}

export type Post = {
  id: number
  title: string
  slug: string
  excerpt: string
  template: string
  category: string
  tags: string
  published: boolean
}

export type Page = {
  id: number
  title: string
  slug: string
  published: boolean
}

export type Category = {
  id: number
  name: string
  slug: string
  displayName: string
}

export type Tag = {
  id: number
  name: string
  slug: string
}

export type Menu = {
  id: number
  name: string
  items: string
}

export type Comment = {
  id: number
  postId: number
  author: string
  email: string
  content: string
}

export type Reply = {
  id: number
  commentId: number
  author: string
  email: string
  content: string
  status: string
}

export type RegistrationPayload = {
  username: string
  email: string
  password: string
}

export type SearchResult = {
  posts: Post[]
  pages: Page[]
}

export type Theme = {
  id: number
  name: string
  displayName: string
  activated: boolean
}

export type ThemeScanItem = {
  name: string
  displayName: string
}

export type PluginScanItem = {
  name: string
  displayName: string
  path: string
}

export type Plugin = {
	 id: number
	 name: string
	 displayName: string
	 path?: string
	 enabled: boolean
}

export type Setting = {
  id: number
  key: string
  value: string
}

export type User = {
	 id: number
	 username: string
	 role: string
}

export type LoginResponse = {
  token: string
  user: User
}

export type Attachment = {
  id: number
  filename: string
  path: string
  url: string
  size: number
}
