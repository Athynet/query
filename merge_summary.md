此次合并主要添加了SM2签名算法支持，扩展了签名功能的多样性，同时优化了代码实现和文档说明。
| 文件 | 变更 |
|------|------|
| README.md | - 更新功能特点，增加SM2算法支持<br>- 更新技术栈，增加github.com/tjfoc/gmsm依赖<br>- 更新使用方法，增加SM2算法使用示例<br>- 更新命令行参数，增加-a参数用于指定签名算法<br>- 更新私钥格式要求，增加SM2私钥说明<br>- 更新生成私钥方法，增加SM2私钥生成方法<br>- 更新核心函数列表，增加LoadSM2PrivateKey和SM2_Sign函数 |
| csv.go | - 修改检查sign-String列的方式，从使用slices.Contains改为手动遍历，提高兼容性 |
| go.mod | - 增加github.com/tjfoc/gmsm v1.4.1依赖，用于SM2签名实现 |
| main.go | - 增加-algorithm参数用于指定签名算法<br>- 增加算法参数验证逻辑<br>- 根据选择的算法加载不同私钥并创建相应签名函数<br>- 更新输出信息，显示当前使用的签名算法 |
| rsa.go | - 增加LoadSM2PrivateKey函数用于加载SM2私钥<br>- 增加SM2_Sign函数用于执行SM2签名<br>- 增加github.com/tjfoc/gmsm/sm2的导入 |
| go.sum | - 新增依赖文件，包含github.com/tjfoc/gmsm及其依赖项的版本信息 |
| output_rsa.csv | - 新增测试输出文件，包含RSA签名结果 |
| rsa.key | - 新增RSA私钥文件，用于测试 |
| rsa.pem | - 新增RSA私钥PEM文件，用于测试 |
| sm2.key | - 新增SM2私钥文件，用于测试 |
| test.csv | - 新增测试输入文件，用于验证签名功能 |