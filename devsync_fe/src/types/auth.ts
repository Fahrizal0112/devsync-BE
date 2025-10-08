export interface User {
    id: number
    github_id: number
    username: string
    email: string
    name: string
    avatar_url: string
    created_at: string
    updated_at: string
}

export interface LoginResponse {
    token: string
    user: User
}

export interface ErrorResponse {
    error: string
}

export interface DevLoginRequest {
    username: string
    email: string
}