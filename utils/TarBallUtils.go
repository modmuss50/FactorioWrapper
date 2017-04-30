package utils

import (
	"github.com/mholt/archiver"
	"fmt"
)

func ExtractTarXZ(tarBall string, dest string){
	fmt.Println("Extracting " + tarBall)
	archiver.TarXZ.Open(tarBall, dest)
}