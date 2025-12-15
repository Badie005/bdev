export default function Home() {
    return (
        <main className="flex min-h-screen flex-col items-center justify-center p-24">
            <div className="text-center">
                <h1 className="text-4xl font-bold mb-4">
                    Welcome to <span className="text-blue-500">B.DEV</span>
                </h1>
                <p className="text-gray-500 mb-8">
                    Your Next.js project is ready to go!
                </p>
                <div className="flex gap-4 justify-center">
                    <a
                        href="https://nextjs.org/docs"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition"
                    >
                        Documentation
                    </a>
                    <a
                        href="https://nextjs.org/learn"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-100 transition"
                    >
                        Learn
                    </a>
                </div>
            </div>
        </main>
    )
}
