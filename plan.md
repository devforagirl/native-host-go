# FlowMeter Native Host (Go) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 重构原生消息宿主为 Go 语言实现，消除 Node.js 依赖，降低资源占用并实现跨平台独立运行。

**Architecture:** 采用分层架构，将通信协议 (`messenger`)、数据采集 (`sampler`) 和系统注册 (`registrar`) 解耦。主程序 `main.go` 负责协调命令行参数和运行主循环。

**Tech Stack:** Go 1.21+, `github.com/shirou/gopsutil/v3/net` (网络统计), `golang.org/x/sys/windows/registry` (Windows 注册表)。

---

## File Structure

- `native-host-go/go.mod`: 项目依赖管理。
- `native-host-go/main.go`: 程序入口，处理 `--register` 标志及主运行循环。
- `native-host-go/internal/messenger/messenger.go`: 实现 Chrome Native Messaging 协议 (4字节长度头 + JSON)。
- `native-host-go/internal/messenger/messenger_test.go`: 通信协议单元测试。
- `native-host-go/internal/sampler/sampler.go`: 网络流量采样及速率计算。
- `native-host-go/internal/sampler/sampler_test.go`: 采样逻辑单元测试。
- `native-host-go/internal/registrar/registrar.go`: 注册逻辑接口定义。
- `native-host-go/internal/registrar/registrar_windows.go`: Windows 注册表写入实现。
- `native-host-go/internal/registrar/registrar_unix.go`: macOS/Linux 配置文件写入实现。

---

## Implementation Tasks

### Task 1: Project Initialization & Messenger Protocol (M1)

**Files:**
- Create: `native-host-go/go.mod`
- Create: `native-host-go/internal/messenger/messenger.go`
- Create: `native-host-go/internal/messenger/messenger_test.go`

- [ ] **Step 1: Initialize Go module**
  Run: `cd native-host-go && go mod init flowmeter-host`

- [ ] **Step 2: Write failing test for `ReadMessage`**
  Implement `TestReadMessage` in `messenger_test.go` that feeds a 4-byte length header + JSON and expects the JSON string.

- [ ] **Step 3: Implement `ReadMessage` in `messenger.go`**
  Use `binary.Read(r, binary.LittleEndian, &length)` to read the header, then `io.ReadFull` for the body.

- [ ] **Step 4: Write failing test for `SendMessage`**
  Implement `TestSendMessage` that checks if the output starts with the correct 4-byte little-endian length.

- [ ] **Step 5: Implement `SendMessage` in `messenger.go`**
  Use `binary.Write(w, binary.LittleEndian, uint32(len(msg)))` followed by the message bytes.

- [ ] **Step 6: Run tests and commit**
  Run: `go test ./internal/messenger/... -v`
  Commit: `git commit -m "feat(messenger): implement chrome native messaging protocol"`

### Task 2: Network Sampler Implementation (M2)

**Files:**
- Create: `native-host-go/internal/sampler/sampler.go`
- Create: `native-host-go/internal/sampler/sampler_test.go`

- [ ] **Step 1: Add `gopsutil` dependency**
  Run: `go get github.com/shirou/gopsutil/v3/net`

- [ ] **Step 2: Implement `GetTotalBytes`**
  Create a function that sums `BytesSent` and `BytesRecv` across all active interfaces.

- [ ] **Step 3: Implement `Sampler` struct and `Sample()` method**
  The `Sampler` should store the previous byte count and timestamp. `Sample()` calculates `(current - prev) / (time_diff)`.

- [ ] **Step 4: Implement JSON formatting**
  Ensure the output matches the original TS structure: `{"sent": 123, "recv": 456, "total": 789}`.

- [ ] **Step 5: Write and run sampler tests**
  Mock the byte counts to verify the rate calculation logic.
  Run: `go test ./internal/sampler/... -v`

- [ ] **Step 6: Commit**
  Commit: `git commit -m "feat(sampler): implement network traffic sampling and rate calculation"`

### Task 3: OS Registrar Implementation (M3)

**Files:**
- Create: `native-host-go/internal/registrar/registrar.go`
- Create: `native-host-go/internal/registrar/registrar_windows.go`
- Create: `native-host-go/internal/registrar/registrar_unix.go`

- [ ] **Step 1: Define `Registrar` interface**
  Method: `Register(hostName, binaryPath string) error`.

- [ ] **Step 2: Implement Windows Registry writer**
  Use `golang.org/x/sys/windows/registry` to create the key `HKCU\Software\Google\Chrome\NativeMessagingHosts\com.flowmeter.host` with the value pointing to the manifest JSON.

- [ ] **Step 3: Implement Unix config writer**
  Write the manifest JSON to `~/.config/google-chrome/NativeMessagingHosts/com.flowmeter.host.json`.

- [ ] **Step 4: Implement Manifest JSON generator**
  Function to create the JSON content containing the `name` and `path` to the binary.

- [ ] **Step 5: Test registration (Manual/Semi-auto)**
  Run the binary with `--register` on Windows and verify registry key existence.

- [ ] **Step 6: Commit**
  Commit: `git commit -m "feat(registrar): implement cross-platform browser host registration"`

### Task 4: Main Entry & Integration (M4)

**Files:**
- Create: `native-host-go/main.go`

- [ ] **Step 1: Implement CLI flag handling**
  Use `os.Args` to detect `--register`. If present, call `registrar.Register` and exit.

- [ ] **Step 2: Implement main loop**
  1. Initialize `messenger` using `os.Stdin` and `os.Stdout`.
  2. Initialize `sampler`.
  3. Start a `time.Ticker(1 * time.Second)`.
  4. On tick: `sampler.Sample()` $\rightarrow$ `messenger.SendMessage(json)`.

- [ ] **Step 3: Implement graceful shutdown**
  Listen for `os.Interrupt` or EOF from `os.Stdin` to exit cleanly.

- [ ] **Step 4: End-to-End Verification**
  1. Run `flowmeter-host --register`.
  2. Restart Chrome.
  3. Verify that the extension connects and receives real-time speed updates.

- [ ] **Step 5: Commit**
  Commit: `git commit -m "feat(main): integrate sampler and messenger into main loop"`

### Task 5: Distribution Wrapper (M5)

**Files:**
- Create: `native-host-go/build.sh` (or a Go build script)

- [ ] **Step 1: Create build script for cross-compilation**
  Implement the matrix: `GOOS=windows GOARCH=amd64`, `GOOS=darwin GOARCH=arm64`, `GOOS=linux GOARCH=amd64`.

- [ ] **Step 2: Create npm package structure**
  Prepare the `package.json` for the `flowmeter-host` wrapper.

- [ ] **Step 3: Implement `postinstall` script**
  A script that detects OS $\rightarrow$ picks the right binary $\rightarrow$ runs `binary --register`.

- [ ] **Step 4: Final verification**
  Install via npm and verify the host is registered and working.

- [ ] **Step 5: Final Commit**
  Commit: `git commit -m "chore: add build scripts and npm distribution wrapper"`
