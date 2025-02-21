import ReactDOM from "react-dom/client"
import App from "./App"
import { BrowserRouter, Routes, Route } from "react-router"

ReactDOM.createRoot(document.getElementById("root")).render(
    <BrowserRouter>
        <Routes>
            <Route path="/" element={<App />} />
        </Routes>
    </BrowserRouter>
    )