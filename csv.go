package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// ReadCSV 读取CSV文件内容
func ReadCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 允许不同行有不同数量的字段

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// countLines 计算文件行数
func countLines(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// ProcessCSVStream 流式处理CSV文件，减少内存占用
func ProcessCSVStream(inputFile, outputFile string, signFunc func([]byte) (string, error), template string, concurrency int) error {
	// 计算文件总行数
	totalLines, err := countLines(inputFile)
	if err != nil {
		return err
	}
	// 减去表头行数
	totalDataLines := totalLines - 1
	if totalDataLines < 0 {
		totalDataLines = 0
	}

	// 打开输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	// 创建输出文件
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 创建CSV读写器
	reader := csv.NewReader(inFile)
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// 读取并处理表头
	header, err := reader.Read()
	if err != nil {
		return err
	}

	// 检查是否已有sign-String列
	hasSignColumn := false
	for _, col := range header {
		if col == "sign-String" {
			hasSignColumn = true
			break
		}
	}

	// 如果没有sign-String列，添加到第二列位置
	if !hasSignColumn {
		newHeader := make([]string, len(header)+1)
		newHeader[0] = header[0]
		newHeader[1] = "sign-String"
		copy(newHeader[2:], header[1:])
		header = newHeader
	}

	// 写入表头
	if err := writer.Write(header); err != nil {
		return err
	}

	// 使用channel进行并发处理
	jobs := make(chan []string, 1000)       // 任务队列
	results := make(chan []string, 1000)    // 结果队列
	errors := make(chan error, concurrency) // 错误通道
	wg := &sync.WaitGroup{}

	// 原子计数器，用于进度跟踪
	processedLines := int64(0)

	// 启动工作goroutine
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for record := range jobs {
				if len(record) == 0 {
					continue // 跳过空行
				}

				// 获取第一列数据
				firstCol := record[0]

				// 替换模板中的占位符
				content := fmt.Sprintf(template, firstCol)

				// 生成签名
				signature, err := signFunc([]byte(content))
				if err != nil {
					errors <- err
					return
				}

				// 确保记录有足够的列
				if len(record) < 2 {
					// 扩展记录到至少2列
					record = append(record, make([]string, 2-len(record))...)
				}

				// 将签名结果写入第二列
				record[1] = signature
				results <- record

				// 更新已处理行数
				current := atomic.AddInt64(&processedLines, 1)

				// 每处理1000行打印一次进度
				if current%1000 == 0 || current == int64(totalDataLines) {
					percentage := float64(current) / float64(totalDataLines) * 100
					fmt.Printf("\r已处理 %d/%d 行 (%.1f%%)", current, totalDataLines, percentage)
				}
			}
		}()
	}

	// 启动写入goroutine
	writeWg := &sync.WaitGroup{}
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()
		flushCount := 0
		for result := range results {
			if err := writer.Write(result); err != nil {
				errors <- err
				return
			}
			flushCount++
			// 每1000行刷新一次，减少内存占用
			if flushCount%1000 == 0 {
				writer.Flush()
				flushCount = 0
			}
		}
	}()

	// 开始计时
	startTime := time.Now()

	// 读取数据并发送到任务队列
	for {
		record, err := reader.Read()
		if err != nil {
			break // 文件读取完毕
		}
		jobs <- record
	}
	close(jobs)

	// 等待所有工作完成
	wg.Wait()
	close(results)
	writeWg.Wait()

	// 打印最终进度
	fmt.Printf("\r已处理 %d/%d 行 (100.0%%)\n", totalDataLines, totalDataLines)

	// 打印总耗时
	totalTime := time.Since(startTime)
	fmt.Printf("处理完成！总耗时: %v\n", totalTime)

	// 检查是否有错误
	select {
	case err := <-errors:
		return err
	default:
		return writer.Error()
	}
}

// WriteCSV 将内容写入CSV文件
func WriteCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return writer.Error()
}

// ProcessCSVRecords 处理CSV记录，对第一列进行签名并添加到第二列
func ProcessCSVRecords(records [][]string, signFunc func([]byte) (string, error), template string) ([][]string, error) {
	if len(records) == 0 {
		return records, nil
	}

	// 处理表头
	header := records[0]
	if len(header) == 0 {
		return nil, fmt.Errorf("CSV文件表头为空")
	}

	// 检查是否已有sign-String列
	hasSignColumn := slices.Contains(header, "sign-String")

	// 如果没有sign-String列，添加到第二列位置
	if !hasSignColumn {
		newHeader := make([]string, len(header)+1)
		newHeader[0] = header[0]
		newHeader[1] = "sign-String"
		copy(newHeader[2:], header[1:])
		records[0] = newHeader
	}

	// 处理数据行
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) == 0 {
			continue // 跳过空行
		}

		// 获取第一列数据
		firstCol := record[0]

		// 替换模板中的占位符
		content := fmt.Sprintf(template, firstCol)

		// 生成签名
		signature, err := signFunc([]byte(content))
		if err != nil {
			return nil, fmt.Errorf("处理第%d行时签名失败: %w", i+1, err)
		}

		// 确保记录有足够的列
		if len(record) < 2 {
			// 扩展记录到至少2列
			record = append(record, make([]string, 2-len(record))...)
		}

		// 将签名结果写入第二列
		record[1] = signature
		records[i] = record
	}

	return records, nil
}
