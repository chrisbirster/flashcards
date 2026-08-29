import { defineConfig } from "vite";
import solid from "@solidjs/vite-plugin";
import stylex from "@stylexjs/unplugin/vite";

export default defineConfig({
  plugins: [stylex(), solid()],
  server: {
    port: 5173,
    strictPort: true,
    proxy: {
      "/api": "http://localhost:8000",
    },
  },
});
