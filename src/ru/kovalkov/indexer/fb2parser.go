package main

import (
    "fmt"
    "strconv"
    "strings"
    "ru/kovalkov/xmltextreader"
    "regexp"
)

const (
    annotation_tag = "annotation"
    author_tag = "author"
    body_tag = "body"
    book_title_tag = "book-title"
    description_tag = "description"
    document_info_tag = "document-info"
    emphasis_tag = "emphasis"
    epigraph_tag = "epigraph"
    fictionbook_tag = "FictionBook"
    first_name_tag = "first-name"
    last_name_tag = "last-name"
    link_tag = "a"
    middle_name_tag = "middle-name"
    nick_name_tag = "nickname"
    paragraph_tag = "p"
    publish_info_tag = "publish-info"
    section_tag = "section"
    src_author_first_name_attr = "src-author-first-name"
    src_author_last_name_attr = "src-author-last-name"
    src_author_middle_name_attr = "src-author-middle-name"
    src_book_title_attr = "src-book-title"
    src_title_info_tag = "src-title-info"
    stanza_tag = "stanza"
    strong_tag = "strong"
    subtitle_tag = "subtitle"
    text_author_tag = "text-author"
    title_info_tag = "title-info"
    title_tag = "title"
    translator_tag = "translator"
    verse_tag = "v"
)

type Word struct {
    text        string
    fb2pointer  string
    weight      float32
}

func processBook(c chan *Word, qchan chan bool) {
    reader, err := xmltextreader.Filename("Panov_V._Ruchnoyi_Privod.fb2")
    if nil != err {
        fmt.Println(err)
        return
    }

    res := reader.Read()
    eventType := reader.NodeType()

    if res == 1 && xmltextreader.XML_START_ELEMENT == eventType && fictionbook_tag == reader.Name() {
        if 1 == processDescription(reader, "/1/1", c) {
            processBody(reader, fictionbook_tag, "/1", c)
        }
    }
    qchan <- true
}

func processBody(reader *xmltextreader.XmlTextReaderPtr, tag string, currentPointer string, c chan *Word) {
    currentElementNum := 0
    res := reader.Read()
    eventType := reader.NodeType()
    if currentPointer == "/1" {
        currentElementNum = 1;
    }

    for ;!(xmltextreader.XML_END_ELEMENT == eventType && tag == reader.Name()) && -1 != res; {
        name := reader.Name()

        if xmltextreader.XML_START_ELEMENT == eventType {
            currentElementNum += 1

            if paragraph_tag == name || text_author_tag == name || subtitle_tag == name || verse_tag == name {
                getParaToTagEnd(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            } else if body_tag == name {
                hasAttr, err := reader.HasAttributes()
                if nil == err && hasAttr == false {
                    processBody(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
                }
            } else if section_tag == name {
                hasAttr, err := reader.HasAttributes()
                if nil == err && hasAttr {
//                  Skip section with id attr, notes definitions
                    getParaToTagEnd(reader, name, "", c)
                } else {
                    processBody(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
                }
            } else if title_tag == name || epigraph_tag == name {
                processBody(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            }
        } else if xmltextreader.XML_END_ELEMENT == eventType {
            name = reader.Name()
            if stanza_tag == name {
                //elementPointer = currentPointer + "/" + strconv.Itoa(currentElementNum)
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
}

func processDescription(reader *xmltextreader.XmlTextReaderPtr, currentPointer string, c chan *Word) int {
    nextTag(reader)
    if reader.Name() != description_tag {
        return -1
    }
    res := reader.Read()
    eventType := reader.NodeType()
    currentElementNum := 0

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && reader.Name() == description_tag) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            currentElementNum += 1
            name := reader.Name()
            if name == title_info_tag || name == src_title_info_tag {
                processTitleInfo(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
    return 1
}

func processTitleInfo(reader *xmltextreader.XmlTextReaderPtr, tag string, currentPointer string, c chan *Word) {
    res := reader.Read()
    eventType := reader.NodeType()
    currentElementNum := 0

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == reader.Name()) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            currentElementNum += 1
            name := reader.Name()
            if author_tag == name {
                processPersonInfo(reader, name, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            } else if book_title_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 2.0, c)
            } else if annotation_tag == name {
                getParagraphsToTagEnd(reader, annotation_tag, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            } else if translator_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 1.1, c)
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
}

func processPersonInfo(reader *xmltextreader.XmlTextReaderPtr, tag string, currentPointer string, c chan *Word) {
    res := reader.Read()
    eventType := reader.NodeType()
    currentElementNum := 0

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == reader.Name()) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            currentElementNum += 1
            name := reader.Name()
            if first_name_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 1.3, c)
            } else if last_name_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 1.3, c)
            } else if middle_name_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 1.3, c)
            } else if nick_name_tag == name {
                sendWords(getText(reader), currentPointer + "/" + strconv.Itoa(currentElementNum), 1.1, c)
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
}

