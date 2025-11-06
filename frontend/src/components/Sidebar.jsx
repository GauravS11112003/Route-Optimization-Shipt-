import { motion, AnimatePresence } from 'framer-motion';
import { Package, User, Clock, ShoppingCart, ChevronRight } from 'lucide-react';

export default function Sidebar({
    orders,
    shoppers,
    assignments,
    onLoadSampleData,
    onOptimize,
    loading
}) {
    return (
        <motion.div
            initial={{ x: -300, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ duration: 0.5, ease: "easeOut" }}
            className="w-80 bg-white h-full shadow-xl overflow-y-auto flex flex-col"
        >
            {/* Header */}
            <div className="p-6 border-b border-gray-200">
                <h2 className="text-xl font-bold text-gray-800 mb-4">Control Panel</h2>

                {/* Action Buttons */}
                <div className="space-y-3">
                    <button
                        onClick={onLoadSampleData}
                        disabled={loading}
                        className="w-full bg-shipt-green hover:bg-green-600 disabled:bg-gray-300 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-200 shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
                    >
                        {loading ? 'Loading...' : 'Load Sample Data'}
                    </button>

                    <button
                        onClick={onOptimize}
                        disabled={loading || orders.length === 0}
                        className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-200 shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
                    >
                        {loading ? 'Optimizing...' : 'Optimize Routes'}
                    </button>
                </div>
            </div>

            {/* Shoppers Section */}
            {shoppers.length > 0 && (
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.2 }}
                    className="p-6 border-b border-gray-200"
                >
                    <div className="flex items-center gap-2 mb-3">
                        <User className="w-5 h-5 text-shipt-green" />
                        <h3 className="font-semibold text-gray-800">
                            Shoppers ({shoppers.length})
                        </h3>
                    </div>

                    <div className="space-y-2">
                        {shoppers.map((shopper) => {
                            const assignment = assignments.find(a => a.shopperId === shopper.id);
                            const assignedCount = assignment ? assignment.route.length : 0;

                            return (
                                <motion.div
                                    key={shopper.id}
                                    initial={{ opacity: 0, x: -20 }}
                                    animate={{ opacity: 1, x: 0 }}
                                    className="bg-emerald-50 rounded-lg p-3 border border-emerald-100"
                                >
                                    <div className="flex justify-between items-start">
                                        <div>
                                            <div className="font-semibold text-gray-800">
                                                {shopper.id}
                                            </div>
                                            <div className="text-xs text-gray-500 mt-1">
                                                Capacity: {shopper.capacity} orders
                                            </div>
                                        </div>
                                        {assignment && (
                                            <div className="bg-shipt-green text-white text-xs font-semibold px-2 py-1 rounded">
                                                {assignedCount} assigned
                                            </div>
                                        )}
                                    </div>

                                    {assignment && assignment.totalDistance && (
                                        <div className="mt-2 text-xs text-gray-600 flex items-center gap-1">
                                            <ChevronRight className="w-3 h-3" />
                                            Route: {assignment.totalDistance.toFixed(1)} km
                                        </div>
                                    )}
                                </motion.div>
                            );
                        })}
                    </div>
                </motion.div>
            )}

            {/* Orders Section */}
            {orders.length > 0 && (
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.3 }}
                    className="p-6 flex-1"
                >
                    <div className="flex items-center gap-2 mb-3">
                        <Package className="w-5 h-5 text-orange-500" />
                        <h3 className="font-semibold text-gray-800">
                            Orders ({orders.length})
                        </h3>
                    </div>

                    <div className="space-y-2 max-h-[400px] overflow-y-auto pr-2">
                        {orders.map((order) => (
                            <motion.div
                                key={order.id}
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                className="bg-orange-50 rounded-lg p-3 border border-orange-100"
                            >
                                <div className="font-semibold text-gray-800">
                                    {order.id}
                                </div>
                                <div className="flex items-center gap-1 text-xs text-gray-600 mt-1">
                                    <ShoppingCart className="w-3 h-3" />
                                    {order.itemCount} items
                                </div>
                                <div className="flex items-center gap-1 text-xs text-gray-600 mt-1">
                                    <Clock className="w-3 h-3" />
                                    {order.deliveryWindow}
                                </div>
                            </motion.div>
                        ))}
                    </div>
                </motion.div>
            )}

            {/* Empty State */}
            {orders.length === 0 && shoppers.length === 0 && (
                <div className="flex-1 flex items-center justify-center p-6">
                    <div className="text-center text-gray-400">
                        <Package className="w-12 h-12 mx-auto mb-3 opacity-50" />
                        <p className="text-sm">Load sample data to get started</p>
                    </div>
                </div>
            )}
        </motion.div>
    );
}

