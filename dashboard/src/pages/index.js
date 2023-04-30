import Conn from "@/components/Conn"
import NavBar from "@/components/NavBar"
import NodeInfo from "@/components/NodeInfo"


export default function Home() {
  return (
    <main
      className=""
    >
      <NavBar />
     <section className="w-full">
      <NodeInfo />

      <section className="w-full">
        <Conn />
      </section>
     </section>
    </main>
  )
}
