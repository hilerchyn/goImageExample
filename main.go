package main

import (
	"syscall"
	"log"
	"os"
	"unsafe"
	"time"
	"strconv"
)



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





	f.Write([]byte("G92\n"))
	time.Sleep(time.Second/2)
	f.Write([]byte("G91\n"))
	time.Sleep(time.Second/2)
	f.Write([]byte("G21\n"))
	time.Sleep(time.Second/2)

	f.Write([]byte("M03 L100\n"))

	time.Sleep(time.Second/2)



	//level := "L100"
	for y :=1; y<= 25; y++ {

		//time.Sleep(time.Microsecond*10)

		f.Write([]byte("G01 Y-" + strconv.Itoa(1) + " F50\n"))

		time.Sleep(time.Second/2)

		log.Println("G0 Y-" + strconv.Itoa(1) + "\n")
	}



	f.Write([]byte("M05\n"))

}
