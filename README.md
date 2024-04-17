# Openivity: Open Activity

Interactive tool to view (with OpenStreetMap view), edit, convert and combine multiple FIT, GPX and TCX activity files. 100% client-site power!

## Getting Started

### Prerequisites

- Go version >= v1.20.x
- NodeJS >= v20.7.x

### Project Setup

- Install node dependencies

  ```sh
  npm install
  ```

- Copy `wasm_exec.js` of your current version of Go in your local machine to `src/assets/wasm/wasm_exec.js`. This file is tightly coupled with a specific version of Go in order to compile the WebAssembly binary.

  ```sh
  cp $GOROOT/misc/wasm/wasm_exec.js src/assets/wasm/wasm_exec.js
  ```

- Compile Go source code into wasm binary.

  ```sh
  npm run build-wasm
  ```

- Finally, now you can run the project

  ```sh
  npm run dev
  ```

- Type-Check, Compile and Minify for Production

  ```sh
  npm run build-wasm
  npm run build-only
  ```
