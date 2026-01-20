import { AlertTriangle, ArrowRight, Clock, Trash2, TrendingUp, Users, Vote } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type {PollOption} from '@/lib/api';
import { pollAPI } from '@/lib/api'
import { Button } from '@/components/ui/button'

interface PollCardProps {
  id: string
  title: string
  options: Array<PollOption>
  createdAt?: string
  expires_at?: string
  creator_id?: string
}

export function PollCard({ id: pollId, title, options, createdAt, expires_at, creator_id }: PollCardProps) {
  const currentUserId = (() => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('access_token') : null
    if (!token) return null
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.sub
    } catch (e) {
      return null
    }
  })()

  const isCreator = creator_id === currentUserId
  const queryClient = useQueryClient()
  const [showConfirm, setShowConfirm] = useState(false)

  const deleteMutation = useMutation({
    mutationFn: () => pollAPI.deletePoll(pollId),
    onSuccess: () => {
      toast.success('Poll deleted')
      queryClient.invalidateQueries({ queryKey: ['polls'] })
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.message || 'Failed to delete poll')
    }
  })

  const totalVotes = options.reduce((sum: number, opt: any) => sum + opt.votes.length, 0)
  const isExpired = expires_at ? new Date(expires_at).getTime() < Date.now() : false

  const optionsWithPercentage = options.map(opt => ({
    ...opt,
    percentage: totalVotes > 0 ? Math.round((opt.votes.length / totalVotes) * 100) : 0
  }))

  return (
    <div className="bg-card rounded-2xl p-6 shadow border border-border hover:shadow-sm hover:border-primary/20 transition-all duration-300 flex flex-col h-full group relative">
      {/* Delete Confirmation Overlay */}
      {showConfirm && (
        <div className="absolute inset-0 z-20 bg-background/95 backdrop-blur-sm rounded-2xl flex flex-col items-center justify-center p-6 text-center animate-in fade-in duration-200">
           <AlertTriangle className="w-10 h-10 text-destructive mb-3" />
           <p className="text-sm font-bold mb-4">Delete this poll permanently?</p>
           <div className="flex flex-col w-full gap-2">
              <Button
                variant="destructive"
                size="sm"
                className="w-full rounded-xl h-12"
                onClick={() => deleteMutation.mutate()}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? 'Deleting...' : 'Confirm Delete'}
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="w-full rounded-xl h-12"
                onClick={() => setShowConfirm(false)}
              >
                Cancel
              </Button>
           </div>
        </div>
      )}

      {/* Header */}
      <div className="flex items-start gap-4 mb-6">
        <div className="p-3 rounded-xl bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300">
          <Vote className="w-5 h-5" />
        </div>
        <div className="flex-1">
          <div className="flex items-start justify-between">
            <h3 className="text-xl font-bold text-card-foreground mb-1 leading-tight line-clamp-2">
              {title}
            </h3>
            {isCreator && (
              <button
                onClick={() => setShowConfirm(true)}
                className="p-1.5 text-muted-foreground hover:text-destructive transition-colors opacity-0 group-hover:opacity-100"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            )}
          </div>
          <div className="flex items-center gap-2">
             {isExpired && (
               <span className="bg-destructive/10 text-destructive text-[10px] font-black uppercase tracking-tighter px-1.5 py-0.5 rounded">
                 Ended
               </span>
             )}
             <div className="flex items-center gap-2 text-muted-foreground">
                <Users className="w-3.5 h-3.5" />
                <span className="text-xs font-medium">{totalVotes} votes cast</span>
             </div>
          </div>
        </div>
      </div>

      <div className="space-y-4 mb-8 flex-1">
        {isCreator ? (
          optionsWithPercentage.slice(0, 3).map((option) => (
            <div key={option.id} className="space-y-1.5">
              <div className="flex items-center justify-between px-1">
                <span className="text-sm font-semibold text-card-foreground">
                  {option.name}
                </span>
                <span className="text-xs font-black text-primary/70">
                  {option.percentage}%
                </span>
              </div>
              <div className="relative h-1.5 bg-muted rounded-full overflow-hidden">
                <div
                  className="h-full bg-primary/40 rounded-full transition-all duration-1000"
                  style={{ width: `${option.percentage}%` }}
                />
              </div>
            </div>
          ))
        ) : (
          <div className="space-y-2">
            {options.slice(0, 4).map((option) => (
              <div key={option.id} className="flex items-center gap-3 p-3 rounded-xl bg-muted/30 border border-transparent">
                <div className="w-1.5 h-1.5 rounded-full bg-primary/50" />
                <span className="text-sm font-medium text-foreground/80 line-clamp-1">
                  {option.name}
                </span>
              </div>
            ))}
          </div>
        )}
        {options.length > (isCreator ? 3 : 4) && (
          <p className="text-xs text-muted-foreground text-center italic">
            + {options.length - (isCreator ? 3 : 4)} more options
          </p>
        )}
      </div>

      {/* Action Buttons */}
      <div className="flex flex-col gap-3 mt-auto">
        <Link
          to="/polls/$pollId"
          params={{ pollId }}
          className="w-full"
        >
          <button className="w-full py-3.5 px-4 bg-muted hover:bg-primary text-muted-foreground hover:text-primary-foreground font-bold rounded-xl transition-all duration-300 flex items-center justify-center gap-2 group/btn">
             <span>{isCreator ? 'Cast Your Vote' : 'Vote Now'}</span>
             <ArrowRight className="w-4 h-4 group-hover/btn:translate-x-1 transition-transform" />
          </button>
        </Link>

        {isCreator && (
          <Link
            to="/polls/$pollId/view"
            params={{ pollId }}
            className="w-full"
          >
            <button className="w-full py-3.5 px-4 bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground font-bold rounded-xl transition-all duration-300 flex items-center justify-center gap-2 group/btn">
               <TrendingUp className="w-4 h-4" />
               <span>Live Results</span>
               <ArrowRight className="w-4 h-4 group-hover/btn:translate-x-1 transition-transform" />
            </button>
          </Link>
        )}
      </div>

      {/* Footer Stats */}
      <div className="flex items-center justify-between pt-5 mt-5 border-t border-border/50 text-muted-foreground/70">
        {createdAt && (
          <div className="flex items-center gap-1.5">
            <Clock className="w-3.5 h-3.5" />
            <span className="text-[10px] font-bold uppercase tracking-wider">{createdAt}</span>
          </div>
        )}
        <div className="flex items-center gap-1.5">
          <TrendingUp className="w-3.5 h-3.5" />
          <span className="text-[10px] font-bold uppercase tracking-wider">Live Results</span>
        </div>
      </div>
    </div>
  )
}

export function PollCardSkeleton() {
  return (
    <div className="bg-card rounded-2xl p-6 shadow border border-border animate-pulse h-[400px] flex flex-col">
      <div className="flex items-start gap-4 mb-6">
        <div className="w-11 h-11 bg-muted rounded-xl" />
        <div className="flex-1">
          <div className="h-6 bg-muted rounded w-3/4 mb-2" />
          <div className="h-3 bg-muted rounded w-1/4" />
        </div>
      </div>
      <div className="space-y-4 flex-1">
        {[1, 2, 3].map((i) => (
          <div key={i} className="space-y-2">
            <div className="h-4 bg-muted rounded w-1/2" />
            <div className="h-1.5 bg-muted rounded-full" />
          </div>
        ))}
      </div>
      <div className="h-12 bg-muted rounded-xl w-full mt-auto" />
    </div>
  )
}
