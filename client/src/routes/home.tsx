import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Calendar, LogOut, Plus, X } from 'lucide-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useForm } from '@tanstack/react-form'
import { z } from 'zod'
import { toast } from 'sonner'
import { Input } from '@/components/ui/input'
import { Field, FieldError, FieldLabel } from '@/components/ui/field'
import { pollAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { PollCard, PollCardSkeleton } from '@/components/PollCard'
import { Spinner } from '@/components/ui/spinner'

export const Route = createFileRoute('/home')({
  component: RouteComponent,
})

const CreatePollSchema = z.object({
  title: z.string().min(5, 'Title must be at least 5 characters'),
  options: z.array(z.string().min(1, 'Option cannot be empty')).min(2, 'At least 2 options required'),
  expires_at: z.string().refine((val) => !isNaN(Date.parse(val)), 'Invalid date'),
})

function RouteComponent() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [isCreating, setIsCreating] = useState(false)

  const { data: polls, isLoading, error } = useQuery({
    queryKey: ['polls'],
    queryFn: pollAPI.getAllPolls,
    retry: false,
  })

  const createPollMutation = useMutation({
    mutationFn: pollAPI.createPoll,
    onSuccess: () => {
      toast.success('Poll created successfully!')
      queryClient.invalidateQueries({ queryKey: ['polls'] })
      setIsCreating(false)
      form.reset()
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to create poll')
    }
  })

  const form = useForm({
    defaultValues: {
      title: '',
      options: ['', ''],
      expires_at: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    },
    validators: {
      onSubmit: CreatePollSchema,
    },
    onSubmit: ({ value }) => {
      // Convert date to ISO with time
      const date = new Date(value.expires_at)
      date.setHours(23, 59, 59)

      createPollMutation.mutate({
        ...value,
        expires_at: date.toISOString(),
      })
    },
  })

  const handleLogout = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    navigate({ to: '/login' })
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-xl font-bold text-foreground">Jille</h1>
              <p className="text-xs text-muted-foreground">
                Your Polling Dashboard
              </p>
            </div>

            <div className="flex items-center gap-3">
              <Button
                size="lg"
                className="rounded-xl flex gap-2"
                onClick={() => setIsCreating(true)}
              >
                <Plus className="w-4 h-4" />
                <p className="max-md:hidden"> Create Poll</p>
              </Button>
              <Button
                variant="outline"
                size="lg"
                onClick={handleLogout}
                className="rounded-xl flex gap-2"
              >
                <LogOut className="w-4 h-4" />
                <p className="max-md:hidden">Logout</p>
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-12">
        <div className="mb-8 flex justify-between items-end">
          <div>
            <h2 className="text-3xl font-bold mb-2 text-foreground">
              Active Polls
            </h2>
            <p className="text-muted-foreground">
              Browse and participate in ongoing polls
            </p>
          </div>
        </div>

        {/* Error State */}
        {error && (
          <div className="text-center py-12 bg-muted/30 rounded-2xl border border-dashed border-border">
            <div className="inline-flex p-4 rounded-full bg-destructive/10 mb-4">
              <X className="w-8 h-8 text-destructive" />
            </div>
            <h3 className="text-xl font-bold mb-2 text-foreground">
              Failed to load polls
            </h3>
            <p className="text-muted-foreground max-w-md mx-auto">
              {error instanceof Error ? error.message : 'An error occurred'}
            </p>
          </div>
        )}

        {/* Loading State */}
        {isLoading && (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map((i) => (
              <PollCardSkeleton key={i} />
            ))}
          </div>
        )}

        {/* Polls Grid */}
        {!isLoading && !error && polls && polls.length > 0 && (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {polls.map((poll: any) => (
              <PollCard key={poll.id} {...poll} />
            ))}
          </div>
        )}

        {/* Empty State */}
        {!isLoading && !error && polls && polls.length === 0 && (
          <div className="text-center py-20 bg-muted/30 rounded-3xl border border-dashed border-border">
            <div className="inline-flex p-6 rounded-full bg-muted mb-6">
              <Plus className="w-12 h-12 text-muted-foreground" />
            </div>
            <h3 className="text-2xl font-bold mb-2 text-foreground">
              No polls yet
            </h3>
            <p className="text-muted-foreground mb-6">
              Create your first poll to get started sharing opinions.
            </p>
            <Button
              size="lg"
              className="rounded-xl"
              onClick={() => setIsCreating(true)}
            >
              <Plus className="w-4 h-4 mr-2" />
              Create Your First Poll
            </Button>
          </div>
        )}
      </main>

      {/* Create Poll Modal Overlay */}
      {isCreating && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-card w-full max-w-md rounded-2xl shadow-2xl border border-border overflow-hidden animate-in zoom-in-95 duration-200">
            <div className="p-6 border-b border-border flex justify-between items-center">
              <h2 className="text-xl font-bold">Create New Poll</h2>
              <button
                onClick={() => setIsCreating(false)}
                className="p-1 rounded-lg hover:bg-muted transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form
              className="p-6 space-y-6"
              onSubmit={(e) => {
                e.preventDefault()
                form.handleSubmit()
              }}
            >
              <form.Field
                name="title"
                children={(field) => (
                  <Field>
                    <FieldLabel>Poll Title</FieldLabel>
                    <Input
                      placeholder="e.g. What's the best fruit?"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                    {field.state.meta.errors && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )}
              />

              <div className="space-y-3">
                <FieldLabel>Options</FieldLabel>
                <form.Field
                  name="options"
                  children={(field) => (
                    <div className="space-y-2">
                      {field.state.value.map((_, i) => (
                        <div key={i} className="flex gap-2">
                          <Input
                            placeholder={`Option ${i + 1}`}
                            value={field.state.value[i]}
                            onChange={(e) => {
                              const newOptions = [...field.state.value]
                              newOptions[i] = e.target.value
                              field.handleChange(newOptions)
                            }}
                          />
                          {field.state.value.length > 2 && (
                            <Button
                              type="button"
                              variant="ghost"
                              size="icon"
                              onClick={() => {
                                const newOptions = [...field.state.value]
                                newOptions.splice(i, 1)
                                field.handleChange(newOptions)
                              }}
                            >
                              <X className="w-4 h-4" />
                            </Button>
                          )}
                        </div>
                      ))}
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        className="w-full mt-1"
                        onClick={() =>
                          field.handleChange([...field.state.value, ''])
                        }
                      >
                        <Plus className="w-4 h-4 mr-2" /> Add Option
                      </Button>
                      {field.state.meta.errors && (
                        <FieldError errors={field.state.meta.errors} />
                      )}
                    </div>
                  )}
                />
              </div>

              <form.Field
                name="expires_at"
                children={(field) => (
                  <Field>
                    <FieldLabel>Expiration Date</FieldLabel>
                    <div className="relative">
                      <Input
                        type="date"
                        className="pl-10"
                        min={new Date().toISOString().split('T')[0]}
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                      />
                      <Calendar className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
                    </div>
                    {field.state.meta.errors && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )}
              />

              <div className="pt-4 flex gap-3">
                <Button
                  type="button"
                  variant="outline"
                  className="flex-1 rounded-xl"
                  onClick={() => setIsCreating(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  className="flex-1 rounded-xl"
                  disabled={createPollMutation.isPending}
                >
                  {createPollMutation.isPending ? <Spinner /> : 'Create Poll'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
