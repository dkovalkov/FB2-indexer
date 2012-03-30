package xmltextreader

/*
#cgo CFLAGS: -I/usr/include/libxml2
#cgo LDFLAGS: -lxml2

#include <libxml/xmlreader.h> 

char* xmlChar2C(xmlChar* x) { return (char *) x; } 
*/
import "C"
import (
	"errors"
	"fmt"
)

type Errno int

func (e Errno) Error() string {
	s := errText[e]
	if s == "" {
		return fmt.Sprintf("errno %d", int(e))
	}
	return s
}

var (
	ErrError        error = Errno(1)
	ErrNoAttributes error = Errno(2)
)

var errText = map[Errno]string{
	1: "Error",
	2: "Has no attributes",
}

const (
	XML_START_ELEMENT      = 1
	XML_ATTRIBUTE_NODE     = 2
	XML_TEXT_NODE          = 3
	XML_CDATA_SECTION_NODE = 4
	XML_ENTITY_REF_NODE    = 5
	XML_ENTITY_NODE        = 6
	XML_PI_NODE            = 7
	XML_COMMENT_NODE       = 8
	XML_DOCUMENT_NODE      = 9
	XML_DOCUMENT_TYPE_NODE = 10
	XML_DOCUMENT_FRAG_NODE = 11
	XML_NOTATION_NODE      = 12
	XML_HTML_DOCUMENT_NODE = 13
	XML_DTD_NODE           = 14
	XML_END_ELEMENT        = 15
	XML_ATTRIBUTE_DECL     = 16
	XML_ENTITY_DECL        = 17
	XML_NAMESPACE_DECL     = 18
	XML_XINCLUDE_START     = 19
	XML_XINCLUDE_END       = 20
	XML_DOCB_DOCUMENT_NODE = 21
)

type XmlTextReaderPtr struct {
	Ptr *C.struct_xmlTextReaderPtr
}

func Filename(filename string) (*XmlTextReaderPtr, error) {
	reader := C.xmlNewTextReaderFilename(C.CString(filename))
	if nil == reader {
		return &XmlTextReaderPtr{}, errors.New("Unable to open " + filename)
	}
	return &XmlTextReaderPtr{reader}, nil
}

func (reader *XmlTextReaderPtr) Read() int {
	return int(C.xmlTextReaderRead(reader.Ptr))
}

func (reader *XmlTextReaderPtr) Name() string {
	name := C.GoString(C.xmlChar2C(C.xmlTextReaderName(reader.Ptr)))
	return name
}

func (reader *XmlTextReaderPtr) NodeType() int {
	return int(C.xmlTextReaderNodeType(reader.Ptr))
}

func (reader *XmlTextReaderPtr) Value() string {
	return C.GoString(C.xmlChar2C(C.xmlTextReaderConstValue(reader.Ptr)))
}

func (reader *XmlTextReaderPtr) HasAttributes() (bool, error) {
	hasAttr := int(C.xmlTextReaderHasAttributes(reader.Ptr))
	if hasAttr == -1 {
		return false, ErrNoAttributes
	} else if hasAttr == 1 {
		return true, nil
	}
	return false, nil
}
