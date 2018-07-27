// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var debug = flag.Bool("debug", false, "")

func main() {
	flag.Parse()

	w := new(bytes.Buffer)
	w.WriteString("// generated by \"go run gen.go\". DO NOT EDIT.\n\n" +
		"package draw\n\nimport (\n" +
		"\"image\"\n" +
		"\"image/color\"\n" +
		"\"math\"\n" +
		"\n" +
		"\"github.com/golang/image/math/f64\"\n" +
		")\n")

	gen(w, "nnInterpolator", codeNNScaleLeaf, codeNNTransformLeaf)
	gen(w, "ablInterpolator", codeABLScaleLeaf, codeABLTransformLeaf)
	genKernel(w)

	if *debug {
		os.Stdout.Write(w.Bytes())
		return
	}
	out, err := format.Source(w.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("impl.go", out, 0660); err != nil {
		log.Fatal(err)
	}
}

var (
	// dsTypes are the (dst image type, src image type) pairs to generate
	// scale_DType_SType implementations for. The last element in the slice
	// should be the fallback pair ("Image", "image.Image").
	//
	// TODO: add *image.CMYK src type after Go 1.5 is released.
	// An *image.CMYK is also alwaysOpaque.
	dsTypes = []struct{ dType, sType string }{
		{"*image.RGBA", "*image.Gray"},
		{"*image.RGBA", "*image.NRGBA"},
		{"*image.RGBA", "*image.RGBA"},
		{"*image.RGBA", "*image.YCbCr"},
		{"*image.RGBA", "image.Image"},
		{"Image", "image.Image"},
	}
	dTypes, sTypes  []string
	sTypesForDType  = map[string][]string{}
	subsampleRatios = []string{
		"444",
		"422",
		"420",
		"440",
	}
	ops = []string{"Over", "Src"}
	// alwaysOpaque are those image.Image implementations that are always
	// opaque. For these types, Over is equivalent to the faster Src, in the
	// absence of a source mask.
	alwaysOpaque = map[string]bool{
		"*image.Gray":  true,
		"*image.YCbCr": true,
	}
)

func init() {
	dTypesSeen := map[string]bool{}
	sTypesSeen := map[string]bool{}
	for _, t := range dsTypes {
		if !sTypesSeen[t.sType] {
			sTypesSeen[t.sType] = true
			sTypes = append(sTypes, t.sType)
		}
		if !dTypesSeen[t.dType] {
			dTypesSeen[t.dType] = true
			dTypes = append(dTypes, t.dType)
		}
		sTypesForDType[t.dType] = append(sTypesForDType[t.dType], t.sType)
	}
	sTypesForDType["anyDType"] = sTypes
}

type data struct {
	dType    string
	sType    string
	sratio   string
	receiver string
	op       string
}

func gen(w *bytes.Buffer, receiver string, codes ...string) {
	expn(w, codeRoot, &data{receiver: receiver})
	for _, code := range codes {
		for _, t := range dsTypes {
			for _, op := range ops {
				if op == "Over" && alwaysOpaque[t.sType] {
					continue
				}
				expn(w, code, &data{
					dType:    t.dType,
					sType:    t.sType,
					receiver: receiver,
					op:       op,
				})
			}
		}
	}
}

func genKernel(w *bytes.Buffer) {
	expn(w, codeKernelRoot, &data{})
	for _, sType := range sTypes {
		expn(w, codeKernelScaleLeafX, &data{
			sType: sType,
		})
	}
	for _, dType := range dTypes {
		for _, op := range ops {
			expn(w, codeKernelScaleLeafY, &data{
				dType: dType,
				op:    op,
			})
		}
	}
	for _, t := range dsTypes {
		for _, op := range ops {
			if op == "Over" && alwaysOpaque[t.sType] {
				continue
			}
			expn(w, codeKernelTransformLeaf, &data{
				dType: t.dType,
				sType: t.sType,
				op:    op,
			})
		}
	}
}

func expn(w *bytes.Buffer, code string, d *data) {
	if d.sType == "*image.YCbCr" && d.sratio == "" {
		for _, sratio := range subsampleRatios {
			e := *d
			e.sratio = sratio
			expn(w, code, &e)
		}
		return
	}

	for _, line := range strings.Split(code, "\n") {
		line = expnLine(line, d)
		if line == ";" {
			continue
		}
		fmt.Fprintln(w, line)
	}
}

func expnLine(line string, d *data) string {
	for {
		i := strings.IndexByte(line, '$')
		if i < 0 {
			break
		}
		prefix, s := line[:i], line[i+1:]

		i = len(s)
		for j, c := range s {
			if !('A' <= c && c <= 'Z' || 'a' <= c && c <= 'z') {
				i = j
				break
			}
		}
		dollar, suffix := s[:i], s[i:]

		e := expnDollar(prefix, dollar, suffix, d)
		if e == "" {
			log.Fatalf("couldn't expand %q", line)
		}
		line = e
	}
	return line
}

// expnDollar expands a "$foo" fragment in a line of generated code. It returns
// the empty string if there was a problem. It returns ";" if the generated
// code is a no-op.
func expnDollar(prefix, dollar, suffix string, d *data) string {
	switch dollar {
	case "dType":
		return prefix + d.dType + suffix
	case "dTypeRN":
		return prefix + relName(d.dType) + suffix
	case "sratio":
		return prefix + d.sratio + suffix
	case "sType":
		return prefix + d.sType + suffix
	case "sTypeRN":
		return prefix + relName(d.sType) + suffix
	case "receiver":
		return prefix + d.receiver + suffix
	case "op":
		return prefix + d.op + suffix

	case "switch":
		return expnSwitch("", "", true, suffix)
	case "switchD":
		return expnSwitch("", "", false, suffix)
	case "switchS":
		return expnSwitch("", "anyDType", false, suffix)

	case "preOuter":
		switch d.dType {
		default:
			return ";"
		case "Image":
			s := ""
			if d.sType == "image.Image" {
				s = "srcMask, smp := opts.SrcMask, opts.SrcMaskP\n"
			}
			return s +
				"dstMask, dmp := opts.DstMask, opts.DstMaskP\n" +
				"dstColorRGBA64 := &color.RGBA64{}\n" +
				"dstColor := color.Color(dstColorRGBA64)"
		}

	case "preInner":
		switch d.dType {
		default:
			return ";"
		case "*image.RGBA":
			return "d := " + pixOffset("dst", "dr.Min.X+adr.Min.X", "dr.Min.Y+int(dy)", "*4", "*dst.Stride")
		}

	case "preKernelOuter":
		switch d.sType {
		default:
			return ";"
		case "image.Image":
			return "srcMask, smp := opts.SrcMask, opts.SrcMaskP"
		}

	case "preKernelInner":
		switch d.dType {
		default:
			return ";"
		case "*image.RGBA":
			return "d := " + pixOffset("dst", "dr.Min.X+int(dx)", "dr.Min.Y+adr.Min.Y", "*4", "*dst.Stride")
		}

	case "blend":
		args, _ := splitArgs(suffix)
		if len(args) != 4 {
			return ""
		}
		switch d.sType {
		default:
			return argf(args, ""+
				"$3r = $0*$1r + $2*$3r\n"+
				"$3g = $0*$1g + $2*$3g\n"+
				"$3b = $0*$1b + $2*$3b\n"+
				"$3a = $0*$1a + $2*$3a",
			)
		case "*image.Gray":
			return argf(args, ""+
				"$3r = $0*$1r + $2*$3r",
			)
		case "*image.YCbCr":
			return argf(args, ""+
				"$3r = $0*$1r + $2*$3r\n"+
				"$3g = $0*$1g + $2*$3g\n"+
				"$3b = $0*$1b + $2*$3b",
			)
		}

	case "clampToAlpha":
		if alwaysOpaque[d.sType] {
			return ";"
		}
		// Go uses alpha-premultiplied color. The naive computation can lead to
		// invalid colors, e.g. red > alpha, when some weights are negative.
		return `
			if pr > pa {
				pr = pa
			}
			if pg > pa {
				pg = pa
			}
			if pb > pa {
				pb = pa
			}
		`

	case "convFtou":
		args, _ := splitArgs(suffix)
		if len(args) != 2 {
			return ""
		}

		switch d.sType {
		default:
			return argf(args, ""+
				"$0r := uint32($1r)\n"+
				"$0g := uint32($1g)\n"+
				"$0b := uint32($1b)\n"+
				"$0a := uint32($1a)",
			)
		case "*image.Gray":
			return argf(args, ""+
				"$0r := uint32($1r)",
			)
		case "*image.YCbCr":
			return argf(args, ""+
				"$0r := uint32($1r)\n"+
				"$0g := uint32($1g)\n"+
				"$0b := uint32($1b)",
			)
		}

	case "outputu":
		args, _ := splitArgs(suffix)
		if len(args) != 3 {
			return ""
		}

		switch d.op {
		case "Over":
			switch d.dType {
			default:
				log.Fatalf("bad dType %q", d.dType)
			case "Image":
				return argf(args, ""+
					"qr, qg, qb, qa := dst.At($0, $1).RGBA()\n"+
					"if dstMask != nil {\n"+
					"	_, _, _, ma := dstMask.At(dmp.X + $0, dmp.Y + $1).RGBA()\n"+
					"	$2r = $2r * ma / 0xffff\n"+
					"	$2g = $2g * ma / 0xffff\n"+
					"	$2b = $2b * ma / 0xffff\n"+
					"	$2a = $2a * ma / 0xffff\n"+
					"}\n"+
					"$2a1 := 0xffff - $2a\n"+
					"dstColorRGBA64.R = uint16(qr*$2a1/0xffff + $2r)\n"+
					"dstColorRGBA64.G = uint16(qg*$2a1/0xffff + $2g)\n"+
					"dstColorRGBA64.B = uint16(qb*$2a1/0xffff + $2b)\n"+
					"dstColorRGBA64.A = uint16(qa*$2a1/0xffff + $2a)\n"+
					"dst.Set($0, $1, dstColor)",
				)
			case "*image.RGBA":
				return argf(args, ""+
					"$2a1 := (0xffff - $2a) * 0x101\n"+
					"dst.Pix[d+0] = uint8((uint32(dst.Pix[d+0])*$2a1/0xffff + $2r) >> 8)\n"+
					"dst.Pix[d+1] = uint8((uint32(dst.Pix[d+1])*$2a1/0xffff + $2g) >> 8)\n"+
					"dst.Pix[d+2] = uint8((uint32(dst.Pix[d+2])*$2a1/0xffff + $2b) >> 8)\n"+
					"dst.Pix[d+3] = uint8((uint32(dst.Pix[d+3])*$2a1/0xffff + $2a) >> 8)",
				)
			}

		case "Src":
			switch d.dType {
			default:
				log.Fatalf("bad dType %q", d.dType)
			case "Image":
				return argf(args, ""+
					"if dstMask != nil {\n"+
					"	qr, qg, qb, qa := dst.At($0, $1).RGBA()\n"+
					"	_, _, _, ma := dstMask.At(dmp.X + $0, dmp.Y + $1).RGBA()\n"+
					"	pr = pr * ma / 0xffff\n"+
					"	pg = pg * ma / 0xffff\n"+
					"	pb = pb * ma / 0xffff\n"+
					"	pa = pa * ma / 0xffff\n"+
					"	$2a1 := 0xffff - ma\n"+ // Note that this is ma, not $2a.
					"	dstColorRGBA64.R = uint16(qr*$2a1/0xffff + $2r)\n"+
					"	dstColorRGBA64.G = uint16(qg*$2a1/0xffff + $2g)\n"+
					"	dstColorRGBA64.B = uint16(qb*$2a1/0xffff + $2b)\n"+
					"	dstColorRGBA64.A = uint16(qa*$2a1/0xffff + $2a)\n"+
					"	dst.Set($0, $1, dstColor)\n"+
					"} else {\n"+
					"	dstColorRGBA64.R = uint16($2r)\n"+
					"	dstColorRGBA64.G = uint16($2g)\n"+
					"	dstColorRGBA64.B = uint16($2b)\n"+
					"	dstColorRGBA64.A = uint16($2a)\n"+
					"	dst.Set($0, $1, dstColor)\n"+
					"}",
				)
			case "*image.RGBA":
				switch d.sType {
				default:
					return argf(args, ""+
						"dst.Pix[d+0] = uint8($2r >> 8)\n"+
						"dst.Pix[d+1] = uint8($2g >> 8)\n"+
						"dst.Pix[d+2] = uint8($2b >> 8)\n"+
						"dst.Pix[d+3] = uint8($2a >> 8)",
					)
				case "*image.Gray":
					return argf(args, ""+
						"out := uint8($2r >> 8)\n"+
						"dst.Pix[d+0] = out\n"+
						"dst.Pix[d+1] = out\n"+
						"dst.Pix[d+2] = out\n"+
						"dst.Pix[d+3] = 0xff",
					)
				case "*image.YCbCr":
					return argf(args, ""+
						"dst.Pix[d+0] = uint8($2r >> 8)\n"+
						"dst.Pix[d+1] = uint8($2g >> 8)\n"+
						"dst.Pix[d+2] = uint8($2b >> 8)\n"+
						"dst.Pix[d+3] = 0xff",
					)
				}
			}
		}

	case "outputf":
		args, _ := splitArgs(suffix)
		if len(args) != 5 {
			return ""
		}
		ret := ""

		switch d.op {
		case "Over":
			switch d.dType {
			default:
				log.Fatalf("bad dType %q", d.dType)
			case "Image":
				ret = argf(args, ""+
					"qr, qg, qb, qa := dst.At($0, $1).RGBA()\n"+
					"$3r0 := uint32($2($3r * $4))\n"+
					"$3g0 := uint32($2($3g * $4))\n"+
					"$3b0 := uint32($2($3b * $4))\n"+
					"$3a0 := uint32($2($3a * $4))\n"+
					"if dstMask != nil {\n"+
					"	_, _, _, ma := dstMask.At(dmp.X + $0, dmp.Y + $1).RGBA()\n"+
					"	$3r0 = $3r0 * ma / 0xffff\n"+
					"	$3g0 = $3g0 * ma / 0xffff\n"+
					"	$3b0 = $3b0 * ma / 0xffff\n"+
					"	$3a0 = $3a0 * ma / 0xffff\n"+
					"}\n"+
					"$3a1 := 0xffff - $3a0\n"+
					"dstColorRGBA64.R = uint16(qr*$3a1/0xffff + $3r0)\n"+
					"dstColorRGBA64.G = uint16(qg*$3a1/0xffff + $3g0)\n"+
					"dstColorRGBA64.B = uint16(qb*$3a1/0xffff + $3b0)\n"+
					"dstColorRGBA64.A = uint16(qa*$3a1/0xffff + $3a0)\n"+
					"dst.Set($0, $1, dstColor)",
				)
			case "*image.RGBA":
				ret = argf(args, ""+
					"$3r0 := uint32($2($3r * $4))\n"+
					"$3g0 := uint32($2($3g * $4))\n"+
					"$3b0 := uint32($2($3b * $4))\n"+
					"$3a0 := uint32($2($3a * $4))\n"+
					"$3a1 := (0xffff - uint32($3a0)) * 0x101\n"+
					"dst.Pix[d+0] = uint8((uint32(dst.Pix[d+0])*$3a1/0xffff + $3r0) >> 8)\n"+
					"dst.Pix[d+1] = uint8((uint32(dst.Pix[d+1])*$3a1/0xffff + $3g0) >> 8)\n"+
					"dst.Pix[d+2] = uint8((uint32(dst.Pix[d+2])*$3a1/0xffff + $3b0) >> 8)\n"+
					"dst.Pix[d+3] = uint8((uint32(dst.Pix[d+3])*$3a1/0xffff + $3a0) >> 8)",
				)
			}

		case "Src":
			switch d.dType {
			default:
				log.Fatalf("bad dType %q", d.dType)
			case "Image":
				ret = argf(args, ""+
					"if dstMask != nil {\n"+
					"	qr, qg, qb, qa := dst.At($0, $1).RGBA()\n"+
					"	_, _, _, ma := dstMask.At(dmp.X + $0, dmp.Y + $1).RGBA()\n"+
					"	pr := uint32($2($3r * $4)) * ma / 0xffff\n"+
					"	pg := uint32($2($3g * $4)) * ma / 0xffff\n"+
					"	pb := uint32($2($3b * $4)) * ma / 0xffff\n"+
					"	pa := uint32($2($3a * $4)) * ma / 0xffff\n"+
					"	pa1 := 0xffff - ma\n"+ // Note that this is ma, not pa.
					"	dstColorRGBA64.R = uint16(qr*pa1/0xffff + pr)\n"+
					"	dstColorRGBA64.G = uint16(qg*pa1/0xffff + pg)\n"+
					"	dstColorRGBA64.B = uint16(qb*pa1/0xffff + pb)\n"+
					"	dstColorRGBA64.A = uint16(qa*pa1/0xffff + pa)\n"+
					"	dst.Set($0, $1, dstColor)\n"+
					"} else {\n"+
					"	dstColorRGBA64.R = $2($3r * $4)\n"+
					"	dstColorRGBA64.G = $2($3g * $4)\n"+
					"	dstColorRGBA64.B = $2($3b * $4)\n"+
					"	dstColorRGBA64.A = $2($3a * $4)\n"+
					"	dst.Set($0, $1, dstColor)\n"+
					"}",
				)
			case "*image.RGBA":
				switch d.sType {
				default:
					ret = argf(args, ""+
						"dst.Pix[d+0] = uint8($2($3r * $4) >> 8)\n"+
						"dst.Pix[d+1] = uint8($2($3g * $4) >> 8)\n"+
						"dst.Pix[d+2] = uint8($2($3b * $4) >> 8)\n"+
						"dst.Pix[d+3] = uint8($2($3a * $4) >> 8)",
					)
				case "*image.Gray":
					ret = argf(args, ""+
						"out := uint8($2($3r * $4) >> 8)\n"+
						"dst.Pix[d+0] = out\n"+
						"dst.Pix[d+1] = out\n"+
						"dst.Pix[d+2] = out\n"+
						"dst.Pix[d+3] = 0xff",
					)
				case "*image.YCbCr":
					ret = argf(args, ""+
						"dst.Pix[d+0] = uint8($2($3r * $4) >> 8)\n"+
						"dst.Pix[d+1] = uint8($2($3g * $4) >> 8)\n"+
						"dst.Pix[d+2] = uint8($2($3b * $4) >> 8)\n"+
						"dst.Pix[d+3] = 0xff",
					)
				}
			}
		}

		return strings.Replace(ret, " * 1)", ")", -1)

	case "srcf", "srcu":
		lhs, eqOp := splitEq(prefix)
		if lhs == "" {
			return ""
		}
		args, extra := splitArgs(suffix)
		if len(args) != 2 {
			return ""
		}

		tmp := ""
		if dollar == "srcf" {
			tmp = "u"
		}

		// TODO: there's no need to multiply by 0x101 in the switch below if
		// the next thing we're going to do is shift right by 8.

		buf := new(bytes.Buffer)
		switch d.sType {
		default:
			log.Fatalf("bad sType %q", d.sType)
		case "image.Image":
			fmt.Fprintf(buf, ""+
				"%sr%s, %sg%s, %sb%s, %sa%s := src.At(%s, %s).RGBA()\n",
				lhs, tmp, lhs, tmp, lhs, tmp, lhs, tmp, args[0], args[1],
			)
			if d.dType == "" || d.dType == "Image" {
				fmt.Fprintf(buf, ""+
					"if srcMask != nil {\n"+
					"	_, _, _, ma := srcMask.At(smp.X+%s, smp.Y+%s).RGBA()\n"+
					"	%sr%s = %sr%s * ma / 0xffff\n"+
					"	%sg%s = %sg%s * ma / 0xffff\n"+
					"	%sb%s = %sb%s * ma / 0xffff\n"+
					"	%sa%s = %sa%s * ma / 0xffff\n"+
					"}\n",
					args[0], args[1],
					lhs, tmp, lhs, tmp,
					lhs, tmp, lhs, tmp,
					lhs, tmp, lhs, tmp,
					lhs, tmp, lhs, tmp,
				)
			}
		case "*image.Gray":
			fmt.Fprintf(buf, ""+
				"%si := %s\n"+
				"%sr%s := uint32(src.Pix[%si]) * 0x101\n",
				lhs, pixOffset("src", args[0], args[1], "", "*src.Stride"),
				lhs, tmp, lhs,
			)
		case "*image.NRGBA":
			fmt.Fprintf(buf, ""+
				"%si := %s\n"+
				"%sa%s := uint32(src.Pix[%si+3]) * 0x101\n"+
				"%sr%s := uint32(src.Pix[%si+0]) * %sa%s / 0xff\n"+
				"%sg%s := uint32(src.Pix[%si+1]) * %sa%s / 0xff\n"+
				"%sb%s := uint32(src.Pix[%si+2]) * %sa%s / 0xff\n",
				lhs, pixOffset("src", args[0], args[1], "*4", "*src.Stride"),
				lhs, tmp, lhs,
				lhs, tmp, lhs, lhs, tmp,
				lhs, tmp, lhs, lhs, tmp,
				lhs, tmp, lhs, lhs, tmp,
			)
		case "*image.RGBA":
			fmt.Fprintf(buf, ""+
				"%si := %s\n"+
				"%sr%s := uint32(src.Pix[%si+0]) * 0x101\n"+
				"%sg%s := uint32(src.Pix[%si+1]) * 0x101\n"+
				"%sb%s := uint32(src.Pix[%si+2]) * 0x101\n"+
				"%sa%s := uint32(src.Pix[%si+3]) * 0x101\n",
				lhs, pixOffset("src", args[0], args[1], "*4", "*src.Stride"),
				lhs, tmp, lhs,
				lhs, tmp, lhs,
				lhs, tmp, lhs,
				lhs, tmp, lhs,
			)
		case "*image.YCbCr":
			fmt.Fprintf(buf, ""+
				"%si := %s\n"+
				"%sj := %s\n"+
				"%s\n",
				lhs, pixOffset("src", args[0], args[1], "", "*src.YStride"),
				lhs, cOffset(args[0], args[1], d.sratio),
				ycbcrToRGB(lhs, tmp),
			)
		}

		if dollar == "srcf" {
			switch d.sType {
			default:
				fmt.Fprintf(buf, ""+
					"%sr %s float64(%sru)%s\n"+
					"%sg %s float64(%sgu)%s\n"+
					"%sb %s float64(%sbu)%s\n"+
					"%sa %s float64(%sau)%s\n",
					lhs, eqOp, lhs, extra,
					lhs, eqOp, lhs, extra,
					lhs, eqOp, lhs, extra,
					lhs, eqOp, lhs, extra,
				)
			case "*image.Gray":
				fmt.Fprintf(buf, ""+
					"%sr %s float64(%sru)%s\n",
					lhs, eqOp, lhs, extra,
				)
			case "*image.YCbCr":
				fmt.Fprintf(buf, ""+
					"%sr %s float64(%sru)%s\n"+
					"%sg %s float64(%sgu)%s\n"+
					"%sb %s float64(%sbu)%s\n",
					lhs, eqOp, lhs, extra,
					lhs, eqOp, lhs, extra,
					lhs, eqOp, lhs, extra,
				)
			}
		}

		return strings.TrimSpace(buf.String())

	case "tweakD":
		if d.dType == "*image.RGBA" {
			return "d += dst.Stride"
		}
		return ";"

	case "tweakDx":
		if d.dType == "*image.RGBA" {
			return strings.Replace(prefix, "dx++", "dx, d = dx+1, d+4", 1)
		}
		return prefix

	case "tweakDy":
		if d.dType == "*image.RGBA" {
			return strings.Replace(prefix, "for dy, s", "for _, s", 1)
		}
		return prefix

	case "tweakP":
		switch d.sType {
		case "*image.Gray":
			if strings.HasPrefix(strings.TrimSpace(prefix), "pa * ") {
				return "1,"
			}
			return "pr,"
		case "*image.YCbCr":
			if strings.HasPrefix(strings.TrimSpace(prefix), "pa * ") {
				return "1,"
			}
		}
		return prefix

	case "tweakPr":
		if d.sType == "*image.Gray" {
			return "pr *= s.invTotalWeightFFFF"
		}
		return ";"

	case "tweakVarP":
		switch d.sType {
		case "*image.Gray":
			return strings.Replace(prefix, "var pr, pg, pb, pa", "var pr", 1)
		case "*image.YCbCr":
			return strings.Replace(prefix, "var pr, pg, pb, pa", "var pr, pg, pb", 1)
		}
		return prefix
	}
	return ""
}

func expnSwitch(op, dType string, expandBoth bool, template string) string {
	if op == "" && dType != "anyDType" {
		lines := []string{"switch op {"}
		for _, op = range ops {
			lines = append(lines,
				fmt.Sprintf("case %s:", op),
				expnSwitch(op, dType, expandBoth, template),
			)
		}
		lines = append(lines, "}")
		return strings.Join(lines, "\n")
	}

	switchVar := "dst"
	if dType != "" {
		switchVar = "src"
	}
	lines := []string{fmt.Sprintf("switch %s := %s.(type) {", switchVar, switchVar)}

	fallback, values := "Image", dTypes
	if dType != "" {
		fallback, values = "image.Image", sTypesForDType[dType]
	}
	for _, v := range values {
		if dType != "" {
			// v is the sType. Skip those always-opaque sTypes, where Over is
			// equivalent to Src.
			if op == "Over" && alwaysOpaque[v] {
				continue
			}
		}

		if v == fallback {
			lines = append(lines, "default:")
		} else {
			lines = append(lines, fmt.Sprintf("case %s:", v))
		}

		if dType != "" {
			if v == "*image.YCbCr" {
				lines = append(lines, expnSwitchYCbCr(op, dType, template))
			} else {
				lines = append(lines, expnLine(template, &data{dType: dType, sType: v, op: op}))
			}
		} else if !expandBoth {
			lines = append(lines, expnLine(template, &data{dType: v, op: op}))
		} else {
			lines = append(lines, expnSwitch(op, v, false, template))
		}
	}

	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func expnSwitchYCbCr(op, dType, template string) string {
	lines := []string{
		"switch src.SubsampleRatio {",
		"default:",
		expnLine(template, &data{dType: dType, sType: "image.Image", op: op}),
	}
	for _, sratio := range subsampleRatios {
		lines = append(lines,
			fmt.Sprintf("case image.YCbCrSubsampleRatio%s:", sratio),
			expnLine(template, &data{dType: dType, sType: "*image.YCbCr", sratio: sratio, op: op}),
		)
	}
	lines = append(lines, "}")
	return strings.Join(lines, "\n")
}

func argf(args []string, s string) string {
	if len(args) > 9 {
		panic("too many args")
	}
	for i, a := range args {
		old := fmt.Sprintf("$%d", i)
		s = strings.Replace(s, old, a, -1)
	}
	return s
}

func pixOffset(m, x, y, xstride, ystride string) string {
	return fmt.Sprintf("(%s-%s.Rect.Min.Y)%s + (%s-%s.Rect.Min.X)%s", y, m, ystride, x, m, xstride)
}

func cOffset(x, y, sratio string) string {
	switch sratio {
	case "444":
		return fmt.Sprintf("( %s    - src.Rect.Min.Y  )*src.CStride + ( %s    - src.Rect.Min.X  )", y, x)
	case "422":
		return fmt.Sprintf("( %s    - src.Rect.Min.Y  )*src.CStride + ((%s)/2 - src.Rect.Min.X/2)", y, x)
	case "420":
		return fmt.Sprintf("((%s)/2 - src.Rect.Min.Y/2)*src.CStride + ((%s)/2 - src.Rect.Min.X/2)", y, x)
	case "440":
		return fmt.Sprintf("((%s)/2 - src.Rect.Min.Y/2)*src.CStride + ( %s    - src.Rect.Min.X  )", y, x)
	}
	return fmt.Sprintf("unsupported sratio %q", sratio)
}

func ycbcrToRGB(lhs, tmp string) string {
	s := `
		// This is an inline version of image/color/ycbcr.go's YCbCr.RGBA method.
		$yy1 := int(src.Y[$i]) * 0x10100
		$cb1 := int(src.Cb[$j]) - 128
		$cr1 := int(src.Cr[$j]) - 128
		$r@ := ($yy1 + 91881*$cr1) >> 8
		$g@ := ($yy1 - 22554*$cb1 - 46802*$cr1) >> 8
		$b@ := ($yy1 + 116130*$cb1) >> 8
		if $r@ < 0 {
			$r@ = 0
		} else if $r@ > 0xffff {
			$r@ = 0xffff
		}
		if $g@ < 0 {
			$g@ = 0
		} else if $g@ > 0xffff {
			$g@ = 0xffff
		}
		if $b@ < 0 {
			$b@ = 0
		} else if $b@ > 0xffff {
			$b@ = 0xffff
		}
	`
	s = strings.Replace(s, "$", lhs, -1)
	s = strings.Replace(s, "@", tmp, -1)
	return s
}

func split(s, sep string) (string, string) {
	if i := strings.Index(s, sep); i >= 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+len(sep):])
	}
	return "", ""
}

