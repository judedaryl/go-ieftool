package internal

import (
	"fmt"
	"io"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

func configureDecoder() (yqlib.Decoder, error) {
	yqlibInputFormat, err := yqlib.InputFormatFromString(inputFormat)
	if err != nil {
		return nil, err
	}
	switch yqlibInputFormat {
	case yqlib.XMLInputFormat:
		return yqlib.NewXMLDecoder(xmlAttributePrefix, xmlContentName, xmlStrictMode, xmlKeepNamespace, xmlUseRawToken), nil
	case yqlib.PropertiesInputFormat:
		return yqlib.NewPropertiesDecoder(), nil
	}

	return yqlib.NewYamlDecoder(), nil
}

func configurePrinterWriter(format yqlib.PrinterOutputFormat, out io.Writer) (yqlib.PrinterWriter, error) {

	var printerWriter yqlib.PrinterWriter

	if splitFileExp != "" {
		colorsEnabled = forceColor
		splitExp, err := yqlib.ExpressionParser.ParseExpression(splitFileExp)
		if err != nil {
			return nil, fmt.Errorf("bad split document expression: %w", err)
		}
		printerWriter = yqlib.NewMultiPrinterWriter(splitExp, format)
	} else {
		printerWriter = yqlib.NewSinglePrinterWriter(out)
	}
	return printerWriter, nil
}

func configureEncoder(format yqlib.PrinterOutputFormat) yqlib.Encoder {
	switch format {
	case yqlib.JSONOutputFormat:
		return yqlib.NewJONEncoder(indent, colorsEnabled)
	case yqlib.PropsOutputFormat:
		return yqlib.NewPropertiesEncoder(unwrapScalar)
	case yqlib.CSVOutputFormat:
		return yqlib.NewCsvEncoder(',')
	case yqlib.TSVOutputFormat:
		return yqlib.NewCsvEncoder('\t')
	case yqlib.YamlOutputFormat:
		return yqlib.NewYamlEncoder(indent, colorsEnabled, !noDocSeparators, unwrapScalar)
	case yqlib.XMLOutputFormat:
		return yqlib.NewXMLEncoder(indent, xmlAttributePrefix, xmlContentName)
	}
	panic("invalid encoder")
}
