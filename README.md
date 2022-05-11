# Unreal Rem | Unreal 开发辅助工具

## 快速开始

```bash
go install github.com/zhiruili/urem
urem --help
```

## 功能

### 刷新工程

```bash
urem gen vs PATH_TO_THE_PROJECT_FILE
urem gen clang PATH_TO_THE_PROJECT_FILE
# Example:
#  urem gen vs projects/MyUeProject/MyUeProject.uproject
#  urem gen clang projects/MyUeProject/MyUeProject.uproject
```

### 新增模块

```bash
urem new mod MODULE_NAME MODULE_OUTPUT_PATH
# Example:
#  urem new mod AnExample projects/MyUeProject/Source
#  urem new mod AnExample projects/MyUeProject/Plugins/MyPlug/Source
```

### 新增 gitignore 模板

```bash
urem new ig PATH_TO_THE_PROJECT_FILE
# Example:
#  urem new ig projects/MyUeProject/MyUeProject.uproject
```

### 新增 clang-format 模板

```bash
urem new fmt PATH_TO_THE_PROJECT_FILE
# Example:
#  urem new fmt projects/MyUeProject/MyUeProject.uproject
```

### 查看工程用的 UE 的版本和安装路径信息

```bash
urem info ue PATH_TO_THE_PROJECT_FILE
# Example:
#  urem info ue projects/MyUeProject/MyUeProject.uproject
```

### 搜索一个接口的定义位置

```bash
urem info def CLASS_NAME_PATTERN
# Example:
#  urem info def projects/MyUeProject/MyUeProject.uproject
```

### 查看枚举值

``` bash
urem info enum ENUM_TYPE
# Example:
#  urem info enum modtype
#  urem info enum loadphase
```
