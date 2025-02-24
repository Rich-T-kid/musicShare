import { Link } from "react-router-dom"

const Navigation = () => {
    return (
      <nav>
        <div className="">
          musicShare
        </div>
        <ul>
          <li><Link to="/">Home</Link></li>
          <li><Link to="/recomendation"> Recomendations</Link></li>
        </ul>
      </nav>
    )
}

export default Navigation