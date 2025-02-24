import { Link } from "react-router-dom"
import MobileMenuButton from "./MobileMenuButton"
import NavButtons from "./NavButtons"
import ProfileButton  from "./ProfileButton"
const Navigation = () => {
  return (
    <nav className="bg-gray-800">
      <div className="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
        <div className="relative flex h-16 items-center justify-between">
          <MobileMenuButton />
          <NavButtons />
          <ProfileButton />
        </div>
      </div>
    </nav>
  )
}

export default Navigation