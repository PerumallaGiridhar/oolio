package index

import (
	"bufio"
	"io"
	"os"

	"github.com/willf/bloom"
)

func countPromos(srcPath string) (int, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Buffer(make([]byte, 1024), 64*1024)

	promoCount := 0
	for scan.Scan() {
		promoCount++
	}

	return promoCount, nil
}

func buildBloomFilterFromFile(srcPath string) (*bloom.BloomFilter, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Buffer(make([]byte, 1024), 64*1024)

	promoCount, err := countPromos(srcPath)
	if err != nil {
		return nil, err
	}
	bloomFilter := bloom.NewWithEstimates(uint(promoCount), 0.001)
	for scan.Scan() {
		promo := scan.Text()
		if promo == "" {
			continue
		}
		bloomFilter.AddString(promo)
	}

	return bloomFilter, scan.Err()

}

func saveBloomFilter(bf *bloom.BloomFilter, dstPath string) error {
	tmp := dstPath + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, writeErr := bf.WriteTo(out)
	closeErr := out.Close()

	if writeErr != nil {
		_ = os.Remove(tmp)
		return writeErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}

	return os.Rename(tmp, dstPath)
}

func loadBloomFilter(srcPath string) (*bloom.BloomFilter, error) {
	file, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bf := &bloom.BloomFilter{}
	if _, err := bf.ReadFrom(file); err != nil && err != io.EOF {
		return nil, err
	}

	return bf, nil
}

func BuildOrLoadBloomFilter(srcPath string) (*bloom.BloomFilter, error) {
	dstPath := srcPath + ".bloom"

	srcStat, err := os.Stat(srcPath)
	if err != nil {
		return nil, err
	}

	if dstStat, err := os.Stat(dstPath); err == nil && srcStat.ModTime().Before(dstStat.ModTime()) {
		bf, err := loadBloomFilter(dstPath)
		return bf, err
	}

	bf, err := buildBloomFilterFromFile(srcPath)
	if err != nil {
		return nil, err
	}

	if err := saveBloomFilter(bf, dstPath); err != nil {
		return nil, err
	}

	return bf, nil
}
