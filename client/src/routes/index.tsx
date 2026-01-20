import { Link, createFileRoute } from '@tanstack/react-router'
import { ChartPie, Shield, Vote } from 'lucide-react'
import { Button } from '@/components/ui/button'

export const Route = createFileRoute('/')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <div className="container mx-auto px-4 py-16 md:py-24">
        <div className="text-center max-w-4xl mx-auto mb-16">
          <h1 className="text-5xl md:text-7xl font-bold mb-6 text-foreground leading-tight">
            Make Every Vote Count
          </h1>

          <p className="text-xl md:text-2xl text-muted-foreground mb-10 leading-relaxed">
            A modern, secure, and minimal voting system.
            Create polls, gather opinions, and make decisions together.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <Link to="/register">
              <Button
                size="lg"
                className="px-8 py-6 text-lg rounded-xl"
              >
                Get Started
              </Button>
            </Link>
            <Link to="/login">
              <Button
                size="lg"
                variant="outline"
                className="px-8 py-6 text-lg rounded-xl"
              >
                Sign In
              </Button>
            </Link>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto mt-24">
          <FeatureCard
            icon={<Vote className="w-8 h-8" />}
            title="Simple & Intuitive"
            description="Create polls in seconds with our streamlined interface."
          />
          <FeatureCard
            icon={<Shield className="w-8 h-8" />}
            title="Secure & Private"
            description="Built with JWT authentication and industry-standard security."
          />
          <FeatureCard
            icon={<ChartPie className="w-8 h-8" />}
            title="Real-time Results"
            description="Watch votes come in live with interactive visualizations."
          />
        </div>
      </div>
    </div>
  )
}

function FeatureCard({ icon, title, description }: {
  icon: React.ReactNode
  title: string
  description: string
}) {
  return (
    <div className="bg-card rounded-xl p-6 border border-border hover:shadow-lg transition-all duration-300">
      <div className="inline-flex p-3 rounded-lg bg-primary text-primary-foreground mb-4">
        {icon}
      </div>
      <h3 className="text-lg font-bold mb-2 text-card-foreground">
        {title}
      </h3>
      <p className="text-muted-foreground text-sm leading-relaxed">
        {description}
      </p>
    </div>
  )
}
