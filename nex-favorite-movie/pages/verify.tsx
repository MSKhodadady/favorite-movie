import { showAlert } from "@/components/AlertProvider"
import { LoginDispatchContext } from "@/components/LoginProvider"
import { serverAddress } from "@/components/serverAddress"
import jwtDecode from "jwt-decode"
import { useRouter } from "next/router"
import { useContext, useEffect, useState } from "react"

export default function Verify() {
  const [verifyMessage, setVerifyMessage] = useState("")
  const [loading, setLoading] = useState(true)

  const router = useRouter()

  const loginDispatch = useContext(LoginDispatchContext)

  useEffect(() => {
    if (!router.isReady) return
    setLoading(true)

    const token = router.query.token

    if (!token || (Array.isArray(token) && token.length == 0)) {
      setLoading(false)
      setVerifyMessage("not any token in url")
      return
    }

    const fetchData = async () => {
      const res = await fetch(serverAddress() + '/verify?token=' + token)

      if (res.ok) {
        const body: { token: string } = await res.json()
        const payload: Payload = jwtDecode(body.token)

        setLoading(false)

        loginDispatch!({
          type: 'login',
          payload: {
            username: payload.username,
            token: body.token,
            exp: payload.exp
          }
        })

        setVerifyMessage("You are verified. Going to your page in 3 seconds ...")
        setTimeout(() => {
          router.push('/u/' + payload.username)
        }, 3000);
        return
      }


      var message = ""

      if (res.status == 406) {
        console.log(await res.text())
        message = "Your token is not verified."
      } else if (res.status == 401) {
        message = "Your token is expired."
      } else {
        console.log(`code ${res.status} - ${await res.text()}`)
        message = "unknown error ocurred while verification!"
      }

      setVerifyMessage(message + " Going to sign up page is 3 seconds ...")

      setTimeout(() => {
        router.push('/sign-up')
      }, 3000);
      setLoading(false)
      return
    }


    fetchData().catch(console.error)

  }, [router.isReady, router.query])

  return (<div>
    {loading
      ? <p>Loading ...</p>
      : <p>{verifyMessage}</p>
        /* ? <p>Your token is verified. Going to your page in 5 seconds.</p>
    : <p>Your token is <em>not verified</em>! Going to sign up page in 5 seconds.</p>*/}
  </div>)
}