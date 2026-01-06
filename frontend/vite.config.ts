import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "../web",   // 👈 输出到 Gin 项目
    emptyOutDir: true
  }
})
