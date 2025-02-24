import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home/Home"
import Recommendation from "./pages/Recommendation/Recommendation"

const App = () => {

  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/reccomendation" element={<Recommendation />} />
      </Routes>
    </>
  )
}

export default App