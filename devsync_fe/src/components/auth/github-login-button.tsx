"use client"

import { Button } from "@/components/ui/button"
import { Github, Loader2 } from "lucide-react"

interface GitHubLoginButtonProps {
  isLoading: boolean
  setIsLoading: (loading: boolean) => void
}

export function GitHubLoginButton({ isLoading, setIsLoading }: GitHubLoginButtonProps) {
  const handleGitHubLogin = async () => {
    setIsLoading(true)
    try {
      // Redirect to backend GitHub OAuth
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
      window.location.href = `${apiUrl}/auth/github`
    } catch (error) {
      console.error('GitHub login error:', error)
      setIsLoading(false)
    }
  }

  return (
    <Button
      variant="default"
      onClick={handleGitHubLogin}
      disabled={isLoading}
      className="w-full"
      size="lg"
    >
      {isLoading ? (
        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
      ) : (
        <Github className="mr-2 h-4 w-4" />
      )}
      Continue with GitHub
    </Button>
  )
}