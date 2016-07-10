package main

import (
    "github.com/lucasb-eyer/go-colorful"
    "github.com/llgcode/draw2d/draw2dimg"
    "github.com/nfnt/resize"
    "image"
    "image/color"
)

type SimpleRenderer struct {
    size int
}

func NewSimpleRenderer(size int) *SimpleRenderer {
    r := new(SimpleRenderer)
    r.size = size
    return r
}

func (r *SimpleRenderer) Render(code uint32, fileName string) {

    // decode the code into parts
    // bit 0-6   : frame patch type
    // bit 7     : center type
    // bit 8     : cross type
    // bit 9     : turn cross type
    // bit 10    : slash type
    // bit 11-13 : inner type
    // bit 14-18 : fill option
    // bit 19-23 : saturation for HSV
    // bit 24-31 : hue for HSV
    frameType     := code & 0x7F
    centerType    := ((code >> 7) & 0x1) != 0
    crossType     := ((code >> 8) & 0x1) != 0
    turnCrossType := ((code >> 9) & 0x1) != 0
    slashType     := ((code >> 10) & 0x1) != 0
    innerType     := (code >> 11) & 0x7
    fillOption    := (code >> 14) & 0x1F
    saturation    := (code >> 19) & 0x1F
    hue           := (code >> 24) & 0xFF


    // baseColor := color.RGBA{uint8(red << 3), uint8(green << 3), uint8(blue << 3), 0xFF}
    c := colorful.Hsv(float64(hue),float64(saturation)/0x1F,1)
    rc,gc,bc := c.RGB255()
    baseColor :=color.RGBA{rc,gc,bc,0xFF}
    secondColor, thirdColor := GetTriadColor(baseColor)
    src := image.NewRGBA(image.Rect(0, 0, 7, 7))
    frameFlip := 0 < (fillOption & 0x10)

    //draw frame
    for i := 0; i < 7; i++ {
        if 0 < ((0x1 << uint8(i)) & frameType) {
            src.SetRGBA(i, 0, baseColor)
            src.SetRGBA(6-i, 6, baseColor)
            if frameFlip {
                src.SetRGBA(0, 6-i, baseColor)
                src.SetRGBA(6, i, baseColor)
            }else{
                src.SetRGBA(0, i, baseColor)
                src.SetRGBA(6, 6-i, baseColor)
            }
        }
    }

    centerTypeTable := [][]image.Point{
                                        {image.Pt(3,3)},
                                        {image.Pt(3,2)},
                                        {image.Pt(3,2)},
                                        {image.Pt(3,3)},
                                        {image.Pt(3,3)},
                                        {image.Pt(3,3),image.Pt(3,2),image.Pt(2,3),image.Pt(4,3),image.Pt(3,4)},
                                        {image.Pt(3,3),image.Pt(3,2),image.Pt(2,3),image.Pt(4,3),image.Pt(3,4)},
                                        {image.Pt(2,2),image.Pt(4,2),image.Pt(2,4),image.Pt(4,4)},
                                        }
    var centerColor = baseColor
    if !centerType {
        centerColor = color.RGBA{0x00,0x00,0x00,0x00}
    }
    for i := 0; i < len(centerTypeTable[innerType]); i++ {
        if 0 < (fillOption & 0x1) {
            src.SetRGBA(centerTypeTable[innerType][i].X, centerTypeTable[innerType][i].Y, centerColor)
        }
    }


    crossTypeTable := [][]image.Point{
                                        {image.Pt(3,3),image.Pt(3,1),image.Pt(3,2),image.Pt(3,4),image.Pt(3,5),image.Pt(1,3),image.Pt(2,3),image.Pt(4,3),image.Pt(5,3),},
                                        {image.Pt(2,1),image.Pt(4,1),image.Pt(2,3),image.Pt(4,3),image.Pt(1,4),image.Pt(5,4),},
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(4,2),image.Pt(3,3),image.Pt(2,4),image.Pt(1,5),image.Pt(5,5),},
                                        {image.Pt(2,1),image.Pt(4,1),image.Pt(1,2),image.Pt(3,2),image.Pt(5,2),image.Pt(1,4),image.Pt(3,4),image.Pt(5,4),image.Pt(2,5),image.Pt(4,5),},
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(2,2),image.Pt(4,2),image.Pt(1,3),image.Pt(5,3),image.Pt(1,4),image.Pt(3,4),image.Pt(5,4),image.Pt(2,5),image.Pt(4,5),},
                                        {image.Pt(3,1),image.Pt(1,3),image.Pt(5,3),image.Pt(3,5),},
                                        {image.Pt(2,2),image.Pt(4,2),image.Pt(1,3),image.Pt(5,3),image.Pt(2,4),image.Pt(4,4),},
                                        {image.Pt(3,3),},
                                        }
    var crossColor = secondColor
    if crossType {
        crossColor = thirdColor
    }
    for i := 0; i < len(crossTypeTable[innerType]); i++ {
        if 0 < (fillOption & 0x2) {
            src.SetRGBA(crossTypeTable[innerType][i].X, crossTypeTable[innerType][i].Y, crossColor)
        }
    }

    turnCrossTypeTable := [][]image.Point{
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(2,2),image.Pt(4,2),image.Pt(2,4),image.Pt(4,4),image.Pt(1,5),image.Pt(5,5),},
                                        {image.Pt(3,1),image.Pt(1,2),image.Pt(2,2),image.Pt(4,2),image.Pt(5,2),image.Pt(3,3),image.Pt(3,4),image.Pt(3,5)},
                                        {image.Pt(3,1),image.Pt(2,2),image.Pt(1,3),image.Pt(5,3),image.Pt(4,4),image.Pt(3,5),},
                                        {image.Pt(3,1),image.Pt(2,2),image.Pt(4,2),image.Pt(1,3),image.Pt(5,3),image.Pt(2,4),image.Pt(4,4),image.Pt(3,5),},
                                        {image.Pt(2,1),image.Pt(4,1),image.Pt(1,2),image.Pt(3,2),image.Pt(5,2),image.Pt(2,4),image.Pt(4,4),image.Pt(1,5),image.Pt(5,5),},
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(2,2),image.Pt(4,2),image.Pt(2,4),image.Pt(4,4),image.Pt(1,5),image.Pt(5,5),},
                                        {image.Pt(1,1),image.Pt(2,1),image.Pt(3,1),image.Pt(4,1),image.Pt(5,1),image.Pt(1,5),image.Pt(2,5),image.Pt(3,5),image.Pt(4,5),image.Pt(5,5),},
                                        {image.Pt(3,2),image.Pt(2,3),image.Pt(4,3),image.Pt(3,4),},
                                        }
    var turnCrossColor = secondColor
    if turnCrossType {
        turnCrossColor = thirdColor
    }
    for i := 0; i < len(turnCrossTypeTable[innerType]); i++ {
        if 0 < (fillOption & 0x4) {
            src.SetRGBA(turnCrossTypeTable[innerType][i].X, turnCrossTypeTable[innerType][i].Y, turnCrossColor)
        }
    }

    slashTypeTable := [][]image.Point{
                                        {image.Pt(2,1),image.Pt(4,1),image.Pt(1,2),image.Pt(5,2),image.Pt(1,4),image.Pt(5,4),image.Pt(2,5),image.Pt(4,5),},
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(1,3),image.Pt(5,3),image.Pt(2,4),image.Pt(4,4),image.Pt(1,5),image.Pt(2,5),image.Pt(4,5),image.Pt(5,5)},
                                        {image.Pt(2,1),image.Pt(4,2),image.Pt(1,2),image.Pt(5,2),image.Pt(2,3),image.Pt(4,3),image.Pt(1,4),image.Pt(3,4),image.Pt(5,4),image.Pt(2,5),image.Pt(4,5),},
                                        {image.Pt(1,1),image.Pt(5,1),image.Pt(2,3),image.Pt(4,3),image.Pt(1,5),image.Pt(5,5),},
                                        {image.Pt(3,1),image.Pt(2,3),image.Pt(4,3),image.Pt(3,5),},
                                        {image.Pt(2,1),image.Pt(4,1),image.Pt(1,2),image.Pt(5,2),image.Pt(1,4),image.Pt(5,4),image.Pt(2,5),image.Pt(4,5),},
                                        {image.Pt(1,2),image.Pt(5,2),image.Pt(1,4),image.Pt(5,4),},
                                        {image.Pt(1,1),image.Pt(2,1),image.Pt(3,1),image.Pt(4,1),image.Pt(5,1),image.Pt(1,2),image.Pt(5,2),image.Pt(1,3),image.Pt(5,3),image.Pt(1,4),image.Pt(5,4),image.Pt(1,5),image.Pt(2,5),image.Pt(3,5),image.Pt(4,5),image.Pt(5,5),},
                                        }
    var slashColor = secondColor
    if slashType {
        slashColor = thirdColor
    }
    for i := 0; i < len(slashTypeTable[innerType]); i++ {
        if 0 < (fillOption & 0x8) {
            src.SetRGBA(slashTypeTable[innerType][i].X, slashTypeTable[innerType][i].Y, slashColor)
        }
    }
    dest := resize.Resize(uint(r.size), 0, src, resize.NearestNeighbor)
    draw2dimg.SaveToPngFile(fileName, dest)
}
