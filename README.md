# 全民k歌下载

练习python异步。顺便把从网易云音乐歌单下架的cover歌曲下载一下。

## install && usage


```bash
# install
pip install -r requirements.txt

# python run
python3 main.py "https://kg.qq.com/node/personal?uid=1243454"


# go run

go mod tidy
go build
./qqkg_download "https://kg.qq.com/node/personal?uid=1243454"

```

## TODO
用不同语言实现

- [x] golang
- [ ] lua

