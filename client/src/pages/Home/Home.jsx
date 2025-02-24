import Navigation from "../../components/shared/Navigation"
import Hero from "../../components/home/Hero"
import Features from "../../components/home/Features"
import Testimonials from "../../components/home/Testimonials"
import Footer from "../../components/shared/Footer"

const Home = () => {

    return (
        <div>
            <Navigation />
            <Hero />
            <Features />
            <Testimonials />
            <Footer />
        </div>
    )
}

export default Home