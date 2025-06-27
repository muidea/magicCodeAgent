# magicCodeAgent

## 📋 前置条件

1. 已安装 `magicCodeAgent` 可执行文件至 `/usr/local/bin/magicCodeAgent`
2. 系统已启用 systemd 用户服务功能
3. 当前用户具有家目录写入权限

## 🛠️ 安装步骤

### 1. 创建配置文件目录

```bash
# 创建配置存储目录
mkdir -p ~/.magicCodeAgent
chmod 700 ~/.magicCodeAgent
```

### 2. 创建 systemd 用户服务文件

```bash
mkdir -p ~/.config/systemd/user
nano ~/.config/systemd/user/magicCodeAgent.service
```

粘贴以下内容(注意替换用户根目录)：

```ini
[Unit]
Description=magicCodeAgent (User Service)

[Service]
Type=simple
ExecStart=/usr/local/bin/magicCodeAgent -config={用户根目录}/.magicCodeAgent/config.json
Restart=always
RestartSec=5

MemoryLimit=512M
CPUQuota=80%

[Install]
WantedBy=default.target
```

### 3. 启动并启用服务

```bash
# 重载服务配置
systemctl --user daemon-reload

# 启用开机自启
systemctl --user enable magicCodeAgent

# 立即启动服务
systemctl --user start magicCodeAgent
```

## ✅ 服务验证

```bash
# 检查服务状态
systemctl --user status magicCodeAgent

# 查看实时日志
journalctl --user-unit=magicCodeAgent -f

```

## 功能使用

http://127.0.0.1:8090/embed/ 

## 🔧 服务管理命令

| 命令 | 功能 |
|------|------|
| `systemctl --user start magicCodeAgent` | 启动服务 |
| `systemctl --user stop magicCodeAgent` | 停止服务 |
| `systemctl --user restart magicCodeAgent` | 重启服务 |
| `systemctl --user status magicCodeAgent` | 查看状态 |
| `systemctl --user enable magicCodeAgent` | 启用开机自启 |
| `systemctl --user disable magicCodeAgent` | 禁用开机自启 |
| `journalctl --user-unit=magicCodeAgent -f` | 查看实时日志 |
