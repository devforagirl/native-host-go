# FlowMeter Native Host (Go 版本) 设计方案

## 1. 项目概述

本项目的目标是将原有的基于 Node.js 的原生消息宿主重构为使用 Go 语言编写的独立二进制文件。旨在消除运行时的依赖，极大地降低内存占用，并提供极简的跨平台安装体验。

- 项目路径：D:\GitHub\FlowMeter\native-host-go
- 核心目标：
    - 独立性：不依赖 Node.js，编译后直接运行。
    - 轻量化：二进制体积 $\approx 10\text{MB}$，运行内存 $\le 15\text{MB}$。
    - 全平台：支持 Windows, macOS, Linux。

## 2. 整体架构设计

### 2.1 数据流向图

Chrome 浏览器 $\xrightarrow{\text{Standard I/O}}$ FlowMeter Go Binary $\xrightarrow{\text{System Call}}$ OS Network Stats

### 2.2 核心模块分解

我们将程序分为三个解耦的模块：

#### A. 采样模块 (sampler)
- 职责：从操作系统获取网络接口的实时流量数据。
- 技术实现：使用 `github.com/shirou/gopsutil/v3/net` 库。
- 逻辑：
    a. 获取所有网络接口的 BytesSent 和 BytesRecv 总量。
    b. 每秒执行一次采样，通过 $\Delta \text{Value} / \Delta \text{Time}$ 计算实时速率。
    c. 将结果封装为与原 TS 版本一致的 JSON 结构。

#### B. 通信模块 (messenger)
- 职责：实现 Chrome Native Messaging 协议。
- 协议细节：
    - 发送 (Write)：在 JSON 字符串前附加 4 字节的长度头（Little-endian）。
    - 接收 (Read)：读取前 4 字节确定消息长度，然后读取对应长度的 JSON 字符串。
- 接口：提供简单的 `Send(message)` 和 `Receive()` 方法。

#### C. 注册模块 (registrar)
- 职责：将 Host 程序注册到浏览器的白名单中。
- 实现细节：
    - Manifest 生成：根据 `HOST_NAME` 和当前二进制文件的绝对路径生成 JSON。
    - 路径适配：
        - Windows $\rightarrow$ 写入 `HKEY_CURRENT_USER\Software\Google\Chrome\NativeMessagingHosts\...`
        - macOS/Linux $\rightarrow$ 写入 `~/Library/...` 或 `~/.config/...` 的 JSON 文件。
- 指令：通过命令行参数 `--register` 触发。

## 3. 接口与命令行定义

程序将通过单一二进制文件提供两种运行模式：

| 命令 | 模式 | 行为 |
| :--- | :--- | :--- |
| `flowmeter-host` | 宿主模式 | 进入死循环 $\rightarrow$ 采样网速 $\rightarrow$ 发送给浏览器 $\rightarrow$ 等待退出信号。 |
| `flowmeter-host --register` | 注册模式 | 检测 OS $\rightarrow$ 写入注册表/配置文件 $\rightarrow$ 打印成功信息 $\rightarrow$ 立即退出。 |

## 4. 跨平台编译与分发策略

### 4.1 编译矩阵

我们将利用 Go 的交叉编译能力，在开发机上一次性生成所有版本：

| 目标 OS | 目标架构 | 命令示例 | 输出文件名 |
| :--- | :--- | :--- | :--- |
| Windows | amd64 | `GOOS=windows GOARCH=amd64 go build` | `flowmeter-host-win.exe` |
| macOS | arm64 | `GOOS=darwin GOARCH=arm64 go build` | `flowmeter-host-macos` |
| Linux | amd64 | `GOOS=linux GOARCH=amd64 go build` | `flowmeter-host-linux` |

### 4.2 分发路径

1. **npm 分发 (flowmeter-host)**:
    - 将上述三个二进制文件放入 npm 包。
    - `postinstall` 脚本 $\rightarrow$ 检测 OS $\rightarrow$ 执行 `binary --register`。
2. **独立文件分发**:
    - 将二进制文件上传至 GitHub Releases。
    - 用户下载 $\rightarrow$ 运行 $\rightarrow$ 自动注册 $\rightarrow$ 完成。

## 5. 关键技术挑战与对策

- **挑战 1：权限问题**。在 macOS/Linux 上写入配置目录可能需要权限。
    - 对策：优先写入用户家目录（User Home）下的配置文件夹，无需 sudo 权限。
- **挑战 2：路径动态化**。注册时需要写入二进制文件的绝对路径，但安装路径各异。
    - 对策：在 `--register` 模式下，使用 `os.Executable()` 获取当前运行程序的绝对路径。
- **挑战 3：数据精度**。确保 Go 的采样频率与原 TS 版本完全一致。
    - 对策：使用 `time.Ticker` 实现严格的 1 秒心跳。

## 6. 开发计划 (Milestones)

1. **M1**：搭建 `native-host-go` 基础结构 $\rightarrow$ 实现 `messenger` 协议。
2. **M2**：实现 `sampler` 采样逻辑 $\rightarrow$ 验证 JSON 数据正确性。
3. **M3**：实现 `registrar` 注册逻辑 $\rightarrow$ 验证三端注册成功。
4. **M4**：集成测试 $\rightarrow$ 浏览器连接 $\rightarrow$ 实时显示网速。
5. **M5**：编写 npm wrapper $\rightarrow$ 完成发布准备。
