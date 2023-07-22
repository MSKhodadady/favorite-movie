import { Dispatch, ReactNode, createContext, useReducer } from "react";

export type LoginState = {
  isLoggedIn: boolean
  username: string
  token: string
  exp: number
}

export interface LoginActions {
  type: 'login' | 'logout'
  payload?: {
    username: string
    token: string
    exp: number
  }
}

const reducer = (s: LoginState, a: LoginActions): LoginState => {
  switch (a.type) {
    case 'login':
      window.localStorage.setItem("token", a.payload!.token)
      window.localStorage.setItem("username", a.payload!.username)
      window.localStorage.setItem("exp", a.payload!.exp.toString())
      return {
        isLoggedIn: true,
        ...a.payload!
      }
    case 'logout':
      window.localStorage.removeItem("token")
      window.localStorage.removeItem("username")
      window.localStorage.removeItem("exp")
      return {
        isLoggedIn: false,
        username: '',
        token: '',
        exp: 0
      }
    default:
      return s
  }
}

const initState: LoginState = {
  isLoggedIn: false,
  username: '',
  token: '',
  exp: 0
}

export const LoginStateContext = createContext<LoginState>(initState)
export const LoginDispatchContext = createContext<Dispatch<LoginActions> | null>(null)

export function LoginStateProvider(props: { children: ReactNode }) {
  const [loginState, loginDispatch] = useReducer(reducer, initState)

  return (<LoginStateContext.Provider value={loginState}>
    <LoginDispatchContext.Provider value={loginDispatch}>
      {props.children}
    </LoginDispatchContext.Provider>
  </LoginStateContext.Provider>)
}