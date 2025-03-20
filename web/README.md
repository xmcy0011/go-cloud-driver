# 云盘

## 开发

1. 安装 node

```bash
# 安装 nvm
$ curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
# 重新加载配置
$ source ~/.bashrc
# 安装 LTS 版本的 Node.js
$ nvm install --lts
# 使用新安装的版本
$ nvm use --lts
# 验证版本
$ node -v
v22.14.0
```
> node 版本须 >= 20.11.0

2. 安装依赖

```bash
npm install
```

3. 启动开发环境

```bash
npm run dev
```

4. 打包

```bash
npm run build
```
