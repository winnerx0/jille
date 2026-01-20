import { useForm } from '@tanstack/react-form'
import { z } from 'zod'
import { Link, useNavigate } from '@tanstack/react-router'
import { useMutation } from '@tanstack/react-query'
import { toast } from 'sonner'
import { Input } from './ui/input'
import { Spinner } from './ui/spinner'
import { Button } from './ui/button'
import { Field, FieldError, FieldGroup, FieldLabel } from './ui/field'
import type { AxiosError } from 'axios'
import { authAPI } from '@/lib/api'

const Auth = ({ type }: { type: 'Register' | 'Login' }) => {
  const navigate = useNavigate()

  const RegisterSchema = z.object({
    email: z.string().email({ message: 'Invalid email' }),
    password: z
      .string()
      .min(8, { message: 'Password must have at least 8 characters' })
      .max(16, { message: 'Password must have at most 16 characters' }),
    username: z.string().min(1, { message: 'Username required' }),
  })

  const LoginSchema = z.object({
    email: z.string().email({ message: 'Invalid email' }),
    password: z
      .string()
      .min(8, { message: 'Password must have at least 8 characters' })
      .max(16, { message: 'Password must have at most 16 characters' }),
  })

  const { mutate, isPending } = useMutation({
    mutationFn: async (
      data: z.infer<typeof RegisterSchema> | z.infer<typeof LoginSchema>,
    ) => {
      if (type === 'Register') {
        return authAPI.register(data as z.infer<typeof RegisterSchema>)
      } else {
        return authAPI.login(data as z.infer<typeof LoginSchema>)
      }
    },
    onSuccess(data) {
      localStorage.setItem('access_token', data.data.accessToken)
      localStorage.setItem('refresh_token', data.data.refreshToken)
      navigate({
        to: '/home',
      })
    },
    onError(error: AxiosError) {
      console.log(error)

      const data = error.response?.data as { message: string }
      return toast(data.message || 'An error occurred')
    },
  })

  const form = useForm({
    defaultValues:
      type === 'Register'
        ? {
            email: '',
            username: '',
            password: '',
          }
        : {
            email: '',
            password: '',
          },
    validators: {
      onSubmit: type === 'Register' ? RegisterSchema : LoginSchema,
    },
    onSubmit: (data) => mutate(data.value),
  })
  return (
    <form
      className="flex flex-col space-y-4 w-full max-w-125"
      onSubmit={(e) => {
        e.preventDefault()
        form.handleSubmit(e)
      }}
    >
      <h3 className="font-semibold text-xl">
        {type === 'Register' ? 'Register for Jille' : 'Welcome back'}
      </h3>
      <FieldGroup>
        {type === 'Register' ? (
          <>
            <form.Field
              name="email"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid
                return (
                  <Field className="space-y-1">
                    <FieldLabel>{field.name}</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      onChange={(e) => field.handleChange(e.target.value)}
                      type="email"
                      placeholder="Enter a valid email address"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )
              }}
            />
            <form.Field
              name="username"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid
                return (
                  <Field className="space-y-1">
                    <FieldLabel>{field.name}</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      onChange={(e) => field.handleChange(e.target.value)}
                      type="text"
                      placeholder="Enter a valid username"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )
              }}
            />

            <form.Field
              name="password"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid
                return (
                  <Field className="space-y-1">
                    <FieldLabel>{field.name}</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      onChange={(e) => field.handleChange(e.target.value)}
                      type="password"
                      placeholder="Enter a valid password"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )
              }}
            />
          </>
        ) : (
          <>
            <form.Field
              name="email"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid
                return (
                  <Field className="space-y-1">
                    <FieldLabel>{field.name}</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      onChange={(e) => field.handleChange(e.target.value)}
                      type="email"
                      placeholder="Enter a valid email address"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )
              }}
            />

            <form.Field
              name="password"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid
                return (
                  <Field className="space-y-1">
                    <FieldLabel>{field.name}</FieldLabel>
                    <Input
                      id={field.name}
                      name={field.name}
                      onChange={(e) => field.handleChange(e.target.value)}
                      type="password"
                      placeholder="Enter a valid password"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                )
              }}
            />
          </>
        )}
        <Button type="submit" className="self-center w-36" disabled={isPending}>
          {isPending && <Spinner />} {type}
        </Button>
      </FieldGroup>
      {type === 'Login' ? (
        <p className="self-center">
          Don&apos;t have an account ?{' '}
          <Link className="text-foreground/60" to="/register">
            Sign up
          </Link>
        </p>
      ) : (
        <p className="self-center">
          Have an account ?{' '}
          <Link className="text-foreground/60" to="/login">
            Sign in
          </Link>
        </p>
      )}
    </form>
  )
}

export default Auth
