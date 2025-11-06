import { useEffect, useRef } from 'react';
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Fix for default marker icons in webpack
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
    iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
    shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
});

// Custom icons
const shopperIcon = new L.Icon({
    iconUrl: 'data:image/svg+xml;base64,' + btoa(`
    <svg width="25" height="41" viewBox="0 0 25 41" xmlns="http://www.w3.org/2000/svg">
      <path d="M12.5 0C5.6 0 0 5.6 0 12.5c0 9.4 12.5 28.5 12.5 28.5S25 21.9 25 12.5C25 5.6 19.4 0 12.5 0z" fill="#00C389"/>
      <circle cx="12.5" cy="12.5" r="6" fill="white"/>
    </svg>
  `),
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
});

const orderIcon = new L.Icon({
    iconUrl: 'data:image/svg+xml;base64,' + btoa(`
    <svg width="25" height="41" viewBox="0 0 25 41" xmlns="http://www.w3.org/2000/svg">
      <path d="M12.5 0C5.6 0 0 5.6 0 12.5c0 9.4 12.5 28.5 12.5 28.5S25 21.9 25 12.5C25 5.6 19.4 0 12.5 0z" fill="#f97316"/>
      <circle cx="12.5" cy="12.5" r="6" fill="white"/>
    </svg>
  `),
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
});

// Route colors for different shoppers
const routeColors = ['#00C389', '#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b'];

// Component to fit map bounds
function MapBounds({ orders, shoppers, assignments }) {
    const map = useMap();

    useEffect(() => {
        const allPoints = [];

        shoppers.forEach(s => allPoints.push([s.lat, s.lng]));
        orders.forEach(o => allPoints.push([o.lat, o.lng]));

        if (allPoints.length > 0) {
            const bounds = L.latLngBounds(allPoints);
            map.fitBounds(bounds, { padding: [50, 50] });
        }
    }, [map, orders, shoppers, assignments]);

    return null;
}

export default function MapView({ orders, shoppers, assignments, routeGeometries = [] }) {
    const defaultCenter = [33.5186, -86.8104]; // Birmingham, AL

    // Build route lines from assignments or use real geometries
    const routeLines = routeGeometries.length > 0
        ? routeGeometries.map((geometry, idx) => ({
            coords: geometry.points,
            color: routeColors[idx % routeColors.length],
            shopperId: geometry.shopperId,
        }))
        : assignments.map((assignment, idx) => {
            const shopper = shoppers.find(s => s.id === assignment.shopperId);
            if (!shopper) return null;

            const routeCoords = [[shopper.lat, shopper.lng]];

            assignment.route.forEach(orderId => {
                const order = orders.find(o => o.id === orderId);
                if (order) {
                    routeCoords.push([order.lat, order.lng]);
                }
            });

            return {
                coords: routeCoords,
                color: routeColors[idx % routeColors.length],
                shopperId: assignment.shopperId,
            };
        }).filter(Boolean);

    return (
        <div className="h-full w-full relative">
            <MapContainer
                center={defaultCenter}
                zoom={11}
                className="h-full w-full"
                style={{ background: '#e5e7eb' }}
            >
                <TileLayer
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                <MapBounds orders={orders} shoppers={shoppers} assignments={assignments} />

                {/* Render shoppers */}
                {shoppers.map(shopper => (
                    <Marker
                        key={shopper.id}
                        position={[shopper.lat, shopper.lng]}
                        icon={shopperIcon}
                    >
                        <Popup>
                            <div className="font-semibold text-shipt-green">Shopper {shopper.id}</div>
                            <div className="text-sm text-gray-600">Capacity: {shopper.capacity} orders</div>
                            <div className="text-xs text-gray-500 mt-1">
                                {shopper.lat.toFixed(4)}, {shopper.lng.toFixed(4)}
                            </div>
                        </Popup>
                    </Marker>
                ))}

                {/* Render orders */}
                {orders.map(order => (
                    <Marker
                        key={order.id}
                        position={[order.lat, order.lng]}
                        icon={orderIcon}
                    >
                        <Popup>
                            <div className="font-semibold text-orange-500">Order {order.id}</div>
                            <div className="text-sm text-gray-600">{order.itemCount} items</div>
                            <div className="text-sm text-gray-600">Window: {order.deliveryWindow}</div>
                            <div className="text-xs text-gray-500 mt-1">
                                {order.lat.toFixed(4)}, {order.lng.toFixed(4)}
                            </div>
                        </Popup>
                    </Marker>
                ))}

                {/* Render route lines */}
                {routeLines.map((route, idx) => (
                    <Polyline
                        key={`route-${idx}-${route.shopperId}`}
                        positions={route.coords}
                        pathOptions={{
                            color: route.color,
                            weight: 4,
                            opacity: 0.8,
                            dashArray: routeGeometries.length > 0 ? null : '10, 10',
                        }}
                    />
                ))}
            </MapContainer>

            {/* Legend */}
            {(shoppers.length > 0 || orders.length > 0) && (
                <div className="absolute bottom-6 right-6 bg-white rounded-lg shadow-lg p-4 z-[1000]">
                    <div className="text-sm font-semibold mb-2 text-gray-700">Legend</div>
                    <div className="flex items-center gap-2 mb-2">
                        <div className="w-4 h-4 rounded-full bg-shipt-green"></div>
                        <span className="text-xs text-gray-600">Shoppers</span>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="w-4 h-4 rounded-full bg-orange-500"></div>
                        <span className="text-xs text-gray-600">Orders</span>
                    </div>
                </div>
            )}
        </div>
    );
}

