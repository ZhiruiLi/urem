# Unreal Rem | Unreal 开发辅助工具

## 快速开始

```bash
go install github.com/zhiruili/urem
urem --help
```

## 功能

### 刷新工程

```bash
urem regensln PATH_TO_THE_PROJECT_FILE
# Example:
#  urem regensln projects/MyUeProject/MyUeProject.uproject
```

### 新增模块

```bash
urem newmod -o MODULE_OUTPUT_PATH MODULE_NAME
# Example:
#  urem newmod -o projects/MyUeProject/Source AnExample
#  urem newmod -o projects/MyUeProject/Plugins/MyPlug/Source AnExample
```
