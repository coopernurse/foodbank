import { defineConfig } from 'vite';

export default defineConfig({
  root: 'frontend',  // Point to the source directory for Vite
  build: {
    outDir: '../static',  // Output bundled files into the static directory
    emptyOutDir: false,    // Clean up the static directory before building
    rollupOptions: {
        input: 'frontend/index.js',  // Explicitly set the entry file to index.js
        output: {
            entryFileNames: 'js/index.js',  // Specify the output filename for the main JS entry point
            chunkFileNames: 'js/[name].js', // Chunk file naming (optional, if you have code splitting)
            assetFileNames: 'assets/[name][extname]', // Non-JS assets like images will go into 'assets/'
          },
    },
  },
  server: {
    port: 3000,  // Vite dev server port
    proxy: {
      '/': 'http://localhost:8080',  // Proxy API requests to Go server
    },
  },
});
