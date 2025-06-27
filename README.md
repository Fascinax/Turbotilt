# 🛠️Turbotilt : Dev Temp Docker + Tilt

> **Zero‑setup local orchestration** for multi‑module Java micro‑services (Spring Boot, Quarkus, Micronaut …), with optional live‑reload via Tilt.

![status-badge](https://img.shields.io/badge/status-alpha-red)

---

## ✨ Features

| What you get                                                               | How it works                                                                 |
| -------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| 🔍 **Auto‑detection** of Maven/Gradle projects under the current directory | Scans for `pom.xml` & `build.gradle` files                                   |
| 🐳 **On‑the‑fly containerisation**                                         | Generates a temporary `Dockerfile.dev.tmp` per service                       |
| 🧩 **Composable runtime**                                                  | Creates a throw‑away `dev-compose.temp.yml` and launches `docker compose up` |
| ⚡ **Live‑reload dev‑loop** (optional)                                      | Pass `--tilt` → script autogenerates a `Tiltfile` and starts `tilt up`       |
| 🧹 **Zero footprint**                                                      | All temp files & containers cleaned on exit (Ctrl‑C)                         |
| 🏗️ **Multi‑JDK ready**                                                    | Future roadmap: detect `<java.version>` and pick 11 / 17 / 21 image          |

---

## 📦 Repository Layout

```
.
├─ dev_temp_docker.sh          # <── main entrypoint (Bash)
├─ README.md                   # you are here
└─ services/                   # example Spring/Quarkus modules (optional)
```

> The script is self‑contained; **no file is committed** to your repo when you run it.

---

## 🚀 Quick Start

```bash
# clone your micro‑services mono‑repo
$ git clone git@github.com:your‑org/your‑repo.git && cd your‑repo

# make the helper executable
$ chmod +x dev_temp_docker.sh

# 1⃣ Classic Docker Compose mode (no live‑reload)
$ ./dev_temp_docker.sh

# 2⃣ Tilt‑powered live‑reload mode
$ ./dev_temp_docker.sh --tilt

# 3⃣ Run a subset of services
$ SERVICES="service‑a service‑b" ./dev_temp_docker.sh --tilt
```

### Default Port Mapping

First service exposed on **8081**, then 8082, 8083 …
Customize easily by editing the generated YAML or submitting a PR 😉.

---

## 🔍 How It Works (Under the Hood)

1. **Scan phase** – find every directory ≤2 levels deep containing a `pom.xml`.
2. **Dockerfile generation** – writes a `Dockerfile.dev.tmp` (multi‑stage Maven build).
3. **Compose + Tiltfile render** – builds `dev-compose.temp.yml`; if `--tilt`, adds live‑update rules.
4. **Launch** – either `docker compose up` or `tilt up`.
5. **Cleanup** – `trap` removes containers & temp files on exit.

Flow diagram:

```
        +-------- scan --------+
        |  *.xml / *.gradle    |
        +----------+-----------+
                   ↓
        +---- generate Dockerfile.dev.tmp ----+
        |   multi‑stage build (17‑JDK)        |
        +----------+-----------+-------------+
                   ↓
        +---- render dev‑compose.temp.yml ----+
        +---- render Tiltfile (optional) -----+
                   ↓
          docker compose ⬆︎   |   tilt ⬆︎
                   ↓         |   ↻ live‑update
              container(s)   |   class hot‑swap
```

---

## 🪛 Configuration & Overrides

| Variable / Flag         | Purpose                                            | Default               |
| ----------------------- | -------------------------------------------------- | --------------------- |
| `SERVICES` env          | Space‑delimited list of sub‑directories to include | All detected services |
| `--tilt`                | Enable Tilt live‑reload workflow                   | Off                   |
| `.env` file per service | Runtime environment variables                      | *Optional*            |

Roadmap items (PRs welcome):

* **`tilt-stack.yaml`** manifest to override port, JDK, buildpack.
* **Quarkus Gradle detection**.
* **Multi‑runtime support** (Node, Go, Python).

---

## 🐞 Troubleshooting

| Symptom                           | Fix                                                                 |
| --------------------------------- | ------------------------------------------------------------------- |
| `❌ No Maven micro‑services found` | Ensure each service has a `pom.xml` (Gradle support WIP).           |
| Port already in use               | Set `SERVICES` env to run fewer modules or edit the generated YAML. |
| `tilt: command not found`         | Install Tilt ≥ 0.33 or run without `--tilt`.                        |

---

## 🤝 Contributing

1. **Fork** the repo & create your feature branch (`git checkout -b feat/detector-gradle`).
2. **Commit** your changes (`git commit -m "Add Gradle detector"`).
3. **Push** to the branch (`git push origin feat/detector-gradle`).
4. **Open** a Pull Request.

Please run `shellcheck dev_temp_docker.sh` and `bats` tests before submitting.

---

## 📜 License

MIT — do what you want, but don’t blame us if your laptop takes off 🚀.

---

## ❤️ Acknowledgements

* [Tilt](https://tilt.dev) for the excellent dev‑loop engine.
* Maven, Eclipse Temurin, Docker, and the OSS community.

> *Made with caffeine by Océan ☕*
