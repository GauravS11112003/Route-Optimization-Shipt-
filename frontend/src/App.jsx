import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Info, X, BarChart3, Route as RouteIcon } from 'lucide-react';
import MapView from './components/MapView';
import Sidebar from './components/Sidebar';
import SummaryPanel from './components/SummaryPanel';
import AnalyticsDashboard from './components/AnalyticsDashboard';
import { getSampleData, optimizeWithAnalytics } from './api/optimizer';

function App() {
    const [orders, setOrders] = useState([]);
    const [shoppers, setShoppers] = useState([]);
    const [assignments, setAssignments] = useState([]);
    const [stats, setStats] = useState(null);
    const [analytics, setAnalytics] = useState(null);
    const [routeGeometries, setRouteGeometries] = useState([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [showAbout, setShowAbout] = useState(false);
    const [showAnalytics, setShowAnalytics] = useState(false);
    const [useRealRoutes, setUseRealRoutes] = useState(false);

    const handleLoadSampleData = async () => {
        setLoading(true);
        setError(null);
        setAssignments([]);
        setStats(null);

        try {
            const data = await getSampleData();
            setOrders(data.orders);
            setShoppers(data.shoppers);
        } catch (err) {
            setError('Failed to load sample data. Make sure the backend is running on port 8080.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleOptimize = async () => {
        if (orders.length === 0 || shoppers.length === 0) {
            setError('Please load sample data first.');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const result = await optimizeWithAnalytics({ orders, shoppers }, useRealRoutes);
            setAssignments(result.optimization.assignments);
            setStats(result.optimization);
            setAnalytics(result.analytics);
            setRouteGeometries(result.analytics.routeGeometries || []);
        } catch (err) {
            setError('Failed to optimize routes. Make sure the backend is running.');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="h-screen w-screen flex flex-col overflow-hidden">
            <motion.header
                initial={{ y: -100 }}
                animate={{ y: 0 }}
                transition={{ duration: 0.5, ease: "easeOut" }}
                className="bg-white shadow-md z-40"
            >
                <div className="px-6 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 bg-shipt-green rounded-lg flex items-center justify-center">
                            <svg viewBox="0 0 24 24" fill="white" className="w-6 h-6">
                                <path d="M9 11.75A.75.75 0 0 1 9.75 11h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1-.75-.75Z" />
                                <path fillRule="evenodd" d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm10-8a8 8 0 1 0 0 16 8 8 0 0 0 0-16Z" clipRule="evenodd" />
                            </svg>
                        </div>
                        <div>
                            <h1 className="text-2xl font-bold text-gray-800">
                                Shipt Route Optimizer
                            </h1>
                            <p className="text-sm text-gray-500">
                                Intelligent delivery route planning
                            </p>
                        </div>
                    </div>

                    <div className="flex items-center gap-3">
                        <label className="flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg cursor-pointer hover:bg-gray-200 transition-colors">
                            <input
                                type="checkbox"
                                checked={useRealRoutes}
                                onChange={(e) => setUseRealRoutes(e.target.checked)}
                                className="w-4 h-4 text-shipt-green rounded focus:ring-shipt-green"
                            />
                            <RouteIcon className="w-4 h-4 text-gray-600" />
                            <span className="text-sm font-medium text-gray-700">Real Routes</span>
                        </label>

                        {analytics && (
                            <button
                                onClick={() => setShowAnalytics(!showAnalytics)}
                                className={
                                    showAnalytics
                                        ? 'flex items-center gap-2 px-4 py-2 rounded-lg transition-colors duration-200 bg-shipt-green text-white'
                                        : 'flex items-center gap-2 px-4 py-2 rounded-lg transition-colors duration-200 bg-gray-100 hover:bg-gray-200 text-gray-700'
                                }
                            >
                                <BarChart3 className="w-4 h-4" />
                                <span className="text-sm font-medium">Analytics</span>
                            </button>
                        )}

                        <button
                            onClick={() => setShowAbout(true)}
                            className="flex items-center gap-2 px-4 py-2 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors duration-200"
                        >
                            <Info className="w-4 h-4 text-gray-600" />
                            <span className="text-sm font-medium text-gray-700">About</span>
                        </button>
                    </div>
                </div>
            </motion.header>

            <AnimatePresence>
                {error && (
                    <motion.div
                        initial={{ height: 0, opacity: 0 }}
                        animate={{ height: 'auto', opacity: 1 }}
                        exit={{ height: 0, opacity: 0 }}
                        className="bg-red-50 border-b border-red-200 z-40"
                    >
                        <div className="px-6 py-3 flex items-center justify-between">
                            <div className="flex items-center gap-2 text-red-800">
                                <span className="text-sm">{error}</span>
                            </div>
                            <button
                                onClick={() => setError(null)}
                                className="text-red-600 hover:text-red-800"
                            >
                                <X className="w-4 h-4" />
                            </button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            <div className="flex-1 flex overflow-hidden relative">
                <Sidebar
                    orders={orders}
                    shoppers={shoppers}
                    assignments={assignments}
                    onLoadSampleData={handleLoadSampleData}
                    onOptimize={handleOptimize}
                    loading={loading}
                />

                <div className="flex-1 relative z-0">
                    <MapView
                        orders={orders}
                        shoppers={shoppers}
                        assignments={assignments}
                        routeGeometries={routeGeometries}
                    />

                    <SummaryPanel
                        stats={stats}
                        visible={assignments.length > 0}
                    />
                </div>

                <AnimatePresence>
                    {showAnalytics && (
                        <AnalyticsDashboard
                            analytics={analytics}
                            visible={showAnalytics}
                        />
                    )}
                </AnimatePresence>
            </div>

            <AnimatePresence>
                {showAbout && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={() => setShowAbout(false)}
                        className="fixed inset-0 bg-black bg-opacity-50 z-[100] flex items-center justify-center p-4"
                    >
                        <motion.div
                            initial={{ scale: 0.9, opacity: 0 }}
                            animate={{ scale: 1, opacity: 1 }}
                            exit={{ scale: 0.9, opacity: 0 }}
                            onClick={(e) => e.stopPropagation()}
                            className="bg-white rounded-2xl shadow-2xl max-w-2xl w-full p-8"
                        >
                            <div className="flex items-start justify-between mb-6">
                                <div>
                                    <h2 className="text-2xl font-bold text-gray-800 mb-2">
                                        About Shipt Route Optimizer
                                    </h2>
                                    <div className="h-1 w-20 bg-shipt-green rounded"></div>
                                </div>
                                <button
                                    onClick={() => setShowAbout(false)}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    <X className="w-6 h-6" />
                                </button>
                            </div>

                            <div className="space-y-4 text-gray-600">
                                <p className="leading-relaxed">
                                    This prototype demonstrates route optimization built with <strong>Go</strong> and <strong>React</strong>,
                                    inspired by Shipt's logistics workflows.
                                </p>

                                <p className="leading-relaxed">
                                    The application uses a <strong>nearest-neighbor clustering algorithm</strong> with
                                    real driving routes (via OpenRouteService) to efficiently assign delivery orders to available
                                    shoppers while minimizing total travel distance and time.
                                </p>

                                <p className="leading-relaxed">
                                    <strong>New features:</strong> Real-time analytics dashboard with shopper performance metrics,
                                    capacity utilization, time estimates, cost projections, and environmental impact calculations.
                                </p>

                                <div className="bg-emerald-50 rounded-lg p-4 border border-emerald-100 mt-6">
                                    <h3 className="font-semibold text-gray-800 mb-2">Tech Stack</h3>
                                    <ul className="space-y-1 text-sm">
                                        <li>• <strong>Backend:</strong> Go with Gin framework</li>
                                        <li>• <strong>Frontend:</strong> React with Vite & TailwindCSS</li>
                                        <li>• <strong>Mapping:</strong> Leaflet.js</li>
                                        <li>• <strong>Animations:</strong> Framer Motion</li>
                                    </ul>
                                </div>

                                <div className="bg-blue-50 rounded-lg p-4 border border-blue-100">
                                    <h3 className="font-semibold text-gray-800 mb-2">Features</h3>
                                    <ul className="space-y-1 text-sm">
                                        <li>• Real driving routes vs. straight-line distance</li>
                                        <li>• Comprehensive analytics dashboard</li>
                                        <li>• Shopper performance & efficiency metrics</li>
                                        <li>• Capacity utilization tracking</li>
                                        <li>• Time estimates & delivery windows</li>
                                        <li>• Cost & environmental impact</li>
                                        <li>• Interactive map with live updates</li>
                                    </ul>
                                </div>
                            </div>

                            <div className="mt-8 flex justify-end">
                                <button
                                    onClick={() => setShowAbout(false)}
                                    className="px-6 py-2 bg-shipt-green hover:bg-green-600 text-white font-semibold rounded-lg transition-colors duration-200"
                                >
                                    Got it!
                                </button>
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

export default App;
