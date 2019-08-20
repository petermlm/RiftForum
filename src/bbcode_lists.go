// From:
// https://gist.githubusercontent.com/xthexder/44f4b9cec3ed7876780d/raw/3420d5c43d5ebd99f910a2f277ec9184a496d9d0/gistfile1.go
// Note: Use of this list tag implementation requires that `autoCloseTags` be set to true.

// Example use:
// [list]
// [*] Item 1
// [*] Item 2
// [/list]

package main

import "github.com/frustra/bbcode"

func addBBCodeLists(compiler *bbcode.Compiler) {
    compiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
        out := bbcode.NewHTMLTag("")
        out.Name = "ul"
        style := node.GetOpeningTag().Value
        switch style {
        case "a":
            out.Attrs["style"] = "list-style-type: lower-alpha;"
        case "A":
            out.Attrs["style"] = "list-style-type: upper-alpha;"
        case "i":
            out.Attrs["style"] = "list-style-type: lower-roman;"
        case "I":
            out.Attrs["style"] = "list-style-type: upper-roman;"
        case "1":
            out.Attrs["style"] = "list-style-type: decimal;"
        default:
            out.Attrs["style"] = "list-style-type: disc;"
        }

        if len(node.Children) == 0 {
            out.AppendChild(bbcode.NewHTMLTag(""))
        } else {
            node.Info = []*bbcode.HTMLTag{out, out}
            tags := node.Info.([]*bbcode.HTMLTag)
            for _, child := range node.Children {
                curr := tags[1]
                curr.AppendChild(node.Compiler.CompileTree(child))
            }
            if len(tags[1].Children) > 0 {
                last := tags[1].Children[len(tags[1].Children)-1]
                if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
                    last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
                }
            } else {
                tags[1].AppendChild(bbcode.NewHTMLTag(""))
            }
        }
        return out, false
    })

    compiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
        parent := node.Parent
        for parent != nil {
            if parent.ID == bbcode.OPENING_TAG && parent.GetOpeningTag().Name == "list" {
                out := bbcode.NewHTMLTag("")
                out.Name = "li"
                tags := parent.Info.([]*bbcode.HTMLTag)
                if len(tags[1].Children) > 0 {
                    last := tags[1].Children[len(tags[1].Children)-1]
                    if len(last.Children) > 0 && last.Children[len(last.Children)-1].Name == "br" {
                        last.Children[len(last.Children)-1] = bbcode.NewHTMLTag("")
                    }
                } else {
                    tags[1].AppendChild(bbcode.NewHTMLTag(""))
                }
                tags[1] = out
                tags[0].AppendChild(out)

                if len(parent.Children) == 0 {
                    out.AppendChild(bbcode.NewHTMLTag(""))
                } else {
                    for _, child := range node.Children {
                        curr := tags[1]
                        curr.AppendChild(node.Compiler.CompileTree(child))
                    }
                }
                if node.ClosingTag != nil {
                    tag := bbcode.NewHTMLTag(node.ClosingTag.Raw)
                    bbcode.InsertNewlines(tag)
                    out.AppendChild(tag)
                }
                return nil, false
            }
            parent = parent.Parent
        }
        return bbcode.DefaultTagCompiler(node)
    })
}
