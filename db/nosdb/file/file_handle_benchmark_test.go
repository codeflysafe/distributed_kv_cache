/*
 * @Author: sjhuang
 * @Date: 2022-04-24 11:45:44
 * @LastEditTime: 2022-04-24 11:54:06
 * @FilePath: /nosdb/file/file_handle_benchmark_test.go
 */
package file

import "testing"

// benchmark 测试
// sync once
func BenchmarkTestFileHandle_IOWriteAt_sync_once(b *testing.B) {
	h, err := OpenFile(STANDARD_IO, "a.txt", "./bench_mack", 1<<30)
	if err != nil {
		b.Error(err)
	}
	defer h.Close()
	var offset int64 = 0
	value := []byte("test benchmark insert sync, adsasdasdasdasdadasdasd")
	for i := 0; i < b.N; i++ {
		offset, err = h.WriteAt(offset, value)
		if err != nil {
			b.Error(err)
		}
	}
	h.Sync()
}

func BenchmarkTestFileHandle_IOWriteAt_sync_every(b *testing.B) {
	h, err := OpenFile(STANDARD_IO, "a.txt", "./bench_mack", 1<<30)
	if err != nil {
		b.Error(err)
	}
	defer h.Close()
	var offset int64 = 0
	value := []byte("test benchmark insert sync, adsasdasdasdasdadasdasd")
	for i := 0; i < b.N; i++ {
		offset, err = h.WriteAt(offset, value)
		if err != nil {
			b.Error(err)
		}
		h.Sync()
	}
}

func BenchmarkTestFileHandle_Mmap_WriteAt_sync_once(b *testing.B) {
	h, err := OpenFile(M_MAP, "a.txt", "./bench_mack", 1<<30)
	if err != nil {
		b.Error(err)
	}
	defer h.Close()
	var offset int64 = 0
	value := []byte("test benchmark insert sync, adsasdasdasdasdadasdasd")
	for i := 0; i < b.N; i++ {
		offset, err = h.WriteAt(offset, value)
		if err != nil {
			b.Error(err)
		}
	}
	h.Sync()
}

func BenchmarkTestFileHandle_Mmap_WriteAt_sync_every(b *testing.B) {
	h, err := OpenFile(M_MAP, "a.txt", "./bench_mack", 1<<30)
	if err != nil {
		b.Error(err)
	}
	defer h.Close()
	var offset int64 = 0
	value := []byte("test benchmark insert sync, adsasdasdasdasdadasdasd")
	for i := 0; i < b.N; i++ {
		offset, err = h.WriteAt(offset, value)
		if err != nil {
			b.Error(err)
		}
		h.Sync()
	}
}
