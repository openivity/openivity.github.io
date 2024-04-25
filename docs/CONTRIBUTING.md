# Contributing

When contributing to this repository, please first discuss the change you wish to make via issue on this repository before making a change.
Please note we have a [code of conduct](CODE_OF_CONDUCT.md), please follow it in all your interactions with the project.

## Development environment setup

### Prerequisites

This project need:

- Go version >= v1.20.x
- NodeJS >= v20.7.x

Project Git Branch:

- `master` release branch
- `dev` dev consolidate branch

### Project Setup

To set up a development environment, please follow these steps:

1. Clone the repo (and/or switch to consolidate `dev` branch)

   ```sh
   git clone https://github.com/openivity/openivity.github.io
   git checkout dev # Optional, dev branch
   ```

2. Install node dependencies

   ```sh
   npm install
   ```

3. Copy `wasm_exec.js` of your current version of Go in your local machine to `src/assets/wasm/wasm_exec.js`. This file is tightly coupled with a specific version of Go in order to compile the WebAssembly binary.

   ```sh
   cp $GOROOT/misc/wasm/wasm_exec.js src/assets/wasm/wasm_exec.js
   ```

4. Compile Go source code into wasm binary.

   ```sh
   npm run build-wasm
   ```

5. Finally, now you can run the project

   ```sh
   npm run dev
   ```

6. Type-Check, Compile and Minify for Production

   ```sh
   npm run build-wasm
   npm run build-only
   ```

## Issues and feature requests

You've found a bug in the source code, a mistake in the documentation or maybe you'd like a new feature? You can help us by [submitting an issue on GitHub](https://github.com/openivity/openivity.github.io/issues). Before you create an issue, make sure to search the issue archive -- your issue may have already been addressed!

Please try to create bug reports that are:

- _Reproducible._ Include steps to reproduce the problem.
- _Specific._ Include as much detail as possible: which version, what environment, etc.
- _Unique._ Do not duplicate existing opened issues.
- _Scoped to a Single Bug._ One bug per report.

**Even better: Submit a pull request with a fix or new feature!**

### How to submit a Pull Request

1. Search our repository for open or closed
   [Pull Requests](https://github.com/openivity/openivity.github.io/pulls)
   that relate to your submission. You don't want to duplicate effort.
2. Fork the project
3. Create your feature branch (`git checkout -b feat/amazing_feature`)
4. Commit your changes (`git commit -m 'feat: add amazing_feature'`)
5. Push to the branch (`git push origin feat/amazing_feature`)
6. [Open a Pull Request](https://github.com/openivity/openivity.github.io/compare?expand=1)
