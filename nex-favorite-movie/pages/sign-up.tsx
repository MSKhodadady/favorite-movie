import SignLayout from "@/components/SignLayout";
import { FieldValue, LiteralUnion, RegisterOptions, UseFormRegister, useForm } from "react-hook-form";

type Inputs = {
  username: string
  password: string
  email: string
}

export default function SignPage() {

  const { register, handleSubmit, formState: { errors }, setError } = useForm<Inputs>()

  const onSubmit = handleSubmit(data => {
    var errored = false
    if (data.email == 'sadeq@email.com') {
      setError('email', { type: 'value', message: 'duplicate email found!' })
      errored = true
    }

    if (data.username == 'sadeq1') {
      setError('username', { type: 'value', message: 'duplicate username found!' })
      errored = true
    }

    if (errored) return

    console.log(data)
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
      })} // TODO:
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