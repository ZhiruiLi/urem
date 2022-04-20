# Unreal Rem | Unreal 开发辅助工具

## 快速开始

```bash
go install github.com/zhiruili/urem
urem --help
```

## 功能

### 刷新工程

```bash
urem genvs PATH_TO_THE_PROJECT_FILE
urem genclang PATH_TO_THE_PROJECT_FILE
# Example:
#  urem genvs projects/MyUeProject/MyUeProject.uproject
#  urem genclang projects/MyUeProject/MyUeProject.uproject
```

### 新增模块

```bash
urem newmod MODULE_NAME MODULE_OUTPUT_PATH
# Example:
#  urem newmod AnExample projects/MyUeProject/Source
#  urem newmod AnExample projects/MyUeProject/Plugins/MyPlug/Source
```

### 新增 gitignore 模板

```bash
urem newignore PATH_TO_THE_PROJECT_FILE
# Example:
#  urem newignore projects/MyUeProject/MyUeProject.uproject
```

### 新增 clang-format 模板

```bash
urem newformat PATH_TO_THE_PROJECT_FILE
# Example:
#  urem newformat projects/MyUeProject/MyUeProject.uproject
```
