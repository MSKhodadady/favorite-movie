import { showAlert } from "@/components/AlertProvider";
import SignLayout from "@/components/SignLayout";
import { serverAddress } from "@/components/serverAddress";
import { useRouter } from "next/router";
import { FieldValue, LiteralUnion, RegisterOptions, UseFormRegister, useForm } from "react-hook-form";

type Inputs = {
  username: string
  password: string
  email: string
}

export default function SignPage() {

  const { register, handleSubmit, formState: { errors }, setError } = useForm<Inputs>()

  const router = useRouter()

  const onSubmit = handleSubmit(async data => {
    const res = await fetch(serverAddress() + '/sign-up', {
      method: 'POST',
      headers: { "Content-Type": "application/json", },
      body: JSON.stringify({
        username: data.username,
        password: data.password,
        email: data.email,
      })
    });

    if (res.status == 409) {
      const body: { username: boolean, email: boolean } = await res.json()

      if (body.username) {
        setError('username', { type: 'value', message: 'chosen before' })
      }
      if (body.email) {
        setError('email', { type: 'value', message: 'chosen before' })
      }

      return
    } else if (res.ok) {
      showAlert('verification email sent. check your email', 'success')
      router.push('/')
    } else {
      console.log(res)
      console.log(await res.text())

      showAlert('unknown error', 'warning')
    }
  })

  return (<SignLayout onSubmit={onSubmit} header="Sign Up">
    {/* EMAIL */}
    <input
      {...register('email', { required: true, pattern: /^[\w-\.]+@([\w-]+\.)+[\w-]{2,10}$/ })}
      placeholder="Email"
      className={`input ${errors.email == null ? '' : 'input-error'}`}
    />
    <span className="text-red-600 text-sm mb-2">{matchCase({
      'required': 'Required',
      'pattern': 'Incorrect',
      'value': errors.email?.message,
      '': ''
    }, errors.email?.type ?? '')}</span>
    {/* USERNAME */}
    <input
      {...register("username", { required: true, pattern: /^[A-Za-z_]\w*$/ })}
      type="text" placeholder="Username"
      className={`input mb-1 ${errors.username != null ? " input-error" : ""}`}
    />
    <span className="text-red-600 text-sm mb-2">{matchCase({
      'required': 'Required',
      'pattern': 'Incorrect',
      'value': errors.username?.message,
      '': ''
    }, errors.username?.type ?? '')}</span>
    {/* PASSWORD */}
    <input
      {...register("password", {
        required: true,
        pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&_])[A-Za-z\d@$!%*?&_]{8,}$/
      })}
      type="password" placeholder="Password"
      className={`input mb-1 ${errors.password != null ? " input-error" : ""}`} />
    <span className="text-red-600 text-sm mb-2">{matchCase({
      'required': 'Required',
      'pattern': '8 char, at least with a lower letter, a upper letter, a number, and one of @$!%*?&_',
      '': ''
    }, errors.password?.type ?? '')}</span>
    {/* SUBMIT */}
    <button
      type="submit"
      className={`btn`}>
      SIGN UP
    </button>
  </SignLayout>)
}

const matchCase = (m: Record<string, any>, s: string) => m[s]