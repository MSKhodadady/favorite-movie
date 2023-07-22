import { FormEventHandler, ReactNode } from "react";

export default function SignLayout(props: {
  onSubmit: FormEventHandler<HTMLFormElement>,
  children?: ReactNode,
  header: string
}) {
  return <main className="bg-gray-300 rounded-md p-5">
    <div className="sm:mx-32 flex md:flex-row flex-col outline outline-black outline-1 rounded-lg p-5">
      <h1 className="text-3xl basis-1/2 pb-2">{props.header}</h1>
      <form className="basis-1/2 flex flex-col" onSubmit={props.onSubmit}>
        {props.children}
      </form>
    </div>
  </main>
}