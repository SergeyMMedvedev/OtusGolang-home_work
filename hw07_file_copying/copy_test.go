package main

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func compare(path1, path2 string) bool {
	f1, err := os.ReadFile(path1)
	if err != nil {
		log.Fatal(err)
	}
	f2, err := os.ReadFile(path2)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Equal(f1, f2)
}

var (
	sourceFile             = path.Join("testdata", "input.txt")
	destFile               = "output.txt"
	outOffset0Limit0       = path.Join("testdata", "out_offset0_limit0.txt")
	outOffset0Limit10      = path.Join("testdata", "out_offset0_limit10.txt")
	outOffset0Limit1000    = path.Join("testdata", "out_offset0_limit1000.txt")
	outOffset0Limit10000   = path.Join("testdata", "out_offset0_limit10000.txt")
	outOffset100Limit1000  = path.Join("testdata", "out_offset100_limit1000.txt")
	outOffset6000Limit1000 = path.Join("testdata", "out_offset6000_limit1000.txt")
)

func TestCopy(t *testing.T) {
	err := Copy(sourceFile, destFile, 0, 0)
	require.NoError(t, err)
	require.True(t, compare(sourceFile, destFile))
	require.True(t, compare(destFile, outOffset0Limit0))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffset0Limit10(t *testing.T) {
	err := Copy(sourceFile, destFile, 0, 10)
	require.NoError(t, err)
	require.True(t, compare(destFile, outOffset0Limit10))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffset0Limit1000(t *testing.T) {
	err := Copy(sourceFile, destFile, 0, 1000)
	require.NoError(t, err)
	require.True(t, compare(destFile, outOffset0Limit1000))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffset0Limit10000(t *testing.T) {
	err := Copy(sourceFile, destFile, 0, 10000)
	require.NoError(t, err)
	require.True(t, compare(destFile, outOffset0Limit10000))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffset100Limit1000(t *testing.T) {
	err := Copy(sourceFile, destFile, 100, 1000)
	require.NoError(t, err)
	require.True(t, compare(destFile, outOffset100Limit1000))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffset6000Limit1000(t *testing.T) {
	err := Copy(sourceFile, destFile, 6000, 1000)
	require.NoError(t, err)
	require.True(t, compare(destFile, outOffset6000Limit1000))
	t.Cleanup(func() { os.Remove(destFile) })
}

func TestCopyOffsetExceeds(t *testing.T) {
	err := Copy(sourceFile, destFile, 60000, 1000)
	require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
}
