# netpbm-parse

解析 Netpbm 图像：P1–P6 六种格式、头注释、P4 的行填充与 P5/P6 的双字节采样。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/netpbm --sample=comment
go run ./cmd/netpbm --sample=pad
go run ./cmd/netpbm --sample=maxval
go run ./cmd/netpbm --sample=rgb
go run ./cmd/netpbm --sample=work
```

## 口径

- **头**：`#` 到行尾是注释，解析头时跳过；头由魔数、宽、高、以及 P1–P3/P5/P6 的 maxval 组成。
- **ASCII**：P1–P3 的样本是十进制数字，P3 每个像素三个样本。
- **P4**：每行按位打包，每行补齐到整字节，位从最高位开始。
- **P5/P6**：maxval 小于 256 时每样本一个字节；否则每样本两字节大端。
- **规模**：解析只扫一遍数据。

## 不变量

- 解析出的样本数与宽、高、每像素样本数一致。
- 头注释不影响宽高与 maxval。
- `scanned` 不随样本数平方级放大：两千个样本的 `scanned` 不超过 6000。

## 输出契约

`Image` 含魔数、宽、高、maxval 与样本列表。场景打印一行 JSON。
`--sample=work` 打印 `{"samples": N, "scanned": N}`。
