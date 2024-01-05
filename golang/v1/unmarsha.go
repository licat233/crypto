package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"io"
	"net/http"
)

// 解析被压缩的request数据，req必须为指针
func UnmarshalReqZipData(r *http.Request, req any) (compressType string, err error) {
	if req == nil {
		panic("req is nil")
	}
	resp := r.Response
	// 读取请求体数据
	compressedData, err := io.ReadAll(r.Body)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		var buf bytes.Buffer
		io.WriteString(&buf, "failed to read request body")
		resp.Write(&buf)
		return
	}
	var decompressedData []byte
	compressType = r.Header.Get("Content-Encoding")
	switch compressType {
	case "deflate":
		decompressedData, err = decompressedFlateData(compressedData)
	case "gzip":
		decompressedData, err = decompressGzipData(compressedData)
	case "zlib":
		decompressedData, err = decompressZlibData(compressedData)
	default:
		err = json.NewDecoder(r.Body).Decode(req)
	}
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		var buf bytes.Buffer
		io.WriteString(&buf, "the requested data format is incorrect")
		resp.Write(&buf)
		return
	}
	err = json.Unmarshal(decompressedData, req)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		var buf bytes.Buffer
		io.WriteString(&buf, "failed to json unmarshal request data")
		resp.Write(&buf)
		return
	}
	return
}

func decompressedFlateData(compressedData []byte) ([]byte, error) {
	// 创建deflate解压缩器
	reader := flate.NewReader(bytes.NewReader(compressedData))
	defer reader.Close()

	return io.ReadAll(reader)
}

func compressFlateData(data any) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	_, err = writer.Write(jsonData)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompressGzipData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

func compressGzipData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	defer writer.Close()

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decompressZlibData(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	zlibReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	return io.ReadAll(zlibReader)
}

func compressZlibData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zlibWriter := zlib.NewWriter(&buf)
	defer zlibWriter.Close()

	_, err := zlibWriter.Write(data)
	if err != nil {
		return nil, err
	}

	err = zlibWriter.Flush()
	if err != nil {
		return nil, err
	}

	err = zlibWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
