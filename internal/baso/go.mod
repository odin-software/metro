module baso

replace internal/model => ../models

replace internal/dbstore => ../dbstore

require github.com/odin-software/metro/internal/dbstore v1.0.0

go 1.22
