import FeaturesCard from "./FeaturesCard"

const Features = () => {

    return (
        <section className="py-24 ">
            <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div className="mb-10 lg:mb-16 flex justify-center items-center flex-col gap-x-0 gap-y-6 lg:gap-y-0 lg:flex-row lg:justify-between max-md:max-w-lg max-md:mx-auto">
                    <div className="relative w-full text-center lg:text-left lg:w-2/4">
                        <h2 className="text-4xl font-bold text-gray-900 leading-[3.25rem] lg:mb-6 mx-auto max-w-max lg:max-w-md lg:mx-0">Get the most out of your listening expierence</h2>
                    </div>
                </div>
                <div className="flex justify-center items-center  gap-x-5 gap-y-8 lg:gap-y-0 flex-wrap md:flex-wrap lg:flex-nowrap lg:flex-row lg:justify-between lg:gap-x-8">
                    <FeaturesCard header="Personalize your music for you!" desc="The more you use musicShare, the better recommendations you'll get, just like or add songs to a playlist." />
                    <FeaturesCard header="Favorite Songs directly in the app" desc="Just click the favorite button, and we'll add the song to your favorites." />
                    <FeaturesCard header="Discover new genres" desc="musicShare will curate a selection of genre's based on what you've been listening to as well as what people with similar tastes listen to." />
                </div>
            </div>
        </section>
    )
}

export default Features