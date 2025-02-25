import Navigation from "../../components/shared/Navigation"
import MusicCard from "../../components/reccomendation/MusicCard"
import Sidebar from "../../components/reccomendation/Sidebar"
import Footer from "../../components/shared/Footer"

const Recomendation = () => {
    return(
        <div>
            <Navigation />
            <div className="grid grid-cols-2">
                <Sidebar />
                <MusicCard />
            </div>
            
            <Footer />
        </div>
    )
}

export default Recomendation