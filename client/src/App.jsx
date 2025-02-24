import { Routes, Route } from "react-router-dom"
import Home from "./pages/Home/Home"
import Recommendation from "./pages/Recommendation/Recommendation"
import "./index.css"

const App = () => {

  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/recomendation" element={<Recommendation />} />
      </Routes>
    </>
  )
}

export default App