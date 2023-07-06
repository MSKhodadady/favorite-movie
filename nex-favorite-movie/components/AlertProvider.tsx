import { Dispatch, ReactNode, createContext, useReducer } from "react";

//: types
export type AlertType = "info" | "" | "success" | "warning" | "error"
export interface AlertState {
  text: string
  show: boolean
  type: AlertType
}
export interface AlertAction {
  type: 'show' | 'hide',
  payload?: AlertState
}

const reducer = (s: AlertState, a: AlertAction): AlertState => {
  switch (a.type) {
    case "hide":
      return {
        text: '',
        show: false,
        type: ''
      }
    case 'show':
      return a.payload!
    default:
      return {
        text: '',
        show: false,
        type: ''
      }
  }
}
const init: AlertState = {
  show: false,
  text: '',
  type: ''
}
export const AlertContext =
  createContext<{ alertState: AlertState, alertDispatch: Dispatch<AlertAction> }>({
    alertState: init,
    alertDispatch: () => { }
  })
let alertDispatchHolder: Dispatch<AlertAction> | null = null;

export function AlertProvider(props: { children: ReactNode }) {
  const [alertState, alertDispatch] = useReducer(reducer, init)

  alertDispatchHolder = alertDispatch

  return (<AlertContext.Provider value={{
    alertState, alertDispatch
  }}>
    {props.children}
  </AlertContext.Provider>)
}

export function showAlert(text: string, alertType: AlertType = 'info', time = 3000) {
  alertDispatchHolder!({
    type: 'show', payload: {
      show: true,
      text,
      type: alertType
    }
  })

  setTimeout(() => {
    alertDispatchHolder!({
      type: 'hide'
    })
  }, time);
}


export function alertTypeClassName(t: AlertType) {
  switch (t) {
    case 'error':
      return 'alert-error'
    case 'info':
      return 'alert-info'
    case 'success':
      return 'alert-success'
    case 'warning':
      return 'alert-warning'
    case '':
      return ''
    default:
      ''
  }
}