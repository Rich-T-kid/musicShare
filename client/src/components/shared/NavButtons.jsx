import { Link } from "react-router-dom"

const NavButtons = () => {

    return (
        <div className="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
            <div className="flex shrink-0 items-center">
                <img className="h-8 w-auto" src="https://img.icons8.com/?size=100&id=RaE1ps1GXivB&format=png&color=000000" alt="musicShare"></img>
            </div>
            <div className="hidden sm:ml-6 sm:block">
                <div className="flex space-x-4">
                    {/* Current: "bg-gray-900 text-white", Default: "text-gray-300 hover:bg-gray-700 hover:text-white" */}
                    <Link to="/" className="rounded-md bg-gray-900 px-3 py-2 text-sm font-medium text-white" aria-current="page">Home</Link>
                    <Link to="/recommendation" className="rounded-md px-3 py-2 text-sm font-medium text-gray-300 hover:bg-gray-700 hover:text-white">Recomendations</Link>
                </div>
            </div>
        </div>
    )
}

export default NavButtons