func splitEq(s string) (lhs, eqOp string) {
	s = strings.TrimSpace(s)
	if lhs, _ = split(s, ":="); lhs != "" {
		return lhs, ":="
	}
	if lhs, _ = split(s, "+="); lhs != "" {
		return lhs, "+="
	}
	return "", ""
}

func splitArgs(s string) (args []string, extra string) {
	s = strings.TrimSpace(s)
	if s == "" || s[0] != '[' {
		return nil, ""
	}
	s = s[1:]

	i := strings.IndexByte(s, ']')
	if i < 0 {
		return nil, ""
	}
	args, extra = strings.Split(s[:i], ","), s[i+1:]
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	return args, extra
}

func relName(s string) string {
	if i := strings.LastIndex(s, "."); i >= 0 {
		return s[i+1:]
	}
	return s
}

const (
	codeRoot = `
		func (z $receiver) Scale(dst Image, dr image.Rectangle, src image.Image, sr image.Rectangle, op Op, opts *Options) {
			// Try to simplify a Scale to a Copy.
			if dr.Size() == sr.Size() {
				Copy(dst, dr.Min, src, sr, op, opts)
				return
			}

			var o Options
			if opts != nil {
				o = *opts
			}

			// adr is the affected destination pixels.
			adr := dst.Bounds().Intersect(dr)
			adr, o.DstMask = clipAffectedDestRect(adr, o.DstMask, o.DstMaskP)
			if adr.Empty() || sr.Empty() {
				return
			}
			// Make adr relative to dr.Min.
			adr = adr.Sub(dr.Min)
			if op == Over && o.SrcMask == nil && opaque(src) {
				op = Src
			}

			// sr is the source pixels. If it extends beyond the src bounds,
			// we cannot use the type-specific fast paths, as they access
			// the Pix fields directly without bounds checking.
			//
			// Similarly, the fast paths assume that the masks are nil.
			if o.DstMask != nil || o.SrcMask != nil || !sr.In(src.Bounds()) {
				switch op {
				case Over:
					z.scale_Image_Image_Over(dst, dr, adr, src, sr, &o)
				case Src:
					z.scale_Image_Image_Src(dst, dr, adr, src, sr, &o)
				}
			} else if _, ok := src.(*image.Uniform); ok {
				Draw(dst, dr, src, src.Bounds().Min, op)
			} else {
				$switch z.scale_$dTypeRN_$sTypeRN$sratio_$op(dst, dr, adr, src, sr, &o)
			}
		}

		func (z $receiver) Transform(dst Image, s2d f64.Aff3, src image.Image, sr image.Rectangle, op Op, opts *Options) {
			// Try to simplify a Transform to a Copy.
			if s2d[0] == 1 && s2d[1] == 0 && s2d[3] == 0 && s2d[4] == 1 {
				dx := int(s2d[2])
				dy := int(s2d[5])
				if float64(dx) == s2d[2] && float64(dy) == s2d[5] {
					Copy(dst, image.Point{X: sr.Min.X + dx, Y: sr.Min.X + dy}, src, sr, op, opts)
					return
				}
			}

			var o Options
			if opts != nil {
				o = *opts
			}

			dr := transformRect(&s2d, &sr)
			// adr is the affected destination pixels.
			adr := dst.Bounds().Intersect(dr)
			adr, o.DstMask = clipAffectedDestRect(adr, o.DstMask, o.DstMaskP)
			if adr.Empty() || sr.Empty() {
				return
			}
			if op == Over && o.SrcMask == nil && opaque(src) {
				op = Src
			}

			d2s := invert(&s2d)
			// bias is a translation of the mapping from dst coordinates to src
			// coordinates such that the latter temporarily have non-negative X
			// and Y coordinates. This allows us to write int(f) instead of
			// int(math.Floor(f)), since "round to zero" and "round down" are
			// equivalent when f >= 0, but the former is much cheaper. The X--
			// and Y-- are because the TransformLeaf methods have a "sx -= 0.5"
			// adjustment.
			bias := transformRect(&d2s, &adr).Min
			bias.X--
			bias.Y--
			d2s[2] -= float64(bias.X)
			d2s[5] -= float64(bias.Y)
			// Make adr relative to dr.Min.
			adr = adr.Sub(dr.Min)
			// sr is the source pixels. If it extends beyond the src bounds,
			// we cannot use the type-specific fast paths, as they access
			// the Pix fields directly without bounds checking.
			//
			// Similarly, the fast paths assume that the masks are nil.
			if o.DstMask != nil || o.SrcMask != nil || !sr.In(src.Bounds()) {
				switch op {
				case Over:
					z.transform_Image_Image_Over(dst, dr, adr, &d2s, src, sr, bias, &o)
				case Src:
					z.transform_Image_Image_Src(dst, dr, adr, &d2s, src, sr, bias, &o)
				}
			} else if u, ok := src.(*image.Uniform); ok {
				transform_Uniform(dst, dr, adr, &d2s, u, sr, bias, op)
			} else {
				$switch z.transform_$dTypeRN_$sTypeRN$sratio_$op(dst, dr, adr, &d2s, src, sr, bias, &o)
			}
		}
	`

	codeNNScaleLeaf = `
		func (nnInterpolator) scale_$dTypeRN_$sTypeRN$sratio_$op(dst $dType, dr, adr image.Rectangle, src $sType, sr image.Rectangle, opts *Options) {
			dw2 := uint64(dr.Dx()) * 2
			dh2 := uint64(dr.Dy()) * 2
			sw := uint64(sr.Dx())
			sh := uint64(sr.Dy())
			$preOuter
			for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
				sy := (2*uint64(dy) + 1) * sh / dh2
				$preInner
				for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ { $tweakDx
					sx := (2*uint64(dx) + 1) * sw / dw2
					p := $srcu[sr.Min.X + int(sx), sr.Min.Y + int(sy)]
					$outputu[dr.Min.X + int(dx), dr.Min.Y + int(dy), p]
				}
			}
		}
	`

	codeNNTransformLeaf = `
		func (nnInterpolator) transform_$dTypeRN_$sTypeRN$sratio_$op(dst $dType, dr, adr image.Rectangle, d2s *f64.Aff3, src $sType, sr image.Rectangle, bias image.Point, opts *Options) {
			$preOuter
			for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
				dyf := float64(dr.Min.Y + int(dy)) + 0.5
				$preInner
				for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ { $tweakDx
					dxf := float64(dr.Min.X + int(dx)) + 0.5
					sx0 := int(d2s[0]*dxf + d2s[1]*dyf + d2s[2]) + bias.X
					sy0 := int(d2s[3]*dxf + d2s[4]*dyf + d2s[5]) + bias.Y
					if !(image.Point{sx0, sy0}).In(sr) {
						continue
					}
					p := $srcu[sx0, sy0]
					$outputu[dr.Min.X + int(dx), dr.Min.Y + int(dy), p]
				}
			}
		}
	`

	codeABLScaleLeaf = `
		func (ablInterpolator) scale_$dTypeRN_$sTypeRN$sratio_$op(dst $dType, dr, adr image.Rectangle, src $sType, sr image.Rectangle, opts *Options) {
			sw := int32(sr.Dx())
			sh := int32(sr.Dy())
			yscale := float64(sh) / float64(dr.Dy())
			xscale := float64(sw) / float64(dr.Dx())
			swMinus1, shMinus1 := sw - 1, sh - 1
			$preOuter

			for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
				sy := (float64(dy)+0.5)*yscale - 0.5
				// If sy < 0, we will clamp sy0 to 0 anyway, so it doesn't matter if
				// we say int32(sy) instead of int32(math.Floor(sy)). Similarly for
				// sx, below.
				sy0 := int32(sy)
				yFrac0 := sy - float64(sy0)
				yFrac1 := 1 - yFrac0
				sy1 := sy0 + 1
				if sy < 0 {
					sy0, sy1 = 0, 0
					yFrac0, yFrac1 = 0, 1
				} else if sy1 > shMinus1 {
					sy0, sy1 = shMinus1, shMinus1
					yFrac0, yFrac1 = 1, 0
				}
				$preInner

				for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ { $tweakDx
					sx := (float64(dx)+0.5)*xscale - 0.5
					sx0 := int32(sx)
					xFrac0 := sx - float64(sx0)
					xFrac1 := 1 - xFrac0
					sx1 := sx0 + 1
					if sx < 0 {
						sx0, sx1 = 0, 0
						xFrac0, xFrac1 = 0, 1
					} else if sx1 > swMinus1 {
						sx0, sx1 = swMinus1, swMinus1
						xFrac0, xFrac1 = 1, 0
					}

					s00 := $srcf[sr.Min.X + int(sx0), sr.Min.Y + int(sy0)]
					s10 := $srcf[sr.Min.X + int(sx1), sr.Min.Y + int(sy0)]
					$blend[xFrac1, s00, xFrac0, s10]
					s01 := $srcf[sr.Min.X + int(sx0), sr.Min.Y + int(sy1)]
					s11 := $srcf[sr.Min.X + int(sx1), sr.Min.Y + int(sy1)]
					$blend[xFrac1, s01, xFrac0, s11]
					$blend[yFrac1, s10, yFrac0, s11]
					$convFtou[p, s11]
					$outputu[dr.Min.X + int(dx), dr.Min.Y + int(dy), p]
				}
			}
		}
	`

	codeABLTransformLeaf = `
		func (ablInterpolator) transform_$dTypeRN_$sTypeRN$sratio_$op(dst $dType, dr, adr image.Rectangle, d2s *f64.Aff3, src $sType, sr image.Rectangle, bias image.Point, opts *Options) {
			$preOuter
			for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
				dyf := float64(dr.Min.Y + int(dy)) + 0.5
				$preInner
				for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ { $tweakDx
					dxf := float64(dr.Min.X + int(dx)) + 0.5
					sx := d2s[0]*dxf + d2s[1]*dyf + d2s[2]
					sy := d2s[3]*dxf + d2s[4]*dyf + d2s[5]
					if !(image.Point{int(sx) + bias.X, int(sy) + bias.Y}).In(sr) {
						continue
					}

					sx -= 0.5
					sx0 := int(sx)
					xFrac0 := sx - float64(sx0)
					xFrac1 := 1 - xFrac0
					sx0 += bias.X
					sx1 := sx0 + 1
					if sx0 < sr.Min.X {
						sx0, sx1 = sr.Min.X, sr.Min.X
						xFrac0, xFrac1 = 0, 1
					} else if sx1 >= sr.Max.X {
						sx0, sx1 = sr.Max.X-1, sr.Max.X-1
						xFrac0, xFrac1 = 1, 0
					}

					sy -= 0.5
					sy0 := int(sy)
					yFrac0 := sy - float64(sy0)
					yFrac1 := 1 - yFrac0
					sy0 += bias.Y
					sy1 := sy0 + 1
					if sy0 < sr.Min.Y {
						sy0, sy1 = sr.Min.Y, sr.Min.Y
						yFrac0, yFrac1 = 0, 1
					} else if sy1 >= sr.Max.Y {
						sy0, sy1 = sr.Max.Y-1, sr.Max.Y-1
						yFrac0, yFrac1 = 1, 0
					}

					s00 := $srcf[sx0, sy0]
					s10 := $srcf[sx1, sy0]
					$blend[xFrac1, s00, xFrac0, s10]
					s01 := $srcf[sx0, sy1]
					s11 := $srcf[sx1, sy1]
					$blend[xFrac1, s01, xFrac0, s11]
					$blend[yFrac1, s10, yFrac0, s11]
					$convFtou[p, s11]
					$outputu[dr.Min.X + int(dx), dr.Min.Y + int(dy), p]
				}
			}
		}
	`

	codeKernelRoot = `
		func (z *kernelScaler) Scale(dst Image, dr image.Rectangle, src image.Image, sr image.Rectangle, op Op, opts *Options) {
			if z.dw != int32(dr.Dx()) || z.dh != int32(dr.Dy()) || z.sw != int32(sr.Dx()) || z.sh != int32(sr.Dy()) {
				z.kernel.Scale(dst, dr, src, sr, op, opts)
				return
			}

			var o Options
			if opts != nil {
				o = *opts
			}

			// adr is the affected destination pixels.
			adr := dst.Bounds().Intersect(dr)
			adr, o.DstMask = clipAffectedDestRect(adr, o.DstMask, o.DstMaskP)
			if adr.Empty() || sr.Empty() {
				return
			}
			// Make adr relative to dr.Min.
			adr = adr.Sub(dr.Min)
			if op == Over && o.SrcMask == nil && opaque(src) {
				op = Src
			}

			if _, ok := src.(*image.Uniform); ok && o.DstMask == nil && o.SrcMask == nil && sr.In(src.Bounds()) {
				Draw(dst, dr, src, src.Bounds().Min, op)
				return
			}

			// Create a temporary buffer:
			// scaleX distributes the source image's columns over the temporary image.
			// scaleY distributes the temporary image's rows over the destination image.
			var tmp [][4]float64
			if z.pool.New != nil {
				tmpp := z.pool.Get().(*[][4]float64)
				defer z.pool.Put(tmpp)
				tmp = *tmpp
			} else {
				tmp = z.makeTmpBuf()
			}

			// sr is the source pixels. If it extends beyond the src bounds,
			// we cannot use the type-specific fast paths, as they access
			// the Pix fields directly without bounds checking.
			//
			// Similarly, the fast paths assume that the masks are nil.
			if o.SrcMask != nil || !sr.In(src.Bounds()) {
				z.scaleX_Image(tmp, src, sr, &o)
			} else {
				$switchS z.scaleX_$sTypeRN$sratio(tmp, src, sr, &o)
			}

			if o.DstMask != nil {
				switch op {
				case Over:
					z.scaleY_Image_Over(dst, dr, adr, tmp, &o)
				case Src:
					z.scaleY_Image_Src(dst, dr, adr, tmp, &o)
				}
			} else {
				$switchD z.scaleY_$dTypeRN_$op(dst, dr, adr, tmp, &o)
			}
		}

		func (q *Kernel) Transform(dst Image, s2d f64.Aff3, src image.Image, sr image.Rectangle, op Op, opts *Options) {
			var o Options
			if opts != nil {
				o = *opts
			}

			dr := transformRect(&s2d, &sr)
			// adr is the affected destination pixels.
			adr := dst.Bounds().Intersect(dr)
			adr, o.DstMask = clipAffectedDestRect(adr, o.DstMask, o.DstMaskP)
			if adr.Empty() || sr.Empty() {
				return
			}
			if op == Over && o.SrcMask == nil && opaque(src) {
				op = Src
			}
			d2s := invert(&s2d)
			// bias is a translation of the mapping from dst coordinates to src
			// coordinates such that the latter temporarily have non-negative X
			// and Y coordinates. This allows us to write int(f) instead of
			// int(math.Floor(f)), since "round to zero" and "round down" are
			// equivalent when f >= 0, but the former is much cheaper. The X--
			// and Y-- are because the TransformLeaf methods have a "sx -= 0.5"
			// adjustment.
			bias := transformRect(&d2s, &adr).Min
			bias.X--
			bias.Y--
			d2s[2] -= float64(bias.X)
			d2s[5] -= float64(bias.Y)
			// Make adr relative to dr.Min.
			adr = adr.Sub(dr.Min)

			if u, ok := src.(*image.Uniform); ok && o.DstMask != nil && o.SrcMask != nil && sr.In(src.Bounds()) {
				transform_Uniform(dst, dr, adr, &d2s, u, sr, bias, op)
				return
			}

			xscale := abs(d2s[0])
			if s := abs(d2s[1]); xscale < s {
				xscale = s
			}
			yscale := abs(d2s[3])
			if s := abs(d2s[4]); yscale < s {
				yscale = s
			}

			// sr is the source pixels. If it extends beyond the src bounds,
			// we cannot use the type-specific fast paths, as they access
			// the Pix fields directly without bounds checking.
			//
			// Similarly, the fast paths assume that the masks are nil.
			if o.DstMask != nil || o.SrcMask != nil || !sr.In(src.Bounds()) {
				switch op {
				case Over:
					q.transform_Image_Image_Over(dst, dr, adr, &d2s, src, sr, bias, xscale, yscale, &o)
				case Src:
					q.transform_Image_Image_Src(dst, dr, adr, &d2s, src, sr, bias, xscale, yscale, &o)
				}
			} else {
				$switch q.transform_$dTypeRN_$sTypeRN$sratio_$op(dst, dr, adr, &d2s, src, sr, bias, xscale, yscale, &o)
			}
		}
	`

	codeKernelScaleLeafX = `
		func (z *kernelScaler) scaleX_$sTypeRN$sratio(tmp [][4]float64, src $sType, sr image.Rectangle, opts *Options) {
			t := 0
			$preKernelOuter
			for y := int32(0); y < z.sh; y++ {
				for _, s := range z.horizontal.sources {
					var pr, pg, pb, pa float64 $tweakVarP
					for _, c := range z.horizontal.contribs[s.i:s.j] {
						p += $srcf[sr.Min.X + int(c.coord), sr.Min.Y + int(y)] * c.weight
					}
					$tweakPr
					tmp[t] = [4]float64{
						pr * s.invTotalWeightFFFF, $tweakP
						pg * s.invTotalWeightFFFF, $tweakP
						pb * s.invTotalWeightFFFF, $tweakP
						pa * s.invTotalWeightFFFF, $tweakP
					}
					t++
				}
			}
		}
	`

	codeKernelScaleLeafY = `
		func (z *kernelScaler) scaleY_$dTypeRN_$op(dst $dType, dr, adr image.Rectangle, tmp [][4]float64, opts *Options) {
			$preOuter
			for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ {
				$preKernelInner
				for dy, s := range z.vertical.sources[adr.Min.Y:adr.Max.Y] { $tweakDy
					var pr, pg, pb, pa float64
					for _, c := range z.vertical.contribs[s.i:s.j] {
						p := &tmp[c.coord*z.dw+dx]
						pr += p[0] * c.weight
						pg += p[1] * c.weight
						pb += p[2] * c.weight
						pa += p[3] * c.weight
					}
					$clampToAlpha
					$outputf[dr.Min.X + int(dx), dr.Min.Y + int(adr.Min.Y + dy), ftou, p, s.invTotalWeight]
					$tweakD
				}
			}
		}
	`

	codeKernelTransformLeaf = `
		func (q *Kernel) transform_$dTypeRN_$sTypeRN$sratio_$op(dst $dType, dr, adr image.Rectangle, d2s *f64.Aff3, src $sType, sr image.Rectangle, bias image.Point, xscale, yscale float64, opts *Options) {
			// When shrinking, broaden the effective kernel support so that we still
			// visit every source pixel.
			xHalfWidth, xKernelArgScale := q.Support, 1.0
			if xscale > 1 {
				xHalfWidth *= xscale
				xKernelArgScale = 1 / xscale
			}
			yHalfWidth, yKernelArgScale := q.Support, 1.0
			if yscale > 1 {
				yHalfWidth *= yscale
				yKernelArgScale = 1 / yscale
			}

			xWeights := make([]float64, 1 + 2*int(math.Ceil(xHalfWidth)))
			yWeights := make([]float64, 1 + 2*int(math.Ceil(yHalfWidth)))

			$preOuter
			for dy := int32(adr.Min.Y); dy < int32(adr.Max.Y); dy++ {
				dyf := float64(dr.Min.Y + int(dy)) + 0.5
				$preInner
				for dx := int32(adr.Min.X); dx < int32(adr.Max.X); dx++ { $tweakDx
					dxf := float64(dr.Min.X + int(dx)) + 0.5
					sx := d2s[0]*dxf + d2s[1]*dyf + d2s[2]
					sy := d2s[3]*dxf + d2s[4]*dyf + d2s[5]
					if !(image.Point{int(sx) + bias.X, int(sy) + bias.Y}).In(sr) {
						continue
					}

					// TODO: adjust the bias so that we can use int(f) instead
					// of math.Floor(f) and math.Ceil(f).
					sx += float64(bias.X)
					sx -= 0.5
					ix := int(math.Floor(sx - xHalfWidth))
					if ix < sr.Min.X {
						ix = sr.Min.X
					}
					jx := int(math.Ceil(sx + xHalfWidth))
					if jx > sr.Max.X {
						jx = sr.Max.X
					}

					totalXWeight := 0.0
					for kx := ix; kx < jx; kx++ {
						xWeight := 0.0
						if t := abs((sx - float64(kx)) * xKernelArgScale); t < q.Support {
							xWeight = q.At(t)
						}
						xWeights[kx - ix] = xWeight
						totalXWeight += xWeight
					}
					for x := range xWeights[:jx-ix] {
						xWeights[x] /= totalXWeight
					}

					sy += float64(bias.Y)
					sy -= 0.5
					iy := int(math.Floor(sy - yHalfWidth))
					if iy < sr.Min.Y {
						iy = sr.Min.Y
					}
					jy := int(math.Ceil(sy + yHalfWidth))
					if jy > sr.Max.Y {
						jy = sr.Max.Y
					}

					totalYWeight := 0.0
					for ky := iy; ky < jy; ky++ {
						yWeight := 0.0
						if t := abs((sy - float64(ky)) * yKernelArgScale); t < q.Support {
							yWeight = q.At(t)
						}
						yWeights[ky - iy] = yWeight
						totalYWeight += yWeight
					}
					for y := range yWeights[:jy-iy] {
						yWeights[y] /= totalYWeight
					}

					var pr, pg, pb, pa float64 $tweakVarP
					for ky := iy; ky < jy; ky++ {
						if yWeight := yWeights[ky - iy]; yWeight != 0 {
							for kx := ix; kx < jx; kx++ {
								if w := xWeights[kx - ix] * yWeight; w != 0 {
									p += $srcf[kx, ky] * w
								}
							}
						}
					}
					$clampToAlpha
					$outputf[dr.Min.X + int(dx), dr.Min.Y + int(dy), fffftou, p, 1]
				}
			}
		}
	`
)
