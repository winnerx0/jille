import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowLeft, Award, Clock, Share2, Trash2, TrendingUp, Users } from 'lucide-react'
import { useEffect, useState } from 'react'
import { pollAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/polls/$pollId/view')({
  component: PollResultsPage,
})

function PollResultsPage() {
  const { pollId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  const { data: poll, isLoading, error } = useQuery({
    queryKey: ['poll-view', pollId],
    queryFn: () => pollAPI.getPollView(pollId),
  })

  const deleteMutation = useMutation({
    mutationFn: () => pollAPI.deletePoll(pollId),
    onSuccess: () => {
      toast.success('Poll deleted successfully')
      queryClient.invalidateQueries({ queryKey: ['polls'] })
      navigate({ to: '/home' })
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message || 'Failed to delete poll')
    }
  })

  const currentUserId = (() => {
    const token = localStorage.getItem('access_token')
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.sub
    } catch (e) {
      return null
    }
  })()

  const isCreator = poll?.creator_id === currentUserId

  // SSE for real-time updates (accessible to creator)
  useEffect(() => {
    if (!isCreator) return

    const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:9000'
    const eventSource = new EventSource(`${backendUrl}/api/v1/sse`)

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'POLL_VOTE' && data.payload.poll_id === pollId) {
          queryClient.invalidateQueries({ queryKey: ['poll-view', pollId] })
        }
      } catch (e) {
        console.error('Error parsing SSE event:', e)
      }
    }

    return () => {
      eventSource.close()
    }
  }, [pollId, queryClient, isCreator])

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.origin + `/polls/${pollId}`)
    toast.success('Voting link copied to clipboard!')
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Spinner className="w-12 h-12 mb-4 mx-auto text-primary" />
          <p className="text-muted-foreground animate-pulse">Loading live results...</p>
        </div>
      </div>
    )
  }

  if (error || !poll) {
    const isUnauthorised = (error as any)?.response?.status === 403
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-background p-4">
        <div className="bg-card p-10 text-center max-w-md w-full relative overflow-hidden">

          <h1 className="text-3xl font-black mb-3 tracking-tight">
            {isUnauthorised ? 'Creator Access Only' : 'Poll Not Found'}
          </h1>
          <p className="text-muted-foreground mb-10 leading-relaxed font-medium">
            {isUnauthorised
              ? "This live view is strictly for the poll creator. Please return to the public voting page to participate."
              : "Oops! The poll you're looking for might have expired, been deleted, or never existed."}
          </p>
          <div className="space-y-3">
             <Button
               size="lg"
               className="w-full rounded-2xl h-12 text-lg font-bold"
               onClick={() => navigate({ to: '/home' })}
             >
               Return to Dashboard
             </Button>
          </div>
        </div>
      </div>
    )
  }

  const totalVotes = poll.options.reduce((sum, opt) => sum + opt.votes.length, 0)
  const maxVotes = Math.max(...poll.options.map(o => o.votes.length))
  const isPollActive = new Date(poll.expires_at).getTime() > Date.now()

  return (
    <div className="min-h-screen bg-background text-foreground py-16 px-4 selection:bg-primary selection:text-primary-foreground">
      <div className="max-w-3xl mx-auto">
        {/* Navigation & Actions */}
        <div className="flex items-center justify-between mb-12">
          <button
            onClick={() => navigate({ to: '/home' })}
            className="flex items-center gap-3 text-muted-foreground hover:text-foreground transition-all duration-300 group"
          >
            <div className="w-10 h-10 rounded-full border border-border flex items-center justify-center group-hover:border-primary group-hover:bg-primary/5 transition-all">
              <ArrowLeft className="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
            </div>
            <span className="font-bold tracking-tight">Dashboard</span>
          </button>

          <div className="flex gap-4">
            <Button
              variant="destructive"
              size="lg"
              onClick={() => setShowDeleteConfirm(true)}
              className="rounded-2xl px-6 group"
            >
              <Trash2 className="w-4 h-4 group-hover:animate-bounce" />
            </Button>
          </div>
        </div>

        {/* Delete Confirmation Overlay */}
        {showDeleteConfirm && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-background/80 backdrop-blur-sm animate-in fade-in duration-300">
            <div className="bg-card p-8 rounded-3xl border border-border shadow-2xl max-w-sm w-full text-center scale-in-center overflow-hidden relative">
              <h2 className="text-2xl font-black mb-2">Delete Poll?</h2>
              <p className="text-muted-foreground mb-8 text-sm leading-relaxed">
                This action is permanent and will remove all cast votes. Are you absolutely sure?
              </p>
              <div className="flex flex-col gap-3">
                <Button
                  variant="destructive"
                  size="lg"
                  className="w-full rounded-2xl font-bold h-12"
                  disabled={deleteMutation.isPending}
                  onClick={() => deleteMutation.mutate()}
                >
                  {deleteMutation.isPending ? 'Deleting...' : 'Yes, Delete Permanently'}
                </Button>
                <Button
                  variant="ghost"
                  size="lg"
                  className="w-full rounded-2xl font-bold h-12"
                  onClick={() => setShowDeleteConfirm(false)}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Main Card */}
        <div className="p-8 md:p-16 mb-10 overflow-hidden relative">
          {/* Decorative background element */}
          <div className="absolute top-0 right-0 p-12 opacity-[0.03] rotate-12 pointer-events-none">
            <TrendingUp className="w-64 h-64" />
          </div>

          <div className="relative z-10">
            {/* Badges */}
            <div className="flex flex-wrap items-center gap-4 mb-10">
              <div className="bg-primary text-primary-foreground px-5 py-2 rounded-full text-xs font-black uppercase tracking-[0.2em]">
                Creator Analytics
              </div>
              <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-border text-muted-foreground">
                <Users className="w-4 h-4" />
                <span className="text-xs font-bold">{totalVotes.toLocaleString()} Total Votes</span>
              </div>
              <div className="flex items-center gap-2 px-4 py-2 rounded-full border border-border text-muted-foreground">
                <Clock className="w-4 h-4" />
                <span className="text-xs font-bold italic">Live Updates Enabled</span>
              </div>
            </div>

            <h1 className="text-4xl md:text-6xl font-black mb-16 tracking-tight leading-[1.1]">
              {poll.title}
            </h1>

            <div className="grid gap-6">
              {poll.options.map((option) => {
                const percentage = totalVotes > 0
                  ? Math.round((option.votes.length / totalVotes) * 100)
                  : 0
                const isWinner = maxVotes > 0 && option.votes.length === maxVotes

                return (
                  <div
                    key={option.id}
                    className={cn(
                      "w-full text-left relative p-8 rounded-4xl border transition-all duration-500 overflow-hidden",
                      isWinner && totalVotes > 0
                        ? "bg-primary/3 border-primary/20 shadow-[0_0_40px_rgba(var(--primary-rgb),0.05)]"
                        : "bg-muted/10 border-transparent"
                    )}
                  >
                    {/* Progress Bar Background */}
                    <div
                      className={cn(
                        "absolute inset-0 opacity-[0.08] transition-all duration-1000 ease-out origin-left",
                        isWinner ? "bg-primary" : "bg-muted-foreground"
                      )}
                      style={{ transform: `scaleX(${percentage / 100})` }}
                    />

                    <div className="relative z-10 flex items-center justify-between gap-6">
                      <div className="flex flex-col gap-2 flex-1">
                        <div className="flex items-center gap-3">
                          <span className="text-xl md:text-2xl font-black">
                            {option.name}
                          </span>
                          {isWinner && totalVotes > 0 && <Award className="w-6 h-6 text-amber-500 animate-in zoom-in duration-500" />}
                        </div>
                        <div className="flex items-center gap-2">
                           <span className="text-sm font-bold text-muted-foreground uppercase tracking-widest">
                             {option.votes.length.toLocaleString()} votes
                           </span>
                        </div>
                      </div>

                      <div className="text-right flex flex-col items-end">
                        <span className={cn(
                          "text-3xl md:text-5xl font-black tabular-nums transition-colors",
                          isWinner ? "text-primary" : "text-muted-foreground/50"
                        )}>
                          {percentage}<span className="text-xl md:text-2xl ml-0.5">%</span>
                        </span>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </div>

        {/* Legend */}
        <div className="grid md:grid-cols-3 gap-8 py-10 border-t border-border/50">
           <div className="text-center md:text-left">
             <h4 className="text-xs font-black uppercase tracking-widest text-muted-foreground mb-2">Total Outreach</h4>
             <p className="text-2xl font-black text-foreground">{totalVotes.toLocaleString()}</p>
           </div>
           <div className="text-center">
             <h4 className="text-xs font-black uppercase tracking-widest text-muted-foreground mb-2">Poll Status</h4>
             <div className="flex items-center justify-center gap-2">
               <div className={cn(
                 "w-2 h-2 rounded-full animate-pulse",
                 isPollActive ? "bg-green-500" : "bg-red-500"
               )} />
               <p className="text-sm font-bold uppercase tracking-tighter">
                 {isPollActive ? "Actively Collecting" : "Voting Concluded"}
               </p>
             </div>
           </div>
           <div className="text-center md:text-right">
             <h4 className="text-xs font-black uppercase tracking-widest text-muted-foreground mb-2 flex items-center justify-center md:justify-end gap-2 cursor-pointer" onClick={handleShare}>
               <Share2 className="w-3.5 h-3.5" />
               Share Link
             </h4>
             <p className="text-sm font-bold text-foreground">You (Creator)</p>
           </div>
        </div>
      </div>
    </div>
  )
}
