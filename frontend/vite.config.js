import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";

export default defineConfig({
    plugins: [react()],
    base: "/",
    build: {
        assetsDir: "",
    },
    server: {
        host: "0.0.0.0",
        proxy: {
            "/simple_upload": {
                target: "http://127.0.0.1:8088",
                changeOrigin: true,
            },
            "/api": {
                target: "http://127.0.0.1:8088",
                changeOrigin: true,
            },
            "/webdav": {
                target: "http://127.0.0.1:8088",
                changeOrigin: true,
            },
        },
    },
});
