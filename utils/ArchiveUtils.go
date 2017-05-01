package utils

import (
	"github.com/mholt/archiver"
	"fmt"
)

func ExtractTarXZ(tarBall string, dest string){
	fmt.Println("Extracting " + tarBall)
	archiver.TarXZ.Open(tarBall, dest)
}

func ExtractZip(zip string, dest string){
	fmt.Println("Extracting " + zip)
	archiver.Zip.Open(zip, dest)
}