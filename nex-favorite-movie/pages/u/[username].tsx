import { useRouter } from "next/router"
import { ChangeEvent, MouseEventHandler, SetStateAction, useCallback, useContext, useEffect, useState } from "react"
import Page404 from "../404"
import debounce from "@/components/debounce"
import { LoginDispatchContext, LoginStateContext } from "@/components/LoginProvider"
import { showAlert } from "@/components/AlertProvider"
import { serverAddress } from "@/components/serverAddress"

export default function UsernamePage() {
  const router = useRouter()

  const [loading, setLoading] = useState(false)
  const [likedMovies, setLikedMovies] = useState<Movie[]>([])
  const [userNotFound, setUserNotFound] = useState(true)
  const [editMode, setEditMode] = useState(false)
  const [isAddNew, setIsAddNew] = useState(false)

  const loginState = useContext(LoginStateContext)
  const loginDispatch = useContext(LoginDispatchContext)

  useEffect(() => {
    if (!router.isReady) return

    const username = router.query.username
    setLoading(true)

    const fetchData = async () => {
      const v = await fetch(serverAddress() + "/u/" + username)
      if (v.ok) {
        const body = await v.json()

        setLikedMovies(body.movies == null ? [] : body.movies)
        setLoading(false)
        setUserNotFound(false)

        if (loginState.isLoggedIn && loginState.username == username) {
          setEditMode(true)
        } else {
          setEditMode(false)
        }

      } else if (v.status == 404) {
        setLoading(false)
        setUserNotFound(true)
      } else {
        console.log(v.status)
        console.log(await v.text)

        showAlert('unknown error', 'warning')
      }
    }

    fetchData().catch(console.error)

  }, [router.isReady, router.query])


  useEffect(() => {
    if (!loginState.isLoggedIn)
      setEditMode(false)
  }, [loginState.isLoggedIn])

  const fetchAddFilm = useCallback(async (m: SuggestedMovie) => {
    try {
      const res = await fetch(serverAddress() + "/movie", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + loginState.token
        },
        body: JSON.stringify(m)
      })

      if (res.status == 401) {
        showAlert('Login Expired', 'warning')
        loginDispatch!({ type: 'logout' })
        router.push("/sign-in")
      } else if (res.status == 409) {
        showAlert('chosen before!', 'error')
      } else if (res.ok) {
        setLikedMovies(ms => [...ms, m])
        showAlert("Movie Added", 'success')
      } else {
        console.log(res.status)
        console.log(await res.text)

        showAlert('unknown error', 'warning')
      }
    } catch (error) {
      console.error(error)
      showAlert('unknown error', 'warning')
    }
  }, [])

  const fetchRemoveFilm = useCallback(async (m: Movie) => {
    try {
      const res = await fetch(serverAddress() + "/movie", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + loginState.token
        },
        body: JSON.stringify(m)
      })
      if (res.status == 401) {
        showAlert('Login Expired', 'warning')
        loginDispatch!({ type: 'logout' })
        router.push("/sign-in")
      } else if (res.ok) {
        setLikedMovies(ms => ms.filter(m_ => m_.name != m.name && m_.year != m.year))
        showAlert("Movie Removed!", 'warning')
      } else {
        console.log(res.status)
        console.log(await res.text)

        showAlert('unknown error', 'warning')
      }
    } catch (error) {
      console.error(error)
      showAlert('unknown error', 'warning')
    }
  }, [])

  return (
    loading ? <h1>Loading</h1> : userNotFound ? <Page404 /> :
      <main className='bg-gray-600 p-4 rounded-md'>
        <p className="text-2xl mb-4 text-white">{router.query.username}</p>
        <ul>
          {likedMovies.map(lm => (<li key={lm.name + lm.year} className='bg-white p-3 rounded-md flex mb-2'>
            <span className='grow p-3 text-lg'>{lm.name}</span>
            <span className='bg-gray-300 p-3 rounded-lg'>{lm.year}</span>
            {editMode && <button
              onClick={() => fetchRemoveFilm(lm)}
              className="btn ml-2">
              <span className="text-xl">&#128465;</span>
            </button>}
          </li>))}
          {isAddNew && <li className='bg-white p-3 rounded-md flex mb-2'>
            <SearchText onChoose={(s) => {
              fetchAddFilm(s)
              setIsAddNew(false)
            }} />
            {/* cancel add film */}
            <button
              className="btn ml-3"
              onClick={() => setIsAddNew(false)}>
              <span className="text-xl">&#128465;</span>
            </button>
          </li>}
          {editMode && !isAddNew && <li
            className="flex">
            <AddMovieButton onClick={() => setIsAddNew(true)} />
          </li>}
        </ul>
      </main>)
}

function SearchText(props: { onChoose: (movie: SuggestedMovie) => void }) {
  const [suggestedFilms, setSuggestedFilms] = useState<SuggestedMovie[]>([])
  const loginState = useContext(LoginStateContext)
  const loginDispatch = useContext(LoginDispatchContext)
  const router = useRouter()

  const fetchSuggestedFilms = useCallback(async (s: string) => {
    if (s == "") {
      setSuggestedFilms([])
      return
    }
    try {
      const res = await fetch(serverAddress() + "/suggest", {
        method: "POST", headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + loginState.token
        },
        body: JSON.stringify({
          text: s
        })
      })
      if (res.status == 401) {
        showAlert('Login Expired')
        loginDispatch!({ type: 'logout' })
        router.push("/sign-in")

      } else if (res.ok) {
        const body: SuggestedMovie[] = await res.json()
        setSuggestedFilms(body)
      } else {
        console.log(res.status)
        console.log(await res.text)
        showAlert('unknown error', 'warning')
      }
    } catch (error) {
      console.error(error)
      showAlert('unknown error', 'warning')
    }
  }, [])

  return (<div className={`dropdown grow dropdown-open `}>
    <input
      tabIndex={0} autoFocus
      className='p-3 text-lg input border-black w-full'
      placeholder="Search Movie" onChange={debounce((e: ChangeEvent<HTMLInputElement>) => {
        fetchSuggestedFilms(e.target.value)
      }, 1000)} />
    <ul
      tabIndex={0}
      className={`dropdown-content menu z-[1] p-2 shadow bg-base-100 rounded-box w-full ${suggestedFilms.length == 0 ? " hidden" : ""}`}>
      {suggestedFilms.map(s => <li key={s.hash} className="w-full"><a className="block" onClick={() => {
        props.onChoose(s)
      }}>{s.name} - {s.year}</a></li>)}
    </ul>
  </div>)
}

const AddMovieButton = (props: { onClick: MouseEventHandler<HTMLButtonElement> }) => (<button onClick={props.onClick} className="
  text-4xl m-1 text-center text-white rounded-md grow
  outline-white outline-dashed outline-2
  focus:ring-gray-400 focus:ring-2 focus:ring-offset-4 focus:ring-offset-gray-600
  active:ring-yellow-100 active:text-yellow-300 active:outline-yellow-300">+</button>)