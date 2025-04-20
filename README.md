<div align="center">
  <a href='' target="_blank">
    <img src="./build/appicon.png" alt="Project Icon" width="100" style="pointer-events: none;">
  </a>
  <h1>AVD Launcher</h1>
</div>

- A lightning-fast, ultra-lightweight native app for launching Android Virtual Devices.  
- Cross-platform, portable, and minimal.
<br/><br/>

## 🔧 Dev Prerequisites

To build or run in live dev mode, follow the [official Wails installation guide](https://wails.io/docs/gettingstarted/installation).  
You'll need Go installed, along with Node and a package manager like `npm`, `yarn`, or `pnpm`.
<br/><br/>

## ⚙️ Live Development
To start the app in live development mode:
```bash
wails dev
```
This runs a Vite-powered dev server with hot reload for the frontend.
<br/><br/>

## 📦 Release Build
To generate a production-ready, standalone binary:
```bash
wails build
```
This compiles the app and outputs a native executable, ready to distribute.
<br/><br/>

## 🚀 Release Workflow

AVD Launcher uses a fully automated release pipeline powered by **GitHub Actions** and a helper script.

To create a new release:

### 📦 Step 1: Run the Release Script
In the project root, run the following npm script:
```bash
npm run release
```

This will:
1. Prompt you to select the version type (`Patch`, `Minor`, or `Major`).
2. Bump the version in `frontend/package.json`.
3. Commit the version bump and create a Git tag.
4. Push the commit and tag to GitHub.

> ℹ️ The version bump uses a conventional commit message like:  
> `chore: bumped version to v1.2.3`

### ⚙️ Step 2: GitHub Actions Kicks In
When a `v*` tag is pushed, the [`release.yml`](.github/workflows/release.yml) GitHub Actions workflow is triggered.

It automatically:
- 🔧 Builds native binaries for:
  - Linux (amd64)
  - Windows (.exe)
  - macOS (.pkg)
- 🗃 Renames and organizes the build artifacts.
- 📝 Creates a new GitHub Release and uploads the binaries with OS-specific labels.

You can view the release process under the repo's **Actions** tab.

> 🧠 _Note: This release pipeline wasn't built overnight — it took a full day of trial, error, and frustration to get it working just right. If you're struggling to set up something similar, you're not alone!_

<br/>

## Built with ❤️
This project is built with passion using:
- [Wails](https://wails.io/)
- [Go](https://go.dev/)
- [Vue 3](https://vuejs.org/)

<div align="center">
  <sub>Made with 💜 by Symon from Belgium</sub>
</div>
<div align="center">
  <sub>Powered by <a href="https://wails.io/">Wails</a></sub>
</div>
