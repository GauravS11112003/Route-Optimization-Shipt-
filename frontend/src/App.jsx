import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Info, X, BarChart3, Route as RouteIcon, Zap, Settings, ExternalLink } from 'lucide-react';
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
    const [showSettings, setShowSettings] = useState(false);
    const [useRealRoutes, setUseRealRoutes] = useState(true);
    const [algorithm, setAlgorithm] = useState('astar'); // 'nearest-neighbor' or 'astar'
    const [apiKey, setApiKey] = useState('');
    const [apiKeyInput, setApiKeyInput] = useState('');

    // Load API key from localStorage on mount
    useEffect(() => {
        const savedKey = localStorage.getItem('openroute_api_key');
        if (savedKey) {
            setApiKey(savedKey);
            setApiKeyInput(savedKey);
        }
    }, []);

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
            // Check if real routes are enabled but no API key is set
            if (useRealRoutes && !apiKey) {
                setError('API key required for real routes. Click Settings to add your OpenRouteService API key.');
                setShowSettings(true);
                setLoading(false);
                return;
            }

            const result = await optimizeWithAnalytics({ orders, shoppers }, useRealRoutes, algorithm, apiKey);

            setAssignments(result.optimization.assignments);
            setStats(result.optimization);
            setAnalytics(result.analytics);
            const geometries = result.analytics.routeGeometries || [];
            if (geometries.length > 0) {
                const firstRoutePoints = geometries[0].points?.length || 0;

                // Alert user if routes are falling back to straight lines
                if (firstRoutePoints < 10 && useRealRoutes) {
                    setError('Real routes unavailable - API key may be invalid. Check your OpenRouteService API key in Settings.');
                } else if (firstRoutePoints >= 10) {
                    setError(null); // Clear any previous errors
                }
            }
            setRouteGeometries(geometries);
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
                        {/* Algorithm Selection */}
                        <div className="flex items-center gap-2 bg-white border-2 border-gray-200 rounded-lg p-1">
                            <button
                                onClick={() => setAlgorithm('nearest-neighbor')}
                                className={
                                    algorithm === 'nearest-neighbor'
                                        ? "flex items-center gap-1.5 px-3 py-1.5 bg-shipt-green text-white rounded-md transition-all text-sm font-medium"
                                        : "flex items-center gap-1.5 px-3 py-1.5 text-gray-600 hover:bg-gray-100 rounded-md transition-all text-sm font-medium"
                                }
                            >
                                <RouteIcon className="w-3.5 h-3.5" />
                                <span>Greedy</span>
                            </button>
                            <button
                                onClick={() => setAlgorithm('astar')}
                                className={
                                    algorithm === 'astar'
                                        ? "flex items-center gap-1.5 px-3 py-1.5 bg-shipt-green text-white rounded-md transition-all text-sm font-medium"
                                        : "flex items-center gap-1.5 px-3 py-1.5 text-gray-600 hover:bg-gray-100 rounded-md transition-all text-sm font-medium"
                                }
                            >
                                <Zap className="w-3.5 h-3.5" />
                                <span>A* Search</span>
                            </button>
                        </div>

                        {/* Real Routes Toggle */}
                        <label className={
                            useRealRoutes
                                ? "flex items-center gap-2 px-3 py-2 bg-green-50 border border-green-200 rounded-lg cursor-pointer hover:bg-green-100 transition-colors"
                                : "flex items-center gap-2 px-3 py-2 bg-gray-100 rounded-lg cursor-pointer hover:bg-gray-200 transition-colors"
                        }>
                            <input
                                type="checkbox"
                                checked={useRealRoutes}
                                onChange={(e) => setUseRealRoutes(e.target.checked)}
                                className="w-4 h-4 text-shipt-green rounded focus:ring-shipt-green"
                            />
                            <RouteIcon className={useRealRoutes ? "w-4 h-4 text-green-600" : "w-4 h-4 text-gray-600"} />
                            <span className={useRealRoutes ? "text-sm font-medium text-green-700" : "text-sm font-medium text-gray-700"}>
                                Real Routes
                            </span>
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
                            onClick={() => setShowSettings(true)}
                            className={
                                apiKey
                                    ? "flex items-center gap-2 px-4 py-2 bg-green-50 border border-green-200 hover:bg-green-100 rounded-lg transition-colors duration-200"
                                    : "flex items-center gap-2 px-4 py-2 bg-orange-50 border border-orange-200 hover:bg-orange-100 rounded-lg transition-colors duration-200"
                            }
                        >
                            <Settings className={apiKey ? "w-4 h-4 text-green-600" : "w-4 h-4 text-orange-600"} />
                            <span className={apiKey ? "text-sm font-medium text-green-700" : "text-sm font-medium text-orange-700"}>
                                {apiKey ? 'API Key Set' : 'Settings'}
                            </span>
                        </button>

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
                                    The application features two optimization algorithms: <strong>Greedy Nearest-Neighbor</strong> (fast, efficient)
                                    and <strong>A* Search</strong> (optimal pathfinding with heuristics). Both work with
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
                                        <li>• <strong>Dual algorithms</strong>: Greedy vs A* Search optimization</li>
                                        <li>• <strong>Real driving routes</strong> following actual roads (enable "Real Routes")</li>
                                        <li>• Comprehensive analytics dashboard</li>
                                        <li>• Shopper performance & efficiency metrics</li>
                                        <li>• Capacity utilization tracking</li>
                                        <li>• Time estimates & delivery windows</li>
                                        <li>• Cost & environmental impact</li>
                                        <li>• Interactive map with live updates</li>
                                    </ul>
                                </div>

                                <div className="bg-amber-50 rounded-lg p-4 border border-amber-100">
                                    <h3 className="font-semibold text-gray-800 mb-2">Real Routes Setup</h3>
                                    <p className="text-xs text-gray-700 leading-relaxed">
                                        The "Real Routes" toggle is now enabled by default. For best results, set up an
                                        OpenRouteService API key (free) in the backend <code className="bg-amber-100 px-1 rounded">.env</code> file.
                                        See <code className="bg-amber-100 px-1 rounded">ROUTING_SETUP.md</code> for details.
                                    </p>
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

            {/* Settings Modal */}
            <AnimatePresence>
                {showSettings && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        onClick={() => setShowSettings(false)}
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
                                        Settings
                                    </h2>
                                    <div className="h-1 w-20 bg-shipt-green rounded"></div>
                                </div>
                                <button
                                    onClick={() => setShowSettings(false)}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    <X className="w-6 h-6" />
                                </button>
                            </div>

                            <div className="space-y-6">
                                {/* API Key Section */}
                                <div className="bg-blue-50 rounded-xl p-6 border border-blue-100">
                                    <div className="flex items-start gap-3 mb-4">
                                        <RouteIcon className="w-6 h-6 text-blue-600 mt-1" />
                                        <div className="flex-1">
                                            <h3 className="text-lg font-semibold text-gray-800 mb-2">
                                                OpenRouteService API Key
                                            </h3>
                                            <p className="text-sm text-gray-600 mb-4">
                                                Required for routes to follow actual roads instead of straight lines.
                                                Get your <strong>FREE</strong> API key from OpenRouteService.
                                            </p>

                                            <a
                                                href="https://openrouteservice.org/dev/#/signup"
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium mb-4"
                                            >
                                                <ExternalLink className="w-4 h-4" />
                                                Get Free API Key
                                            </a>

                                            <div className="mt-4">
                                                <label className="block text-sm font-medium text-gray-700 mb-2">
                                                    API Key
                                                </label>
                                                <input
                                                    type="password"
                                                    value={apiKeyInput}
                                                    onChange={(e) => setApiKeyInput(e.target.value)}
                                                    placeholder="Paste your API key here (e.g., eyJvcmci...)"
                                                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-shipt-green focus:border-transparent text-sm font-mono"
                                                />
                                                {apiKeyInput && (
                                                    <p className="text-xs text-gray-500 mt-2">
                                                        Key length: {apiKeyInput.length} characters
                                                    </p>
                                                )}
                                            </div>

                                            <div className="flex gap-3 mt-4">
                                                <button
                                                    onClick={() => {
                                                        setApiKey(apiKeyInput);
                                                        localStorage.setItem('openroute_api_key', apiKeyInput);
                                                        setError(null);
                                                        setShowSettings(false);
                                                    }}
                                                    disabled={!apiKeyInput}
                                                    className="px-6 py-2 bg-shipt-green hover:bg-green-600 text-white font-semibold rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                                >
                                                    Save API Key
                                                </button>
                                                {apiKey && (
                                                    <button
                                                        onClick={() => {
                                                            setApiKey('');
                                                            setApiKeyInput('');
                                                            localStorage.removeItem('openroute_api_key');
                                                        }}
                                                        className="px-6 py-2 bg-red-100 hover:bg-red-200 text-red-700 font-semibold rounded-lg transition-colors"
                                                    >
                                                        Clear Key
                                                    </button>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                {/* Instructions */}
                                <div className="bg-gray-50 rounded-xl p-6 border border-gray-200">
                                    <h4 className="font-semibold text-gray-800 mb-3">How to use:</h4>
                                    <ol className="space-y-2 text-sm text-gray-600">
                                        <li className="flex gap-2">
                                            <span className="font-semibold text-shipt-green">1.</span>
                                            <span>Click "Get Free API Key" above to sign up at OpenRouteService</span>
                                        </li>
                                        <li className="flex gap-2">
                                            <span className="font-semibold text-shipt-green">2.</span>
                                            <span>Copy your API key (starts with "eyJ" and is ~600 characters long)</span>
                                        </li>
                                        <li className="flex gap-2">
                                            <span className="font-semibold text-shipt-green">3.</span>
                                            <span>Paste it in the field above and click "Save API Key"</span>
                                        </li>
                                        <li className="flex gap-2">
                                            <span className="font-semibold text-shipt-green">4.</span>
                                            <span>Enable "Real Routes" toggle and optimize - routes will follow roads!</span>
                                        </li>
                                    </ol>

                                    <div className="mt-4 p-3 bg-amber-50 border border-amber-200 rounded-lg">
                                        <p className="text-xs text-amber-800">
                                            <strong>Note:</strong> Your API key is stored locally in your browser and sent directly to OpenRouteService.
                                            It never passes through our servers.
                                        </p>
                                    </div>
                                </div>
                            </div>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
}

export default App;
