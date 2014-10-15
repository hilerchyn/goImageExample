package main

import (
	"syscall"
	"log"
	"os"
	"unsafe"
	"time"
	"strconv"

	"image"
	"image/png" // register the PNG format with the image package
	"image/color"
	"strings"
	"fmt"
)


var flagX = 76
var flagY = 80


func main(){
	deviceName := "/dev/ttyUSB0"
	rate := uint32(syscall.B9600)

	f, err := os.OpenFile(deviceName, syscall.O_RDWR, 0660 | syscall.S_IFCHR)
	if err != nil {
		log.Println("open device file error:", err)
	}

	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()

	fd := f.Fd()
	t := syscall.Termios{
		Iflag:  syscall.IGNPAR | syscall.IGNBRK | syscall.IXON | syscall.IXOFF,
		Oflag:  syscall.IGNPAR | syscall.IGNBRK | syscall.IXON | syscall.IXOFF,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | rate,
		Cc:     [32]uint8{syscall.VMIN: 1},
		Ispeed: rate,
		Ospeed: rate,
	}

	if _, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		//return nil, errno
	}

	if err = syscall.SetNonblock(int(fd), false); err != nil {
		//return
	}




	//
	infile, err := os.Open("./resources/demo.png")
	if err != nil {
		// replace this with real error handling
		log.Println(err)
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, err := png.Decode(infile)
	if err != nil {
		// replace this with real error handling
		log.Println(err)
	}


	//
	f.Write([]byte("G92\n"))
	time.Sleep(time.Second/2)
	f.Write([]byte("G91\n"))
	time.Sleep(time.Second/2)
	f.Write([]byte("G21\n"))
	time.Sleep(time.Second/2)

	f.Write([]byte("M03 L130\n"))

	//time.Sleep(time.Second * 5)

	//f.Write([]byte("M05\n"))

	//os.Exit(0)


	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)

	//time.Sleep(time.Second*5)
	for y := 1; y <= h; y++ {

		if y>1 {
			f.Write([]byte("G00 X-"+strconv.Itoa(w)+" Y-1 F0\n"))
			time.Sleep(time.Second/2)
		}

		log.Println("X:", y)


		for x := 1; x <= w; x++ {
			oldColor := src.At(x, y)

			grayColor := color.GrayModel.Convert(oldColor)

			gray.Set(x, y, grayColor)

			grayVal := strings.Replace(fmt.Sprint(grayColor), "{", "", 1)
			grayVal = strings.Replace(fmt.Sprint(grayVal), "}", "", 1)

			val, _ := strconv.Atoi(grayVal)
			log.Println("M03 L"+strconv.Itoa(255-val)+"\n")

			//f.Write([]byte("M03 L"+strconv.Itoa(255-val)+"\n"))
			if val == 0 {
				f.Write([]byte("M03 L200\n"))
			} else {
				f.Write([]byte("M03 L1\n"))
			}

			/*
			if  val <=50 {
				f.Write([]byte("M03 L5\n"))
			} else {
				f.Write([]byte("M03 L100\n"))
			}
			*/

			f.Write([]byte("G01 X" + strconv.Itoa(1) + " \n"))
			time.Sleep(time.Second)
			//log.Println("G01 Y-" + strconv.Itoa(1) + " F50\n")
		}

	}






	time.Sleep(time.Second/2)

	f.Write([]byte("M03 F50\n"))
	f.Write([]byte("G00 X-"+strconv.Itoa(w)+" Y"+strconv.Itoa(h)+"\n"))
	/*

	//level := "L100"
	for y :=1; y<= 25; y++ {

		//time.Sleep(time.Microsecond*10)

		f.Write([]byte("G01 Y-" + strconv.Itoa(1) + " F50\n"))

		time.Sleep(time.Second/2)

		log.Println("G0 Y-" + strconv.Itoa(1) + "\n")
	}
	*/



	f.Write([]byte("M05\n"))

}
