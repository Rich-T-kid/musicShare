import { Routes, Route } from "react-router"
import Home from "./pages/Home"
import Reccomendation from "./pages/Reccomendation"

const App = () => {

  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/reccomendations" element={<Reccomendation />} />
    </Routes>
  )
}

export default App