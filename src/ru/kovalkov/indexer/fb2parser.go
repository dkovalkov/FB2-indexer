package main

import (
    "fmt"
    "strconv"
    "ru/kovalkov/xmltextreader"
)

const (
    annotation_tag = "annotation"
    author_tag = "author"
    body_tag = "body"
    book_title_tag = "book-title"
    description_tag = "description"
    document_info_tag = "document-info"
    emphasis_tag = "emphasis"
    fictionbook_tag = "FictionBook"
    first_name_tag = "first-name"
    last_name_tag = "last-name"
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

func main() {
    reader, err := xmltextreader.Filename("Panov_V._Ruchnoyi_Privod.fb2")
    if nil != err {
        fmt.Println(err)
        return
    }

    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)

    if res == 1 && xmltextreader.XML_START_ELEMENT == eventType && fictionbook_tag == xmltextreader.Name(reader) {
        if 1 == processDescription(reader) {
            processBody(reader, fictionbook_tag, "", "/1")
        }
    }
}

func processBody(reader xmltextreader.XmlTextReaderPtr, tag string, parentTag string, currentPointer string) int {
    currentElementNum := 0
    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)
    if currentPointer == "/1" {
        currentElementNum = 1;
    }

    for ;!(xmltextreader.XML_END_ELEMENT == eventType && tag == xmltextreader.Name(reader)) && -1 != res;
           res = xmltextreader.Read(reader) {
        elementPointer := currentPointer
        name := xmltextreader.Name(reader)

        if xmltextreader.XML_START_ELEMENT == eventType {
            currentElementNum += 1

            if paragraph_tag == name || text_author_tag == name || subtitle_tag == name || verse_tag == name {
                elementPointer = currentPointer + "/" + strconv.Itoa(currentElementNum)
            }
        } else if xmltextreader.XML_END_ELEMENT == eventType {
            if stanza_tag == name {
                elementPointer = currentPointer + "/" + strconv.Itoa(currentElementNum)
            }
        }
        fmt.Println("elementPointer", elementPointer)
    }
    return 1
}

func processDescription(reader xmltextreader.XmlTextReaderPtr) int {
    nextTag(reader)
    if xmltextreader.Name(reader) != description_tag {
        return -1
    }
    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)
    for ;!(eventType == xmltextreader.XML_END_ELEMENT && xmltextreader.Name(reader) == description_tag) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            name := xmltextreader.Name(reader)
            if name == title_info_tag || name == src_title_info_tag {
                processTitleInfo(reader, name)
            }
        }
        res = xmltextreader.Read(reader)
        eventType = xmltextreader.NodeType(reader)
    }
    return 1
}

func processTitleInfo(reader xmltextreader.XmlTextReaderPtr, tag string) {
    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == xmltextreader.Name(reader)) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            name := xmltextreader.Name(reader)
            if author_tag == name {
                //fmt.Println("author")
                processPersonInfo(reader, name)
            } else if book_title_tag == name {
                fmt.Println("book title", getText(reader))
            } else if annotation_tag == name {
                fmt.Println("annotation", getText(reader))
            } else if translator_tag == name {
                fmt.Println("translator", getText(reader))
            }
        }
        res = xmltextreader.Read(reader)
        eventType = xmltextreader.NodeType(reader)
    }
}

func processPersonInfo(reader xmltextreader.XmlTextReaderPtr, tag string) {
    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)

    for ;!(eventType == xmltextreader.XML_END_ELEMENT && tag == xmltextreader.Name(reader)) && res != -1; {
        if eventType == xmltextreader.XML_START_ELEMENT {
            name := xmltextreader.Name(reader)
            if first_name_tag == name {
                fmt.Println("first name", getText(reader))
            } else if last_name_tag == name {
                fmt.Println("last name", getText(reader))
            } else if middle_name_tag == name {
                fmt.Println("middle name", getText(reader))
            } else if nick_name_tag == name {
                fmt.Println("nick name", getText(reader))
            }
        }
        res = xmltextreader.Read(reader)
        eventType = xmltextreader.NodeType(reader)
    }
}

func getText(reader xmltextreader.XmlTextReaderPtr) string {
    res := xmltextreader.Read(reader)
    eventType := xmltextreader.NodeType(reader)

    for ;eventType != xmltextreader.XML_END_ELEMENT && xmltextreader.XML_TEXT_NODE != eventType && res != -1; {
        res = xmltextreader.Read(reader)
        eventType = xmltextreader.NodeType(reader)
    }

    if xmltextreader.XML_TEXT_NODE == eventType {
        return xmltextreader.Value(reader)
    }
    return ""
}

func nextTag(reader xmltextreader.XmlTextReaderPtr) int {
    res := xmltextreader.Read(reader)
    if res == 1 {
        for nodeType := xmltextreader.NodeType(reader);
                nodeType != -1 && nodeType != xmltextreader.XML_START_ELEMENT && res == 1; {
            res = xmltextreader.Read(reader)
            nodeType = xmltextreader.NodeType(reader)
        }
    }
    return res
}
