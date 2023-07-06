import { AlertContext, AlertProvider, AlertType, alertTypeClassName } from '@/components/AlertProvider'
import { LoginActions, LoginDispatchContext, LoginStateContext, LoginStateProvider } from '@/components/LoginProvider'
import '@/styles/globals.css'
import type { AppProps } from 'next/app'
import Link from 'next/link'
import { NextRouter } from 'next/router'
import { Dispatch, useContext, useEffect } from 'react'

export default function App({ Component, pageProps }: AppProps) {
  return (<AlertProvider>
    <LoginStateProvider>
      <div className='p-3'>
        <AppNavBar />
        <div><Component {...pageProps} /></div>
      </div>
    </LoginStateProvider>
  </AlertProvider>)
}


function LoginButton(props: { text: string, link: string }) {
  return <Link href={props.link} className='
  block
  bg-white p-2 rounded-md me-2
  shadow-inner
  active:shadow-2xl'>{props.text}</Link>
}

function AppNavBar() {

  const loginState = useContext(LoginStateContext)
  const loginDispatch = useContext(LoginDispatchContext)

  const alertProvider = useContext(AlertContext)



  const logOut = () => {
    loginDispatch!({ type: 'logout' })
  }

  useEffect(() => {
    const token = localStorage.getItem("token")
    const exp = Number(localStorage.getItem("exp"))

    if (token) {
      const nowEpochSeconds = Math.floor(new Date().valueOf() / 1000);
      if (exp >= nowEpochSeconds)
        loginDispatch!({
          type: 'login',
          payload: {
            username: localStorage.getItem("username")!,
            token: localStorage.getItem("token")!,
            exp: Number(localStorage.getItem("exp")!)
          }
        })
      else
        loginDispatch!({ type: 'logout' })
    }
  }, [])

  return (<nav className='mb-3'>
    <ul className='flex bg-gray-200 p-2 rounded-md justify-center'>
      <li className='grow text-4xl'>
        <Link href={'/'}>
          &#128253; My Favorite Movie
        </Link>
      </li>
      {loginState.isLoggedIn ?
        <div className="dropdown dropdown-hover dropdown-left p-2">
          <span className='text-lg'>
            Hello <span className='inline font-bold'>{loginState.username}</span>
          </span>
          <ul tabIndex={0} className="dropdown-content z-[1] menu p-2 shadow bg-base-100 rounded-box w-52">
            <li><Link href={`/u/${loginState.username}`}>User Page</Link></li>
            <li><a onClick={logOut}>Logout</a></li>
          </ul>
        </div>
        :
        <>
          <li><LoginButton text='Log in' link='/login' /></li>
          <li>
            <LoginButton text='Sign in' link='/sign-in'></LoginButton>
          </li>
        </>}
    </ul>
    {alertProvider.alertState.show && <div
      className={"alert z-10 absolute left-2 bottom-2 w-fit " + alertTypeClassName(alertProvider.alertState.type)}>
      <p>{alertProvider.alertState.text}</p>
    </div>}
  </nav>)
}