<h1 align="center">
  <a href="https://github.com/openivity/openivity.github.io">
    <img src="docs/images/logo.png" alt="Logo" height="35">
  </a>
</h1>

<div align="center">
</div>

<div align="center">
<br />

![GitHub Workflow Status](https://github.com/openivity/openivity.github.io/actions/workflows/main.yml/badge.svg)
[![Project license](https://img.shields.io/github/license/openivity/openivity.github.io.svg?style=flat-square)](LICENSE)
[![Pull Requests welcome](https://img.shields.io/badge/PRs-welcome-33CC56.svg?style=flat-square)](https://github.com/openivity/openivity.github.io/issues?q=is%3Aissue+is%3Aopen)

</div>

<details open="open">
<summary>Table of Contents</summary>

- [About](#about)
- [Getting Started](#getting-started)
- [Features](#features)
- [Roadmap](#roadmap)
- [Support](#support)
- [Project assistance](#project-assistance)
- [Contributing](#contributing)
- [Contributors](#contributors)
- [Security](#security)
- [License](#license)

</details>

---

## About

Interactive tool to view (with OpenStreetMap), edit, convert and combine multiple FIT, GPX and TCX activity files. 100% client-site power!

<img src="docs/images/sample.jpg" title="App Page" width="100%">

### Built With

- [Go](https://go.dev) - [WebAssembly](https://github.com/golang/go/wiki/WebAssembly)
- [FIT SDK for Go](https://github.com/muktihari/fit)
- [NodeJS](https://nodejs.org) - [Vite](https://vitejs.dev) - [Vue](https://vuejs.org) - [Typescript](https://www.typescriptlang.org)
- [OpenLayers](https://openlayers.org)
- [OpenStreetMap](https://www.openstreetmap.org)
- [Bootstrap](https://getbootstrap.com)
- [d3.js](https://d3js.org)

## Getting Started

Please see [Development environment setup](docs/CONTRIBUTING.md#development-environment-setup) to set-up.

## Features

- Supported files: **\*.fit**, **\*.gpx**, and **\*.tcx**
- Support opening Single or Multiple files
- Support multiple sport session in Single or Multiple files
- Map view (with OpenStreetMap)
- Graphs:
  - Elevation
  - Heart Rate Zone
  - Splits Pace
  - Pace
  - Speed
  - Cadence
  - Heart Rate
  - Power
  - Temperature
- Laps & Sessions
- **Tools**
  - Export to FIT, GPX, or TCX
  - Edit Relevant Data
    - Change Sport Type
    - Change Device
    - Trim Track
    - Conceal GPS Positions
    - Remove Fields: Cadence, Heart Rate, Power, and Temperature
  - Combine Multiple Activities
  - Split Activity per Session

## Roadmap

See the [open issues](https://github.com/openivity/openivity.github.io/issues) for a list of proposed features (and known issues).

## Support

If you have any questions or encounter any issues, feel free to open an [issue](https://github.com/openivity/openivity.github.io/issues/new) and we will assist you in resolving them.

## Project assistance

If you want to say **thank you** or/and support active development of Open Activity:

- Add a [GitHub Star](https://github.com/openivity/openivity.github.io) to the project.
- Tweet about the Open Activity.
- Write interesting articles about the project on [Dev.to](https://dev.to/), [Medium](https://medium.com/) or your personal blog.

## Contributing

First off, thanks for taking the time to contribute! Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make will benefit everybody else and are **greatly appreciated**.

Please read [our contribution guidelines](docs/CONTRIBUTING.md), and thank you for being involved!

## Contributors

The original setup of this repository is by [Openivity](https://github.com/openivity).

For a full list of all authors and contributors, see [the contributors page](https://github.com/openivity/openivity.github.io/contributors).

## Security

Openivity - Open Activity runs 100% client-side, and we don't collect any data since all data reside on the user's machine, ensuring user privacy. We strive to follow security best practices; however, we cannot assure 100% free from security breaches. Please [Report a vulnerability](https://github.com/openivity/openivity.github.io/security/advisories/new) if you find any security issue in this repository.

## License

This project is licensed under the **GPL 3**.

See [LICENSE](LICENSE) for more information.

Openivity - Open Activity is provided **"as is"** without any **warranty**. Use at your own risk.
