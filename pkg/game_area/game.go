package gamearea

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"pixelbattle/pkg/utils"
)

type PixelBattleArea struct {
	area []byte
	Width int
	Heigth int
}

type PixelBattle struct {
	PixelBattleArea
	file *os.File
}

func NewPixelBattle(width, heigth int, filename string) (*PixelBattle, error) {
	var file *os.File
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		file, err =  os.Create(filename)
		if err != nil{
			return nil, fmt.Errorf("error creating pixel battle file: %s", err)
		}

		game := PixelBattle{
			PixelBattleArea: PixelBattleArea{make([]byte, width * heigth * 3), width, heigth},
			file: file,
		}
		err := game.SaveAll()
		if err != nil{
			return nil, fmt.Errorf("error saving initial area data: %s", err)
		}
		return &game, nil
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR, 0644)
		if err != nil{
			return nil, fmt.Errorf("error opening pixel battle file: %s", err)
		}
		area, err := readArea(file)
		if err != nil{
			return nil, fmt.Errorf("error reading game area data: %s", err)
		}
		file.Seek(0, 0)
		return &PixelBattle{
			PixelBattleArea: *area,
			file: file,
		}, nil
	}

}

func (p *PixelBattle) SetPixel(x, y int, color utils.RGB) error {
	i := y * p.Width + x
	if i > len(p.area) {
		return fmt.Errorf("coordinates %d %d out of bounce", x, y)
	}
	p.area[i] = color.R
	p.area[i + 1] = color.G
	p.area[i + 2] = color.B
	n, err := p.file.WriteAt(p.area[i:i+3], int64(i) + 8)
	fmt.Println(n)
	return err
}

func (p *PixelBattle) GetPixel(x, y int) (*utils.RGB, error) {
	i := y * p.Width + x
	if i > len(p.area) {
		return nil, fmt.Errorf("coordinates %d %d out of bounce", x, y)
	}
	return &utils.RGB{R: p.area[i], G: p.area[i + 1], B: p.area[i + 2]}, nil
}


func (p *PixelBattle) SaveAll() error {
	if err := p.file.Truncate(0); err != nil{
		return err
	}
	if _, err := p.file.Seek(0, 0); err != nil{
		return err
	}
	if err := binary.Write(p.file, binary.LittleEndian, uint32(p.Width)); err != nil{
		return err
	}
	
	if err := binary.Write(p.file, binary.LittleEndian, uint32(p.Heigth)); err != nil{
		return err
	}

	_, err := p.file.Write(p.area)
	return err
}

func (p *PixelBattle) Close() {
	p.file.Close()
}

func readArea(file *os.File) (*PixelBattleArea, error) {
	var width uint32
	err := binary.Read(file, binary.LittleEndian, &width)
	if err != nil {
		return nil, err
	}
	var height uint32
	err = binary.Read(file, binary.LittleEndian, &height)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err !=  nil {
		return nil, err
	}
	size := int(width) * int(height) * 3
	if len(data) != size {
		return nil, fmt.Errorf("invalid count of elements in file array: need %d, got %d", size, len(data))
	} 
	return &PixelBattleArea{data, int(width), int(height)}, nil
}