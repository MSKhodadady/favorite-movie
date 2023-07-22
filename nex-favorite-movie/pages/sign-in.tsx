import { showAlert } from "@/components/AlertProvider"
import { LoginDispatchContext } from "@/components/LoginProvider"
import SignLayout from "@/components/SignLayout"
import { serverAddress } from "@/components/serverAddress"
import jwtDecode from "jwt-decode"
import { useRouter } from "next/router"
import { useContext, useState } from "react"
import { SubmitHandler, useForm } from "react-hook-form"

export default function LoginPage() {
  type Inputs = {
    username: string,
    password: string,
  }

  const { register, handleSubmit, formState: { errors } } = useForm<Inputs>()
  const [userNotFound, setUserNotFound] = useState(false)
  const router = useRouter()
  const loginDispatch = useContext(LoginDispatchContext)

  const submit: SubmitHandler<Inputs> = async ({ username, password }) => {
    const res = await fetch(serverAddress() + "/sign-in", {
      method: "POST", headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: username,
        password: password
      })
    })

    if (res.ok) {
      const body: { token: string } = await res.json()
      const payload: Payload = jwtDecode(body.token)

      loginDispatch!({
        type: 'login',
        payload: {
          ...payload,
          token: body.token
        }
      })

      if (router.isReady) router.push("/u/" + payload.username)

    } else if (res.status == 404) {
      setUserNotFound(true)
      setTimeout(() => {
        setUserNotFound(false)
      }, 2000)
    } else {
      console.log(res.status)
      console.log(await res.text())

      showAlert('unknown error', 'warning')
    }
  }

  return (<SignLayout onSubmit={handleSubmit(submit)} header="Sign In">
    <input
      {...register("username", { required: true, pattern: /^[A-Za-z_]\w*$/ })}
      type="text" placeholder="Username"
      className={`input mb-1 ${errors.username != null ? " input-error" : ""}`}
    />
    <span className="text-red-600 text-sm mb-2">{errors.username?.type == "required" ?
      "Required" :
      errors.username?.type == 'pattern' ?
        "Incorrect Username" : ""}</span>
    <input
      {...register("password", { required: true })}
      type="password" placeholder="Password"
      className={`input mb-1 ${errors.password != null ? " input-error" : ""}`} />
    <span className="text-red-600 text-sm mb-2">{errors.password?.type == "required" ?
      "Required" :
      errors.password?.type == 'pattern' ?
        "Incorrect Username" : ""}</span>
    <button
      type="submit"
      className={`btn ${userNotFound ? " btn-error" : ""}`}>
      {userNotFound ? "USER OR PASSWORD IS INCORRECT" : "SIGN IN"}
    </button>
  </SignLayout>)
}