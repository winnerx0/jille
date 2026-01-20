import { createFileRoute, useRouter } from '@tanstack/react-router'
import { ArrowLeft } from 'lucide-react'
import Auth from '@/components/Auth'
import { Button } from '@/components/ui/button'

export const Route = createFileRoute('/register')({
  component: RouteComponent,
})

function RouteComponent() {
  const router = useRouter()

  return (
    <div className='flex items-center justify-center min-h-screen bg-background relative'>
      <Button
        variant="ghost"
        className="absolute top-6 left-6"
        onClick={() => router.history.back()}
      >
        <ArrowLeft className="w-4 h-4 mr-2" />
        Back
      </Button>

      <div className="w-full max-w-md mx-4">
        <div className="bg-card rounded-2xl p-8 shadow-lg border border-border">
          <div className="mb-6 text-center">
            <h2 className="text-2xl font-bold text-foreground">
              Jille
            </h2>
            <p className="text-sm text-muted-foreground mt-1">Create your account</p>
          </div>
          <Auth type="Register"/>
        </div>
      </div>
    </div>
  )
}
