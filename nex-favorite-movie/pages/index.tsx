import { serverAddress } from "@/components/serverAddress"
import Link from "next/link"
import { useEffect, useState } from "react"

type UsernameMovieCount = { name: string, movieCount: number }

export default function Home() {


  const [isLoading, setIsLoading] = useState(false)
  const [usernameList, setUsernameList] = useState<UsernameMovieCount[]>([])

  useEffect(() => {
    setIsLoading(true)

    const fetchData = async () => {
      const res = await fetch(serverAddress() + "/user-list")

      if (res.ok) {
        const body: UsernameMovieCount[] = await res.json()

        setUsernameList(body)
      } else {
        console.error(res.status)
        console.error(await res.text())
      }

      setIsLoading(false)
    }

    fetchData().catch(console.error)
  }, [])

  return (<div>
    <div className="flex flex-col text-center py-10">
      <h1 className="text-5xl mb-5">Welcome!</h1>
      <h2 className="text-3xl"><em>Let your friends see your favorite movies!</em></h2>
    </div>
    {isLoading ? <p>Loading</p> : <div className="bg-gray-200 p-3 rounded-md">
      <p className="text-2xl mb-5">Some of users</p>
      <ul>
        {usernameList.map(u => <Link
          href={"/u/" + u.name} key={u.name}>
          <li
            className="flex bg-gray-300 p-3 rounded-md mb-2">
            <span className="grow">{u.name}</span>
            <span>{u.movieCount} Favorite Movie</span>
          </li></Link>)}

      </ul>
    </div>}
  </div>)
}
