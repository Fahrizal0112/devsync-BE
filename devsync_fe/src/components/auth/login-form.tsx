"use client"

import { useState } from "react"
import { GitHubLoginButton } from "./github-login-button"
import { DevLoginForm } from "./dev-login-form"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"

export function LoginForm() {
  const [isLoading, setIsLoading] = useState(false)
  const [showDevLogin, setShowDevLogin] = useState(false)

  return (
    <div className="grid gap-6">
      <GitHubLoginButton isLoading={isLoading} setIsLoading={setIsLoading} />
      
      {process.env.NODE_ENV === 'development' && (
        <>
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <Separator className="w-full" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">
                Or continue with
              </span>
            </div>
          </div>
          
          <Button
            variant="outline"
            onClick={() => setShowDevLogin(!showDevLogin)}
            disabled={isLoading}
          >
            {showDevLogin ? 'Hide' : 'Show'} Development Login
          </Button>
          
          {showDevLogin && <DevLoginForm />}
        </>
      )}
    </div>
  )
}