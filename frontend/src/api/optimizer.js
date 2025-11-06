const BASE_URL = "http://localhost:8080/api";

export async function getSampleData() {
    try {
        const res = await fetch(`${BASE_URL}/sample-data`);
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error fetching sample data:", error);
        throw error;
    }
}

export async function optimizeRoutes(data) {
    try {
        const res = await fetch(`${BASE_URL}/optimize`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error optimizing routes:", error);
        throw error;
    }
}

export async function optimizeWithAnalytics(data, useRealRoutes = false) {
    try {
        const res = await fetch(`${BASE_URL}/optimize-analytics`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                orders: data.orders,
                shoppers: data.shoppers,
                useRealRoutes: useRealRoutes
            }),
        });
        if (!res.ok) {
            throw new Error(`HTTP error! status: ${res.status}`);
        }
        return await res.json();
    } catch (error) {
        console.error("Error optimizing with analytics:", error);
        throw error;
    }
}

export async function healthCheck() {
    try {
        const res = await fetch(`${BASE_URL}/health`);
        return await res.json();
    } catch (error) {
        console.error("Error checking health:", error);
        throw error;
    }
}
