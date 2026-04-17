import axios from 'axios'

export const api = axios.create({
  baseURL: '/api'
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

export type Theme = {
  id: number
  name: string
  displayName: string
  activated: boolean
}

export type Plugin = {
	 id: number
	 name: string
	 displayName: string
	 enabled: boolean
}

export type User = {
	 id: number
	 username: string
	 role: string
}