func getParagraphsToTagEnd(reader *xmltextreader.XmlTextReaderPtr, tag string, currentPointer string, c chan *Word) {
    res := reader.Read()
    eventType := reader.NodeType()
    currentElementNum := 0

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == reader.Name()) && res != -1; {
        if xmltextreader.XML_START_ELEMENT == eventType {
            currentElementNum += 1
            name := reader.Name()
            if name == paragraph_tag {
                getParaToTagEnd(reader, paragraph_tag, currentPointer + "/" + strconv.Itoa(currentElementNum), c)
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
}

func getParaToTagEnd(reader *xmltextreader.XmlTextReaderPtr, tag string, fb2Pointer string, c chan *Word) {
    res := reader.Read()
    eventType := reader.NodeType()

    var (
        strong      byte = 1 << 0
        emphasis    byte = 1 << 1
        link        byte = 1 << 2
    )

    var name string
    var styles byte = 0

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == reader.Name()) && res != -1; {
        if xmltextreader.XML_TEXT_NODE == eventType {
            var weight float32 = 1.0
            if styles & strong != 0 {
                weight += 0.2
            } else if styles & emphasis != 0 {
                weight += 0.1
            } else if styles & link != 0 {
                weight += 0.1
            }
            if len(fb2Pointer) > 0 {
                sendWords(reader.Value(), fb2Pointer, weight, c)
            }
        } else if xmltextreader.XML_START_ELEMENT == eventType {
            name = reader.Name()
            if emphasis_tag == name {
                styles |= emphasis
            } else if strong_tag == name {
                styles |= strong
            } else if link_tag == name {
                styles |= link
            }
        } else if xmltextreader.XML_END_ELEMENT == eventType {
            name = reader.Name()
            if emphasis_tag == name {
                styles &^= emphasis
            } else if strong_tag == name {
                styles &^= strong
            } else if link_tag == name {
                styles &^= link
            }
        }
        res = reader.Read()
        eventType = reader.NodeType()
    }
}

func getText(reader *xmltextreader.XmlTextReaderPtr) string {
    res := reader.Read()
    eventType := reader.NodeType()

    for ;xmltextreader.XML_END_ELEMENT != eventType && xmltextreader.XML_TEXT_NODE != eventType && res != -1; {
        res = reader.Read()
        eventType = reader.NodeType()
    }

    if xmltextreader.XML_TEXT_NODE == eventType {
        return reader.Value()
    }
    return ""
}

func sendWords(text string, pointer string, weight float32, c chan *Word) {
	patrn, _ := regexp.Compile("(\\pL+|\\pN+)")
    for _, w := range patrn.FindAllString(text, -1) {
        if len([]rune(w)) > 1 {
            c <- &Word{strings.ToLower(w), pointer, weight}
        }
    }
}

func nextTag(reader *xmltextreader.XmlTextReaderPtr) int {
    res := reader.Read()
    if res == 1 {
        for nodeType := reader.NodeType();
                nodeType != -1 && nodeType != xmltextreader.XML_START_ELEMENT && res == 1; {
            res = reader.Read()
            nodeType = reader.NodeType()
        }
    }
    return res
}
