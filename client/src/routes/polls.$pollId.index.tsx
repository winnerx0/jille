import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ArrowLeft, Award, CheckCircle2, Clock, Share2, TrendingUp, Users, Vote } from 'lucide-react'
import { useState } from 'react'
import { pollAPI, voteAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/lib/utils'

export const Route = createFileRoute('/polls/$pollId/')({
  component: PollDetailsPage,
})

function PollDetailsPage() {
  const { pollId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [votedOptionId, setVotedOptionId] = useState<string | null>(null)

  const { data: poll, isLoading, error } = useQuery({
    queryKey: ['poll', pollId],
    queryFn: () => pollAPI.getPoll(pollId),
  })

  const voteMutation = useMutation({
    mutationFn: (optionId: string) =>
      voteAPI.vote({ poll_id: pollId, option_id: optionId }),
    onSuccess: (_, optionId) => {
      toast.success('Vote cast successfully!')
      setVotedOptionId(optionId)
      queryClient.invalidateQueries({ queryKey: ['poll', pollId] })
      queryClient.invalidateQueries({ queryKey: ['polls'] })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to cast vote')
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

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href)
    toast.success('Link copied to clipboard!')
  }

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <Spinner className="w-12 h-12 mb-4 mx-auto text-primary" />
          <p className="text-muted-foreground animate-pulse">Entering the voting booth...</p>
        </div>
      </div>
    )
  }

  if (error || !poll) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-background p-4">
        <div className="bg-card p-10 text-center max-w-md w-full relative overflow-hidden">
          <h1 className="text-3xl font-black mb-3 tracking-tight">Poll Not Found</h1>
          <p className="text-muted-foreground mb-10 leading-relaxed">
            Oops! The poll you're looking for might have expired, been deleted, or never existed.
          </p>
          <Button
            size="lg"
            className="w-full rounded-2xl h-14 text-lg font-bold shadow-lg shadow-primary/20"
            onClick={() => navigate({ to: '/home' })}
          >
            Back to Safety
          </Button>
        </div>
      </div>
    )
  }

  const isPollActive = new Date(poll.expires_at).getTime() > Date.now()

  return (
    <div className="min-h-screen bg-background text-foreground py-16 px-4">
      <div className="max-w-2xl mx-auto">
        {/* Simple Header */}
        <div className="flex items-center justify-between mb-16">
          <button
            onClick={() => navigate({ to: '/home' })}
            className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-all group font-bold"
          >
            <ArrowLeft className="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
            Back
          </button>

          <div className="flex gap-4">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleShare}
              className="rounded-xl text-muted-foreground"
            >
              <Share2 className="w-4 h-4 mr-2" />
              Share
            </Button>
          </div>
        </div>

        {/* Voting Booth Layout */}
        <div className="text-center mb-16">
          <div className="w-20 h-20 bg-primary/10 text-primary rounded-3xl flex items-center justify-center mx-auto mb-8 animate-in zoom-in duration-500">
            <Vote className="w-10 h-10" />
          </div>
          <h1 className="text-4xl md:text-5xl font-black mb-4 tracking-tight leading-tight">
            {poll.title}
          </h1>
          <p className="text-muted-foreground font-medium text-lg max-w-md mx-auto">
            {votedOptionId
              ? 'Thank you! Your choice has been securely recorded.'
              : 'Every vote counts. Make your choice below to participate in this poll.'}
          </p>
        </div>

        {!poll.voted && (
          <div className="grid gap-4 mb-16">
            {poll.options.map((option) => {
              const isSelected = votedOptionId === option.id

              return (
                <button
                  key={option.id}
                  disabled={
                    voteMutation.isPending || !!votedOptionId || !isPollActive
                  }
                  onClick={() => voteMutation.mutate(option.id)}
                  className={cn(
                    'w-full text-left group relative p-6 md:p-8 rounded-3xl border-2 transition-all duration-300 flex items-center justify-between',
                    isSelected
                      ? 'bg-primary border-primary text-primary-foreground shadow-xl shadow-primary/20 scale-[1.02]'
                      : votedOptionId || !isPollActive
                        ? 'bg-muted/5 border-border/30 opacity-50 cursor-not-allowed'
                        : 'bg-card border-border hover:border-primary hover:shadow-lg',
                  )}
                >
                  <span className="text-xl font-bold">{option.name}</span>

                  {isSelected ? (
                    <CheckCircle2 className="w-6 h-6 text-primary-foreground animate-in zoom-in" />
                  ) : (
                    !votedOptionId &&
                    isPollActive && (
                      <div className="w-6 h-6 rounded-full border-2 border-border group-hover:border-primary transition-colors" />
                    )
                  )}
                </button>
              )
            })}
          </div>
        )}

        {(votedOptionId || poll.voted) && (
          <div className="p-8 rounded-4xl bg-green-500/5 border border-green-500/10 text-center animate-in slide-in-from-bottom-8 duration-1000">
            <h3 className="text-xl font-bold text-green-600 mb-2">
              Vote Recorded
            </h3>
            <p className="text-muted-foreground text-sm">
              You can close this page now or head back to the dashboard.
            </p>
            <Button
              variant="outline"
              className="mt-6 rounded-xl"
              onClick={() => navigate({ to: '/home' })}
            >
              Go to Dashboard
            </Button>
          </div>
        )}

        {/* Footer Info */}
        <div className="flex flex-col items-center justify-center gap-6 pt-12 border-t border-border/50 text-muted-foreground grayscale opacity-50">
          <div className="flex items-center gap-8">
            <div className="flex items-center gap-2">
              <Users className="w-4 h-4" />
              <span className="text-xs font-bold uppercase tracking-widest">
                Anonymous
              </span>
            </div>
            <div className="flex items-center gap-2">
              <Award className="w-4 h-4" />
              <span className="text-xs font-bold uppercase tracking-widest">
                Verified
              </span>
            </div>
          </div>
          <p className="text-[10px] font-medium tracking-tighter">
            JILLE SECURE VOTING PROTOCOL • 2026
          </p>
        </div>
      </div>
    </div>
  )
}
