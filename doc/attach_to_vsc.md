# 使用 Kubeconfig 连接到HPCGame Kubernetes 集群

本节将指导您如何使用 Kubeconfig 文件在本地开发环境中连接到 Kubernetes 集群。

## 前提条件

1. **Kubeconfig 文件**：确保您已经从 Kubernetes 集群管理员处获取了 Kubeconfig 文件。
2. **VSCode**：确保您已经安装了 Visual Studio Code（VSCode）。
3. **Kubernetes 插件**：确保您已经安装了 VSCode 的 Kubernetes 插件。
4. **Dev Containers 插件**：确保您已经安装了 VSCode 的 Dev Containers 插件。

## 步骤 1：获取 Kubeconfig 文件

1. 请前往 [HPCGame平台](https://hpcgame.pku.edu.cn/) 上开通账号。
2. 访问[此处](hpcgame.pku.edu.cn/kube/_/ui/) 获得 kubeconfig。
3. 将 Kubeconfig 文件保存到本地。

## 步骤 2：安装 VSCode

1. 访问 [VSCode 官方网站](https://code.visualstudio.com/) 下载并安装 VSCode。
2. 安装完成后，启动 VSCode。

## 步骤 3：安装 Kubernetes 和 Dev Containers 插件

1. 打开 VSCode，点击左侧的扩展图标（或按 `Ctrl+Shift+X` / `Command+Shift+X`(对于MacOS使用者)）。
2. 在搜索栏中输入 `Kubernetes`，找到并安装 `Kubernetes` 插件。
3. 同样地，搜索并安装 `Dev Containers` 插件。

## 步骤 4：导入 Kubeconfig 文件

1. 打开 VSCode，点击左侧的 Kubernetes 图标（或按 `Ctrl+Shift+P` / `Command+Shift+P`(对于MacOS使用者) 并输入 `Kubernetes`）。
2. 在 Kubernetes 插件界面中，点击 `Set Kubeconfig`(可能需要查看`CLUSTERS`选项卡旁边的`···`来找到这个选项)。或者，按 `Ctrl+Shift+P` / `Command+Shift+P`(对于MacOS使用者)并输入`Kubernetes: Set Kubeconfig`，点击。
3. 选择 `Add new kubeconfig`，然后浏览并选择您之前下载的 Kubeconfig 文件。
4. 确认导入后，Kubernetes 插件将显示您的集群信息。

## 步骤 5：连接容器进行开发

1. 在连接的集群中，找到`Workloads`，在其中找到`Pods`，在该栏中将列出诸多Pod。
2. 选择您想要连接的容器名称，右键并点击`Attach VSCode`选项。