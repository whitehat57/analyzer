#!/bin/bash

# Nama file Go
GO_FILE="analyzer.go"
BINARY_NAME="analyzer"

# 1. Mengecek apakah Go terinstall
if ! [ -x "$(command -v go)" ]; then
  echo "Go tidak ditemukan! Silakan install Go terlebih dahulu." >&2
  exit 1
fi

# 2. Membuat module baru dan melakukan inisialisasi
echo "Inisialisasi Go module..."
go mod init analyzer || { echo "Gagal inisialisasi Go module."; exit 1; }

# 3. Menambahkan dependensi yang diperlukan
echo "Menambahkan dependensi..."
go mod tidy || { echo "Gagal menambahkan dependensi."; exit 1; }

# 4. Kompilasi file Go menjadi binary
echo "Mengkompilasi file Go menjadi binary..."
go build -o $BINARY_NAME $GO_FILE || { echo "Gagal mengkompilasi."; exit 1; }

# 5. Menghapus file .go asli dan file module Go
echo "Menghapus file .go dan Go mod files..."
rm -f $GO_FILE go.mod go.sum

# 6. Menampilkan pesan sukses
echo "Instalasi berhasil. Binary '$BINARY_NAME' telah dibuat."
echo "Jalankan program dengan perintah ./$BINARY_NAME"
