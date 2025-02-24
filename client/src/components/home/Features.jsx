import FeaturesCard from "./FeaturesCard"

const Features = () => {

    return (
        <section className="py-24 ">
            <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                <div className="mb-10 lg:mb-16 flex justify-center items-center flex-col gap-x-0 gap-y-6 lg:gap-y-0 lg:flex-row lg:justify-between max-md:max-w-lg max-md:mx-auto">
                    <div className="relative w-full text-center lg:text-left lg:w-2/4">
                        <h2 className="text-4xl font-bold text-gray-900 leading-[3.25rem] lg:mb-6 mx-auto max-w-max lg:max-w-md lg:mx-0">Enjoy the finest features with our products</h2>
                    </div>
                    <div className="relative w-full text-center  lg:text-left lg:w-2/4">
                        <p className="text-lg font-normal text-gray-500 mb-5">We provide all the advantages that can simplify all your financial transactions without any further requirements</p>
                        <a href="#" className="flex flex-row items-center justify-center gap-2 text-base font-semibold text-indigo-600 lg:justify-start hover:text-indigo-700 ">Button CTA <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M7.5 15L11.0858 11.4142C11.7525 10.7475 12.0858 10.4142 12.0858 10C12.0858 9.58579 11.7525 9.25245 11.0858 8.58579L7.5 5" stroke="#4F46E5" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"></path>
                        </svg>
                        </a>
                    </div>
                </div>
                <div className="flex justify-center items-center  gap-x-5 gap-y-8 lg:gap-y-0 flex-wrap md:flex-wrap lg:flex-nowrap lg:flex-row lg:justify-between lg:gap-x-8">
                    <FeaturesCard header="header" desc="We Provide Various Methods For You To Carry Out All Transactions Related To Your Finances" />
                    <FeaturesCard header="Safe Transaction" desc="We have the most up-to-date security to support the security of all our customers in carrying out all transactions." />
                    <FeaturesCard header="Fast Customer Service" desc="Provide Customer Service For Those Of You Who Have Problems 24 Hours A Week" />
                    <FeaturesCard header="Quick Transaction" desc="We provide faster transaction speeds than competitors, so money arrives and is received faster." />
                </div>
            </div>
        </section>
    )
}

export default Features